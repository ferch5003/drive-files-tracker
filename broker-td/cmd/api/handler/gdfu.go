package handler

import (
	"github.com/gofiber/fiber/v2"
)

type GDriveFamilyHandler struct {
}

func NewGDriveFamilyHandler() *GDriveFamilyHandler {
	return &GDriveFamilyHandler{}
}

func (h *GDriveFamilyHandler) Post(c *fiber.Ctx) error {
	var data any
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(data)
}
