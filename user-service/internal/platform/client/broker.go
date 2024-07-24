package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type BrokerClient struct {
	baseURL string
}

func NewBrokerClient(baseURL string) *BrokerClient {
	return &BrokerClient{
		baseURL: baseURL,
	}
}

func (c *BrokerClient) PostFolderParentsCreations(b []byte) (map[string]any, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		c.baseURL+"/gdrive-folder-creator",
		bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Using a client to obtain the response of the service.
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

	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	dataErr, ok := data["error"]
	if ok {
		return nil, errors.New(fmt.Sprintf("%s", dataErr))
	}

	return data, nil
}
