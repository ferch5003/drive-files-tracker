package handler

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	agent := fiber.Get("http://user-service/users")
	statusCode, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": errs,
		})
	}

	var users []fiber.Map
	err := json.Unmarshal(body, &users)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err,
		})
	}

	return c.Status(statusCode).JSON(users)
}
