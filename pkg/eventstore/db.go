package eventstore

import (
	"context"
	"database/sql"

	"github.com/dmitrymomot/go-utils"
)

type dbTx interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func newQueries(db dbTx, eventStreamName string) *queries {
	eventStreamName = utils.ToSnakeCase(eventStreamName)
	return &queries{
		db:                db,
		snapshotTableName: eventStreamName + "_snapshots",
		eventTableName:    eventStreamName + "_events",
	}
}

type queries struct {
	db                dbTx
	snapshotTableName string
	eventTableName    string
}

func (q *queries) WithTx(tx *sql.Tx) *queries {
	return &queries{
		db:                tx,
		snapshotTableName: q.snapshotTableName,
		eventTableName:    q.eventTableName,
	}
}
