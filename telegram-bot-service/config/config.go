package config

import (
	"github.com/joho/godotenv"
	"os"
	"telegram-bot-service/internal/platform/files"
)

type EnvVars struct {
	// Telegram Data
	GDUnityFamilyToken  string
	GDGardenFamilyToken string
	GDOSCommercialToken string

	// BaseURLs
	BrokerTDBaseURL string
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
	gdGardenFamilyToken := os.Getenv("GD_GARDEN_FAMILY_TOKEN")
	gdOSCommercialToken := os.Getenv("GD_OS_COMMERCIAL_TOKEN")
	brokerTDBaseURL := os.Getenv("BROKER_TD_BASE_URL")

	environment := &EnvVars{
		GDUnityFamilyToken:  gdUnityFamilyToken,
		GDGardenFamilyToken: gdGardenFamilyToken,
		GDOSCommercialToken: gdOSCommercialToken,

		BrokerTDBaseURL: brokerTDBaseURL,
	}

	return environment, nil
}
