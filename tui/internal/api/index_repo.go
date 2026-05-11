package api

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"tui/internal/service"
)

var indexBaseUrl = "http://localhost:8080/index_repo"
var client = &http.Client{
	Timeout: 10 * time.Second,
}

func IndexRepo(username, repoName string) error {
	payload, err := service.MarshalRequestBody(username, repoName)
	if err != nil {
		return err
	}

	resp, err := client.Post(fmt.Sprintf("%s/add", indexBaseUrl), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println("error closing index repo resp body")
		}
	}()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected error sending req: %d", resp.StatusCode)
	}

	return nil
}

// used to poll the job status for the progress bar
func GetJobStatus(username, repoName string) (string, error) {
	payload, err := service.MarshalRequestBody(username, repoName)
	if err != nil {
		return "", err
	}

	resp, err := client.Post(fmt.Sprintf("%s/job_status", indexBaseUrl), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println("error closing get job status resp body")
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected error sending req: %d", resp.StatusCode)
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	return string(res), err
}
