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
	navBtnStyle := styles.NavBtnStyle
	navBtnTextStyle := styles.NavBtnTextStyle

	if focusRepoList {
		scrollText := navBtnStyle.Render(navBtnTextStyle.Render("[↑↓] change selected repo"))
		switchFocusBtn := navBtnStyle.Render(navBtnTextStyle.Render("[s] back to search"))
		selectRepoBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↵] select repo"))

		return lipgloss.JoinHorizontal(lipgloss.Left, scrollText, switchFocusBtn, selectRepoBtn)
	}

	if isSearching {
		focusBtn := navBtnStyle.Render(navBtnTextStyle.Render("[/] cancel search"))
		clearSearchBtn := navBtnStyle.Render(navBtnTextStyle.Render("[esc] clear search query"))
		searchBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↵] search"))

		return lipgloss.JoinHorizontal(lipgloss.Left, focusBtn, clearSearchBtn, searchBtn)
	}

	searchBtn := navBtnStyle.Render(navBtnTextStyle.Render("[/] new query"))
	scrollSourceBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↑↓] scroll sources"))
	switchFocusBtn := navBtnStyle.Render(navBtnTextStyle.Render("[s] switch selected repo"))
	openBrowserBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↵] open in browser"))

	return lipgloss.JoinHorizontal(lipgloss.Left, searchBtn, scrollSourceBtn, switchFocusBtn, openBrowserBtn)
}
