package user_repo

import (
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ActionBarModel struct {
	ctx *context.App
}

func NewActionBar(ctx *context.App) *ActionBarModel {
	return &ActionBarModel{ctx: ctx}
}

func (m ActionBarModel) Init() tea.Msg {
	return nil
}

type ToggleFocusMsg struct {}
type IndexRepoRequestMsg struct {}

func (m *ActionBarModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "s":
			return func() tea.Msg { return ToggleFocusMsg{} }
		case "r": 
			return func() tea.Msg { return IndexRepoRequestMsg{} }
		}
	}

	return nil
}

func (m ActionBarModel) View(isSidebarFocused bool) tea.View {
	return tea.NewView(lipgloss.NewStyle().
		BorderBottom(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderBottomForeground(styles.Divider).
		Width(m.ctx.WindowWidth - 2).
		Render(m.actionBarBtns(isSidebarFocused)))
}

func (m ActionBarModel) actionBarBtns(isSidebarFocused bool) string {
	navBtnStyle := lipgloss.NewStyle().PaddingLeft(2)

	navBtnTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444444")).
		Bold(true)

	if !isSidebarFocused {
		navBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↑↓] scroll repos"))
		searchBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↵] search repos"))
		clearSearchBtn := navBtnStyle.Render(navBtnTextStyle.Render("[esc] clear search"))
		swapFocusBtn := navBtnStyle.Render(navBtnTextStyle.Render("[s] focus sidebar"))
		indexBtn := navBtnStyle.Render(navBtnTextStyle.Render("[r] index"))

		return lipgloss.JoinHorizontal(lipgloss.Left, navBtn, searchBtn, clearSearchBtn, swapFocusBtn, indexBtn)
	}

	navBtn := navBtnStyle.Render(navBtnTextStyle.Render("[↑↓] scroll dependencies"))
	swapFocusBtn := navBtnStyle.Render(navBtnTextStyle.Render("[s] focus repo list"))
	indexBtn := navBtnStyle.Render(navBtnTextStyle.Render("[r] reindex"))

	return lipgloss.JoinHorizontal(lipgloss.Left, navBtn, swapFocusBtn, indexBtn)
}