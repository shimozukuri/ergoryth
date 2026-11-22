package outbox

import (
	"context"
	"fmt"

	appoutbox "github.com/shimozukuri/ergoryth/backend/internal/application/outbox"
)

const query = `
	INSERT INTO ergoryth.email_outbox (
		id,
		lead_id,
		status,
		attempts,
		last_error,
		created_at,
		available_at,
		sent_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8
	)
`

func (r *Repository) Save(
	ctx context.Context,
	message appoutbox.Message,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.executor.OpTimeout())
	defer cancel()

	_, err := r.executor.Exec(
		ctx,
		query,
		message.ID.String(),
		message.LeadID.String(),
		message.Status.String(),
		message.Attempts,
		message.LastError,
		message.CreatedAt,
		message.AvailableAt,
		message.SentAt,
	)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}
