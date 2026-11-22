package transaction

import (
	poolpgx "github.com/shimozukuri/ergoryth/backend/internal/infrastructure/persistence/postgres/pool/pgx"

	"github.com/shimozukuri/ergoryth/backend/internal/application/service/lead/command/createlead"

	leadrepo "github.com/shimozukuri/ergoryth/backend/internal/infrastructure/persistence/postgres/pool/pgx/lead"

	outboxrepo "github.com/shimozukuri/ergoryth/backend/internal/infrastructure/persistence/postgres/pool/pgx/outbox"
)

type Transaction struct {
	leads  createlead.LeadRepository
	outbox createlead.OutboxRepository
}

func NewTransaction(pool *poolpgx.Executor) *Transaction {
	return &Transaction{
		leads:  leadrepo.NewRepository(pool),
		outbox: outboxrepo.NewRepository(pool),
	}
}

func (t *Transaction) Leads() createlead.LeadRepository {
	return t.leads
}

func (t *Transaction) Outbox() createlead.OutboxRepository {
	return t.outbox
}

var _ createlead.Transaction = (*Transaction)(nil)
