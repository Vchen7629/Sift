package common

import (
	"tui/internal/api"
	"tui/internal/types"

	tea "charm.land/bubbletea/v2"
)

type FetchIndexedRepoMsg struct { IndexedRepos []types.IndexedRepo }

func FetchIndexedRepo(username string) tea.Cmd {
	return func() tea.Msg {
		indexRepos, err := api.GetAllIndexedRepos(username)
		if err != nil {
			return err
		}

		return FetchIndexedRepoMsg{ IndexedRepos: indexRepos }
	}
}

type FetchIndexJobStatusMsg struct { Status string }

func FetchIndexJobStatus(username, repoName string) tea.Cmd {
	return func() tea.Msg {
		status, err := api.GetJobStatus(username, repoName)
		if err != nil {
			return err
		}

		return FetchIndexJobStatusMsg{ Status: status}
	}
}