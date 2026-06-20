package svc

import (
	"context"
	"database/sql"
	"math/big"
	"sync"
	"time"

	"github.com/KyoUK4n/etherscan/internal/model"
	"github.com/KyoUK4n/etherscan/pkg/ethx"
	"github.com/cenkalti/backoff/v5"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/zeromicro/go-zero/core/logx"
)

// Scanner 批量拉取区块信息，并将区块内交易、事件同步到DB中
type Scanner struct {
	client           *ethx.EthClient // ETH客户端
	abis             map[string]*abi.ABI // 合约地址->abi映射，用于交易、事件解析
	addresses        []common.Address // 要监听的合约地址，为空则监听全网
	step             uint64           // 每次批量扫描的区块步长（太大会被 RPC 节点拒绝，推荐 1000-5000）
	eventLogModel    model.EventLogModel // 交易事件的DB model实例
	transactionModel model.TransactionModel// 交易的DB model实例
	wg               sync.WaitGroup
	lock             sync.Mutex
	finished         bool // 控制拉取任务串行执行
}

func NewScanner(client *ethx.EthClient, eventLogModel model.EventLogModel, transactionModel model.TransactionModel, subs []*ethx.SubscriptionConfig, step uint64) *Scanner {

	scanner := &Scanner{
		client:           client,
		step:             step,
		abis:             make(map[string]*abi.ABI),
		eventLogModel:    eventLogModel,
		transactionModel: transactionModel,
		wg:               sync.WaitGroup{},
		lock:             sync.Mutex{},
	}

	addrs := make([]common.Address, len(subs))
	for i, sub := range subs {
		addrs[i] = common.HexToAddress(sub.Address)
		parseABI, err := ethx.ParseABI(sub.AbiPath)
		if err != nil {
			logx.Errorf("ABI %s parse error: %s", sub.AbiPath, err)
			continue
		}
		scanner.abis[sub.Address] = parseABI
	}
	scanner.addresses = addrs

	return scanner
}

func (s *Scanner) IsFinished() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.finished
}

// ReplayRange 回放扫描核心函数，扫描[start - end]区块的交易、事件
// 如果start & end未传递，或不合格，将默认扫描[latest - replayCount]区间的区块
// TODO:
// 1. 记录扫描过的区块，自动自增区块号，每次启动都接着上一次扫描进度
// 2. 记录扫描失败的区块，每次启动自动重试
func (s *Scanner) ReplayRange(ctx context.Context, startBlock uint64, endBlock uint64, replayCount uint64) {

	s.lock.Lock()
	s.finished = false
	s.lock.Unlock()
	defer func() {
		s.lock.Lock()
		s.finished = true
		s.lock.Unlock()
	}()

	if replayCount <= 0 {
		replayCount = 1000
	}

	if startBlock > endBlock || startBlock <= 0 || endBlock <= 0 {
		logx.Info("⏳ fetching latest block...")
		// 先获取latest区块
		latestBlock, err := backoff.Retry(ctx, func() (*ethtypes.Block, error) {
			latestBlock, err := s.client.BlockByNumber(ctx, big.NewInt(int64(rpc.LatestBlockNumber)))
			if err != nil {
				logx.Errorf("attempting to get latest block err: %v", err)
				return nil, err
			}
			return latestBlock, nil
		}, backoff.WithMaxElapsedTime(time.Minute*5))
		if err != nil {
			logx.Errorf("🛑 get latest block failed, replay canceled: %v", err)
			return
		}

		endBlock = latestBlock.Number().Uint64()
		if startBlock <= 0 {
			startBlock = endBlock - replayCount
		}
		if startBlock < 0 {
			startBlock = 0
		}
	}

	logx.Infof("🚀 start replay blocks: [%d] ---> [%d]", startBlock, endBlock)

	current := startBlock
	for current <= endBlock {
		select {
		case <-ctx.Done():
			logx.Error("🛑 scanning canceled")
			return
		default:
		}

		// 计算当前批次的结束高度
		toBlock := current + s.step - 1
		if toBlock > endBlock {
			toBlock = endBlock
		}

		logx.Infof("⏳ scanning blocks: [%d - %d]...", current, toBlock)

		// 执行扫描（带重试机制）
		s.scanBatchWithRetry(ctx, current, toBlock)

		// 步进到下一批次
		current = toBlock + 1
	}

	logx.Info("🎉 history block replay finished")
}

// scanBatchWithRetry 带重试的单批次扫描
func (s *Scanner) scanBatchWithRetry(ctx context.Context, from, to uint64) {

	s.wg.Add(2)
	go func() {
		defer s.wg.Done()

		retries := 0
		var maxRetries uint = 5
		exponentialBackOff := backoff.NewExponentialBackOff()
		events, err := backoff.Retry(ctx, func() ([]*model.EventLog, error) {
			events, err := s.scanLogsBatch(ctx, from, to)
			if err == nil {
				return events, nil
			}
			logx.Errorf(
				"⚠️ [%d/%d ,next: %v] scanning event logs at [%d - %d] failed: %v",
				retries+1,
				maxRetries,
				exponentialBackOff.NextBackOff(),
				from,
				to,
				err,
			)
			return nil, err
		}, backoff.WithBackOff(exponentialBackOff), backoff.WithMaxTries(maxRetries))
		if err != nil {
			logx.Errorf(
				"⚠️ scanning event logs at [%d - %d] failed: %v",
				from,
				to,
				err,
			)
			return
		}
		logx.Infof("⭐ got %d event logs [%d - %d]", len(events), from, to)

		if len(events) > 0 {
			err = s.eventLogModel.BulkInsert(ctx, events)
			if err != nil {
				logx.Errorf("❌ bulk insert event logs failed: %v", err)
				return
			}
		}
	}()

	go func() {
		defer s.wg.Done()
		transactions, err := s.scanTransactionsBatch(ctx, from, to)
		if err != nil {
			logx.Errorf(
				"⚠️ scanning transactions at [%d - %d] failed: %v",
				from,
				to,
				err,
			)
			return
		}
		logx.Infof("⭐ got %d transactions [%d - %d]", len(transactions), from, to)
		if len(transactions) > 0 {
			err = s.transactionModel.BulkInsert(ctx, transactions)
			if err != nil {
				logx.Errorf("❌ bulk insert transactions failed: %v", err)
				return
			}
		}

	}()
	s.wg.Wait()
}

