package handlers

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"os"
	"telegram-bot-service/config"
	"telegram-bot-service/internal/platform/files"
	"time"
)

const (
	RFC3339OnlyDateFormat = "2006-01-02"
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

func (g *GDFamilyUnityBot) GetPhoto(c telebot.Context) error {
	// Get the photo information
	photo := c.Message().Photo

	tmpDir, err := files.GetDir("tmp")
	if err != nil {
		return fmt.Errorf("error downloading photo: %w", err)
	}

	// Check If username directory exist, if not, created.
	username := c.Message().Sender.Username
	userDir, err := files.GetDir(fmt.Sprintf("%s/%s", tmpDir, username))
	if os.IsNotExist(err) {
		userDir, err = files.CreateDir(tmpDir, username)
		if err != nil {
			return fmt.Errorf("error downloading photo: %w", err)
		}
	}

	// Generate a unique filename
	messageTime := c.Message().Time()
	filename := fmt.Sprintf("%s_%d.jpg", messageTime.Format(RFC3339OnlyDateFormat), c.Message().ID)

	// Generate filepath to download.
	filepath := fmt.Sprintf("%s/%s", userDir, filename)

	// Download the photo using File and local filename
	err = g.TelegramBot.Download(photo.MediaFile(), filepath)
	if err != nil {
		return fmt.Errorf("error downloading photo: %w", err)
	}

	return c.Send(fmt.Sprintf("¡Imagen guardada! Nombre guardado de la imagen: **%s**", filename))
}
