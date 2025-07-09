package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	sqlc "github.com/iamsuudi/digital-id-server/database/sqlc"
)

// WithTx starts a db transaction and runs fn with a new Queries instance
func WithTx(ctx context.Context, db sqlc.DBTX, fn func(*sqlc.Queries) error) error {
	tx, ok := db.(pgx.Tx)
	if ok {
		// Already inside a transaction
		q := sqlc.New(tx)
		return fn(q)
	}

	pool, ok := db.(interface {
		BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	})
	if !ok {
		return fmt.Errorf("provided DBTX does not support BeginTx")
	}

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := sqlc.New(tx)

	err = fn(q)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
