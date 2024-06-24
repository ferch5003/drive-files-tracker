package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
	"user-service/config"
)

var _testConfigs = &config.EnvVars{}

type mockRouter struct {
	mock.Mock
}

func (m *mockRouter) Register() {
	m.Called()
}

func TestRegister_Successful(t *testing.T) {
	// Given
	app := fiber.New()

	mur := new(mockRouter)
	mur.On("Register")

	mtr := new(mockRouter)
	mtr.On("Register")

	router := NewRouter(app, _testConfigs, mur, mtr) // Always have the /health endpoint.
	expectedRoute := "/health"
	expectedStatusCode := fiber.StatusOK

	// When
	router.Register() // Register routes.

	req := httptest.NewRequest("GET", expectedRoute, nil)
	resp, err := app.Test(req, -1)

	// Then
	require.NoError(t, err)
	require.Equal(t, expectedStatusCode, resp.StatusCode)
}
