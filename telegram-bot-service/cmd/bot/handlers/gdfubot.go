package handlers

import (
	"bytes"
	"encoding/json"
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
	_errProcessingService = "Hubo un error procesando los datos del servicio..."
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

	// Send date in order to identify the folder.
	err = writer.WriteField("date", messageTime.Format(_RFC3339OnlyDateFormat))
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("err: ", err)
		return c.Send(_errProcessingService)
	}

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println("err: ", err)
		return c.Send(_errProcessingService)
	}

	dataErr, ok := data["error"]
	if ok {
		log.Println("err: ", dataErr)
		return c.Send(
			"La imagen no puede ser procesada debido a un error del servicio, por favor contactarse con el dev...",
		)
	}

	return c.Send(fmt.Sprintf("¡Imagen guardada! Nombre guardado de la imagen: **%s**", filename))
}
