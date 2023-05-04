package eventstore

import (
	"context"
	"database/sql"
)

type dbTx interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func newQueries(db dbTx) *queries {
	return &queries{db: db}
}

type queries struct {
	db dbTx
}

func (q *queries) WithTx(tx *sql.Tx) *queries {
	return &queries{
		db: tx,
	}
}
