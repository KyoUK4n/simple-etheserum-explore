// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package transaction

import (
	"context"

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

	transactions, err := l.svcCtx.TransactionModel.ListOfAccount(l.ctx, req.Address, req.PageSize, req.PageIndex)
	switch {
	case err == nil:
		return logic.OutSuccess(transactions)
	case errors.Is(err, model.ErrNotFound):
		return logic.OutSuccess(transactions)
	default:
		return logic.OutFailedWithErr(err, "query transactions failed")
	}
}
