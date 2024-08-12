package config

import (
	"os"
)

type EnvVars struct {
	// Server Environment.
	Port          string
	IsDevelopment bool
	ActivateCRON  bool

	// Database Environment.
	PostgreDSN string

	// BaseURLs
	BrokerTDBaseURL string
}

func NewConfigurations() (*EnvVars, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	isDevelopment := os.Getenv("IS_DEVELOPMENT") == "true"
	activateCRON := os.Getenv("ACTIVATE_CRON") == "true"

	postgreDSN := os.Getenv("POSTGRE_DSN")
	brokerTDBaseURL := os.Getenv("BROKER_TD_BASE_URL")

	environment := &EnvVars{
		Port:          port,
		IsDevelopment: isDevelopment,
		ActivateCRON:  activateCRON,
		PostgreDSN:    postgreDSN,

		BrokerTDBaseURL: brokerTDBaseURL,
	}

	return environment, nil
}
