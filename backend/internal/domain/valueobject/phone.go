package valueobject

import (
	"fmt"
	"regexp"
	"strings"

	internalerrors "github.com/shimozukuri/ergoryth/backend/internal/errors"
)

type Phone struct {
	value string
}

func NewPhone(value string) (Phone, error) {
	phone := strings.TrimSpace(value)

	if err := validatePhone(phone); err != nil {
		return Phone{}, fmt.Errorf(
			"validate phone value: %w",
			err,
		)
	}

	return Phone{value: phone}, nil
}

func validatePhone(value string) error {
	if value == "" {
		return fmt.Errorf(
			"'Phone' must not be empty: %w",
			internalerrors.ErrInvalidArgument,
		)
	}

	re := regexp.MustCompile("^7[0-9]{10}$")
	if !re.MatchString(value) {
		return fmt.Errorf(
			"invalid 'Phone' format: %w: %s",
			internalerrors.ErrInvalidArgument,
			value,
		)
	}

	return nil
}

func (p Phone) String() string {
	return p.value
}
