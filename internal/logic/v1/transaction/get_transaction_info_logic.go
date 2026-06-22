// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package transaction

import (
	"context"
	"strings"

	"github.com/KyoUK4n/etherscan/internal/logic"
	"github.com/KyoUK4n/etherscan/internal/svc"
	"github.com/KyoUK4n/etherscan/internal/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetTransactionInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTransactionInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTransactionInfoLogic {
	return &GetTransactionInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTransactionInfoLogic) GetTransactionInfo(req *types.GetTransactionInfoReq) (*types.Response, error) {
	req.Hash = strings.TrimSpace(req.Hash)
	if len(req.Hash) <= 0 {
		return logic.OutFailed("transaction hash is required")
	}

	txHash := common.HexToHash(req.Hash)
	tx, isPending, err := l.svcCtx.EthClient.TransactionByHash(l.ctx, txHash)
	if err != nil {
		return logic.OutFailedWithErr(err, "query transaction failed")
	}

	txInfo := &types.Transaction{
		Hash:      tx.Hash().Hex(),
		To:        tx.To().Hex(),
		Value:     tx.Value().Uint64(),
		DataLen:   len(tx.Data()),
		Gas:       tx.Gas(),
		GasPrice:  tx.GasPrice().Uint64(),
		Nonce:     tx.Nonce(),
		IsPending: isPending,
		Timestamp: uint64(tx.Time().UnixMilli() / 1000),
	}

	// 通过signer获取sender
	chainID, err := l.svcCtx.EthClient.ChainID(l.ctx)
	if err != nil {
		logx.Errorf("get chain id failed: %v", err)
	} else {
		signer := ethtypes.LatestSignerForChainID(chainID)

		sender, err := ethtypes.Sender(signer, tx)
		if err == nil {
			txInfo.From = sender.Hex()
		} else {
			logx.Errorf("decode sender on [%s] failed: %v", tx.Hash().Hex(), err)
		}
	}

	if !isPending {
		// 已结束的交易，补充信息
		receipt, err := l.svcCtx.EthClient.TransactionReceipt(l.ctx, txHash)
		if err != nil {
			return logic.OutFailedWithErr(err, "get transaction receipt failed")
		}
		txInfo.Status = receipt.Status
		txInfo.BlockNumber = receipt.BlockNumber.Uint64()
		txInfo.BlockHash = receipt.BlockHash.Hex()
		txInfo.TxIndex = uint64(receipt.TransactionIndex)
		txInfo.GasUsed = receipt.GasUsed
		txInfo.Logs = uint64(len(receipt.Logs))

		// 通过区块头获取交易时间
		blockHeader, err := l.svcCtx.EthClient.HeaderByNumber(l.ctx, receipt.BlockNumber)
		if err != nil {
			logx.Errorf("get block header failed: %v", err)
		} else {
			txInfo.Timestamp = blockHeader.Time
		}

	}

	return logic.OutSuccess(txInfo)
}
