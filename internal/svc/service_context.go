// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/KyoUK4n/etherscan/internal/config"
	"github.com/KyoUK4n/etherscan/internal/model"
	"github.com/KyoUK4n/etherscan/pkg/ethx"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config             config.Config
	EthClient          *ethx.EthClient
	EventSubscriptions sync.Map
	eventLogsCh        chan types.Log
	EventLogModel      model.EventLogModel
	TransactionModel   model.TransactionModel
	Scanner            *Scanner
}

func NewServiceContext(c config.Config) *ServiceContext {

	// 单次查询复用rpc连接
	ethRpcClient, err := ethx.NewEthClient(c.Eth.RpcUrl)
	if err != nil {
		panic(err)
	}
	if c.Eth.RequestPeriod <= 0 {
		c.Eth.RequestPeriod = 100 * time.Millisecond
	}

	mysqlConn := sqlx.NewMysql(c.Mysql.Datasource)
	eventLogModel := model.NewEventLogModel(mysqlConn)
	transactionModel := model.NewTransactionModel(mysqlConn)
	svcCtx := &ServiceContext{
		Config:             c,
		EthClient:          ethRpcClient,
		EventSubscriptions: sync.Map{},
		eventLogsCh:        make(chan types.Log, 10),
		EventLogModel:      eventLogModel,
		TransactionModel:   transactionModel,
		Scanner:            NewScanner(ethRpcClient, eventLogModel, transactionModel, c.Eth.Subscriptions, c.Eth.ReplayStep),
	}

	go svcCtx.startEventSubscription()
	go svcCtx.handleEventLogs()

	go svcCtx.Scanner.ReplayRange(context.Background(), c.Eth.ReplayStartBlock, c.Eth.ReplayEndBlock, c.Eth.ReplayCount)

	return svcCtx
}

func (c *ServiceContext) Close() {
	c.StopEventSubscription()
	c.EthClient.Close()
}

func (c *ServiceContext) startEventSubscription() {

	for _, subscription := range c.Config.Eth.Subscriptions {

		// 订阅模式需要用WS协议
		ethWSClient, err := ethx.NewEthClient(c.Config.Eth.WSUrl)
		if err != nil {
			logx.Errorf("subscribe on [%s] failed: %s", subscription.Address, err)
			continue
		}

		subInfo, err := ethx.NewSubscribe(ethWSClient, subscription, c.eventLogsCh)
		if err != nil {
			logx.Errorf("subscribe on [%s] failed: %s", subscription.Address, err)
			continue
		}
		logx.Infof("Subscribe on [%s] with abi at [%s]", subscription.Address, subscription.AbiPath)
		c.EventSubscriptions.Store(subscription.Address, subInfo)
	}
}

func (c *ServiceContext) StopEventSubscription() {
	c.EventSubscriptions.Range(func(key, value interface{}) bool {
		subInfo, ok := value.(*ethx.SubInfo)
		if !ok {
			logx.Errorf("invalid subscription [%s] type: %T", key, value)
			return false
		}
		subInfo.Unsubscribe()
		logx.Infof("unsubscribed [%s]", key)
		return true
	})
}

func (c *ServiceContext) handleEventLogs() {
	for {
		select {
		case log, open := <-c.eventLogsCh:
			if !open {
				logx.Errorf("Event logs channel closed")
				return
			}
			address := log.Address.Hex()
			v, ok := c.EventSubscriptions.Load(address)
			if !ok {
				logx.Errorf("subscribe on [%s] not found", address)
				continue
			}
			subInfo, ok := v.(*ethx.SubInfo)
			if !ok {
				logx.Errorf("subscribe on [%s] type assert failed: %T", address, v)
				continue
			}
			eventLog := ethx.ParseLogEvent(&log, *subInfo.ABI)
			_, err := c.EventLogModel.Insert(context.Background(), &model.EventLog{
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
			if err != nil {
				logx.Errorf("insert event log to db failed: %s", err)
			}
		}
	}
}
