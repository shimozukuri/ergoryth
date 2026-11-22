package pgx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	Begin(ctx context.Context) (pgx.Tx, error)

	Exec(
		ctx context.Context,
		query string,
		args ...any,
	) (pgconn.CommandTag, error)
}
