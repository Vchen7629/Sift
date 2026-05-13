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
	RepoName 	 string 	   `json:"repoName"`
	NumSources   int 		   `json:"numSources"`
	IssueSources []IssueSource `json:"issues"`
	Summary		 string 	   `json:"summary"`
}

type IssueSource struct {
	Url  		   string  `json:"url"`
	Title 		   string  `json:"title"`
	Body 		   string  `json:"body"`
	RelevanceScore float32 `json:"rerankScore"`
}

var searchBaseUrl = "http://localhost:8080/search"

func Search(username, repoName, searchQuery string) (SearchRes, error) {
	repoNameFmt := fmt.Sprintf("%s/%s", username, repoName)

	payload := SearchReq{UserId: username, RepoName: repoNameFmt, SearchQuery: searchQuery}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return SearchRes{}, err
	}

	resp, err := client.Post(fmt.Sprintf("%s/new", searchBaseUrl), "application/json", bytes.NewBuffer(jsonData))
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
