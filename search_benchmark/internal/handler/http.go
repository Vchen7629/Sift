package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// field from api res
type SearchResult struct {
	URL string `json:"url"`
}

func CallSearchEndpoint(endpoint, query, userId string) ([]SearchResult, error) {
	bodyBytes, err := requestBodyBuilder(userId, query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	err = json.Unmarshal(bodyBytes, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func requestBodyBuilder(userId, query string) ([]byte, error) {
	reqBody := struct{
		UserID string `json:"userId"`
		Query  string `json:"searchQuery"`
	} {
		UserID: userId,
		Query: query,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}