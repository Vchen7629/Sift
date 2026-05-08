package views

import (
	"strings"
	"tui/internal/ui/components/user_repo"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type UserRepoModel struct {
	Ctx 		      *context.App
	ActionBar         *user_repo.ActionBarModel
	SearchBar 	      *user_repo.SearchBarModel
	RepoList	      *user_repo.ListModel
	Sidebar 	      *user_repo.Sidebar
	focusedRepo       user_repo.FocusedRepo
	isSidebarFocused bool
}

func (m UserRepoModel) Init() tea.Cmd {
	return nil
}

func (m *UserRepoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case user_repo.ToggleFocusMsg:
		if !m.SearchBar.IsSearching() {
			m.isSidebarFocused = !m.isSidebarFocused
		}
		return m, nil
	}

	actionBarCmd := m.ActionBar.Update(msg)
	repoListCmd := m.RepoList.Update(msg, m.isSidebarFocused)
	searchBarCmd := m.SearchBar.Update(msg, m.isSidebarFocused)
	sidebarCmd := m.Sidebar.Update(msg, m.isSidebarFocused)

	m.focusedRepo = m.RepoList.Focused
	m.Sidebar.FocusedRepo = m.focusedRepo

	return m, tea.Batch(actionBarCmd, repoListCmd, searchBarCmd, sidebarCmd)
}

func (m *UserRepoModel) View() tea.View {
	if m.Ctx.WindowWidth == 0 {
		return tea.NewView("")
	}

	dividerLine := strings.Repeat("│\n", m.Ctx.MainHeight - 1) + "│"
	divider := lipgloss.NewStyle().Foreground(styles.Divider).Render(dividerLine)

	repoListContent := lipgloss.JoinVertical(lipgloss.Top, m.SearchBar.View(), m.RepoList.View().Content)
	content := lipgloss.JoinHorizontal(lipgloss.Left, repoListContent, divider, m.Sidebar.View().Content)

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, m.ActionBar.View(m.isSidebarFocused).Content, content))
}