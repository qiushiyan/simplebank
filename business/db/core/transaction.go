package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	db_generated "github.com/qiushiyan/simplebank/business/db/generated"
)

func (s *PostgresStore) ExecuteInTransaction(ctx context.Context, fn QueryFunc) error {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := db_generated.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
