package handlers

import (
	"bytes"
	"fmt"
	"gopkg.in/telebot.v3"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"telegram-bot-service/config"
	"time"
)

const (
	_RFC3339OnlyDateFormat = "2006-01-02"
)

const (
	_errReadingImage      = "Hubo un error leyendo la imagen..."
	_errProcessingMessage = "Hubo un error procesando la imagen..."
	_errIdentifyingUser   = "Hubo un error identificando al usuario..."
	_errConnectingService = "Hubo un error conectandose al servicio..."
)

type GDFamilyUnityBot struct {
	TelegramBot *telebot.Bot
}

func NewGDFamilyUnityBot(configs *config.EnvVars) (*GDFamilyUnityBot, error) {
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

	if err = writer.Close(); err != nil {
		log.Println("err: ", err)
		return c.Send(_errProcessingMessage)
	}

	// Making the request to outsource service to process.
	req, err := http.NewRequest(
		http.MethodPost,
		"http://broker-td/gdrive-family-uploader",
		bytes.NewReader(buffer.Bytes()))
	if err != nil {
		log.Println("err: ", err)
		return c.Send(_errConnectingService)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Using a client to obtain the response of the service.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("err: ", err)
		return c.Send(_errProcessingMessage)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("err: ", err)
		}
	}(resp.Body)

	return c.Send(fmt.Sprintf("¡Imagen guardada! Nombre guardado de la imagen: **%s**", filename))
}
