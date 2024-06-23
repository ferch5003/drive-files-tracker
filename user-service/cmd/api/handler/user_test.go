package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const _usersPath = "/users"

func createUserServer() *fiber.App {
	app := fiber.New()

	userHandler := NewUserHandler()

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

func TestGDriveFamilyHandlerPost_Successful(t *testing.T) {
	// Given
	server := createUserServer()

	req, err := createUserRequest(fiber.MethodGet, _usersPath, "")
	require.NoError(t, err)

	// When
	resp, err := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	require.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var data any
	err = json.Unmarshal(body, &data)
	require.NoError(t, err)

	require.NotEmpty(t, data)
}
