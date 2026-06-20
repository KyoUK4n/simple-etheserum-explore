// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package blocks

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/KyoUK4n/etherscan/internal/logic"
	"github.com/KyoUK4n/etherscan/internal/svc"
	"github.com/KyoUK4n/etherscan/internal/types"
	"github.com/cenkalti/backoff/v5"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetBlocksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBlocksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBlocksLogic {
	return &GetBlocksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBlocksLogic) GetBlocks(req *types.GetBlockInfoReq) (*types.Response, error) {

	if req.Start <= 0 || req.End <= 0 {
		return logic.OutFailed("invalid start or end number")
	}
	blockRange := req.End - req.Start
	if blockRange > l.svcCtx.Config.Eth.MaxBlockRange {
		return logic.OutFailed(fmt.Sprintf("query range too large: %d vs %d", l.svcCtx.Config.Eth.MaxBlockRange, blockRange))
	}
	if blockRange <= 0 {
		return logic.OutFailed("negative block range")
	}

	blocks, err := l.fetchBlocks(req)
	if err != nil {
		return logic.OutFailedWithErr(err, fmt.Sprintf("fetch blocks [%d - %d] failed", req.Start, req.End))
	}

	return logic.OutSuccess(blocks)
}

func (l *GetBlocksLogic) fetchBlocks(req *types.GetBlockInfoReq) (*types.Response, error) {

	var finalizedBlockNumber uint64
	finalizedBlock, err := l.svcCtx.EthClient.BlockByNumber(l.ctx, big.NewInt(int64(ethrpc.FinalizedBlockNumber)))
	if err != nil {
		logx.Error("fetch finalized block failed")
	} else {
		finalizedBlockNumber = finalizedBlock.NumberU64()
	}
	blocks := make([]*types.BlockInfo, 0)

	exponentialBackOff := backoff.NewExponentialBackOff()
	exponentialBackOff.InitialInterval = 200 * time.Millisecond

	for number := req.Start; number <= req.End; number++ {
		exponentialBackOff.Reset()
		// 限速
		time.Sleep(l.svcCtx.Config.Eth.RequestPeriod)
		rawBlock, err := backoff.Retry(l.ctx, func() (*ethtypes.Block, error) {
			rawBlock, err := l.svcCtx.EthClient.BlockByNumber(l.ctx, big.NewInt(number))
			if err != nil {
				logx.Errorf("attempting to fetch block %d failed: %v", number, err)
				return nil, err
			}
			return rawBlock, nil
		}, backoff.WithBackOff(exponentialBackOff), backoff.WithMaxElapsedTime(time.Minute))
		if err != nil {
			logx.Errorf("fetch block %d failed: %s", number, err)
			continue
		}

		var confirmations uint64
		if finalizedBlockNumber > 0 {
			confirmations = rawBlock.NumberU64() - finalizedBlockNumber
		}

		transactions := rawBlock.Transactions()
		blocks = append(blocks, &types.BlockInfo{
			Number:           rawBlock.NumberU64(),
			Hash:             rawBlock.Hash().Hex(),
			ParentHash:       rawBlock.ParentHash().Hex(),
			Time:             rawBlock.Time(),
			GasLimit:         rawBlock.GasLimit(),
			GasUsed:          rawBlock.GasUsed(),
			BaseFee:          rawBlock.BaseFee().Uint64(),
			Confirmations:    confirmations,
			TransactionCount: len(transactions),
		})
	}

	return logic.OutSuccess(blocks)
}
