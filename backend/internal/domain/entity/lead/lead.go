package lead

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/shimozukuri/ergoryth/backend/internal/domain/valueobject"
	internalerrors "github.com/shimozukuri/ergoryth/backend/internal/errors"
)

type Lead struct {
	id        ID
	name      string
	contact   valueobject.Contact
	company   string
	message   string
	status    Status
	createdAt time.Time
}

func NewLead(
	id ID,
	name string,
	contact valueobject.Contact,
	company string,
	message string,
	status Status,
	createdAt time.Time,
) (Lead, error) {
	entity := Lead{
		id:        id,
		name:      strings.TrimSpace(name),
		contact:   contact,
		company:   strings.TrimSpace(company),
		message:   strings.TrimSpace(message),
		status:    status,
		createdAt: createdAt.UTC(),
	}

	if err := entity.Validate(); err != nil {
		return Lead{}, fmt.Errorf(
			"validate lead domain: %w",
			err,
		)
	}

	return entity, nil
}

func NewLeadInitialized(
	name string,
	contact valueobject.Contact,
	company string,
	message string,
) (Lead, error) {
	LeadID, err := NewID(uuid.New())
	if err != nil {
		return Lead{}, fmt.Errorf(
			"create new ID: %w",
			err,
		)
	}

	return NewLead(
		LeadID,
		name,
		contact,
		company,
		message,
		StatusNew,
		time.Now().UTC(),
	)
}

func (l *Lead) Validate() error {
	nameLen := utf8.RuneCountInString(l.name)
	if nameLen < 2 || nameLen > 100 {
		return fmt.Errorf(
			"name length must be between 2 and 100, got %d: %w",
			nameLen,
			internalerrors.ErrInvalidArgument,
		)
	}

	companyLen := utf8.RuneCountInString(l.company)
	if companyLen > 150 {
		return fmt.Errorf(
			"company length must be less 150, got %d: %w",
			companyLen,
			internalerrors.ErrInvalidArgument,
		)
	}

	messageLen := utf8.RuneCountInString(l.message)
	if messageLen > 2000 {
		return fmt.Errorf(
			"message length must be less 2000, got %d: %w",
			messageLen,
			internalerrors.ErrInvalidArgument,
		)
	}

	if !l.id.isValid() {
		return fmt.Errorf(
			"invalid lead id: %w",
			internalerrors.ErrInvalidArgument,
		)
	}

	if l.contact.IsEmpty() {
		return fmt.Errorf(
			"contact is required: %w",
			internalerrors.ErrInvalidArgument,
		)
	}

	if !l.status.IsValid() {
		return fmt.Errorf(
			"invalid lead status: %w",
			internalerrors.ErrInvalidArgument,
		)
	}

	if l.createdAt.IsZero() {
		return fmt.Errorf(
			"created at time is required: %w",
			internalerrors.ErrInvalidArgument,
		)
	}

	return nil
}

func (l *Lead) ID() ID {
	return l.id
}

func (l *Lead) Name() string {
	return l.name
}

func (l *Lead) Contact() valueobject.Contact {
	return l.contact
}

func (l *Lead) Company() string {
	return l.company
}

func (l *Lead) Message() string {
	return l.message
}

func (l *Lead) Status() Status {
	return l.status
}

func (l *Lead) CreatedAt() time.Time {
	return l.createdAt
}
