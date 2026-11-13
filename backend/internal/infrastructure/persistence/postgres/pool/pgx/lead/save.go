package lead

import (
	"context"
	"fmt"

	"github.com/shimozukuri/ergoryth/backend/internal/domain/entity/lead"
)

const saveQuery = `
		INSERT INTO ergoryth.leads 
			(id, name, email, phone, company, message, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

func (r *Repository) Save(
	ctx context.Context,
	lead lead.Lead,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var email *string
	var phone *string

	if value, ok := lead.Contact().Email(); ok {
		email = new(value.String())
	}

	if value, ok := lead.Contact().Phone(); ok {
		phone = new(value.String())
	}

	_, err := r.pool.Exec(
		ctx,
		saveQuery,
		lead.ID,
		lead.Name,
		email,
		phone,
		lead.Company(),
		lead.Message(),
		lead.Status(),
		lead.CreatedAt(),
	)
	if err != nil {
		return fmt.Errorf(
			"exec query: %w",
			err,
		)
	}

	return nil
}
