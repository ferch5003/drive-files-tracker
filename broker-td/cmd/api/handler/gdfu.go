package handler

import (
	"encoding/json"
	"fmt"
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
	FolderID string
	Filename string
	Username string
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
	botName := c.FormValue("bot_name")
	date := c.FormValue("date")
	filename := c.FormValue("filename")

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

	folderIDURL := fmt.Sprintf(
		"http://user-service/users/%s/bot/%s?date=%s", username, botName, date)
	agent := fiber.Get(folderIDURL)
	_, body, errs := agent.Bytes()
	if len(errs) > 0 {
		log.Println(errs)
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

	client, err := rpc.Dial("tcp", "drive-service:5001")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	var result string
	if err := client.Call("Server.UploadDriveFile", familyPayload, &result); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": result,
	})
}
