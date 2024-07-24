package config

import (
	"github.com/joho/godotenv"
	"os"
	"user-service/internal/platform/files"
)

type EnvVars struct {
	// Server Environment.
	Port          string
	IsDevelopment bool

	// Database Environment.
	PostgreDSN string

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	isDevelopment := os.Getenv("IS_DEVELOPMENT") == "true"

	postgreDSN := os.Getenv("POSTGRE_DSN")
	brokerTDBaseURL := os.Getenv("BROKER_TD_BASE_URL")

	environment := &EnvVars{
		Port:          port,
		IsDevelopment: isDevelopment,
		PostgreDSN:    postgreDSN,

		BrokerTDBaseURL: brokerTDBaseURL,
	}

	return environment, nil
}
