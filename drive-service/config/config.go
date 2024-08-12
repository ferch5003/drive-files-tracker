package config

import (
	"os"
)

type EnvVars struct {
	// Server Environment.
	Port string

	SAPrivateKeyID      string
	SAPrivateKey        string
	SAClientEmail       string
	SAClientID          string
	SAClientX509CertURL string
}

func NewConfigurations() (*EnvVars, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	saPrivateKeyID := os.Getenv("SA_PRIVATE_KEY_ID")
	saPrivateKey := os.Getenv("SA_PRIVATE_KEY")
	saClientEmail := os.Getenv("SA_CLIENT_EMAIL")
	saClientID := os.Getenv("SA_CLIENT_ID")
	saClientX509CertURL := os.Getenv("SA_CLIENT_X509_CERT_URL")

	environment := &EnvVars{
		Port: port,

		SAPrivateKeyID:      saPrivateKeyID,
		SAPrivateKey:        saPrivateKey,
		SAClientEmail:       saClientEmail,
		SAClientID:          saClientID,
		SAClientX509CertURL: saClientX509CertURL,
	}

	return environment, nil
}
