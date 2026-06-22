// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package transaction

import (
	"context"
	"math/big"

	"github.com/KyoUK4n/etherscan/internal/logic"
	"github.com/KyoUK4n/etherscan/internal/model"
	"github.com/KyoUK4n/etherscan/internal/svc"
	"github.com/KyoUK4n/etherscan/internal/types"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetTransactionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTransactionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTransactionsLogic {
	return &GetTransactionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTransactionsLogic) GetTransactions(req *types.GetTransactionsReq) (*types.Response, error) {

	txs := make([]*types.Transaction, 0)

	transactions, err := l.svcCtx.TransactionModel.ListOfAccount(l.ctx, req.Address, req.PageSize, req.PageIndex)
	switch {
	case err == nil:

		for _, tx := range transactions {
			value, _ := new(big.Int).SetString(tx.Value, 10)
			txs = append(txs, &types.Transaction{
				Hash:      tx.TxHash,
				From:      tx.From,
				To:        tx.To,
				Value:     value.Uint64(),
				Gas:       uint64(tx.GasLimit),
				GasPrice:  uint64(tx.GasPrice),
				Nonce:     uint64(tx.Nonce),
				Timestamp: uint64(tx.TxTimestamp.Time.UnixMilli() / 1000),
				TransactionReceipt: types.TransactionReceipt{
					Status:      uint64(tx.Status),
					BlockNumber: uint64(tx.BlockNumber),
					BlockHash:   tx.BlockHash,
					TxIndex:     uint64(tx.TxIndex),
					GasUsed:     uint64(tx.GasUsed),
				},
			})
		}

		return logic.OutSuccess(txs)
	case errors.Is(err, model.ErrNotFound):
		return logic.OutSuccess(txs)
	default:
		return logic.OutFailedWithErr(err, "query transactions failed")
	}
}
