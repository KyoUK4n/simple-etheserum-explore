// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package transaction

import (
	"context"

	"github.com/KyoUK4n/etherscan/internal/logic"
	"github.com/KyoUK4n/etherscan/internal/svc"
	"github.com/KyoUK4n/etherscan/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PullFromBlockLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPullFromBlockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PullFromBlockLogic {
	return &PullFromBlockLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PullFromBlockLogic) PullFromBlock(req *types.PullTransactionsReq) (*types.Response, error) {

	if req.EndBlock < req.StartBlock {
		return logic.OutFailed("start block must before end block")
	}

	if !l.svcCtx.Scanner.IsFinished() {
		return logic.OutFailed("previous scanning haven't finished yet")
	}
	go l.svcCtx.Scanner.ReplayRange(context.Background(), req.StartBlock, req.EndBlock, l.svcCtx.Config.Eth.ReplayCount)
	return logic.OutSuccess(nil, "submit scanning task")
}
