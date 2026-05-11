package main

import (
	"fmt"
	"os"

	"tui/internal/ui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	app, err := ui.New()
	if err != nil {
		fmt.Printf("error initializing: %v", err)
		os.Exit(1)
	}

	p := tea.NewProgram(&app)
	_, err = p.Run()
	if err != nil {
		fmt.Printf("error launching: %v", err)
		os.Exit(1)
	}
}
