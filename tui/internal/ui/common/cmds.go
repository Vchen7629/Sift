package common

import (
	"errors"
	"os"
	"os/exec"
	"tui/internal/api"
	"tui/internal/types"

	tea "charm.land/bubbletea/v2"
	"github.com/pkg/browser"
)

type ToggleFocusMsg struct{}

type FetchIndexedRepoMsg struct {
	IndexedRepos    []types.IndexedRepo
	NewSessionToken string
}
type FetchIndexedRepoErr struct{ Err error }

func FetchIndexedRepo(sessionToken string) tea.Cmd {
	return func() tea.Msg {
		indexRepos, err := api.GetAllIndexedRepos(sessionToken)

		if errors.Is(err, api.ErrUnauthorized) {
			ghToken := api.GithubPatToken()

			newSessionToken, err := api.NewSession(ghToken)
			if err != nil {
				return FetchIndexedRepoErr{Err: err}
			}

			indexRepos, err = api.GetAllIndexedRepos(newSessionToken)
			if err != nil {
				return FetchIndexedRepoErr{Err: err}
			}

			return FetchIndexedRepoMsg{IndexedRepos: indexRepos, NewSessionToken: newSessionToken}
		}
		if err != nil {
			return FetchIndexedRepoErr{Err: err}
		}

		return FetchIndexedRepoMsg{IndexedRepos: indexRepos, NewSessionToken: sessionToken}
	}
}

type BrowserOpenedMsg struct{ Err error }

func OpenInBrowser(url string) tea.Cmd {
	return func() tea.Msg {
		// need this workaround for wsl2 to windows case, this requires wslu installed aswell
		// sudo apt install wslu
		if _, wslErr := os.Stat("/proc/sys/fs/binfmt_misc/WSLInterop"); wslErr == nil {
			if path, err := exec.LookPath("wslview"); err == nil {
				err = exec.Command(path, url).Start()
				return BrowserOpenedMsg{Err: err}
			}
		}
		err := browser.OpenURL(url)
		return BrowserOpenedMsg{Err: err}
	}
}
