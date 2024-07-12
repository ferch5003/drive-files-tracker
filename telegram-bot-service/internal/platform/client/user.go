package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

type UserServiceClient struct {
	baseURL string
}

func NewUserServiceClient(baseURL string) *UserServiceClient {
	return &UserServiceClient{
		baseURL: baseURL,
	}
}

func (c *UserServiceClient) GetUsernames() ([]string, error) {
	// Making the request to outsource service to process.
	req, err := http.NewRequest(
		http.MethodGet,
		c.baseURL+"/users",
		nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("err: ", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var users []struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	}

	if err := json.Unmarshal(body, &users); err != nil {
		return nil, err
	}

	usernames := make([]string, 0)
	for _, user := range users {
		usernames = append(usernames, user.Username)
	}

	return usernames, nil
}

func (c *UserServiceClient) PostPhoto(buffer bytes.Buffer, writer *multipart.Writer) error {
	req, err := http.NewRequest(
		http.MethodPost,
		c.baseURL+"/gdrive-family-uploader",
		bytes.NewReader(buffer.Bytes()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Using a client to obtain the response of the service.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("err: ", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	dataErr, ok := data["error"]
	if ok {
		return errors.New(fmt.Sprintf("%s", dataErr))
	}

	return nil
}
