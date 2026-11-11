package lead

import (
	"fmt"

	"github.com/google/uuid"
	internalerrors "github.com/shimozukuri/ergoryth/backend/internal/errors"
)

type ID struct {
	value uuid.UUID
}

func NewID(value uuid.UUID) (ID, error) {
	if value == uuid.Nil {
		return ID{}, fmt.Errorf(
			"'ID' must not be an empty: %w",
			internalerrors.ErrInvalidArgument,
		)
	}

	return ID{value: value}, nil
}

func (id ID) IsValid() bool {
	return id.value != uuid.Nil
}

func (id ID) String() string {
	return id.value.String()
}

func (id ID) UUID() uuid.UUID {
	return id.value
}
