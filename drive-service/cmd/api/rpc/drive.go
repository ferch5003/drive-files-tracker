package rpc

import (
	"drive-service/internal/platform/driveaccount"
	"drive-service/internal/platform/files"
	"fmt"
	"os"
)

type FamilyPayload struct {
	Photo    []byte
	FolderID string
	Filename string
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

	*resp = fmt.Sprintf("File '%s' uploaded successfully", file.Name)

	return nil
}
