package config

import (
	"broker-td/internal/platform/files"
	"github.com/joho/godotenv"
	"os"
)

type EnvVars struct {
	Port string
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

	environment := &EnvVars{
		Port: port,
	}

	return environment, nil
}
