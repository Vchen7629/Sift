package views

import tea "charm.land/bubbletea/v2"

type QueryModel struct {
	width, height int
}

func (m QueryModel) Init() tea.Cmd {
	return nil
}

// user actions
func (m QueryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m QueryModel) View() tea.View {
	return tea.NewView("IF OUR LOVEEEEEEE IS TRAGEDY WHY ARE YOU MY REMEDYYYYYYYYY")
}