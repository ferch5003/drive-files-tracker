package rpc

import (
	"log"
)

type FamilyPayload struct {
	Photo    []byte
	Username string
	Date     string
}

// UploadDriveFile uploads a file into Drive given the file and the file ID.
func (s *Server) UploadDriveFile(payload FamilyPayload, resp *string) error {
	log.Printf("Payload: %+v \n", payload)

	*resp = "Processed payload via RPC:"

	return nil
}
