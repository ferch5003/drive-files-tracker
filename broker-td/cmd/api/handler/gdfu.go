package handler

import (
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"net/rpc"
)

type GDriveFamilyHandler struct {
}

func NewGDriveFamilyHandler() *GDriveFamilyHandler {
	return &GDriveFamilyHandler{}
}

type FamilyPayload struct {
	Photo    []byte
	Username string
	Date     string
}

func (h *GDriveFamilyHandler) Post(c *fiber.Ctx) error {
	photo, err := c.FormFile("tg-bot-file")
	if err != nil {
		log.Println("form file", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	username := c.FormValue("username")
	date := c.FormValue("date")

	reader, err := photo.Open()
	if err != nil {
		log.Println("reader", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	photoBytes, err := io.ReadAll(reader)
	if err != nil {
		log.Println("photoBytes", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	familyPayload := FamilyPayload{
		Photo:    photoBytes,
		Username: username,
		Date:     date,
	}

	client, err := rpc.Dial("tcp", "drive-service:5001")
	if err != nil {
		log.Println("rpc client", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	var result string
	if err := client.Call("Server.UploadDriveFile", familyPayload, &result); err != nil {
		log.Println(err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": result,
	})
}
