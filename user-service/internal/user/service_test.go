package user

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"user-service/internal/domain"
)

type mockRepository struct {
	mock.Mock
}

func (mr *mockRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	args := mr.Called(ctx)
	return args.Get(0).([]domain.User), args.Error(1)
}

func TestServiceGetAll_Successful(t *testing.T) {
	// Given
	expectedUsers := []domain.User{
		{
			ID:       1,
			Username: "john",
		},
		{
			ID:       2,
			Username: "jane",
		},
	}

	mr := new(mockRepository)
	mr.On("GetAll", mock.Anything).Return(expectedUsers, nil)

	service := NewService(mr)

	// When
	users, err := service.GetAll(context.Background())

	// Then
	require.NoError(t, err)
	require.Len(t, users, len(expectedUsers))
	require.EqualValues(t, expectedUsers, users)
}

func TestServiceGetAll_SuccessfulWithZeroUsers(t *testing.T) {
	// Given
	expectedUsers := make([]domain.User, 0)

	mr := new(mockRepository)
	mr.On("GetAll", mock.Anything).Return(expectedUsers, nil)

	service := NewService(mr)

	// When
	users, err := service.GetAll(context.Background())

	// Then
	require.NoError(t, err)
	require.Len(t, users, len(expectedUsers))
	require.EqualValues(t, expectedUsers, users)
}

func TestServiceGetAll_FailsDueToRepositoryError(t *testing.T) {
	// Given
	expectedUsers := make([]domain.User, 0)
	expectedError := errors.New("Error Code: 1054. Unknown column 'wrong' in 'field list'")

	mr := new(mockRepository)
	mr.On("GetAll", mock.Anything).Return(expectedUsers, expectedError)

	service := NewService(mr)

	// When
	users, err := service.GetAll(context.Background())

	// Then
	require.ErrorContains(t, err, "Error Code: 1054")
	require.ErrorContains(t, err, "Unknown column 'wrong' in 'field list'")
	require.Len(t, users, 0)
	require.EqualValues(t, expectedUsers, users)
}
