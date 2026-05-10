package api

import (
	"tui/internal/types"

	gh "github.com/cli/go-gh/v2/pkg/api"
)

type GithubClient struct {
	rest *gh.RESTClient
}

func NewGithubClient() (*GithubClient, error) {
	rest, err := gh.DefaultRESTClient()
	if err != nil {
		return nil, err
	}

	return &GithubClient{rest: rest}, nil
}

// requires the gh cli tool installed and authenticated
func (c *GithubClient) GithubUsername() (string, error) {
	var response struct {
		Login string `json:"login"`
	}
	err := c.rest.Get("user", &response)

	return response.Login, err
}

type RepoApiRes struct {  
	Name        string `json:"name"`
	Description string `json:"description"`
	LastCommit  string `json:"pushed_at"`                                                                                                             
}

func (c *GithubClient) GithubUserRepositories() ([]types.GHRepository, error){
	var apiRes []RepoApiRes
	err := c.rest.Get("user/repos?affiliation=owner&per_page=100", &apiRes)
	if err != nil {
		return nil, err
	}

	res := make([]types.GHRepository, 0, len(apiRes))
	for i, repo := range apiRes {
		res = append(res, types.GHRepository{Id: i, Name: repo.Name, Description: repo.Description, LastCommit: repo.LastCommit})
	}

	return res, err
}