package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var indexBaseUrl = "http://localhost:8080/index_repo"
var client = &http.Client{
	Timeout: 10 * time.Second,
}

// todo: update the userId in the tui and api to be username for clarity
type IndexRepoReq struct {
	UserId   string `json:"userId"`
	RepoName string `json:"repoName"` 
}

func IndexRepo(username, repoName string) error {
	payload := IndexRepoReq{UserId: username, RepoName: repoName}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := client.Post(fmt.Sprintf("%s/add", indexBaseUrl), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected error sending req: %d", resp.StatusCode)
	}

	return nil
}

// used to poll the job status for the progress bar
func GetJobStatus(username, repoName string) (string, error) {
	resp, err := client.Get(fmt.Sprintf("%s/get_status/%s/%s", indexBaseUrl, username, repoName))
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


