package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/internal/domain"
)

const _usersPath = "/users"

type userServiceMock struct {
	mock.Mock
}

func (usm *userServiceMock) GetAll(ctx context.Context) ([]domain.User, error) {
	args := usm.Called(ctx)
	return args.Get(0).([]domain.User), args.Error(1)
}

func createUserServer(usm *userServiceMock) *fiber.App {
	app := fiber.New()

	userHandler := NewUserHandler(usm)

	app.Route("/users", func(api fiber.Router) {
		api.Get("/", userHandler.GetAll).Name("get_all")

	}, "users.")

	return app
}

func createUserRequest(method string, url string, body string) (*http.Request, error) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func TestUserHandlerGetAll_Successful(t *testing.T) {
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

	usm := new(userServiceMock)
	usm.On("GetAll", mock.Anything).Return(expectedUsers, nil)

	server := createUserServer(usm)

	req, err := createUserRequest(fiber.MethodGet, _usersPath, "")
	require.NoError(t, err)

	// When
	resp, err := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var users []domain.User
	err = json.Unmarshal(body, &users)
	require.NoError(t, err)

	require.EqualValues(t, expectedUsers, users)
}

func TestUserHandlerGetAll_FailsDueToServiceError(t *testing.T) {
	// Given
	expectedErr := errors.New("sql: no rows in result set")

	usm := new(userServiceMock)
	usm.On("GetAll", mock.Anything).Return([]domain.User{}, expectedErr)

	server := createUserServer(usm)

	req, err := createUserRequest(
		fiber.MethodGet,
		_usersPath,
		"")
	require.NoError(t, err)

	// When
	resp, _ := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response errorResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	require.Equal(t, expectedErr.Error(), response.Error)
}
