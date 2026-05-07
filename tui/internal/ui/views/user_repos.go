package views

import (
	"strings"
	"tui/internal/ui/components/user_repo"
	"tui/internal/ui/context"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type UserRepoModel struct {
	Ctx 		   *context.App
	SearchBar 	   *user_repo.SearchBarModel
	RepoList	   *user_repo.ListModel
	FocusedSidebar *user_repo.FocusedSidebar
	FocusedRepo    user_repo.FocusedRepo
}

func (m UserRepoModel) Init() tea.Cmd {
	return m.SearchBar.Init()
}

func (m *UserRepoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	_, repoListCmd := m.RepoList.Update(msg)
	searchBarCmd := m.SearchBar.Update(msg)
	
	m.FocusedRepo = m.RepoList.FocusedRepo
	m.FocusedSidebar.FocusedRepo = m.FocusedRepo

	return m, tea.Batch(repoListCmd, searchBarCmd)
}

func (m *UserRepoModel) View() tea.View {
	if m.Ctx.Width == 0 {
		return tea.NewView("")
	}

	dividerLine := strings.Repeat("│\n", m.Ctx.RepoListHeight-1) + "│"
	divider := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444444")).
		Render(dividerLine)

	repoListContent := lipgloss.JoinVertical(lipgloss.Top, m.SearchBar.View(), m.RepoList.View().Content)
	content := lipgloss.JoinHorizontal(lipgloss.Left, repoListContent, divider, m.FocusedSidebar.View().Content)

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, 
		m.userActionBar().Content, 
		content,
	))
}

func (m UserRepoModel) userActionBar() tea.View {
	navBtnStyle := lipgloss.NewStyle().PaddingLeft(2)

	navBtnTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444444")).
		Bold(true)

	navBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↑↓] navigate"))
	searchBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↵] search"))
	focusDetailsBtn := navBtnStyle.Render(navBtnTextStyle.Render("[d] focus details"))
	reindexBtn := navBtnStyle.Render(navBtnTextStyle.Render("[r] reindex"))

	return tea.NewView(lipgloss.NewStyle().
		BorderBottom(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderBottomForeground(lipgloss.Color("#444444")).
		Width(m.Ctx.Width - 2).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, navBtn, searchBtn, focusDetailsBtn, reindexBtn)))
}