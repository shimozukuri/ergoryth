package outbox

import (
	"github.com/shimozukuri/ergoryth/backend/internal/infrastructure/persistence/postgres/pool/pgx"
)

type Repository struct {
	executor *pgx.Executor
}

func NewRepository(
	executor *pgx.Executor,
) *Repository {
	return &Repository{
		executor: executor,
	}
}
