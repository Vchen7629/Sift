package components

import (
	"image/color"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tui/internal/ui/context"
	"tui/internal/ui/styles"
)

type UserRepoListModel struct {
	ctx 		 *context.App
	FetchedRepos []UserRepo
	FocusedRepo  FocusedRepo
	viewport 	 viewport.Model
}

type FocusedRepo struct {
	index    int
	userRepo UserRepo
}

type UserRepo struct {
	id, Name, Status, LastIndexed, TotalLibs string
}

func NewUserRepoList(ctx *context.App) *UserRepoListModel {
	m := &UserRepoListModel{
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
	m.FocusedRepo = FocusedRepo{index: 0, userRepo: m.FetchedRepos[0]}
	return m
}

func (m UserRepoListModel) Init() tea.Cmd {
	return nil
}

func (m *UserRepoListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "down":
			if m.FocusedRepo.index < len(m.FetchedRepos) - 1 {
				m.FocusedRepo.index++
				m.FocusedRepo.userRepo = m.FetchedRepos[m.FocusedRepo.index]
				m.scrollToFocused()
			}
		case "up":
			if m.FocusedRepo.index > 0 {
				m.FocusedRepo.index--
				m.FocusedRepo.userRepo = m.FetchedRepos[m.FocusedRepo.index]
				m.scrollToFocused()
			}
		}
	}

	return m, nil
}

func (m *UserRepoListModel) View() tea.View {
	var cards []string

	for _, repo := range m.FetchedRepos {
		cards = append(cards, m.repoCard(repo))
	}

	m.viewport.SetWidth(m.ctx.RepoListWidth)
	m.viewport.SetHeight(m.ctx.RepoListHeight - 4)
	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, cards...))
	return tea.NewView(m.viewport.View())
}

func (m *UserRepoListModel) repoCard(repo UserRepo) string {
	header := m.repoCardHeader(repo)

	borderColor, _ := m.focusedStyle(repo)
		
	card := lipgloss.NewStyle().
		Width(m.ctx.RepoListWidth).
		PaddingLeft(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		Render(header)

	return card
}

func (m *UserRepoListModel) repoCardHeader(repo UserRepo) string {
	_, textColor := m.focusedStyle(repo)

	repoName := lipgloss.NewStyle().Foreground(textColor).Render(repo.Name)

	indexStatus := lipgloss.NewStyle().Width(8).Align(lipgloss.Right).Render(repo.Status)                                                                       
	lastIndexed := lipgloss.NewStyle().Width(6).Align(lipgloss.Right).Render(repo.LastIndexed)
	totalLibs   := lipgloss.NewStyle().Width(5).Align(lipgloss.Right).Render(repo.TotalLibs)

	right := lipgloss.JoinHorizontal(lipgloss.Top, indexStatus, lastIndexed, totalLibs)

	spacer := lipgloss.NewStyle().
		Width(m.ctx.RepoListWidth - lipgloss.Width(repoName) - lipgloss.Width(right) - 4).
		Render("")

	return lipgloss.JoinHorizontal(lipgloss.Top, repoName, spacer, right)
}

func (m *UserRepoListModel) scrollToFocused() {
	card := m.repoCard(m.FocusedRepo.userRepo)                                                                                                              
    cardHeight := lipgloss.Height(card)

	// this is to calculate the lines the curr focused card occupies
	itemTop    := m.FocusedRepo.index * cardHeight                                                                                                       
	itemBottom := itemTop + cardHeight

	viewTop    := m.viewport.YOffset()
	viewBottom := viewTop + m.viewport.Height()

	if itemBottom > viewBottom {
		// item went below visible area, scroll down just enough
		m.viewport.SetYOffset(itemBottom - m.viewport.Height())
	} else if itemTop < viewTop {
		// item went above visible area, scroll up just enough
		m.viewport.SetYOffset(itemTop)
	}
}

func (m *UserRepoListModel) focusedStyle(repo UserRepo) (color.Color, color.Color) {
	borderColor := lipgloss.Color("#444444")
	textColor := lipgloss.Color("#ffffff")

	if m.FocusedRepo.userRepo.id == repo.id {
		borderColor = styles.Warm.AccentMid
		textColor = styles.Warm.AccentBright
	}

	return borderColor, textColor
}