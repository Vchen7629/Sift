package service

import "encoding/json"

// todo: update the userId in the tui and api to be username for clarity
type RequestBody struct {
	UserId   string `json:"userId"`
	RepoName string `json:"repoName"`
}

func MarshalRequestBody(username, repoName string) ([]byte, error) {
	payload := RequestBody{UserId: username, RepoName: repoName}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
