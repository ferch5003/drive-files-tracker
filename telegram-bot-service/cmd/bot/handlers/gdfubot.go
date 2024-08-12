package handlers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"gopkg.in/telebot.v3"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"telegram-bot-service/config"
	"telegram-bot-service/internal/platform/client"
	"time"
)

const _gdUnityFamilyBotName = "GDUnityFamilyBot"

type GDFamilyUnityBot struct {
	TelegramBot   *telebot.Bot
	userClient    *client.UserServiceClient
	orionStateURL string
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
		TelegramBot:   bot,
		userClient:    userClient,
		orionStateURL: configs.OrionStateURL,
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
	err = writer.WriteField("bot_name", _gdUnityFamilyBotName)
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

func (g *GDFamilyUnityBot) GetPropertyAccountStatement(c telebot.Context) error {
	// Default 30 seconds of timeout, after that returns error for context deadline.
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Grouping Headless, No first run, disable GPU and No sandbox option.
	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoSandbox,
	)
	allocCtx, cancel := chromedp.NewExecAllocator(ctxWithTimeout, options...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	messageTime := c.Message().Time()
	now := messageTime.Format(_RFC3339OnlyDateFormat)

	log.Printf("Processing image at: %s\n", now)

	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(g.orionStateURL, "body", &buf)); err != nil {
		log.Println("err: ", err)
		return c.Send("El servicio esta presentando inconvenientes, por favor intente mas tarde...")
	}

	reader := bytes.NewReader(buf)

	return c.Send(&telebot.Photo{
		File:    telebot.FromReader(reader),
		Caption: fmt.Sprintf("Estado de cuenta predial generado el %s", now),
	})
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(url, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(sel, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("Element visible:", sel)
			return nil
		}),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}
