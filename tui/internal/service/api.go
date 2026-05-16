package service

import "encoding/json"

// todo: update the userId in the tui and api to be username for clarity
type RequestBody struct {
	RepoName string `json:"repoName"`
}

func MarshalRequestBody(repoName string) ([]byte, error) {
	payload := RequestBody{RepoName: repoName}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
