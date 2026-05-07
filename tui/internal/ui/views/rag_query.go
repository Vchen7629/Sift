package views

import (
	tea "charm.land/bubbletea/v2"

	"tui/internal/ui/context"
)

type RagQueryModel struct {
	Ctx *context.App
}

func (m RagQueryModel) Init() tea.Cmd {
	return nil
}

// user actions
func (m RagQueryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m RagQueryModel) View() tea.View {
	return tea.NewView("IF OUR LOVEEEEEEE IS TRAGEDY WHY ARE YOU MY REMEDYYYYYYYYY")
}