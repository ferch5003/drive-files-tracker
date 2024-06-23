package client

import (
	"encoding/json"
	"io"
	"log"
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
