package transaction

import (
	"context"
	"fmt"
	"time"

	"github.com/shimozukuri/ergoryth/backend/internal/application/service/lead/command/createlead"
	"github.com/shimozukuri/ergoryth/backend/internal/infrastructure/persistence/postgres/pool/pgx"
)

type Manager struct {
	executor *pgx.Executor
}

func NewManager(executor *pgx.Executor) *Manager {
	return &Manager{
		executor: executor,
	}
}

func (m *Manager) WithinTransaction(
	ctx context.Context,
	fn func(createlead.Transaction) error,
) error {
	tx, err := m.executor.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		rollbackCtx, cancel := context.WithTimeout(
			context.WithoutCancel(ctx),
			5*time.Second,
		)
		defer cancel()

		_ = tx.Rollback(rollbackCtx)
	}()

	txPool := m.executor.WithTx(tx)

	transaction := NewTransaction(txPool)

	if err = fn(transaction); err != nil {
		return fmt.Errorf("execute transaction: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

var _ createlead.TransactionManager = (*Manager)(nil)
