package driveaccount

import (
	"context"
	"drive-service/config"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
	"io"
	"net/http"
	"os"
)

type ServiceAccount struct {
	secretFile string
	configs    *config.EnvVars
}

func NewServiceAccount(secretFile string, configs *config.EnvVars) ServiceAccount {
	return ServiceAccount{
		secretFile: secretFile,
		configs:    configs,
	}
}

// Get the Service Account Client to operate.
func (s ServiceAccount) Get() (*http.Client, error) {
	b, err := os.ReadFile(s.secretFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	var credentials = struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
	}{}
	if err := json.Unmarshal(b, &credentials); err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	credentials.Email = s.configs.SAClientEmail
	credentials.PrivateKey = s.configs.SAPrivateKey

	jwtConfigs := &jwt.Config{
		Email:      credentials.Email,
		PrivateKey: []byte(credentials.PrivateKey),
		Scopes: []string{
			drive.DriveScope,
		},
		TokenURL: google.JWTTokenURL,
	}

	return jwtConfigs.Client(context.Background()), nil
}

func (s ServiceAccount) CreateFile(
	service *drive.Service,
	name string,
	mimeType string,
	content io.Reader,
	parentID string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentID},
	}

	file, err := service.Files.Create(f).Media(content).Do()
	if err != nil {
		return nil, err
	}

	return file, nil
}
