package valueobject

import (
	"fmt"
	"net/mail"
	"strings"

	internalerrors "github.com/shimozukuri/ergoryth/backend/internal/errors"
)

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	email := strings.TrimSpace(value)

	if err := validateEmail(email); err != nil {
		return Email{}, fmt.Errorf(
			"validate email value: %w",
			err,
		)
	}

	return Email{value: email}, nil
}

func validateEmail(value string) error {
	if value == "" {
		return fmt.Errorf(
			"'Email' must not be empty: %w'",
			internalerrors.ErrInvalidArgument,
		)
	}

	addr, err := mail.ParseAddress(value)
	if err != nil || addr.Address != value {
		return fmt.Errorf(
			"invalid 'Email' format: %w: %v'",
			internalerrors.ErrInvalidArgument,
			err,
		)
	}

	return nil
}

func (e Email) String() string {
	return e.value
}
