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
	SidebarWidth 		      int
	Username	  			  string
	CurrentPage   			  Page
	ThemeSelectorOpen		  bool
	SelectedTheme			  styles.Theme
	GithubApiClient 		  *api.GithubClient
}

func NewApp() (*App, error) {
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
		Username: username,
	}, nil
}