package pgx

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

func NewPool(
	ctx context.Context,
	conf Config,
) (*Pool, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@localhost:5432/%s?sslmode=disable",
		conf.User,
		conf.Password,
		conf.Database,
	)

	pgxconf, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf(
			"parse pgxconfig: %w",
			err,
		)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconf)
	if err != nil {
		return nil, fmt.Errorf(
			"create pgxpool: %w",
			err,
		)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf(
			"ping pgxpool: %w",
			err,
		)
	}

	return &Pool{
		pool,
		conf.Timeout,
	}, nil
}

func (p *Pool) OpTimeout() time.Duration {
	return p.opTimeout
}
