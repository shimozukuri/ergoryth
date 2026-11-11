package createlead

import (
	"context"

	"github.com/shimozukuri/ergoryth/backend/internal/domain/entity/lead"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.8 --name=LeadRepository
type LeadRepository interface {
	Save(ctx context.Context, entity lead.Lead) error
}
