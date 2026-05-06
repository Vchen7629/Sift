package views

import tea "charm.land/bubbletea/v2"

type UserRepoModel struct {
	width, height int
}

func (m UserRepoModel) Init() tea.Cmd {
	return nil
}

func (m UserRepoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m UserRepoModel) View() tea.View {
	return tea.NewView("IF OUR LOVEEEEEEE IS TRAGEDY WHY ARE YOU MY REMEDYYYYYYYYY")
}