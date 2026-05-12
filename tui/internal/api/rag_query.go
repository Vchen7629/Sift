package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SearchReq struct {
	UserId      string `json:"userId"`
	RepoName 	string `json:"repoName"`
	SearchQuery string `json:"searchQuery"`
}

type SearchRes struct {
	IssueSources []issueSource `json:"issues"`
	Summary		 string 	   `json:"summary"`
}

type issueSource struct {
	Url  		   string  `json:"url"`
	Title 		   string  `json:"title"`
	Body 		   string  `json:"body"`
	RelevanceScore float32 `json:"rerankScore"`
}

var searchBaseUrl = "http://localhost:8080/search"

func Search(username, repoName, searchQuery string) (SearchRes, error) {
	payload := SearchReq{UserId: username, RepoName: repoName, SearchQuery: searchQuery}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return SearchRes{}, err
	}

	resp, err := client.Post(fmt.Sprintf("%s/add", searchBaseUrl), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return SearchRes{}, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("error closing search resp body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return SearchRes{}, fmt.Errorf("unexpected error sending req: %d", resp.StatusCode)
	}

	var result SearchRes
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return SearchRes{}, fmt.Errorf("error decoding search response: %w", err)
	}

	return result, err
}
