package repository

import (
	"context"
	"database/sql"
)

type (
	TxQuerier interface {
		Querier // nolint:revive // Use the same method names as in the Querier interface

		BeginTx(ctx context.Context) (TxQuerier, error)
		Commit() error
		Rollback() error
	}

	queries struct {
		*Queries // nolint:structcheck // Embed all the queries into the struct
		db       *sql.DB
		tx       *sql.Tx
	}
)

// NewQuerier returns a new instance of the Querier interface implementation.
func NewQuerier(db *sql.DB) TxQuerier {
	return &queries{
		Queries: New(db), // nolint:exhaustivestruct // All the queries are embedded into the struct
		db:      db,
	}
}

// BeginTx starts a transaction and returns a new instance of the TxQuerier interface implementation
// which can be used to execute queries in the scope of the transaction.
// If the transaction is committed or rolled back, the TxQuerier is no longer usable.
// The TxQuerier must be closed after the transaction is committed or rolled back.
func (q *queries) BeginTx(ctx context.Context) (TxQuerier, error) {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &queries{
		Queries: New(tx), // nolint:exhaustivestruct // All the queries are embedded into the struct
		tx:      tx,
		db:      q.db,
	}, nil
}

// Commit commits the transaction.
func (q *queries) Commit() error {
	return q.tx.Commit()
}

// Rollback rolls back the transaction.
func (q *queries) Rollback() error {
	return q.tx.Rollback()
}
