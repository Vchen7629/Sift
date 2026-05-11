package api

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"tui/internal/service"
)

type SearchReq struct {
	UserId      string `json:"userId"`
	SearchQuery string `json:"searchQuery"`
}

var searchBaseUrl = "http://localhost:8080/search"

func Search(username, searchQuery string) (string, error) {
	payload, err := service.MarshalRequestBody(username, searchQuery)
	if err != nil {
		return "", err
	}

	resp, err := client.Post(fmt.Sprintf("%s/add", searchBaseUrl), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println("error closing search resp body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected error sending req: %d", resp.StatusCode)
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading search response: %w", err)
	}

	return string(res), err
}
