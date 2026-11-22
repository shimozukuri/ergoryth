package lead_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shimozukuri/ergoryth/backend/internal/infrastructure/persistence/postgres/pool/pgx/lead"

	"github.com/shimozukuri/ergoryth/backend/internal/infrastructure/persistence/postgres/pool/pgx/testutil"
)

const saveQuery = `
				SELECT
					id, name, email, phone,
					company, message, status, created_at
				FROM ergoryth.leads
				WHERE id = $1
			`

func TestRepository_Save(t *testing.T) {
	executor, pool := testutil.Open(t)
	repo := lead.NewRepository(executor)

	email := "test@example.com"
	phone := "79991234567"

	tests := []struct {
		name  string
		email *string
		phone *string
	}{
		{
			name:  "email and phone",
			email: &email,
			phone: &phone,
		},
		{
			name:  "email only",
			email: &email,
		},
		{
			name:  "phone only",
			phone: &phone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := testutil.NewLead(t, tt.email, tt.phone)

			t.Cleanup(func() {
				testutil.DeleteLead(t, pool, entity.ID().UUID())
			})

			ctx, cancel := context.WithTimeout(
				context.Background(),
				5*time.Second,
			)
			defer cancel()

			if err := repo.Save(ctx, entity); err != nil {
				t.Fatalf("save lead: %v", err)
			}

			var (
				gotID     uuid.UUID
				name      string
				gotEmail  pgtype.Text
				gotPhone  pgtype.Text
				company   string
				message   string
				status    string
				createdAt time.Time
			)

			err := pool.QueryRow(ctx, saveQuery, entity.ID().UUID()).Scan(
				&gotID,
				&name,
				&gotEmail,
				&gotPhone,
				&company,
				&message,
				&status,
				&createdAt,
			)
			if err != nil {
				t.Fatalf("select saved lead: %v", err)
			}

			if gotID != entity.ID().UUID() {
				t.Errorf("unexpected ID: %v", gotID)
			}
			if name != entity.Name() {
				t.Errorf("unexpected name: %q", name)
			}
			if company != entity.Company() {
				t.Errorf("unexpected company: %q", company)
			}
			if message != entity.Message() {
				t.Errorf("unexpected message: %q", message)
			}
			if status != string(entity.Status()) {
				t.Errorf("unexpected status: %q", status)
			}
			if createdAt.IsZero() {
				t.Error("created_at must not be zero")
			}

			if tt.email == nil {
				if gotEmail.Valid {
					t.Errorf("expected NULL email, got %q", gotEmail.String)
				}
			} else if !gotEmail.Valid || gotEmail.String != *tt.email {
				t.Errorf("unexpected email: %+v", gotEmail)
			}

			if tt.phone == nil {
				if gotPhone.Valid {
					t.Errorf("expected NULL phone, got %q", gotPhone.String)
				}
			} else if !gotPhone.Valid || gotPhone.String != *tt.phone {
				t.Errorf("unexpected phone: %+v", gotPhone)
			}
		})
	}
}

func TestRepository_SaveDuplicateID(t *testing.T) {
	executor, pool := testutil.Open(t)
	repo := lead.NewRepository(executor)

	email := "duplicate@example.com"
	entity := testutil.NewLead(t, &email, nil)

	t.Cleanup(func() {
		testutil.DeleteLead(t, pool, entity.ID().UUID())
	})

	ctx := context.Background()

	if err := repo.Save(ctx, entity); err != nil {
		t.Fatalf("first save: %v", err)
	}

	err := repo.Save(ctx, entity)
	if err == nil {
		t.Fatal("expected duplicate ID error")
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		t.Fatalf("expected PostgreSQL error, got: %T: %v", err, err)
	}
	if pgErr.Code != "23505" {
		t.Errorf("expected unique violation (23505), got %s", pgErr.Code)
	}
}
