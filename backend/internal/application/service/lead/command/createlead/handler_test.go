package createlead_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shimozukuri/ergoryth/backend/internal/application/service/lead/command/createlead"
	"github.com/shimozukuri/ergoryth/backend/internal/application/service/lead/command/createlead/mocks"
	"github.com/shimozukuri/ergoryth/backend/internal/domain/entity/lead"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/shimozukuri/ergoryth/backend/internal/application/outbox"
)

func strPtr(value string) *string {
	return &value
}

func newTestMocks(t *testing.T) (
	*createlead.Handler,
	*mocks.TransactionManager,
	*mocks.Transaction,
	*mocks.LeadRepository,
	*mocks.OutboxRepository,
) {
	t.Helper()

	tm := mocks.NewTransactionManager(t)
	tx := mocks.NewTransaction(t)
	leadRepo := mocks.NewLeadRepository(t)
	outboxRepo := mocks.NewOutboxRepository(t)

	handler := createlead.NewHandler(tm)

	return handler, tm, tx, leadRepo, outboxRepo
}

func setupSuccessfulTransaction(
	t *testing.T,
	tm *mocks.TransactionManager,
	tx *mocks.Transaction,
	leadRepo *mocks.LeadRepository,
	outboxRepo *mocks.OutboxRepository,
) {
	t.Helper()

	tm.On(
		"WithinTransaction",
		mock.Anything,
		mock.Anything,
	).Return(
		func(_ context.Context, fn func(createlead.Transaction) error) error {
			return fn(tx)
		},
	).Once()

	tx.On("Leads").Return(leadRepo).Once()
	if outboxRepo != nil {
		tx.On("Outbox").Return(outboxRepo).Once()
	}

	leadRepo.On(
		"Save",
		mock.Anything,
		mock.Anything,
	).Return(nil).Once()

	if outboxRepo != nil {
		outboxRepo.On(
			"Save",
			mock.Anything,
			mock.Anything,
		).Return(nil).Once()
	}
}

func TestHandler_Handle_Success(t *testing.T) {
	tests := []struct {
		name  string
		email *string
		phone *string
	}{
		{
			name:  "email and phone",
			email: strPtr("test@example.com"),
			phone: strPtr("79991234567"),
		},
		{
			name:  "email only",
			email: strPtr("test@example.com"),
		},
		{
			name:  "phone only",
			phone: strPtr("79991234567"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, tm, tx, leadRepo, outboxRepo := newTestMocks(t)

			setupSuccessfulTransaction(
				t,
				tm,
				tx,
				leadRepo,
				outboxRepo,
			)

			cmd := createlead.Command{
				Name:    "John Doe",
				Email:   tt.email,
				Phone:   tt.phone,
				Company: "",
				Message: "",
			}

			id, err := handler.Handle(context.Background(), cmd)

			require.NoError(t, err)
			require.NotEqual(t, lead.ID{}, id)

			savedMessage := outboxRepo.Calls[0].Arguments.Get(1).(outbox.Message)
			require.Equal(t, id, savedMessage.LeadID)
		})
	}
}

func TestHandler_Handle_InvalidInput(t *testing.T) {
	tests := []struct {
		name string
		cmd  createlead.Command
	}{
		{
			name: "without contact",
			cmd: createlead.Command{
				Name: "John Doe",
			},
		},
		{
			name: "invalid email",
			cmd: createlead.Command{
				Name:  "John Doe",
				Email: strPtr("invalid-email"),
			},
		},
		{
			name: "invalid phone",
			cmd: createlead.Command{
				Name:  "John Doe",
				Phone: strPtr("123"),
			},
		},
		{
			name: "empty name",
			cmd: createlead.Command{
				Email: strPtr("test@example.com"),
			},
		},
		{
			name: "name contains numbers",
			cmd: createlead.Command{
				Name:  "John123",
				Email: strPtr("test@example.com"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, tm, _, _, _ := newTestMocks(t)

			_, err := handler.Handle(context.Background(), tt.cmd)

			require.Error(t, err)

			tm.AssertNotCalled(
				t,
				"WithinTransaction",
				mock.Anything,
				mock.Anything,
			)
		})
	}
}

func TestHandler_Handle_LeadRepositoryError(t *testing.T) {
	expectedErr := errors.New("lead repository error")

	handler, tm, tx, leadRepo, outboxRepo := newTestMocks(t)

	tm.On(
		"WithinTransaction",
		mock.Anything,
		mock.Anything,
	).Return(
		func(_ context.Context, fn func(createlead.Transaction) error) error {
			return fn(tx)
		},
	).Once()

	tx.On("Leads").Return(leadRepo).Once()

	leadRepo.On(
		"Save",
		mock.Anything,
		mock.Anything,
	).Return(expectedErr).Once()

	cmd := createlead.Command{
		Name:  "John Doe",
		Email: strPtr("test@example.com"),
	}

	_, err := handler.Handle(context.Background(), cmd)

	require.ErrorIs(t, err, expectedErr)

	outboxRepo.AssertNotCalled(
		t,
		"Save",
		mock.Anything,
		mock.Anything,
	)
}

func TestHandler_Handle_OutboxRepositoryError(t *testing.T) {
	expectedErr := errors.New("outbox repository error")

	handler, tm, tx, leadRepo, outboxRepo := newTestMocks(t)
	setupSuccessfulTransaction(t, tm, tx, leadRepo, nil)

	tx.On("Outbox").Return(outboxRepo).Once()

	outboxRepo.On(
		"Save",
		mock.Anything,
		mock.Anything,
	).Return(expectedErr).Once()

	cmd := createlead.Command{
		Name:  "John Doe",
		Email: strPtr("test@example.com"),
	}

	_, err := handler.Handle(context.Background(), cmd)

	require.ErrorIs(t, err, expectedErr)
}

func TestHandler_Handle_TransactionManagerError(t *testing.T) {
	expectedErr := errors.New("transaction manager error")

	handler, tm, _, _, _ := newTestMocks(t)

	tm.On(
		"WithinTransaction",
		mock.Anything,
		mock.Anything,
	).Return(expectedErr).Once()

	cmd := createlead.Command{
		Name:  "John Doe",
		Email: strPtr("test@example.com"),
	}

	_, err := handler.Handle(context.Background(), cmd)

	require.ErrorIs(t, err, expectedErr)
}
