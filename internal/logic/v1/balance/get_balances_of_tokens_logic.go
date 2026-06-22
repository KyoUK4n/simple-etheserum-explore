// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package balance

import (
	"context"
	"time"

	"github.com/KyoUK4n/etherscan/internal/logic"
	"github.com/KyoUK4n/etherscan/internal/svc"
	"github.com/KyoUK4n/etherscan/internal/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBalancesOfTokensLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取某个地址的所有代币余额
func NewGetBalancesOfTokensLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBalancesOfTokensLogic {
	return &GetBalancesOfTokensLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBalancesOfTokensLogic) GetBalancesOfTokens(req *types.GetBalanceReq) (*types.Response, error) {

	// eth余额
	ethBalance, err := l.svcCtx.EthClient.BalanceAt(l.ctx, common.HexToAddress(req.Address), nil)
	if err != nil {
		return logic.OutFailedWithErr(err, "get eth balance failed")
	}

	bl := newBalanceLogic(l.ctx, l.svcCtx)
	tokenBalances := make([]*types.BalanceInfo, 0)
	for _, token := range l.svcCtx.Config.Eth.Tokens {
		time.Sleep(l.svcCtx.Config.Eth.RequestPeriod)
		blanceInfo, err := bl.getERC20BalanceV1(req.Address, token)
		if err != nil {
			return logic.OutFailedWithErr(err, "get tokens balance failed")
		}
		tokenBalances = append(tokenBalances, blanceInfo)
	}

	return logic.OutSuccess(&types.TokenBalances{
		EthBalance: &types.BalanceInfo{
			Amount:   ethBalance.String(),
			Decimals: 18,
			Name:     "ETH",
			Symbol:   "ETH",
		},
		TokenBalances: tokenBalances,
	})
}
