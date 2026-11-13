package lead

import (
	"github.com/shimozukuri/ergoryth/backend/internal/infrastructure/persistence/postgres/pool/pgx"
)

type Repository struct {
	pool pgx.Pool
}

func NewRepository(
	pool pgx.Pool,
) *Repository {
	return &Repository{
		pool: pool,
	}
}
