package config

import (
	"github.com/joho/godotenv"
	"os"
	"telegram-bot-service/internal/platform/files"
)

type EnvVars struct {
	// Telegram Data
	GDUnityFamilyToken string
}

func NewConfigurations() (*EnvVars, error) {
	area := os.Getenv("AREA")

	if area == "" {
		envFilepath, err := files.GetFile(".env")
		if err != nil {
			return nil, err
		}

		if err := godotenv.Load(envFilepath); err != nil {
			return nil, err
		}
	}

	gdUnityFamilyToken := os.Getenv("GD_UNITY_FAMILY_TOKEN")

	environment := &EnvVars{
		GDUnityFamilyToken: gdUnityFamilyToken,
	}

	return environment, nil
}
