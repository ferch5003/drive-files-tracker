package config

import (
	"broker-td/internal/platform/files"
	"github.com/joho/godotenv"
	"os"
)

type EnvVars struct {
	Port string

	// BaseURLs
	UserServiceBaseURL  string
	DriveServiceBaseRPC string
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

	userServiceBaseURL := os.Getenv("USER_SERVICE_BASE_URL")
	driveServiceBaseRPC := os.Getenv("DRIVE_SERVICE_BASE_RPC")

	environment := &EnvVars{
		Port: port,

		UserServiceBaseURL:  userServiceBaseURL,
		DriveServiceBaseRPC: driveServiceBaseRPC,
	}

	return environment, nil
}
