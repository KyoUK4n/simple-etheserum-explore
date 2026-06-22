// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package transaction

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/KyoUK4n/etherscan/internal/logic"
	"github.com/KyoUK4n/etherscan/internal/svc"
	"github.com/KyoUK4n/etherscan/internal/types"
	"github.com/KyoUK4n/etherscan/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryEventLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryEventLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryEventLogLogic {
	return &QueryEventLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryEventLogLogic) QueryEventLog(req *types.QueryEventLogReq) (*types.Response, error) {

	req.Address = strings.TrimSpace(req.Address)
	req.TxHash = strings.TrimSpace(req.TxHash)

	eventLogs, total, err := l.svcCtx.EventLogModel.List(l.ctx, req.TxHash, req.Address, req.PageIndex, req.PageSize)
	if err != nil {
		return logic.OutFailedWithErr(err, "query event logs failed")
	}

	res := make([]*types.EventLog, len(eventLogs))
	for i, eventLog := range eventLogs {
		e := &types.EventLog{
			TxHash:      eventLog.TxHash,
			Address:     eventLog.Address,
			EventName:   eventLog.EventName,
			BlockNumber: eventLog.BlockNumber,
			LogIndex:    eventLog.LogIndex,
			TxTimestamp: eventLog.TxTimestamp.Format(time.DateTime),
		}

		topics := make([]types.EventLogData, 0)
		if err = json.Unmarshal(utils.StringToBytes(eventLog.Topics.String), &topics); err != nil {
			logx.Errorf("unmarshal topics failed: %v\n%s", err, eventLog.Topics.String)
		}
		e.Topics = topics

		data := make([]types.EventLogData, 0)
		if err = json.Unmarshal(utils.StringToBytes(eventLog.Data.String), &data); err != nil {
			logx.Errorf("unmarshal data failed: %v\n%s", err, eventLog.Data.String)
		}
		e.Data = data
		res[i] = e
	}

	return logic.OutSuccess(map[string]any{
		"total":     total,
		"list":      res,
		"pageIndex": req.PageIndex,
		"pageSize":  req.PageSize,
	})
}
