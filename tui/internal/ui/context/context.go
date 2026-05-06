package context

import (
	"tui/internal/ui/styles"
)

type Page int

const (
	AuthPage Page = iota
	QueryPage
	UserReposPage
)

type App struct {
	Width, Height int
	ViewPortWidth, ViewPortHeight int
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
