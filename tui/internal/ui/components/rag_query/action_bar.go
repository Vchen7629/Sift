package rag_query

import (
	"fmt"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ActionBarModel struct {
	ctx *context.App
	selectedRepo string
}

func NewActionBar(ctx *context.App) *ActionBarModel {
	return &ActionBarModel{ctx: ctx, selectedRepo: ""}
}

func (m ActionBarModel) Init() tea.Msg {
	return nil
}

func (m *ActionBarModel) Update() tea.Msg {
	return nil
}

func (m ActionBarModel) View() tea.View {
	navBtnStyle := lipgloss.NewStyle().PaddingLeft(2)

	navBtnTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444444")).
		Bold(true)

	selectedRepoName := fmt.Sprintf("Selected Repo: %s", m.selectedRepo)
	if m.selectedRepo == "" {
		selectedRepoName = "No Repo Selected"
	}
	
	selectedRepoText := lipgloss.NewStyle().
		PaddingRight(1).Foreground(m.ctx.SelectedTheme.AccentBright).
		Render(selectedRepoName)

	selectedRepo := navBtnStyle.Render(selectedRepoText)

	searchBtn := navBtnStyle.Render(navBtnTextStyle.Render("[/] new query"))
	scrollBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↑↓] scroll sources"))
	openBrowserBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↵] open in browser"))
	switchRepoBtn := navBtnStyle.Render(navBtnTextStyle.Render("[s] switch repo"))

	return tea.NewView(lipgloss.NewStyle().
		BorderBottom(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderBottomForeground(styles.Divider).
		Width(m.ctx.WindowWidth - 2).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, selectedRepo, searchBtn, scrollBtn, switchRepoBtn, openBrowserBtn)))
}