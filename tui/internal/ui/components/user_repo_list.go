package components

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tui/internal/ui/context"
)

type UserRepoListModel struct {
	Ctx 		 *context.App
	FetchedRepos []UserRepo
	viewport 	 viewport.Model
}

type UserRepo struct {
	Name, Status, LastIndexed, TotalLibs string
}

func (m UserRepoListModel) Init() tea.Cmd {
	return nil
}

func (m *UserRepoListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "down":
			m.viewport.ScrollDown(2)
		case "up":
			m.viewport.ScrollUp(2)
		}
	}

	return m, nil
}

func (m *UserRepoListModel) View() tea.View {
	var cards []string

	for _, repo := range m.FetchedRepos {
		cards = append(cards, m.repoCard(repo))
	}

	m.viewport.SetWidth(m.Ctx.ViewPortWidth)
	m.viewport.SetHeight(m.Ctx.ViewPortHeight)
	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, cards...))
	return tea.NewView(m.viewport.View())
}

func (m *UserRepoListModel) repoCard(repo UserRepo) string {
	header := m.repoCardHeader(repo)

		
	card := lipgloss.NewStyle().
		Width(m.Ctx.ViewPortWidth).
		PaddingLeft(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(0, 1).
		Render(header)

	return card
}

func (m *UserRepoListModel) repoCardHeader(repo UserRepo) string {
	repoName := lipgloss.NewStyle().Render(repo.Name)

	indexStatus := lipgloss.NewStyle().Width(8).Align(lipgloss.Right).Render(repo.Status)                                                                       
	lastIndexed := lipgloss.NewStyle().Width(6).Align(lipgloss.Right).Render(repo.LastIndexed)
	totalLibs   := lipgloss.NewStyle().Width(5).Align(lipgloss.Right).Render(repo.TotalLibs)

	right := lipgloss.JoinHorizontal(lipgloss.Top, indexStatus, lastIndexed, totalLibs)

	spacer := lipgloss.NewStyle().
		Width(m.Ctx.ViewPortWidth - lipgloss.Width(repoName) - lipgloss.Width(right) - 4).
		Render("")

	return lipgloss.JoinHorizontal(lipgloss.Top, repoName, spacer, right)
}