package config

import (
	"os"
)

type EnvVars struct {
	// Telegram Data
	GDUnityFamilyToken  string
	GDGardenFamilyToken string
	GDOSCommercialToken string

	// BaseURLs
	BrokerTDBaseURL string
	OrionStateURL   string
}

func NewConfigurations() (*EnvVars, error) {
	gdUnityFamilyToken := os.Getenv("GD_UNITY_FAMILY_TOKEN")
	gdGardenFamilyToken := os.Getenv("GD_GARDEN_FAMILY_TOKEN")
	gdOSCommercialToken := os.Getenv("GD_OS_COMMERCIAL_TOKEN")

	brokerTDBaseURL := os.Getenv("BROKER_TD_BASE_URL")
	orionStateURL := os.Getenv("ORION_STATE_URL")

	environment := &EnvVars{
		GDUnityFamilyToken:  gdUnityFamilyToken,
		GDGardenFamilyToken: gdGardenFamilyToken,
		GDOSCommercialToken: gdOSCommercialToken,

		BrokerTDBaseURL: brokerTDBaseURL,
		OrionStateURL:   orionStateURL,
	}

	return environment, nil
}
