package rag_query

import (
	"fmt"
	"tui/internal/ui/common"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ActionBarModel struct {
	ctx 			*context.App
	hasSearchResult bool
}

func NewActionBar(ctx *context.App) *ActionBarModel {
	return &ActionBarModel{
		ctx: ctx, hasSearchResult: false,
	}
}

func (m *ActionBarModel) Init() tea.Cmd {
	return nil
}

func (m *ActionBarModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "s":
			return func() tea.Msg { return common.ToggleFocusMsg{} }
		}

	case NewSearchQueryMsg:
		m.hasSearchResult = true
	}
	return nil
}

func (m *ActionBarModel) View(isRepoListFocused, isSearching bool, selectedRepo string) tea.View {
	selectedRepoName := fmt.Sprintf("Selected Repo: %s", selectedRepo)
	if selectedRepo == "" {
		selectedRepoName = "No Repo Selected"
	}

	selectedRepoText := lipgloss.NewStyle().
		PaddingRight(1).Foreground(m.ctx.SelectedTheme.AccentBright).PaddingLeft(2).
		Render(selectedRepoName)

	return tea.NewView(styles.ActionBarBorder.
		Width(m.ctx.WindowWidth - 2).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, selectedRepoText, m.actionBarBtns(isRepoListFocused, isSearching))),
	)
}

func (m *ActionBarModel) actionBarBtns(focusRepoList, isSearching bool) string {
	switch {
	case !m.hasSearchResult:
		return lipgloss.JoinHorizontal(lipgloss.Left, styles.NavBtn("[/] new query"))
	case focusRepoList:
		return lipgloss.JoinHorizontal(lipgloss.Left,
			styles.NavBtn("[/] new query"), styles.NavBtn("[↑↓] change selected repo"), styles.NavBtn("[s] back to query card"),
			styles.NavBtn("[↵] select repo"),
		)
	case isSearching:
		return lipgloss.JoinHorizontal(lipgloss.Left,
			styles.NavBtn("[/] cancel search"), styles.NavBtn("[esc] clear search query"), styles.NavBtn("[↵] search"),
		)
	default:
		return lipgloss.JoinHorizontal(lipgloss.Left,
			styles.NavBtn("[/] new query"), styles.NavBtn("[↑↓] scroll sources"), styles.NavBtn("[s] switch selected repo"),
			styles.NavBtn("[↵] open in browser"),
		)
	}
}
