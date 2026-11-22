package createlead

import (
	"context"

	"github.com/shimozukuri/ergoryth/backend/internal/application/outbox"
	"github.com/shimozukuri/ergoryth/backend/internal/domain/entity/lead"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.8 --name=LeadRepository
type LeadRepository interface {
	Save(ctx context.Context, entity lead.Lead) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.8 --name=OutboxRepository
type OutboxRepository interface {
	Save(ctx context.Context, message outbox.Message) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.8 --name=Transaction
type Transaction interface {
	Leads() LeadRepository
	Outbox() OutboxRepository
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.8 --name=TransactionManager
type TransactionManager interface {
	WithinTransaction(
		ctx context.Context,
		fn func(Transaction) error,
	) error
}
