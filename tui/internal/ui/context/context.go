package context

import (
	"tui/internal/api"
	"tui/internal/ui/styles"
)

type Page int

const (
	QueryPage Page = iota
	UserReposPage
)

type App struct {
	WindowWidth, WindowHeight int
	MainWidth, MainHeight     int
	SidebarWidth              int
	SessionToken              string
	Username                  string
	CurrentPage               Page
	ThemeSelectorOpen         bool
	SelectedTheme             styles.Theme
	GithubApiClient           *api.GithubClient
}

func NewApp() (*App, error) {
	ghToken := api.GithubPatToken()

	SessionToken, err := api.NewSession(ghToken)
	if err != nil {
		return nil, err
	}

	client, err := api.NewGithubClient()
	if err != nil {
		return nil, err
	}

	username, err := client.GithubUsername()
	if err != nil {
		return nil, err
	}

	return &App{
		CurrentPage:     UserReposPage,
		SelectedTheme:   styles.Warm,
		GithubApiClient: client,
		Username:        username,
		SessionToken:    SessionToken,
	}, nil
}
