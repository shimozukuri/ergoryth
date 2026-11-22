package testutil

import (
	"context"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	domainlead "github.com/shimozukuri/ergoryth/backend/internal/domain/entity/lead"
	"github.com/shimozukuri/ergoryth/backend/internal/domain/valueobject"
	dbpgx "github.com/shimozukuri/ergoryth/backend/internal/infrastructure/persistence/postgres/pool/pgx"
)

func Open(t *testing.T) (*dbpgx.Executor, *pgxpool.Pool) {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("parse TEST_DATABASE_URL: %v", err)
	}

	if u.User == nil {
		t.Fatal("TEST_DATABASE_URL must contain username and password")
	}

	password, passwordOK := u.User.Password()
	if !passwordOK {
		t.Fatal("TEST_DATABASE_URL must contain password")
	}

	database := strings.TrimPrefix(u.Path, "/")
	if database == "" {
		t.Fatal("TEST_DATABASE_URL must contain database name")
	}

	port := 5432
	if rawPort := u.Port(); rawPort != "" {
		port, err = strconv.Atoi(rawPort)
		if err != nil || port < 1 || port > 65535 {
			t.Fatalf("invalid database port: %q", rawPort)
		}
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	executor, err := dbpgx.NewExecutor(ctx, dbpgx.Config{
		Host:     u.Hostname(),
		Port:     int16(port),
		User:     u.User.Username(),
		Password: password,
		Database: database,
		Timeout:  5 * time.Second,
	})
	if err != nil {
		t.Fatalf("create postgres executor: %v", err)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		if closer, ok := executor.DBTX.(interface{ Close() }); ok {
			closer.Close()
		}
		t.Fatalf("create assertion pool: %v", err)
	}

	var leadsExists, outboxExists bool
	err = pool.QueryRow(ctx, `
		SELECT
			to_regclass('ergoryth.leads') IS NOT NULL,
			to_regclass('ergoryth.email_outbox') IS NOT NULL
	`).Scan(&leadsExists, &outboxExists)
	if err != nil {
		t.Fatalf("check database schema: %v", err)
	}
	if !leadsExists || !outboxExists {
		t.Fatal("database migrations have not been applied")
	}

	t.Cleanup(func() {
		pool.Close()

		if closer, ok := executor.DBTX.(interface{ Close() }); ok {
			closer.Close()
		}
	})

	return executor, pool
}

func NewLead(
	t *testing.T,
	emailValue *string,
	phoneValue *string,
) domainlead.Lead {
	t.Helper()

	var email *valueobject.Email
	if emailValue != nil {
		value, err := valueobject.NewEmail(*emailValue)
		if err != nil {
			t.Fatalf("create email: %v", err)
		}
		email = &value
	}

	var phone *valueobject.Phone
	if phoneValue != nil {
		value, err := valueobject.NewPhone(*phoneValue)
		if err != nil {
			t.Fatalf("create phone: %v", err)
		}
		phone = &value
	}

	contact, err := valueobject.NewContact(email, phone)
	if err != nil {
		t.Fatalf("create contact: %v", err)
	}

	entity, err := domainlead.NewLeadInitialized(
		"Test User",
		contact,
		"Test Company",
		"Test message",
	)
	if err != nil {
		t.Fatalf("create lead: %v", err)
	}

	return entity
}

func DeleteLead(t *testing.T, pool *pgxpool.Pool, id uuid.UUID) {
	t.Helper()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	_, err := pool.Exec(ctx,
		`DELETE FROM ergoryth.email_outbox WHERE lead_id = $1`,
		id,
	)
	if err != nil {
		t.Errorf("delete outbox messages: %v", err)
	}

	_, err = pool.Exec(ctx,
		`DELETE FROM ergoryth.leads WHERE id = $1`,
		id,
	)
	if err != nil {
		t.Errorf("delete lead: %v", err)
	}
}
