package user_repo

import (
	"image/color"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tui/internal/ui/context"
	"tui/internal/service"
)

type ListModel struct {
	ctx 		 *context.App
	FetchedRepos []UserRepo
	Focused      FocusedRepo
	viewport 	 viewport.Model
}

type FocusedRepo struct {
	index    int
	userRepo UserRepo
}

type UserRepo struct {
	id, Name, Status, LastIndexed, TotalDependencies, Description string
	Dependencies []DependencyStatus
}

func NewUserRepoList(ctx *context.App) *ListModel {
	m := &ListModel{
		ctx: ctx,
		FetchedRepos: dummyData,
	}
	m.Focused = FocusedRepo{index: 0, userRepo: m.FetchedRepos[0]}
	return m
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cardHeight := lipgloss.Height(m.repoCard(m.Focused.userRepo))                                                                                                             

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "down":
			if m.Focused.index < len(m.FetchedRepos) - 1 {
				m.Focused.index++
				m.Focused.userRepo = m.FetchedRepos[m.Focused.index]
    			service.ScrollToFocused(&m.viewport, m.Focused.index, cardHeight)
			}
		case "up":
			if m.Focused.index > 0 {
				m.Focused.index--
				m.Focused.userRepo = m.FetchedRepos[m.Focused.index]
				service.ScrollToFocused(&m.viewport, m.Focused.index, cardHeight)
			}
		}
	}

	return m, nil
}

func (m *ListModel) View() tea.View {
	var cards []string

	for _, repo := range m.FetchedRepos {
		cards = append(cards, m.repoCard(repo))
	}

	m.viewport.SetWidth(m.ctx.MainWidth)
	m.viewport.SetHeight(m.ctx.MainHeight - 4)
	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, cards...))
	return tea.NewView(m.viewport.View())
}

func (m *ListModel) repoCard(repo UserRepo) string {
	header := m.repoCardHeader(repo)

	borderColor, _ := m.focusedStyle(repo)
		
	card := lipgloss.NewStyle().
		Width(m.ctx.MainWidth).
		PaddingLeft(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		Render(header)

	return card
}

func (m *ListModel) repoCardHeader(repo UserRepo) string {
	_, textColor := m.focusedStyle(repo)

	repoName := lipgloss.NewStyle().Foreground(textColor).Render(repo.Name)

	indexStatus := lipgloss.NewStyle().Width(8).Align(lipgloss.Right).Render(repo.Status)                                                                       
	lastIndexed := lipgloss.NewStyle().Width(6).Align(lipgloss.Right).Render(repo.LastIndexed)
	totalDependencies := lipgloss.NewStyle().Width(5).Align(lipgloss.Right).Render(repo.TotalDependencies)

	right := lipgloss.JoinHorizontal(lipgloss.Top, indexStatus, lastIndexed, totalDependencies)

	spacer := lipgloss.NewStyle().
		Width(m.ctx.MainWidth - lipgloss.Width(repoName) - lipgloss.Width(right) - 4).
		Render("")

	return lipgloss.JoinHorizontal(lipgloss.Top, repoName, spacer, right)
}

func (m *ListModel) focusedStyle(repo UserRepo) (color.Color, color.Color) {
	borderColor := lipgloss.Color("#444444")
	textColor := lipgloss.Color("#ffffff")

	if m.Focused.userRepo.id == repo.id {
		borderColor = m.ctx.SelectedTheme.AccentMid
		textColor = m.ctx.SelectedTheme.AccentBright
	}

	return borderColor, textColor
}