package handler

import (
	"broker-td/config"
	"github.com/gofiber/fiber/v2"
	"net/rpc"
)

type GDriveFolderCreatorHandler struct {
	RPCClient *rpc.Client
}

func NewGDriveFolderCreatorHandler(configs *config.EnvVars) (*GDriveFolderCreatorHandler, error) {
	rpcClient, err := rpc.Dial("tcp", configs.DriveServiceBaseRPC)
	if err != nil {
		return nil, err
	}

	return &GDriveFolderCreatorHandler{
		RPCClient: rpcClient,
	}, nil
}

type botUser struct {
	BotID    int    `json:"bot_id"`
	UserID   int    `json:"user_id"`
	Date     string `json:"date"`
	FolderID string `json:"folder_id"`
	IsParent bool   `json:"is_parent"`
}

type BotUsersPayload struct {
	BotUsers []botUser
}

func (h *GDriveFolderCreatorHandler) Post(c *fiber.Ctx) error {
	var botUsers []botUser
	if err := c.BodyParser(&botUsers); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err,
		})
	}

	payload := BotUsersPayload{
		BotUsers: botUsers,
	}

	var result BotUsersPayload
	if err := h.RPCClient.Call("Server.CreateYearlyFolders", payload, &result); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": result.BotUsers,
	})
}
