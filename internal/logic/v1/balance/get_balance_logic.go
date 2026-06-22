// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package balance

import (
	"context"

	"github.com/KyoUK4n/etherscan/internal/logic"
	"github.com/KyoUK4n/etherscan/internal/svc"
	"github.com/KyoUK4n/etherscan/internal/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBalanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBalanceLogic {
	return &GetBalanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBalanceLogic) GetBalance(req *types.GetBalanceReq) (*types.Response, error) {

	if req.TokenAddress != "" {
		bl := newBalanceLogic(l.ctx, l.svcCtx)
		blanceInfo, err := bl.getERC20BalanceV1(req.Address, req.TokenAddress)
		if err != nil {
			return logic.OutFailedWithErr(err, "get balance failed")
		}
		return logic.OutSuccess(blanceInfo)
	}

	balanceWei, err := l.svcCtx.EthClient.BalanceAt(l.ctx, common.HexToAddress(req.Address), nil)
	if err != nil {
		return logic.OutFailedWithErr(err, "get balance failed")
	}

	return logic.OutSuccess(&types.BalanceInfo{
		Amount:   balanceWei.String(),
		Decimals: 18,
		Name:     "ETH",
		Symbol:   "ETH",
	})
}
