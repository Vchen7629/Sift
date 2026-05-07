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
	Theme		  	   	 	  styles.Theme
	CurrentPage   			  Page
	ThemeSelectorOpen		  bool
}

func NewApp() *App {
	return &App{
		CurrentPage: UserReposPage,
		Theme:		 styles.Warm,
	}
}
