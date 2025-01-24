package handler

import (
	"broker-td/config"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/rpc"
)

type GDriveFamilyHandler struct {
	userServiceBaseURL string
	Client             *fiber.Agent
	RPCClient          *rpc.Client
}

func NewGDriveFamilyHandler(configs *config.EnvVars) (*GDriveFamilyHandler, error) {
	rpcClient, err := rpc.Dial("tcp", configs.DriveServiceBaseRPC)
	if err != nil {
		return nil, err
	}

	return &GDriveFamilyHandler{
		userServiceBaseURL: configs.UserServiceBaseURL,
		Client:             fiber.AcquireAgent(),
		RPCClient:          rpcClient,
	}, nil
}

type FamilyPayload struct {
	Photo             []byte
	FolderID          string
	Filename          string
	Username          string
	SpreadsheetID     string
	SpreadsheetGID    string
	SpreadsheetColumn string
}

func (h *GDriveFamilyHandler) Post(c *fiber.Ctx) error {
	photo, err := c.FormFile("tg-bot-file")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	username := c.FormValue("username")
	botName := c.FormValue("bot_name")
	date := c.FormValue("date")
	filename := c.FormValue("filename")

	reader, err := photo.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	photoBytes, err := io.ReadAll(reader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	folderIDURL := fmt.Sprintf(
		"%s/users/%s/bot/%s?date=%s", h.userServiceBaseURL, username, botName, date)
	folderIDBody, err := makeClientRequest(fiber.MethodGet, folderIDURL, h.Client)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var folderData struct {
		FolderID string `json:"folder_id"`
	}
	err = json.Unmarshal(folderIDBody, &folderData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if folderData.FolderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "folder id is empty",
		})
	}

	spreadsheetsURL := fmt.Sprintf(
		"%s/users/%s/bot/%s/date/%s/spreadsheets", h.userServiceBaseURL, username, botName, date)
	spreadsheetsBody, err := makeClientRequest(fiber.MethodGet, spreadsheetsURL, h.Client)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var spreadSheetData struct {
		SpreadsheetID     string `json:"spreadsheet_id"`
		SpreadsheetGID    string `json:"spreadsheet_gid"`
		SpreadsheetColumn string `json:"spreadsheet_column"`
	}
	err = json.Unmarshal(spreadsheetsBody, &spreadSheetData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	familyPayload := FamilyPayload{
		Photo:             photoBytes,
		FolderID:          folderData.FolderID,
		Filename:          filename,
		Username:          username,
		SpreadsheetID:     spreadSheetData.SpreadsheetID,
		SpreadsheetGID:    spreadSheetData.SpreadsheetGID,
		SpreadsheetColumn: spreadSheetData.SpreadsheetColumn,
	}

	var result string
	if err := h.RPCClient.Call("Server.UploadDriveFile", familyPayload, &result); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": result,
	})
}

func makeClientRequest(method, url string, client *fiber.Agent) ([]byte, error) {
	req := client.Request()
	req.Header.SetMethod(method)
	req.SetRequestURI(url)
	if err := client.Parse(); err != nil {
		return nil, err
	}

	_, body, errs := client.Bytes()
	if len(errs) > 0 {
		return nil, errs[0]
	}

	return body, nil
}
