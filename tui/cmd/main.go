package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"tui/internal/ui"
)


func main() {
	p := tea.NewProgram(ui.New())
	_, err := p.Run()
	if err != nil {
		fmt.Printf("error launching: %v", err)
		os.Exit(1)
	}
}