package views

import (
	"strings"
	"tui/internal/ui/components/user_repo"
	"tui/internal/ui/context"

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
	FocusedSidebar *user_repo.FocusedSidebar
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
		_, sidebarCmd = m.FocusedSidebar.Update(msg)
	}

	m.focusedRepo = m.RepoList.FocusedRepo
	m.FocusedSidebar.FocusedRepo = m.focusedRepo

	return m, tea.Batch(repoListCmd, searchBarCmd, sidebarCmd)
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
		BorderBottomForeground(lipgloss.Color("#444444")).
		Width(m.Ctx.Width - 2).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, navBtn, searchBtn, clearSearchBtn, swapFocusBtn, reindexBtn)))
}