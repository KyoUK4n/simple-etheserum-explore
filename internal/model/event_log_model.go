package model

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var _ EventLogModel = (*customEventLogModel)(nil)

type (
	// EventLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customEventLogModel.
	EventLogModel interface {
		eventLogModel
		withSession(session sqlx.Session) EventLogModel
		List(ctx context.Context, txHash string, address string, pageIndex int, pageSize int) ([]EventLog, uint64, error)
		BulkInsert(ctx context.Context, logs []*EventLog) error
	}

	customEventLogModel struct {
		*defaultEventLogModel
	}
)

// NewEventLogModel returns a model for the database table.
func NewEventLogModel(conn sqlx.SqlConn) EventLogModel {
	return &customEventLogModel{
		defaultEventLogModel: newEventLogModel(conn),
	}
}

func (m *customEventLogModel) withSession(session sqlx.Session) EventLogModel {
	return NewEventLogModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customEventLogModel) List(ctx context.Context, txHash string, address string, pageIndex int, pageSize int) ([]EventLog, uint64, error) {
	var eventLogs []EventLog

	where := ""
	var args []any
	if txHash != "" {
		where += " and `tx_hash` = ? "
		args = append(args, txHash)
	}
	if address != "" {
		where += " and `address` = ? "
		args = append(args, address)
	}

	// 查总数
	var total uint64
	countQuery := fmt.Sprintf("select count(1) from %s where 1 = 1 %s", m.table, where)
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	args = append(args, (pageIndex-1)*pageSize, pageSize)
	query := fmt.Sprintf("select %s from %s where 1 = 1 %s limit ?,?", eventLogRows, m.table, where)
	err = m.conn.QueryRowsCtx(ctx, &eventLogs, query, args...)
	switch err {
	case nil:
		return eventLogs, total, nil
	case sqlx.ErrNotFound:
		return nil, 0, ErrNotFound
	default:
		return nil, 0, err
	}
}

func (m *customEventLogModel) getVariables() string {
	eventLogRowsExpect := stringx.Remove(eventLogFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`")
	n := len(eventLogRowsExpect)
	if n <= 0 {
		return "()"
	}
	return "(" + strings.Repeat("?, ", n-1) + "?)"
}

func (m *customEventLogModel) BulkInsert(ctx context.Context, logs []*EventLog) error {

	insertStmt := fmt.Sprintf("insert ignore into %s (%s) values %s", m.table, eventLogRowsExpectAutoSet, m.getVariables())

	bulkInserter, err := sqlx.NewBulkInserter(m.conn, insertStmt)
	if err != nil {
		return err
	}
	for _, log := range logs {
		err = bulkInserter.Insert(
			log.TxHash,
			log.Address,
			log.EventName,
			log.BlockNumber,
			log.LogIndex,
			log.Topics.String,
			log.Data.String,
			log.TxTimestamp.Format(time.DateTime),
		)
		if err != nil {
			return err
		}
	}
	bulkInserter.Flush()
	return nil
}
