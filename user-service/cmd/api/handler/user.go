package handler

import (
	"github.com/gofiber/fiber/v2"
	"user-service/internal/user"
)

type UserHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	users, err := h.userService.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func (h *UserHandler) FindFolderID(c *fiber.Ctx) error {
	userUsername := c.Params("user_username")
	botName := c.Params("bot_name")
	date := c.Query("date")

	folderID, err := h.userService.FindFolderID(c.Context(), userUsername, botName, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"folder_id": folderID,
	})
}
