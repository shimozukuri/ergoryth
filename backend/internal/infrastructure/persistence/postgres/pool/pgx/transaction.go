package pgx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{pool: pool}
}

func (m *TransactionManager) WithinTransaction(
	ctx context.Context,
	fn func( /* application.Transaction */ ) error,
) (err error) {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Создаём репозитории, использующие tx.
	// transaction := ...

	// err = fn(transaction)
	// if err != nil {
	//     return err
	// }

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
