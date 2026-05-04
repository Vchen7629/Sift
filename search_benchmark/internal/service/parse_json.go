package service

import (
	"encoding/json"
	"os"
)

type IssueRecord struct {
	ID		string `json:"id"`
	Queries []string `json:"queries"`
}

// extract the eval json to return a hashmap with ID key Queries Value 
func ExtractIssueQueries(evalFilePath string) (map[string][]string, error) {
	var issues []IssueRecord

	f, err := os.ReadFile(evalFilePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(f, &issues)
	if err != nil {
		return nil, err
	}

	res := make(map[string][]string)

	for _, r := range issues {
		res[r.ID] = r.Queries;
	}

	return res, nil
}