package driveaccount

import (
	"context"
	"drive-service/config"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func (s ServiceAccount) CreateFolder(
	service *drive.Service,
	name string,
	mimeType string,
	parentID string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentID},
	}

	folder, err := service.Files.Create(f).Do()
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func (s ServiceAccount) RemoveFolder(
	service *drive.Service,
	folderID string) error {
	return service.Files.Delete(folderID).Do()
}

func (s ServiceAccount) WriteOnSheet(
	service *sheets.Service,
	spreadsheetID,
	spreadsheetGID,
	spreadsheetColumn,
	monthRow,
	text string) error {
	// Get spreadsheet metadata
	spreadsheet, err := service.Spreadsheets.Get(spreadsheetID).Do()
	if err != nil {
		return err
	}

	// Loop through tabs to find the tab name by GID
	var tabName string
	for _, sheet := range spreadsheet.Sheets {
		if fmt.Sprintf("%d", sheet.Properties.SheetId) == spreadsheetGID {
			tabName = sheet.Properties.Title
			break
		}
	}

	if tabName == "" {
		return err
	}

	cell := fmt.Sprintf("%s!%s", tabName, spreadsheetColumn+monthRow)
	resp, err := service.Spreadsheets.Values.Get(spreadsheetID, cell).Do()
	if err != nil {
		return err
	}

	cellValue := fmt.Sprintf("%s", resp.Values[0][0])
	if cellValue != "" {
		pastValue, err := getMoneyFromTextOnFloat64(cellValue)
		if err != nil {
			return err
		}

		actualValue, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return err
		}

		text = strings.ReplaceAll(
			fmt.Sprintf("%f", pastValue+actualValue),
			".",
			",",
		)
	}

	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{
			{text}, // This writes the value to the specified cell
		},
	}

	_, err = service.Spreadsheets.Values.Update(spreadsheetID, cell, valueRange).
		ValueInputOption("USER_ENTERED").
		Do()
	if err != nil {
		return err
	}

	return nil
}

func getMoneyFromTextOnFloat64(text string) (float64, error) {
	dollarIndex := strings.Index(text, "$") // Find the index of the dollar sign

	if dollarIndex != -1 {
		// Extract the substring starting from the character after the dollar sign
		money := strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(text[dollarIndex+1:], " ", ""),
				".", "",
			),
			",",
			".",
		)

		value, err := strconv.ParseFloat(money, 64)
		if err != nil {
			return 0, err
		}

		return value, nil
	}

	return 0, errors.New("no Money found")
}
