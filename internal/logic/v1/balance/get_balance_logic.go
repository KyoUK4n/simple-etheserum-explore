// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package balance

import (
	"context"

	"github.com/KyoUK4n/etherscan/internal/logic"
	"github.com/KyoUK4n/etherscan/internal/svc"
	"github.com/KyoUK4n/etherscan/internal/types"
	contractsv1 "github.com/KyoUK4n/etherscan/pkg/contracts/v1"
	contractsv2 "github.com/KyoUK4n/etherscan/pkg/contracts/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
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
		return l.getERC20BalanceV1(req)
		//return l.getERC20BalanceV2(req)
	}

	balanceWei, err := l.svcCtx.EthClient.BalanceAt(l.ctx, common.HexToAddress(req.Address), nil)
	if err != nil {
		return logic.OutFailedWithErr(err, "get balance failed")
	}

	return logic.OutSuccess(&types.BalanceInfo{
		Amount:   balanceWei.String(),
		Decimals: 18,
	})
}

func (l *GetBalanceLogic) getERC20BalanceV1(req *types.GetBalanceReq) (*types.Response, error) {
	erc20, err := contractsv1.NewERC20(common.HexToAddress(req.TokenAddress), l.svcCtx.EthClient)
	if err != nil {
		return logic.OutFailedWithErr(err, "get erc20 instance failed")
	}
	callOpts := &bind.CallOpts{Context: l.ctx}
	balance, err := erc20.BalanceOf(callOpts, common.HexToAddress(req.Address))
	if err != nil {
		return logic.OutFailedWithErr(err, "call balanceOf() failed")
	}
	decimals, err := erc20.Decimals(callOpts)
	if err != nil {
		return logic.OutFailedWithErr(err, "call decimals() failed")
	}
	return logic.OutSuccess(&types.BalanceInfo{
		Amount:       balance.String(),
		Decimals:     int(decimals),
		TokenAddress: req.TokenAddress,
	})
}

func (l *GetBalanceLogic) getERC20BalanceV2(req *types.GetBalanceReq) (*types.Response, error) {
	erc20 := contractsv2.NewERC20()
	erc20Instance := erc20.Instance(l.svcCtx.EthClient, common.HexToAddress(req.TokenAddress))
	res := make([]any, 0)
	err := erc20Instance.Call(&bind.CallOpts{Context: l.ctx}, &res, "balanceOf", req.Address)
	if err != nil {
		return logic.OutFailedWithErr(err, "call balanceOf() failed")
	}
	err = erc20Instance.Call(&bind.CallOpts{Context: l.ctx}, &res, "decimals", req.Address)
	if err != nil {
		return logic.OutFailedWithErr(err, "call decimals() failed")
	}
	return logic.OutSuccess(&types.BalanceInfo{
		Amount:       "",
		Decimals:     0,
		TokenAddress: "",
	})
}
