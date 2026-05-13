package rag_query

import (
	"fmt"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ActionBarModel struct {
	ctx 		*context.App
}

type ToggleFocusMsg struct{}

func NewActionBar(ctx *context.App) *ActionBarModel {
	return &ActionBarModel{ctx: ctx}
}

func (m *ActionBarModel) Init() tea.Cmd {
	return nil
}

func (m *ActionBarModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "s":
			return func() tea.Msg { return ToggleFocusMsg{} }
		}
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
	btn := func(text string) string {
		return styles.NavBtnStyle.Render(styles.NavBtnTextStyle.Render(text))
	}

	switch {
	case focusRepoList:
		return lipgloss.JoinHorizontal(lipgloss.Left, 
			btn("[/] new query"), btn("[↑↓] change selected repo"), btn("[s] back to query card"), btn("[↵] select repo"),
		)
	case isSearching:
		return lipgloss.JoinHorizontal(lipgloss.Left,
			btn("[/] cancel search"), btn("[esc] clear search query"), btn("[↵] search"),
		)
	default:
		return lipgloss.JoinHorizontal(lipgloss.Left,
			btn("[/] new query"), btn("[↑↓] scroll sources"), btn("[s] switch selected repo"), btn("[↵] open in browser"),
		)
	}
}