func (s *Scanner) isSubscribeAddress(addr string) bool {
	if len(s.addresses) <= 0 {
		return true
	}
	_, ok := s.abis[addr]
	return ok
}

// scanTransactionsBatch 批量获取blocks中的transactions
func (s *Scanner) scanTransactionsBatch(ctx context.Context, from, to uint64) ([]*model.Transaction, error) {
	transactions := make([]*model.Transaction, 0)
	exponentialBackOff := backoff.NewExponentialBackOff()
	chainID, err := s.client.ChainID(ctx)
	if err != nil {
		logx.Errorf("get chain id failed: %v", err)
		return nil, err
	}
	signer := ethtypes.LatestSignerForChainID(chainID)
	for i := from; i <= to; i++ {
		var maxTries uint = 5
		tries := 0
		block, err := backoff.Retry(ctx, func() (*ethtypes.Block, error) {
			tries++
			block, err := s.client.BlockByNumber(ctx, big.NewInt(int64(i)))
			if err != nil {
				logx.Errorf("[%d/%d, next:%v] get block [%d] failed: %v", tries, maxTries, exponentialBackOff.NextBackOff(), i, err)
				return nil, err
			}
			return block, nil
		}, backoff.WithBackOff(exponentialBackOff), backoff.WithMaxTries(maxTries))
		if err != nil {
			logx.Errorf("get block [%d] failed: %v", i, err)
			continue
		}
		for _, tx := range block.Transactions() {
			var toAddr, fromAddr string
			if tx.To() != nil {
				toAddr = tx.To().Hex()
			}
			sender, err := ethtypes.Sender(signer, tx)
			if err == nil {
				fromAddr = sender.Hex()
			} else {
				logx.Errorf("decode sender on [%s] failed: %v", tx.Hash().Hex(), err)
			}

			// 订阅的合约交易才入库
			if s.isSubscribeAddress(fromAddr) || s.isSubscribeAddress(toAddr) {
				transactions = append(transactions, &model.Transaction{
					TxHash:      tx.Hash().Hex(),
					From:        fromAddr,
					To:          toAddr,
					Value:       tx.Value().String(),
					GasLimit:    int64(tx.Gas()),
					GasPrice:    tx.GasPrice().Int64(),
					Nonce:       int64(tx.Nonce()),
					BlockNumber: block.Number().Int64(),
					BlockHash:   block.Hash().Hex(),
					TxTimestamp: sql.NullTime{
						Time:  tx.Time(),
						Valid: true,
					},
				})
			}
		}
	}

	return transactions, nil
}

// scanLogsBatch 核心逻辑：获取 Logs 并处理数据
func (s *Scanner) scanLogsBatch(ctx context.Context, from, to uint64) ([]*model.EventLog, error) {
	// 1. 构建过滤器
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(from)),
		ToBlock:   big.NewInt(int64(to)),
		Addresses: s.addresses,
		// Topics: [][]common.Hash{{common.HexToHash("0xddf252ad...")}}, // 如果只要特定事件（如Transfer），在这里指定
	}

	// 2. 批量拉取 Logs（走 eth_getLogs 接口）
	logs, err := s.client.FilterLogs(ctx, query)
	if err != nil {
		logx.Errorf("filter logs failed: %v", err)
		return nil, err
	}

	events := make([]*model.EventLog, 0)

	// 3. 处理解析出来的事件
	for _, vLog := range logs {
		parsedABI, ok := s.abis[vLog.Address.Hex()]
		if ok && parsedABI != nil {
			eventLog := ethx.ParseLogEvent(&vLog, *parsedABI)
			if eventLog == nil {
				continue
			}
			events = append(events, &model.EventLog{
				TxHash:      eventLog.TxHash,
				Address:     eventLog.Address,
				EventName:   eventLog.EventName,
				BlockNumber: eventLog.BlockNumber,
				LogIndex:    eventLog.LogIndex,
				Topics: sql.NullString{
					String: eventLog.Topics,
					Valid:  true,
				},
				Data: sql.NullString{
					String: eventLog.Data,
					Valid:  true,
				},
				TxTimestamp: eventLog.TxTimestamp,
			})
		} else {
			events = append(events, &model.EventLog{
				TxHash:      vLog.TxHash.Hex(),
				Address:     vLog.Address.Hex(),
				EventName:   "",
				BlockNumber: int64(vLog.BlockNumber),
				LogIndex:    int64(vLog.Index),
				TxTimestamp: time.Unix(int64(vLog.BlockTimestamp), 0),
			})
		}

	}

	return events, nil
}
