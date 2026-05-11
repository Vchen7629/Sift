package common

import (
	"tui/internal/api"
	"tui/internal/types"

	tea "charm.land/bubbletea/v2"
)

type FetchIndexedRepoMsg struct{ IndexedRepos []types.IndexedRepo }
type FetchIndexedRepoErr struct{ Err error }

func FetchIndexedRepo(username string) tea.Cmd {
	return func() tea.Msg {
		indexRepos, err := api.GetAllIndexedRepos(username)
		if err != nil {
			return FetchIndexedRepoErr{Err: err}
		}

		return FetchIndexedRepoMsg{IndexedRepos: indexRepos}
	}
}
