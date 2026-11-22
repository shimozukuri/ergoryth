package outbox

import (
	"time"

	"github.com/google/uuid"
	"github.com/shimozukuri/ergoryth/backend/internal/domain/entity/lead"
)

type Message struct {
	ID          uuid.UUID
	LeadID      lead.ID
	Status      Status
	Attempts    int
	LastError   *string
	CreatedAt   time.Time
	AvailableAt time.Time
	SentAt      *time.Time
}

func NewMessage(leadID lead.ID) Message {
	now := time.Now().UTC()

	return Message{
		ID:          uuid.New(),
		LeadID:      leadID,
		Status:      StatusPending,
		Attempts:    0,
		CreatedAt:   now,
		AvailableAt: now,
	}
}
