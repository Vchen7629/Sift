package api

import gh "github.com/cli/go-gh/v2/pkg/api"

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

func (c *GithubClient) GithubUserRepositories() ([]RepoApiRes, error) {
	var apiRes []RepoApiRes
	err := c.rest.Get("user/repos?affiliation=owner&per_page=100", &apiRes)
	if err != nil {
		return nil, err
	}

	return apiRes, err
}
