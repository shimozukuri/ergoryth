package outbox_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shimozukuri/ergoryth/backend/internal/infrastructure/persistence/postgres/pool/pgx/outbox"

	appoutbox "github.com/shimozukuri/ergoryth/backend/internal/application/outbox"
	domainlead "github.com/shimozukuri/ergoryth/backend/internal/domain/entity/lead"
	leadrepo "github.com/shimozukuri/ergoryth/backend/internal/infrastructure/persistence/postgres/pool/pgx/lead"
	"github.com/shimozukuri/ergoryth/backend/internal/infrastructure/persistence/postgres/pool/pgx/testutil"
)

const saveQuery = `
		SELECT
			id, lead_id, status, attempts,
			last_error, created_at, available_at, sent_at
		FROM ergoryth.email_outbox
		WHERE id = $1
	`

func TestRepository_Save(t *testing.T) {
	executor, pool := testutil.Open(t)

	leads := leadrepo.NewRepository(executor)
	repo := outbox.NewRepository(executor)

	email := "outbox@example.com"
	entity := testutil.NewLead(t, &email, nil)

	t.Cleanup(func() {
		testutil.DeleteLead(t, pool, entity.ID().UUID())
	})

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := leads.Save(ctx, entity); err != nil {
		t.Fatalf("save lead: %v", err)
	}

	message := appoutbox.NewMessage(entity.ID())

	if err := repo.Save(ctx, message); err != nil {
		t.Fatalf("save outbox message: %v", err)
	}

	var (
		gotID       uuid.UUID
		gotLeadID   uuid.UUID
		status      string
		attempts    int
		lastError   pgtype.Text
		createdAt   time.Time
		availableAt time.Time
		sentAt      pgtype.Timestamptz
	)

	err := pool.QueryRow(ctx, saveQuery, message.ID).Scan(
		&gotID,
		&gotLeadID,
		&status,
		&attempts,
		&lastError,
		&createdAt,
		&availableAt,
		&sentAt,
	)
	if err != nil {
		t.Fatalf("select outbox message: %v", err)
	}

	if gotID != message.ID {
		t.Errorf("unexpected message ID: %v", gotID)
	}
	if gotLeadID != entity.ID().UUID() {
		t.Errorf("unexpected lead ID: %v", gotLeadID)
	}
	if status != message.Status.String() {
		t.Errorf("unexpected status: %q", status)
	}
	if attempts != 0 {
		t.Errorf("expected 0 attempts, got %d", attempts)
	}
	if lastError.Valid {
		t.Errorf("expected NULL last_error, got %q", lastError.String)
	}
	if sentAt.Valid {
		t.Error("expected NULL sent_at")
	}
	if createdAt.IsZero() || availableAt.IsZero() {
		t.Error("timestamps must not be zero")
	}
	if !createdAt.Equal(availableAt) {
		t.Errorf("created_at and available_at differ: %v / %v",
			createdAt, availableAt)
	}
}

func TestRepository_SaveMissingLead(t *testing.T) {
	executor, _ := testutil.Open(t)
	repo := outbox.NewRepository(executor)

	fakeLeadID, err := domainlead.NewID(uuid.New())
	if err != nil {
		t.Fatalf("create fake lead ID: %v", err)
	}

	message := appoutbox.NewMessage(fakeLeadID)

	err = repo.Save(context.Background(), message)
	if err == nil {
		t.Fatal("expected foreign key violation")
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		t.Fatalf("expected PostgreSQL error, got: %T: %v", err, err)
	}
	if pgErr.Code != "23503" {
		t.Errorf("expected foreign key violation (23503), got %s", pgErr.Code)
	}
}
