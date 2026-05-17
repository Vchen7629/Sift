package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

var authBaseUrl = "http://localhost:8080/auth"

var ErrUnauthorized = errors.New("unauthorized")

// this calls the backend to generate a new session for authentication
func NewSession(ghToken string) (string, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/token", authBaseUrl), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ghToken))
	req.Header.Set("User-Agent", "Sift-tui/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println("error closing auth user response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected error sending req: %d", resp.StatusCode)
	}

	var sessionToken string
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "JSESSIONID" {
			sessionToken = cookie.Value
			break
		}
	}

	if sessionToken == "" {
		return "", fmt.Errorf("no session token in response")
	}

	return sessionToken, nil
}
