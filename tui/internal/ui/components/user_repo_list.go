package components

import (
	"image/color"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tui/internal/ui/context"
)

type UserRepoListModel struct {
	ctx 		 *context.App
	FetchedRepos []UserRepo
	viewport 	 viewport.Model
}

type UserRepo struct {
	id, Name, Status, LastIndexed, TotalLibs string
}

func NewUserRepoList(ctx *context.App) *UserRepoListModel {
	return &UserRepoListModel{
		ctx: ctx,
		FetchedRepos: []UserRepo{
			{id: "id1", Name: "react", Status: "indexed", LastIndexed: "1746", TotalLibs: "42"},
			{id: "id2", Name: "next.js", Status: "indexed", LastIndexed: "1745", TotalLibs: "18"},
			{id: "id3", Name: "tailwindcss", Status: "pending", LastIndexed: "0", TotalLibs: "5"},
			{id: "id4", Name: "react", Status: "indexed", LastIndexed: "1746", TotalLibs: "420"},
			{id: "id5", Name: "next.js", Status: "indexed", LastIndexed: "1745", TotalLibs: "18"},
			{id: "id6", Name: "tailwindcss", Status: "pending", LastIndexed: "0", TotalLibs: "5"},
			{id: "id7", Name: "react", Status: "indexed", LastIndexed: "1746", TotalLibs: "42"},
			{id: "id8", Name: "next.js", Status: "indexed", LastIndexed: "1745", TotalLibs: "18"},
			{id: "id9", Name: "tailwindcss", Status: "pending", LastIndexed: "0", TotalLibs: "5"},
			{id: "id10", Name: "react", Status: "indexed", LastIndexed: "1746", TotalLibs: "42"},
			{id: "id11", Name: "next.js", Status: "indexed", LastIndexed: "1745", TotalLibs: "18"},
			{id: "id12", Name: "tailwindcss", Status: "pending", LastIndexed: "0", TotalLibs: "5"},
			{id: "id13", Name: "react", Status: "indexed", LastIndexed: "1746", TotalLibs: "42"},
			{id: "id14", Name: "next.js", Status: "indexed", LastIndexed: "1745", TotalLibs: "18"},
			{id: "id15", Name: "tailwindcss", Status: "pending", LastIndexed: "0", TotalLibs: "5"},
		},
	}
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

	m.viewport.SetWidth(m.ctx.ViewPortWidth)
	m.viewport.SetHeight(m.ctx.ViewPortHeight - 4)
	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, cards...))
	return tea.NewView(m.viewport.View())
}

func (m *UserRepoListModel) repoCard(repo UserRepo) string {
	header := m.repoCardHeader(repo)

	background := lipgloss.NewStyle().Background(color.Transparent)
		
	card := background.Render(lipgloss.NewStyle().
		Width(m.ctx.ViewPortWidth).
		PaddingLeft(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Padding(0, 1).
		Render(header))

	return card
}

func (m *UserRepoListModel) repoCardHeader(repo UserRepo) string {
	repoName := lipgloss.NewStyle().Render(repo.Name)

	indexStatus := lipgloss.NewStyle().Width(8).Align(lipgloss.Right).Render(repo.Status)                                                                       
	lastIndexed := lipgloss.NewStyle().Width(6).Align(lipgloss.Right).Render(repo.LastIndexed)
	totalLibs   := lipgloss.NewStyle().Width(5).Align(lipgloss.Right).Render(repo.TotalLibs)

	right := lipgloss.JoinHorizontal(lipgloss.Top, indexStatus, lastIndexed, totalLibs)

	spacer := lipgloss.NewStyle().
		Width(m.ctx.ViewPortWidth - lipgloss.Width(repoName) - lipgloss.Width(right) - 4).
		Render("")

	return lipgloss.JoinHorizontal(lipgloss.Top, repoName, spacer, right)
}