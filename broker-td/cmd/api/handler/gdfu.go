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
	Photo    []byte
	FolderID string
	Filename string
	Username string
}

func (h *GDriveFamilyHandler) Post(c *fiber.Ctx) error {
	photo, err := c.FormFile("tg-bot-file")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	username := c.FormValue("username")
	botName := c.FormValue("bot_name")
	date := c.FormValue("date")
	filename := c.FormValue("filename")

	reader, err := photo.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	photoBytes, err := io.ReadAll(reader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	folderIDURL := fmt.Sprintf(
		"%s/users/%s/bot/%s?date=%s", h.userServiceBaseURL, username, botName, date)
	req := h.Client.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(folderIDURL)
	if err := h.Client.Parse(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	_, body, errs := h.Client.Bytes()
	if len(errs) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errs,
		})
	}

	var data struct {
		FolderID string `json:"folder_id"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	if data.FolderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	familyPayload := FamilyPayload{
		Photo:    photoBytes,
		FolderID: data.FolderID,
		Filename: filename,
		Username: username,
	}

	var result string
	if err := h.RPCClient.Call("Server.UploadDriveFile", familyPayload, &result); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": result,
	})
}
