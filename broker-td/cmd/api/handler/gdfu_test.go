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

const _gdfuPath = "/gdrive-family-uploader"

type errorResponse struct {
	Error string `json:"error"`
}

func createGDriveFamilyServer() *fiber.App {
	app := fiber.New()

	gdfuHandler := NewGDriveFamilyHandler()

	app.Route("/gdrive-family-uploader", func(api fiber.Router) {
		api.Post("/", gdfuHandler.Post).Name("post")

	}, "gdrive-family-uploader.")

	return app
}

func createGDriveFamilyRequest(method string, url string, body string) (*http.Request, error) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func TestGDriveFamilyHandlerPost_Successful(t *testing.T) {
	// Given
	server := createGDriveFamilyServer()

	req, err := createGDriveFamilyRequest(fiber.MethodPost, _gdfuPath, `{
																	"email": "john@example.com",
																	"password": "12345678"
																}`)
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

func TestUserHandlerLoginUser_FailsDueToInvalidJSONBodyParse(t *testing.T) {
	// Given
	server := createGDriveFamilyServer()

	req, err := createGDriveFamilyRequest(fiber.MethodPost, _gdfuPath, `{invalid_format}`)
	require.NoError(t, err)

	// When
	resp, _ := server.Test(req)

	// Then
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response errorResponse
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	require.Contains(t, response.Error, "invalid character 'i'")
	require.Contains(t, response.Error, "looking for beginning of object key string")
}
