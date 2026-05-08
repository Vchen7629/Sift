package views

import (
	"strings"
	"tui/internal/ui/components/user_repo"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type FocusedPanel int

const (
	FocusList FocusedPanel = iota
	FocusSidebar
)

type UserRepoModel struct {
	Ctx 		   *context.App
	SearchBar 	   *user_repo.SearchBarModel
	RepoList	   *user_repo.ListModel
	Sidebar 	   *user_repo.Sidebar
	focusedRepo    user_repo.FocusedRepo
	panelFocus 	   FocusedPanel
}

func (m UserRepoModel) Init() tea.Cmd {
	return nil
}

func (m *UserRepoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "s" && !m.SearchBar.IsSearching() {
			if m.panelFocus == FocusList {
				m.panelFocus = FocusSidebar
			} else {
				m.panelFocus = FocusList
			}

			return m, nil
		}
	}
	
	var repoListCmd, searchBarCmd, sidebarCmd tea.Cmd
	if m.panelFocus == FocusList {
		_, repoListCmd = m.RepoList.Update(msg)
		searchBarCmd = m.SearchBar.Update(msg)
	} else {
		_, sidebarCmd = m.Sidebar.Update(msg)
	}

	m.focusedRepo = m.RepoList.Focused
	m.Sidebar.FocusedRepo = m.focusedRepo

	return m, tea.Batch(repoListCmd, searchBarCmd, sidebarCmd)
}

func (m *UserRepoModel) View() tea.View {
	if m.Ctx.WindowWidth == 0 {
		return tea.NewView("")
	}

	dividerLine := strings.Repeat("│\n", m.Ctx.MainHeight - 1) + "│"
	divider := lipgloss.NewStyle().Foreground(styles.Divider).Render(dividerLine)

	repoListContent := lipgloss.JoinVertical(lipgloss.Top, m.SearchBar.View(), m.RepoList.View().Content)
	content := lipgloss.JoinHorizontal(lipgloss.Left, repoListContent, divider, m.Sidebar.View().Content)

	return tea.NewView(lipgloss.JoinVertical(lipgloss.Left, 
		m.userRepoListActionBar().Content, 
		content,
	))
}

// todo: maybe rename it to something other than btn since the user doesnt really click it
func (m UserRepoModel) userRepoListActionBar() tea.View {
	navBtnStyle := lipgloss.NewStyle().PaddingLeft(2)

	navBtnTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444444")).
		Bold(true)

	navBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↑↓] navigate"))
	searchBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↵] search"))
	clearSearchBtn := navBtnStyle.Render(navBtnTextStyle.Render("[esc] clear"))
	swapFocusBtn := navBtnStyle.Render(navBtnTextStyle.Render("[s] swap focus"))
	reindexBtn := navBtnStyle.Render(navBtnTextStyle.Render("[r] reindex"))

	return tea.NewView(lipgloss.NewStyle().
		BorderBottom(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderBottomForeground(styles.Divider).
		Width(m.Ctx.WindowWidth - 2).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, navBtn, searchBtn, clearSearchBtn, swapFocusBtn, reindexBtn)))
}