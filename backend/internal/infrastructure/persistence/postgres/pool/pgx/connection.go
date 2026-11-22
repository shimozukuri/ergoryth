package pgx

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Executor struct {
	DBTX
	opTimeout time.Duration
}

func NewExecutor(
	ctx context.Context,
	conf Config,
) (*Executor, error) {
	connectionUrl := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(conf.User, conf.Password),
		Host:   fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Path:   "/" + conf.Database,
	}

	connectionParams := url.Values{}
	connectionParams.Set("sslmode", "disable")

	connectionUrl.RawQuery = connectionParams.Encode()

	connectionString := connectionUrl.String()

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
		pool.Close()

		return nil, fmt.Errorf(
			"ping pgxpool: %w",
			err,
		)
	}

	return &Executor{
		pool,
		conf.Timeout,
	}, nil
}

func (p *Executor) WithTx(tx pgx.Tx) *Executor {
	return &Executor{
		tx,
		p.opTimeout,
	}
}

func (p *Executor) OpTimeout() time.Duration {
	return p.opTimeout
}
