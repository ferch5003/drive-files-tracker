package botuser

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

func (mr *mockRepository) GetAllParents(ctx context.Context) ([]domain.BotUser, error) {
	args := mr.Called(ctx)
	return args.Get(0).([]domain.BotUser), args.Error(1)
}

func (mr *mockRepository) SaveMany(ctx context.Context, botUsers []domain.BotUser) error {
	args := mr.Called(ctx, botUsers)
	return args.Error(0)
}

func TestServiceGetAll_Successful(t *testing.T) {
	// Given
	expectedBotUsers := []domain.BotUser{
		{
			BotID:    1,
			UserID:   1,
			Date:     "test",
			FolderID: "test",
			IsParent: false,
		},
		{
			BotID:    1,
			UserID:   2,
			Date:     "test",
			FolderID: "test",
			IsParent: true,
		},
	}

	mr := new(mockRepository)
	mr.On("GetAllParents", mock.Anything).Return(expectedBotUsers, nil)

	service := NewService(mr)

	// When
	botUsers, err := service.GetAllParents(context.Background())

	// Then
	require.NoError(t, err)
	require.Len(t, botUsers, len(expectedBotUsers))
	require.EqualValues(t, expectedBotUsers, botUsers)
}

func TestServiceGetAll_SuccessfulWithZeroUsers(t *testing.T) {
	// Given
	expectedBotUsers := make([]domain.BotUser, 0)

	mr := new(mockRepository)
	mr.On("GetAllParents", mock.Anything).Return(expectedBotUsers, nil)

	service := NewService(mr)

	// When
	botUsers, err := service.GetAllParents(context.Background())

	// Then
	require.NoError(t, err)
	require.Len(t, botUsers, len(expectedBotUsers))
	require.EqualValues(t, expectedBotUsers, botUsers)
}

func TestServiceGetAll_FailsDueToRepositoryError(t *testing.T) {
	// Given
	expectedBotUsers := make([]domain.BotUser, 0)
	expectedError := errors.New("Error Code: 1054. Unknown column 'wrong' in 'field list'")

	mr := new(mockRepository)
	mr.On("GetAllParents", mock.Anything).Return(expectedBotUsers, expectedError)

	service := NewService(mr)

	// When
	botUsers, err := service.GetAllParents(context.Background())

	// Then
	require.ErrorContains(t, err, "Error Code: 1054")
	require.ErrorContains(t, err, "Unknown column 'wrong' in 'field list'")
	require.Len(t, botUsers, 0)
	require.EqualValues(t, expectedBotUsers, botUsers)
}

func TestServiceSaveMany_Successful(t *testing.T) {
	// Given
	botUsers := []domain.BotUser{
		{
			BotID:    1,
			UserID:   1,
			Date:     "test",
			FolderID: "test",
			IsParent: false,
		},
		{
			BotID:    1,
			UserID:   2,
			Date:     "test",
			FolderID: "test",
			IsParent: true,
		},
	}

	mr := new(mockRepository)
	mr.On("SaveMany", mock.Anything, botUsers).Return(nil)

	service := NewService(mr)

	// When
	err := service.SaveMany(context.Background(), botUsers)

	// Then
	require.NoError(t, err)
}

func TestServiceSaveMany_FailsDueToRepositoryError(t *testing.T) {
	// Given
	var botUsers []domain.BotUser
	expectedError := errors.New("Error Code: 1054. Unknown column 'wrong' in 'field list'")

	mr := new(mockRepository)
	mr.On("SaveMany", mock.Anything, botUsers).Return(expectedError)

	service := NewService(mr)

	// When
	err := service.SaveMany(context.Background(), botUsers)

	// Then
	require.ErrorContains(t, err, "Error Code: 1054")
	require.ErrorContains(t, err, "Unknown column 'wrong' in 'field list'")
}
