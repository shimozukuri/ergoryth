package createlead

import (
	"context"
	"fmt"
	"strings"

	"github.com/shimozukuri/ergoryth/backend/internal/domain/entity/lead"
	"github.com/shimozukuri/ergoryth/backend/internal/domain/valueobject"
)

type Handler struct {
	repository LeadRepository
}

func NewHandler(repository LeadRepository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (h *Handler) Handle(
	ctx context.Context,
	cmd Command,
) (lead.ID, error) {
	var email *valueobject.Email
	var phone *valueobject.Phone

	if cmd.Email != nil && strings.TrimSpace(*cmd.Email) != "" {
		value, err := valueobject.NewEmail(*cmd.Email)
		if err != nil {
			return lead.ID{}, fmt.Errorf(
				"create email: %w",
				err,
			)
		}

		email = &value
	}

	if cmd.Phone != nil && strings.TrimSpace(*cmd.Phone) != "" {
		value, err := valueobject.NewPhone(*cmd.Phone)
		if err != nil {
			return lead.ID{}, fmt.Errorf(
				"create phone: %w",
				err,
			)
		}

		phone = &value
	}

	contact, err := valueobject.NewContact(
		email,
		phone,
	)
	if err != nil {
		return lead.ID{}, fmt.Errorf(
			"create contact: %w",
			err,
		)
	}

	leadDomain, err := lead.NewLeadInitialized(
		cmd.Name,
		contact,
		cmd.Company,
		cmd.Message,
	)
	if err != nil {
		return lead.ID{}, fmt.Errorf(
			"create lead domain: %w",
			err,
		)
	}

	if err = h.repository.Save(ctx, leadDomain); err != nil {
		return lead.ID{}, fmt.Errorf(
			"save lead: %w",
			err,
		)
	}

	return leadDomain.ID(), nil
}
