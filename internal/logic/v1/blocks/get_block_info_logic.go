// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package blocks

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/KyoUK4n/etherscan/internal/logic"
	"github.com/KyoUK4n/etherscan/internal/svc"
	"github.com/KyoUK4n/etherscan/internal/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBlockInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBlockInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBlockInfoLogic {
	return &GetBlockInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBlockInfoLogic) GetBlockInfo(req *types.GetBlockInfoReq) (*types.Response, error) {
	var blockInfo *types.BlockInfo
	var rawBlock *ethtypes.Block
	var err error

	if req.Number > 0 {
		// 使用block number查询区块信息
		rawBlock, err = l.svcCtx.EthClient.BlockByNumber(l.ctx, big.NewInt(req.Number))
		if err != nil {
			return logic.OutFailedWithErr(err, fmt.Sprintf("fetch block by number [%d] failed", req.Number))
		}
	} else if len(strings.TrimSpace(req.Hash)) > 0 {
		// 使用block hash查询区块信息
		rawBlock, err = l.svcCtx.EthClient.BlockByHash(l.ctx, common.HexToHash(req.Hash))
		if err != nil {
			return logic.OutFailedWithErr(err, fmt.Sprintf("fetch block by hash [%s] failed", req.Hash))
		}
	} else if len(strings.TrimSpace(req.Tag)) > 0 {
		// 通过tag查询特殊区块
		rawBlock, err = l.getBlockByTag(req.Tag)
		if err != nil {
			return logic.OutFailedWithErr(err, fmt.Sprintf("fetch block by tag [%s] failed", req.Tag))
		}
	} else {
		return logic.OutFailedWithErr(err, "at least provided one of [number, hash, tag]")
	}

	if rawBlock == nil {
		return logic.OutFailed("block not found")
	}
	transactions := rawBlock.Transactions()
	blockInfo = &types.BlockInfo{
		Number:           rawBlock.NumberU64(),
		Hash:             rawBlock.Hash().Hex(),
		ParentHash:       rawBlock.ParentHash().Hex(),
		Time:             rawBlock.Time(),
		GasLimit:         rawBlock.GasLimit(),
		GasUsed:          rawBlock.GasUsed(),
		BaseFee:          rawBlock.BaseFee().Uint64(),
		TransactionCount: len(transactions),
	}

	if strings.ToLower(req.Tag) != "finalized" {
		finalizedBlock, err := l.getBlockByTag("finalized")
		if err != nil {
			return logic.OutFailedWithErr(err, "fetch finalized block failed")
		}
		blockInfo.Confirmations = blockInfo.Number - finalizedBlock.NumberU64()
	}

	return logic.OutSuccess(blockInfo)
}

var tagToBlockNumber = map[string]*big.Int{
	"earliest":  big.NewInt(int64(ethrpc.EarliestBlockNumber)),
	"safe":      big.NewInt(int64(ethrpc.SafeBlockNumber)),
	"finalized": big.NewInt(int64(ethrpc.FinalizedBlockNumber)),
	"latest":    big.NewInt(int64(ethrpc.LatestBlockNumber)),
	"pending":   big.NewInt(int64(ethrpc.PendingBlockNumber)),
}

// getBlockByTag tag查询在sdk内部使用特殊的number代表几个特殊区块
// 为了可读性以及API定义，传递string然后与特殊block number映射
func (l *GetBlockInfoLogic) getBlockByTag(tag string) (*ethtypes.Block, error) {
	tagInt, ok := tagToBlockNumber[strings.TrimSpace(strings.ToLower(tag))]
	if !ok {
		return nil, errors.Errorf("unknown tag [%s]", tag)
	}

	rawBlock, err := l.svcCtx.EthClient.BlockByNumber(l.ctx, tagInt)
	if err != nil {
		return nil, err
	}

	return rawBlock, nil
}
