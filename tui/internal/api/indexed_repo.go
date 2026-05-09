package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var indexedRepoBaseUrl = "http://localhost:8081/user_repo"

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

func GetAllIndexedRepos(username string) ([]string, error) {
	resp, err := client.Get(fmt.Sprintf("%s/list/%s", indexedRepoBaseUrl, username))
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected error sending req: %d", resp.StatusCode)
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var repos []string
	err = json.Unmarshal(res, &repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}