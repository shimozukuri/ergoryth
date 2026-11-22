package transaction_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shimozukuri/ergoryth/backend/internal/application/service/lead/command/createlead"
	"github.com/shimozukuri/ergoryth/backend/internal/infrastructure/persistence/postgres/pool/pgx/transaction"

	appoutbox "github.com/shimozukuri/ergoryth/backend/internal/application/outbox"
	domainlead "github.com/shimozukuri/ergoryth/backend/internal/domain/entity/lead"
	"github.com/shimozukuri/ergoryth/backend/internal/infrastructure/persistence/postgres/pool/pgx/testutil"
)

const (
	checkExistLeadQuery = `
		SELECT EXISTS (
			SELECT 1 FROM ergoryth.leads WHERE id = $1
		)
	`

	checkExistOutboxQuery = `
		SELECT EXISTS (
			SELECT 1 FROM ergoryth.email_outbox WHERE id = $1
		)
	`
)

func TestManager_Commit(t *testing.T) {
	executor, pool := testutil.Open(t)
	manager := transaction.NewManager(executor)

	entity := testutil.NewLead(t, new("commit@example.com"), nil)
	message := appoutbox.NewMessage(entity.ID())

	t.Cleanup(func() {
		testutil.DeleteLead(t, pool, entity.ID().UUID())
	})

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	err := manager.WithinTransaction(ctx, func(tx createlead.Transaction) error {
		if err := tx.Leads().Save(ctx, entity); err != nil {
			return err
		}

		return tx.Outbox().Save(ctx, message)
	})
	if err != nil {
		t.Fatalf("transaction: %v", err)
	}

	assertRowsExist(t, ctx, pool, entity.ID().UUID(), message.ID)
}

func TestManager_RollbackOnCallbackError(t *testing.T) {
	executor, pool := testutil.Open(t)
	manager := transaction.NewManager(executor)

	entity := testutil.NewLead(t, new("rollback@example.com"), nil)
	message := appoutbox.NewMessage(entity.ID())

	expectedErr := errors.New("callback failed")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	err := manager.WithinTransaction(ctx, func(tx createlead.Transaction) error {
		if err := tx.Leads().Save(ctx, entity); err != nil {
			return err
		}
		if err := tx.Outbox().Save(ctx, message); err != nil {
			return err
		}
		return expectedErr
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected callback error, got: %v", err)
	}

	assertRowsAbsent(t, ctx, pool, entity.ID().UUID(), message.ID)
}

func TestManager_RollbackWhenOutboxSaveFails(t *testing.T) {
	executor, pool := testutil.Open(t)
	manager := transaction.NewManager(executor)

	entity := testutil.NewLead(t, new("atomicity@example.com"), nil)

	fakeLeadID, err := domainlead.NewID(uuid.New())
	if err != nil {
		t.Fatalf("create fake lead ID: %v", err)
	}
	message := appoutbox.NewMessage(fakeLeadID)

	t.Cleanup(func() {
		testutil.DeleteLead(t, pool, entity.ID().UUID())
	})

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	err = manager.WithinTransaction(ctx, func(tx createlead.Transaction) error {
		if err = tx.Leads().Save(ctx, entity); err != nil {
			return err
		}
		return tx.Outbox().Save(ctx, message)
	})

	if err == nil {
		t.Fatal("expected outbox foreign key error")
	}

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); !ok || pgErr.Code != "23503" {
		t.Fatalf("expected foreign key violation, got: %v", err)
	}

	assertRowsAbsent(t, ctx, pool, entity.ID().UUID(), message.ID)
}

func assertRowsExist(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	leadID uuid.UUID,
	messageID uuid.UUID,
) {
	t.Helper()

	var leadExists, outboxExists bool

	err := pool.QueryRow(ctx, checkExistLeadQuery, leadID).Scan(&leadExists)
	if err != nil {
		t.Fatalf("check lead: %v", err)
	}

	err = pool.QueryRow(ctx, checkExistOutboxQuery, messageID).Scan(&outboxExists)
	if err != nil {
		t.Fatalf("check outbox: %v", err)
	}

	if !leadExists || !outboxExists {
		t.Fatalf("expected both rows to exist, lead=%v outbox=%v",
			leadExists, outboxExists)
	}
}

func assertRowsAbsent(
	t *testing.T,
	ctx context.Context,
	pool *pgxpool.Pool,
	leadID uuid.UUID,
	messageID uuid.UUID,
) {
	t.Helper()

	var leadExists, outboxExists bool

	err := pool.QueryRow(ctx, checkExistLeadQuery, leadID).Scan(&leadExists)
	if err != nil {
		t.Fatalf("check lead: %v", err)
	}

	err = pool.QueryRow(ctx, checkExistOutboxQuery, messageID).Scan(&outboxExists)
	if err != nil {
		t.Fatalf("check outbox: %v", err)
	}

	if leadExists || outboxExists {
		t.Fatalf("expected both rows to be absent, lead=%v outbox=%v",
			leadExists, outboxExists)
	}
}
