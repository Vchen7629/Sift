package views

import (
	tea "charm.land/bubbletea/v2"

	"tui/internal/ui/context"
)

type AuthModel struct {
	Ctx *context.App
}

func (m AuthModel) Init() tea.Cmd {
	return nil
}

// user actions
func (m AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m AuthModel) View() tea.View {
	return tea.NewView("IF OUR LOVEEEEEEE IS TRAGEDY WHY ARE YOU MY REMEDYYYYYYYYY")
}