package handlers

import (
	"bytes"
	"fmt"
	"gopkg.in/telebot.v3"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"telegram-bot-service/config"
	"telegram-bot-service/internal/platform/client"
	"time"
)

const (
	_RFC3339OnlyDateFormat = "2006-01-02"
	_RFC3339OnlyYearFormat = "2006"
)

const (
	_errReadingImage      = "Hubo un error leyendo la imagen..."
	_errProcessingMessage = "Hubo un error procesando la imagen..."
	_errIdentifyingUser   = "Hubo un error identificando al usuario..."
)

type GDFamilyUnityBot struct {
	TelegramBot *telebot.Bot
	userClient  *client.UserServiceClient
}

func NewGDFamilyUnityBot(configs *config.EnvVars, userClient *client.UserServiceClient) (*GDFamilyUnityBot, error) {
	pref := telebot.Settings{
		Token:  configs.GDUnityFamilyToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return &GDFamilyUnityBot{
		TelegramBot: bot,
		userClient:  userClient,
	}, nil
}

func (g *GDFamilyUnityBot) UploadImage(c telebot.Context) error {
	// Get the photo information
	photo := c.Message().Photo

	// Generate a unique filename
	messageTime := c.Message().Time()
	filename := fmt.Sprintf("%s_%d.jpg", messageTime.Format(_RFC3339OnlyDateFormat), c.Message().ID)

	// Open the file.
	file, err := g.TelegramBot.File(photo.MediaFile())
	if err != nil {
		log.Println("err: ", err)
		return c.Send(_errReadingImage)
	}

	// Create a new multipart writer.
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Create a multipart form file.
	formFile, err := writer.CreateFormFile("tg-bot-file", filepath.Base(filename))
	if err != nil {
		log.Println("err: ", err)
		return c.Send(_errProcessingMessage)
	}

	// Copy the file content to the form file.
	_, err = io.Copy(formFile, file)
	if err != nil {
		log.Println("err: ", err)
		return c.Send(_errProcessingMessage)
	}

	// Send username in order to identify the sender.
	err = writer.WriteField("username", c.Message().Sender.Username)
	if err != nil {
		log.Println("err: ", err)
		return c.Send(_errIdentifyingUser)
	}

	// Send bot name in order to identify the bot sender.
	err = writer.WriteField("bot_name", "GDUnityFamilyBot")
	if err != nil {
		log.Println("err: ", err)
		return c.Send(_errIdentifyingUser)
	}

	// Send date in order to identify the folder.
	err = writer.WriteField("date", messageTime.Format(_RFC3339OnlyYearFormat))
	if err != nil {
		log.Println("err: ", err)
		return c.Send(_errIdentifyingUser)
	}

	// Send date in order to identify the folder.
	err = writer.WriteField("filename", filename)
	if err != nil {
		log.Println("err: ", err)
		return c.Send(_errIdentifyingUser)
	}

	if err = writer.Close(); err != nil {
		log.Println("err: ", err)
		return c.Send(_errProcessingMessage)
	}

	// Using a client to obtain the response of the service.
	if err = g.userClient.PostPhoto(buffer, writer); err != nil {
		log.Println("err: ", err)
		return c.Send(
			"La imagen no puede ser procesada debido a un error del servicio, por favor contactarse con el dev...",
		)
	}

	return c.Send(fmt.Sprintf("¡Imagen guardada! Nombre guardado de la imagen: **%s**", filename))
}
