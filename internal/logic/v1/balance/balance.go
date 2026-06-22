package balance

import (
	"context"

	"github.com/KyoUK4n/etherscan/internal/svc"
	"github.com/KyoUK4n/etherscan/internal/types"
	contractsv1 "github.com/KyoUK4n/etherscan/pkg/contracts/v1"
	contractsv2 "github.com/KyoUK4n/etherscan/pkg/contracts/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
)

type balanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func newBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *balanceLogic {
	return &balanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *balanceLogic) getERC20BalanceV1(address string, tokenAddress string) (*types.BalanceInfo, error) {
	erc20, err := contractsv1.NewERC20(common.HexToAddress(tokenAddress), l.svcCtx.EthClient)
	if err != nil {
		return nil, err
	}
	callOpts := &bind.CallOpts{Context: l.ctx}
	balance, err := erc20.BalanceOf(callOpts, common.HexToAddress(address))
	if err != nil {
		return nil, err
	}
	decimals, err := erc20.Decimals(callOpts)
	if err != nil {
		return nil, err
	}
	name, err := erc20.Name(callOpts)
	if err != nil {
		return nil, err
	}
	symbol, err := erc20.Symbol(callOpts)
	if err != nil {
		return nil, err
	}
	return &types.BalanceInfo{
		Amount:   balance.String(),
		Decimals: int(decimals),
		Address:  tokenAddress,
		Name:     name,
		Symbol:   symbol,
	}, nil
}

func (l *balanceLogic) getERC20BalanceV2(address string, tokenAddress string) (*types.BalanceInfo, error) {
	erc20 := contractsv2.NewERC20()
	erc20Instance := erc20.Instance(l.svcCtx.EthClient, common.HexToAddress(tokenAddress))
	res := make([]any, 0)
	err := erc20Instance.Call(&bind.CallOpts{Context: l.ctx}, &res, "balanceOf", address)
	if err != nil {
		return nil, err
	}
	err = erc20Instance.Call(&bind.CallOpts{Context: l.ctx}, &res, "decimals", address)
	if err != nil {
		return nil, err
	}
	// TODO: 解析返回数据
	return &types.BalanceInfo{
		Amount:   "",
		Decimals: 0,
		Address:  "",
	}, nil
}
