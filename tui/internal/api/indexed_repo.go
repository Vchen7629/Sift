package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"tui/internal/service"
	"tui/internal/types"
)

var indexedRepoBaseUrl = "http://localhost:8080/user_repo"

func DeleteIndexedRepo(sessionToken, repoName string) error {
	payload, err := service.MarshalRequestBody(repoName)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/delete", indexedRepoBaseUrl), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Sift-tui/1.0")
	req.AddCookie(&http.Cookie{
		Name:  "JSESSIONID",
		Value: sessionToken,
	})

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println("error closing delete indexed repo resp body")
		}
	}()

	switch resp.StatusCode {
	case http.StatusForbidden:
		return ErrUnauthorized
	case http.StatusNoContent: // do nothing when no content resp
	default:
		return fmt.Errorf("unexpected error sending req: %d", resp.StatusCode)
	}

	return nil
}

type fetchIndexedRepoResp struct {
	Name              string             `json:"repoName"`
	LastIndexed       string             `json:"lastIndexed"`
	TotalDependencies int                `json:"totalDependencies"`
	Dependencies      []types.Dependency `json:"dependencies"`
}

func GetAllIndexedRepos(sessionToken string) ([]types.IndexedRepo, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/list", indexedRepoBaseUrl), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Sift-tui/1.0")
	req.AddCookie(&http.Cookie{
		Name:  "JSESSIONID",
		Value: sessionToken,
	})

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println("error closing get all indexed repos resp body")
		}
	}()

	switch resp.StatusCode {
	case http.StatusForbidden:
		return nil, ErrUnauthorized
	case http.StatusOK: // do nothing when ok
	default:
		return nil, fmt.Errorf("unexpected error sending req: %d", resp.StatusCode)
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var repos []fetchIndexedRepoResp
	err = json.Unmarshal(res, &repos)
	if err != nil {
		return nil, err
	}

	indexedRepos := make([]types.IndexedRepo, 0, len(repos))
	for _, repos := range repos {
		for j := range repos.Dependencies {
			repos.Dependencies[j].Id = j
		}

		_, name, _ := strings.Cut(repos.Name, "/")

		indexedRepo := types.IndexedRepo{
			TotalDependencies: repos.TotalDependencies,
			Name:              name,
			LastIndexed:       service.FormatRelativeDate(repos.LastIndexed),
			Dependencies:      repos.Dependencies,
		}
		indexedRepos = append(indexedRepos, indexedRepo)
	}

	return indexedRepos, nil
}
