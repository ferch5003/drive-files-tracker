package user

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"user-service/internal/domain"
)

type userMockRepository struct {
	mock.Mock
}

func (mr *userMockRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	args := mr.Called(ctx)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (mr *userMockRepository) Get(ctx context.Context, username string) (domain.User, error) {
	args := mr.Called(ctx, username)
	return args.Get(0).(domain.User), args.Error(1)
}

func (mr *userMockRepository) FindFolderID(ctx context.Context, userID, botID int, date string) (string, error) {
	args := mr.Called(ctx, userID, botID, date)
	return args.Get(0).(string), args.Error(1)
}

type botMockRepository struct {
	mock.Mock
}

func (mr *botMockRepository) Get(ctx context.Context, name string) (domain.Bot, error) {
	args := mr.Called(ctx, name)
	return args.Get(0).(domain.Bot), args.Error(1)
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

	mr := new(userMockRepository)
	mr.On("GetAll", mock.Anything).Return(expectedUsers, nil)

	service := NewService(mr, nil)

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

	mr := new(userMockRepository)
	mr.On("GetAll", mock.Anything).Return(expectedUsers, nil)

	service := NewService(mr, nil)

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

	mr := new(userMockRepository)
	mr.On("GetAll", mock.Anything).Return(expectedUsers, expectedError)

	service := NewService(mr, nil)

	// When
	users, err := service.GetAll(context.Background())

	// Then
	require.ErrorContains(t, err, "Error Code: 1054")
	require.ErrorContains(t, err, "Unknown column 'wrong' in 'field list'")
	require.Len(t, users, 0)
	require.EqualValues(t, expectedUsers, users)
}

func TestServiceFindFolderID_Successful(t *testing.T) {
	// Given
	expectedFolderID := "Test"

	user := domain.User{
		ID:       1,
		Username: "Test User",
	}

	bot := domain.Bot{
		ID:   1,
		Name: "Test Bot",
	}

	umr := new(userMockRepository)
	umr.On("Get", mock.Anything, user.Username).Return(user, nil)

	bmr := new(botMockRepository)
	bmr.On("Get", mock.Anything, bot.Name).Return(bot, nil)

	umr.On("FindFolderID", mock.Anything, user.ID, bot.ID, "").Return(expectedFolderID, nil)

	service := NewService(umr, bmr)

	// When
	folderID, err := service.FindFolderID(context.Background(), user.Username, bot.Name, "")

	// Then
	require.NoError(t, err)
	require.Equal(t, expectedFolderID, folderID)
}

func TestServiceFindFolderID_FailsDueToUserRepositoryGetError(t *testing.T) {
	// Given
	user := domain.User{}
	expectedError := errors.New("Error Code: 1054. Unknown column 'wrong' in 'field list'")

	umr := new(userMockRepository)
	umr.On("Get", mock.Anything, user.Username).Return(domain.User{}, expectedError)

	service := NewService(umr, nil)

	// When
	folderID, err := service.FindFolderID(context.Background(), user.Username, "", "")

	// Then
	require.ErrorContains(t, err, "Error Code: 1054")
	require.ErrorContains(t, err, "Unknown column 'wrong' in 'field list'")
	require.Equal(t, "", folderID)
}

func TestServiceFindFolderID_FailsDueToBotRepositoryGetError(t *testing.T) {
	// Given
	user := domain.User{
		ID:       1,
		Username: "Test User",
	}

	bot := domain.Bot{}
	expectedError := errors.New("Error Code: 1054. Unknown column 'wrong' in 'field list'")

	umr := new(userMockRepository)
	umr.On("Get", mock.Anything, user.Username).Return(user, nil)

	bmr := new(botMockRepository)
	bmr.On("Get", mock.Anything, bot.Name).Return(domain.Bot{}, expectedError)

	service := NewService(umr, bmr)

	// When
	folderID, err := service.FindFolderID(context.Background(), user.Username, bot.Name, "")

	// Then
	require.ErrorContains(t, err, "Error Code: 1054")
	require.ErrorContains(t, err, "Unknown column 'wrong' in 'field list'")
	require.Equal(t, "", folderID)
}

func TestServiceFindFolderID_FailsDueToUserRepositoryFindFolderIDError(t *testing.T) {
	// Given
	user := domain.User{
		ID:       1,
		Username: "Test User",
	}

	bot := domain.Bot{
		ID:   1,
		Name: "Test Bot",
	}
	expectedError := errors.New("Error Code: 1054. Unknown column 'wrong' in 'field list'")

	umr := new(userMockRepository)
	umr.On("Get", mock.Anything, user.Username).Return(user, nil)

	bmr := new(botMockRepository)
	bmr.On("Get", mock.Anything, bot.Name).Return(bot, nil)

	umr.On("FindFolderID", mock.Anything, user.ID, bot.ID, "").Return("", expectedError)

	service := NewService(umr, bmr)

	// When
	folderID, err := service.FindFolderID(context.Background(), user.Username, bot.Name, "")

	// Then
	require.ErrorContains(t, err, "Error Code: 1054")
	require.ErrorContains(t, err, "Unknown column 'wrong' in 'field list'")
	require.Equal(t, "", folderID)
}
