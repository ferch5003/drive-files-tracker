package handler

import (
	"github.com/gofiber/fiber/v2"
)

type GDriveFolderCreatorHandler struct {
}

func NewGDriveFolderCreatorHandler() (*GDriveFolderCreatorHandler, error) {
	return &GDriveFolderCreatorHandler{}, nil
}

type botUser struct {
	BotID    int    `json:"bot_id"`
	UserID   int    `json:"user_id"`
	Date     string `json:"date"`
	FolderID string `json:"folder_id"`
	IsParent bool   `json:"is_parent"`
}

func (h *GDriveFolderCreatorHandler) Post(c *fiber.Ctx) error {
	var botUsers []botUser
	if err := c.BodyParser(&botUsers); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": botUsers,
	})
}
