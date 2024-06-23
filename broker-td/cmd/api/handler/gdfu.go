package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type GDriveFamilyHandler struct {
}

func NewGDriveFamilyHandler() *GDriveFamilyHandler {
	return &GDriveFamilyHandler{}
}

func (h *GDriveFamilyHandler) Post(c *fiber.Ctx) error {
	if form, err := c.MultipartForm(); err == nil {
		// Get all files from "documents" key:
		image := form.File["tg-bot-file"][0]
		username := form.Value["username"][0]

		fmt.Println(image.Filename, image.Size, username)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": "any",
	})
}
