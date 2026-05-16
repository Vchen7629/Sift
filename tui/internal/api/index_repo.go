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

func IndexRepo(sessionToken, repoName string) error {
	log.Printf("got session token %s for index repo", sessionToken)
	payload, err := service.MarshalRequestBody(repoName)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/add", indexBaseUrl), bytes.NewBuffer(payload))
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
			log.Println("error closing index repo resp body")
		}
	}()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected error sending req: %d", resp.StatusCode)
	}

	return nil
}

// used to poll the job status for the progress bar
func GetJobStatus(sessionToken, repoName string) (string, error) {
	payload, err := service.MarshalRequestBody(repoName)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/job_status", indexBaseUrl), bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Sift-tui/1.0")
	req.AddCookie(&http.Cookie{
		Name:  "JSESSIONID",
		Value: sessionToken,
	})

	resp, err := client.Do(req)
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
		return "", fmt.Errorf("reading job status response: %w", err)
	}

	return string(res), err
}
