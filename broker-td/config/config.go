package config

import (
	"os"
)

type EnvVars struct {
	Port string

	// BaseURLs
	UserServiceBaseURL  string
	DriveServiceBaseRPC string
}

func NewConfigurations() (*EnvVars, error) {
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
