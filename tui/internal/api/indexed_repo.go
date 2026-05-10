package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"tui/internal/service"
	"tui/internal/types"
)

var indexedRepoBaseUrl = "http://localhost:8080/user_repo"

type DeleteIndexedRepoReq struct {
	UserId   string `json:"userId"`
	RepoName string `json:"repoName"`
}

func DeleteIndexedRepo(username, repoName string) error {
	payload := DeleteIndexedRepoReq{UserId: username, RepoName: repoName}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/delete", indexedRepoBaseUrl), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
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

func GetAllIndexedRepos(username string) ([]types.IndexedRepo, error) {
	resp, err := client.Get(fmt.Sprintf("%s/list/%s", indexedRepoBaseUrl, username))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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

	var indexedRepos []types.IndexedRepo
	for i, repos := range repos {
		for j := range repos.Dependencies {
			repos.Dependencies[j].Id = j
		}
		indexedRepo := types.IndexedRepo{
			Id: i, TotalDependencies: repos.TotalDependencies,
			Name:         strings.Split(repos.Name, "/")[1],
			LastIndexed:  service.FormatRelativeDate(repos.LastIndexed),
			Dependencies: repos.Dependencies,
		}
		indexedRepos = append(indexedRepos, indexedRepo)
	}

	return indexedRepos, nil
}