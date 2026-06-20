package model

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var _ TransactionModel = (*customTransactionModel)(nil)

type (
	// TransactionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTransactionModel.
	TransactionModel interface {
		transactionModel
		withSession(session sqlx.Session) TransactionModel
		BulkInsert(ctx context.Context, transactions []*Transaction) error
		ListOfAccount(ctx context.Context, address string, pageSize int, pageIndex int) ([]*Transaction, error)
	}

	customTransactionModel struct {
		*defaultTransactionModel
	}
)

// NewTransactionModel returns a model for the database table.
func NewTransactionModel(conn sqlx.SqlConn) TransactionModel {
	return &customTransactionModel{
		defaultTransactionModel: newTransactionModel(conn),
	}
}

func (m *customTransactionModel) withSession(session sqlx.Session) TransactionModel {
	return NewTransactionModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customTransactionModel) BulkInsert(ctx context.Context, transactions []*Transaction) error {
	insertStmt := fmt.Sprintf("insert ignore into %s (%s) values %s", m.table, transactionRowsExpectAutoSet, m.getVariables())

	bulkInserter, err := sqlx.NewBulkInserter(m.conn, insertStmt)
	if err != nil {
		return err
	}
	for _, tx := range transactions {
		err = bulkInserter.Insert(
			tx.TxHash,
			tx.TxIndex,
			tx.From,
			tx.To,
			tx.Value,
			tx.GasLimit,
			tx.GasPrice,
			tx.GasUsed,
			tx.Nonce,
			tx.Status,
			tx.BlockNumber,
			tx.BlockHash,
			tx.TxTimestamp.Time.Format(time.DateTime),
		)
		if err != nil {
			return err
		}
	}
	bulkInserter.Flush()
	return nil
}

func (m *customTransactionModel) getVariables() string {
	transactionRowsExpect := stringx.Remove(transactionFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`")
	n := len(transactionRowsExpect)
	if n <= 0 {
		return "()"
	}
	return "(" + strings.Repeat("?, ", n-1) + "?)"
}

func (m *customTransactionModel) ListOfAccount(ctx context.Context, address string, pageSize int, pageIndex int) ([]*Transaction, error) {
	stmt := fmt.Sprintf("select %s from %s where `from` = ? or `to` = ? limit ?,?", transactionRows, m.table)
	var resp []*Transaction
	err := m.conn.QueryRowsCtx(ctx, &resp, stmt, address, address, (pageIndex-1)*pageSize, pageSize)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
