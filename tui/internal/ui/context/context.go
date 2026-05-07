package context

import (
	"tui/internal/ui/styles"
)

type Page int

const (
	QueryPage Page = iota
	UserReposPage
	ThemePage
)

type App struct {
	WindowWidth, WindowHeight, MainWidth, MainHeight, SidebarWidth int
	Username	  string
	Theme		  styles.Theme
	CurrentPage   Page
}

func NewApp() *App {
	return &App{
		CurrentPage: UserReposPage,
		Theme:		 styles.Warm,
	}
}
