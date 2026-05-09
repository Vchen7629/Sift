package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SearchReq struct {
	UserId      string `json:"userId"`
	SearchQuery string `json:"searchQuery"`
}

var searchBaseUrl = "http://localhost:8080/search"

func Search(username, searchQuery string) (string, error) {
	payload := SearchReq{UserId: username, SearchQuery: searchQuery}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := client.Post(fmt.Sprintf("%s/add", searchBaseUrl), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected error sending req: %d", resp.StatusCode)
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	return string(res), err
}