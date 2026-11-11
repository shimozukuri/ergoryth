package valueobject

import (
	"fmt"

	internalerrors "github.com/shimozukuri/ergoryth/backend/internal/errors"
)

type Contact struct {
	email *Email
	phone *Phone
}

func NewContact(email *Email, phone *Phone) (Contact, error) {
	if email == nil && phone == nil {
		return Contact{}, fmt.Errorf(
			"'Email' or 'Phone' must not be empty: %w",
			internalerrors.ErrInvalidArgument,
		)
	}

	c := Contact{}

	if email != nil {
		c.email = new(*email)
	}

	if phone != nil {
		c.phone = new(*phone)
	}

	return c, nil
}

func (c Contact) Email() (Email, bool) {
	if c.email == nil {
		return Email{}, false
	}

	return *c.email, true
}

func (c Contact) Phone() (Phone, bool) {
	if c.phone == nil {
		return Phone{}, false
	}

	return *c.phone, true
}

func (c Contact) IsEmpty() bool {
	return c.email == nil && c.phone == nil
}
