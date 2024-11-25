package rpc

import (
	"drive-service/internal/platform/driveaccount"
	"drive-service/internal/platform/files"
	"fmt"
	"log"
	"os"
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
	filePath, err := files.CreateFile("tmp", payload.Filename, payload.Photo)
	if err != nil {
		return err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	file, err := s.ServiceAccount.CreateFile(
		s.DriveService,
		payload.Filename,
		driveaccount.BaseApplicationOctetStream,
		f,
		payload.FolderID)
	if err != nil {
		return err
	}

	if err := os.Remove(filePath); err != nil {
		return err
	}

	*resp = fmt.Sprintf("File '%s' uploaded successfully", file.Name)

	return nil
}

type BotUser struct {
	BotID    int
	UserID   int
	Date     string
	FolderID string
	IsParent bool
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
	for _, botUser := range payload.BotUsers {
		folder, err := s.ServiceAccount.CreateFolder(
			s.DriveService,
			currentYear,
			driveaccount.BaseApplicationVNDGoogleAppsFolder,
			botUser.FolderID)
		if err != nil {
			removeFolderErr := s.removeFolders(newBotUsers)
			if removeFolderErr != nil {
				return removeFolderErr
			}

			return err
		}

		newBotUser := BotUser{
			BotID:    botUser.BotID,
			UserID:   botUser.UserID,
			Date:     currentYear,
			FolderID: folder.Id,
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
