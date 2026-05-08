package context

import (
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
}

func NewApp() *App {
	return &App{
		CurrentPage:   UserReposPage,
		SelectedTheme: styles.Warm,
	}
}
