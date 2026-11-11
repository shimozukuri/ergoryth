package createlead_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/shimozukuri/ergoryth/backend/internal/application/service/lead/command/createlead"
	"github.com/shimozukuri/ergoryth/backend/internal/application/service/lead/command/createlead/mocks"
	"github.com/shimozukuri/ergoryth/backend/internal/domain/entity/lead"
	internalerrors "github.com/shimozukuri/ergoryth/backend/internal/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLead_CreateLead(t *testing.T) {
	name := gofakeit.Name()
	email := gofakeit.Email()
	phone := "79998887766"
	company := gofakeit.Company()
	message := gofakeit.Comment()

	repositoryError := errors.New("unexpected error")

	testCases := []struct {
		testName string
		name     string
		email    string
		phone    string
		company  string
		message  string

		respErr error
		mockErr error
	}{
		{
			testName: "success with email and phone",
			name:     name,
			email:    email,
			phone:    phone,
			company:  company,
			message:  message,
		},
		{
			testName: "success with empty phone",
			name:     name,
			email:    email,
			phone:    "",
			company:  company,
			message:  message,
		},
		{
			testName: "success with nil phone",
			name:     name,
			email:    email,
			company:  company,
			message:  message,
		},
		{
			testName: "success with empty email",
			name:     name,
			email:    "",
			phone:    phone,
			company:  company,
			message:  message,
		},
		{
			testName: "success with nil email",
			name:     name,
			phone:    phone,
			company:  company,
			message:  message,
		},
		{
			testName: "success with empty company",
			name:     name,
			email:    email,
			phone:    phone,
			company:  "",
			message:  message,
		},
		{
			testName: "success with empty message",
			name:     name,
			email:    email,
			phone:    phone,
			company:  company,
			message:  "",
		},
		{
			testName: "error with empty email and phone",
			name:     name,
			email:    "",
			phone:    "",
			company:  company,
			message:  message,

			respErr: internalerrors.ErrInvalidArgument,
		},
		{
			testName: "error with invalid email",
			name:     name,
			email:    "ivanmail.com",
			phone:    phone,
			company:  company,
			message:  message,

			respErr: internalerrors.ErrInvalidArgument,
		},
		{
			testName: "error with invalid phone",
			name:     name,
			email:    email,
			phone:    "7999888776",
			company:  company,
			message:  message,

			respErr: internalerrors.ErrInvalidArgument,
		},
		{
			testName: "error with empty name",
			name:     "",
			email:    email,
			phone:    phone,
			company:  company,
			message:  message,

			respErr: internalerrors.ErrInvalidArgument,
		},
		{
			testName: "error with short name",
			name:     "I",
			email:    email,
			phone:    phone,
			company:  company,
			message:  message,

			respErr: internalerrors.ErrInvalidArgument,
		},
		{
			testName: "error with long name",
			name:     strings.Repeat("a", 101),
			email:    email,
			phone:    phone,
			company:  company,
			message:  message,

			respErr: internalerrors.ErrInvalidArgument,
		},
		{
			testName: "error with numbers in name",
			name:     name + "1",
			email:    email,
			phone:    phone,
			company:  company,
			message:  message,

			respErr: internalerrors.ErrInvalidArgument,
		},
		{
			testName: "error with long company",
			name:     name,
			email:    email,
			phone:    phone,
			company:  strings.Repeat("a", 151),
			message:  message,

			respErr: internalerrors.ErrInvalidArgument,
		},
		{
			testName: "error with long message",
			name:     name,
			email:    email,
			phone:    phone,
			company:  company,
			message:  strings.Repeat("a", 2001),

			respErr: internalerrors.ErrInvalidArgument,
		},
		{
			testName: "save error",
			name:     name,
			email:    email,
			phone:    phone,
			company:  company,
			message:  message,

			respErr: repositoryError,
			mockErr: repositoryError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			leadRepositoryMock := mocks.NewLeadRepository(t)
			h := createlead.NewHandler(leadRepositoryMock)

			ctx := context.Background()

			if tc.respErr == nil || tc.mockErr != nil {
				leadRepositoryMock.
					On(
						"Save",
						ctx,
						mock.MatchedBy(func(entity lead.Lead) bool {
							return entity.ID().IsValid() &&
								entity.Name() == tc.name &&
								entity.Company() == tc.company &&
								entity.Message() == tc.message &&
								entity.Status() == lead.StatusNew &&
								!entity.CreatedAt().IsZero()
						}),
					).
					Return(tc.mockErr).
					Once()
			}

			leadID, err := h.Handle(ctx, createlead.Command{
				Name:    tc.name,
				Company: tc.company,
				Email:   &tc.email,
				Phone:   &tc.phone,
				Message: tc.message,
			})

			if tc.respErr == nil {
				require.NoError(t, err)
				require.True(t, leadID.IsValid())
				return
			}

			require.Error(t, err)
			require.ErrorIs(t, err, tc.respErr)
			require.False(t, leadID.IsValid())
		})
	}
}
