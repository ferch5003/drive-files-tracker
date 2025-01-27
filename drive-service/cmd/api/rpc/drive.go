package rpc

import (
	"drive-service/internal/platform/driveaccount"
	"drive-service/internal/platform/files"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

const _RFC3339OnlyYearFormat = "2006"

type FamilyPayload struct {
	Photo             []byte
	FolderID          string
	Filename          string
	Username          string
	SpreadsheetID     string
	SpreadsheetGID    string
	SpreadsheetColumn string
}

// UploadDriveFile uploads a file into Drive given the file and the file ID.
func (s *Server) UploadDriveFile(payload FamilyPayload, resp *string) error {
	// Uploading file to Drive Process.
	filePath, err := files.CreateFile("tmp", payload.Filename, payload.Photo)
	if err != nil {
		log.Println(err)
		return err
	}

	f, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return err
	}

	file, err := s.ServiceAccount.CreateFile(
		s.DriveService,
		payload.Filename,
		driveaccount.BaseApplicationOctetStream,
		f,
		payload.FolderID)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := os.Remove(filePath); err != nil {
		log.Println(err)
		return err
	}

	*resp = fmt.Sprintf("File '%s' uploaded successfully", file.Name)

	// Write value of payment in Spreadsheet Process.
	if !slices.Contains([]string{
		payload.SpreadsheetID,
		payload.SpreadsheetGID,
		payload.SpreadsheetColumn,
	}, "") {
		location, err := time.LoadLocation("America/Bogota")
		if err != nil {
			log.Println(err)
			return err
		}

		if err := s.OCRClient.SetImageFromBytes(payload.Photo); err != nil {
			log.Println(err)
			return err
		}

		text, err := s.OCRClient.Text()
		if err != nil {
			log.Println(err)
			return fmt.Errorf("text: %v", err)
		}

		actualMoney, err := getMoneyFromText(text)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("getMoneyFromText: %v", err)
		}

		currentMonth := time.Now().In(location).Month()
		if err := s.ServiceAccount.WriteOnSheet(
			s.SheetService,
			payload.SpreadsheetID,
			payload.SpreadsheetGID,
			payload.SpreadsheetColumn,
			fmt.Sprintf("%d", 2+int(currentMonth)), // Default space to row is 2, so it begin on 3 cells.
			actualMoney); err != nil {
			log.Println(err)
			return fmt.Errorf("WriteOnSheet: %v", err)
		}
	}

	return nil
}

func getMoneyFromText(text string) (string, error) {
	dollarIndex := strings.Index(text, "$") // Find the index of the dollar sign

	if dollarIndex != -1 {
		fromMoneySign := text[dollarIndex+1:]
		textsFromMoneySign := strings.Split(fromMoneySign, "\n")

		// Extract the substring starting from the character after the dollar sign
		money := strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(textsFromMoneySign[0], " ", ""),
				".", "",
			),
			",",
			".",
		)

		value, err := strconv.ParseFloat(money, 64)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%f", value), nil
	}

	return "", errors.New("no Money found")
}

type BotUser struct {
	BotID              int
	UserID             int
	Date               string
	FolderID           string
	IsParent           bool
	SpreadsheetID      string
	SpreadsheetGID     string
	SpreadsheetBaseGID string
	SpreadsheetColumn  string
}

type BotUsersPayload struct {
	BotUsers []BotUser
}

// CreateYearlyFolders creates folders with a year format in the indicated parent folder and returns a slice of BotUser
// with the new folder created.
func (s *Server) CreateYearlyFolders(payload BotUsersPayload, resp *BotUsersPayload) error {
	location, err := time.LoadLocation("America/Bogota")
	if err != nil {
		log.Println(err)
		return err
	}

	newBotUsers := make([]BotUser, 0)
	currentYear := time.Now().In(location).Format(_RFC3339OnlyYearFormat)

	botSheetsGID := make(map[string]string)
	for _, botUser := range payload.BotUsers {
		folder, err := s.ServiceAccount.CreateFolder(
			s.DriveService,
			currentYear,
			driveaccount.BaseApplicationVNDGoogleAppsFolder,
			botUser.FolderID)
		if err != nil {
			log.Println(err)
			return removeData(s, newBotUsers, err)
		}

		// Write value of payment in Spreadsheet Process.
		newSpreadsheetGID, ok := botSheetsGID[botUser.SpreadsheetID]
		if !ok && !slices.Contains([]string{
			botUser.SpreadsheetID,
			botUser.SpreadsheetBaseGID,
			botUser.SpreadsheetColumn,
		}, "") {
			newSpreadsheetGID, err = s.ServiceAccount.CopySheetToSameSpreadsheet(
				s.SheetService,
				botUser.SpreadsheetID,
				botUser.SpreadsheetBaseGID,
				currentYear,
			)
			if err != nil {
				log.Println(err)
				return removeData(s, newBotUsers, err)
			}

			// Save saved sheet to not duplicate previous one.
			botSheetsGID[botUser.SpreadsheetID] = newSpreadsheetGID
		}

		newBotUser := BotUser{
			BotID:             botUser.BotID,
			UserID:            botUser.UserID,
			Date:              currentYear,
			FolderID:          folder.Id,
			SpreadsheetID:     botUser.SpreadsheetID,
			SpreadsheetGID:    newSpreadsheetGID,
			SpreadsheetColumn: botUser.SpreadsheetColumn,
		}

		newBotUsers = append(newBotUsers, newBotUser)
	}

	response := BotUsersPayload{
		BotUsers: newBotUsers,
	}

	*resp = response

	return nil
}

func (s *Server) removeFolders(botUsers []BotUser) error {
	for _, botUser := range botUsers {
		if err := s.ServiceAccount.RemoveFolder(s.DriveService, botUser.FolderID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) removeSheets(botUsers []BotUser) error {
	for _, botUser := range botUsers {
		if err := s.ServiceAccount.RemoveSheetOnSpreadsheet(
			s.SheetService, botUser.SpreadsheetID, botUser.SpreadsheetGID); err != nil {
			return err
		}
	}

	return nil
}

func removeData(s *Server, botUsers []BotUser, originalErr error) error {
	removeFolderErr := s.removeFolders(botUsers)
	if removeFolderErr != nil {
		return removeFolderErr
	}

	removeSheetErr := s.removeSheets(botUsers)
	if removeSheetErr != nil {
		return removeSheetErr
	}

	return originalErr
}
