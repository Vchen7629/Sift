package user_repo

import (
	"image/color"
	"strconv"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"tui/internal/service"
	"tui/internal/types"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"
)

type ListModel struct {
	ctx 		 *context.App
	FetchedRepos []types.Repository
	FocusedIdx   int
	viewport 	 viewport.Model
}

func NewUserRepoList(ctx *context.App) *ListModel {
	m := &ListModel{
		ctx: ctx,
		FetchedRepos: []types.Repository{},
	}
	m.FocusedIdx = 0
	return m
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m *ListModel) Update(msg tea.Msg, isSidebarFocused bool) tea.Cmd {                                                                                                            
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if isSidebarFocused || len(m.FetchedRepos) == 0 {
			break
		}
		cardHeight := lipgloss.Height(m.repoCard(m.FocusedIdx, m.FetchedRepos[m.FocusedIdx])) 
		switch msg.String() {
		case "down":
			if m.FocusedIdx < len(m.FetchedRepos) - 1 {
				m.FocusedIdx++
    			service.ScrollToFocused(&m.viewport, m.FocusedIdx, cardHeight)
			}
		case "up":
			if m.FocusedIdx > 0 {
				m.FocusedIdx--
				service.ScrollToFocused(&m.viewport, m.FocusedIdx, cardHeight)
			}
		}
	}

	return nil
}

func (m *ListModel) View() tea.View {
	var cards []string
	
	if len(m.FetchedRepos) > 0 {
		for i, repo := range m.FetchedRepos {
			cards = append(cards, m.repoCard(i, repo))
		}

		m.viewport.SetWidth(m.ctx.MainWidth)
		m.viewport.SetHeight(m.ctx.MainHeight - 4)
		m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, cards...))

		return tea.NewView(m.viewport.View())
	}

	fetchingPlaceholder := lipgloss.NewStyle().Padding(1, 2).Render("Loading your repos...")
	return tea.NewView(fetchingPlaceholder)
}

func (m *ListModel) repoCard(idx int, repo types.Repository) string {
	header := m.repoCardHeader(idx, repo)

	borderColor, _ := m.focusedStyle(idx)
		
	card := lipgloss.NewStyle().
		Width(m.ctx.MainWidth).
		PaddingLeft(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		Render(header)

	return card
}

func (m *ListModel) repoCardHeader(idx int, repo types.Repository) string {
	_, textColor := m.focusedStyle(idx)

	repoName := lipgloss.NewStyle().Foreground(textColor).Render(repo.Name)

	indexStatus := lipgloss.NewStyle().Width(8).Align(lipgloss.Right).Render(repo.Status)                                                                       
	lastIndexed := lipgloss.NewStyle().Width(6).Align(lipgloss.Right).Render(repo.LastIndexed)
	totalDependencies := lipgloss.NewStyle().Width(5).Align(lipgloss.Right).Render(strconv.Itoa(repo.TotalDependencies))

	right := lipgloss.JoinHorizontal(lipgloss.Top, indexStatus, lastIndexed, totalDependencies)

	spacer := lipgloss.NewStyle().
		Width(m.ctx.MainWidth - lipgloss.Width(repoName) - lipgloss.Width(right) - 4).
		Render("")

	return lipgloss.JoinHorizontal(lipgloss.Top, repoName, spacer, right)
}

func (m *ListModel) focusedStyle(idx int) (color.Color, color.Color) {
	if m.FocusedIdx == idx {
		return m.ctx.SelectedTheme.AccentMid, m.ctx.SelectedTheme.AccentBright
	}

	return styles.Divider, lipgloss.Color("#ffffff")
}
