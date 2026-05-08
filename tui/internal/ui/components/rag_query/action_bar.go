package rag_query

import (
	"fmt"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ActionBarModel struct {
	ctx 		 *context.App
}

type ToggleFocusMsg struct{}

func NewActionBar(ctx *context.App) *ActionBarModel {
	return &ActionBarModel{ctx: ctx}
}

func (m ActionBarModel) Init() tea.Msg {
	return nil
}

func (m *ActionBarModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "s":
			return func () tea.Msg { return ToggleFocusMsg{} }
		}
	}
	return nil
}

func (m ActionBarModel) View(isRepoListFocused bool, selectedRepo string) tea.View {
	selectedRepoName := fmt.Sprintf("Selected Repo: %s", selectedRepo)
	if selectedRepo == "" {
		selectedRepoName = "No Repo Selected"
	}
	
	selectedRepoText := lipgloss.NewStyle().
		PaddingRight(1).Foreground(m.ctx.SelectedTheme.AccentBright).PaddingLeft(2).
		Render(selectedRepoName)

	return tea.NewView(lipgloss.NewStyle().
		BorderBottom(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderBottomForeground(styles.Divider).
		Width(m.ctx.WindowWidth - 2).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, selectedRepoText, m.actionBarBtns(isRepoListFocused))))
}

func (m ActionBarModel) actionBarBtns(focusRepoList bool) string {
	navBtnStyle := lipgloss.NewStyle().PaddingLeft(2)
	navBtnTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444444")).
		Bold(true)
	

	if focusRepoList {
		scrollText := navBtnStyle.Render(navBtnTextStyle.Render("[↑↓] change selected repo"))
		switchFocusBtn := navBtnStyle.Render(navBtnTextStyle.Render("[s] back to search"))
		selectRepoBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↵] select repo"))

		return lipgloss.JoinHorizontal(lipgloss.Left, scrollText, switchFocusBtn, selectRepoBtn)
	}

	searchBtn := navBtnStyle.Render(navBtnTextStyle.Render("[/] new query"))
	scrollSourceBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↑↓] scroll sources"))
	switchFocusBtn := navBtnStyle.Render(navBtnTextStyle.Render("[s] switch selected repo"))
	openBrowserBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↵] open in browser"))

	return lipgloss.JoinHorizontal(lipgloss.Left, searchBtn, scrollSourceBtn, switchFocusBtn, openBrowserBtn)
}