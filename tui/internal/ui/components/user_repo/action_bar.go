package user_repo

import (
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ActionBarModel struct {
	ctx              *context.App
	IndexRepoApiDown bool
}

func NewActionBar(ctx *context.App) *ActionBarModel {
	return &ActionBarModel{ctx: ctx, IndexRepoApiDown: false}
}

func (m *ActionBarModel) Init() tea.Cmd {
	return nil
}

type ToggleFocusMsg struct{}
type IndexRepoRequestMsg struct{}

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

func (m *ActionBarModel) View(isSidebarFocused, isIndexed bool) tea.View {
	widthStyle := styles.ActionBarBorder.Width(m.ctx.WindowWidth - 2)

	if m.IndexRepoApiDown {
		statusText := lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Red).Render("Index Repo api is unavailable")

		return tea.NewView(widthStyle.Render(
			lipgloss.JoinHorizontal(lipgloss.Left, statusText, m.actionBarBtns(isSidebarFocused, isIndexed))),
		)
	}
	return tea.NewView(widthStyle.Render(m.actionBarBtns(isSidebarFocused, isIndexed)))
}

func (m *ActionBarModel) actionBarBtns(isSidebarFocused, isIndexed bool) string {
	btn := func(text string) string {
		return styles.NavBtnStyle.Render(styles.NavBtnTextStyle.Render(text))
	}

	switch {
	case isSidebarFocused:
		return lipgloss.JoinHorizontal(lipgloss.Left,
			btn("[↑↓] scroll dependencies"), btn("[s] focus repo list"), btn("[↵] open in browser"),
		)
	case !isIndexed:
		return lipgloss.JoinHorizontal(lipgloss.Left,
			btn("[↑↓] scroll repos"), btn("[/] search repos"), btn("[esc] clear search"), btn("[s] focus sidebar"), btn("[r] index"),
		)
	default:
		return lipgloss.JoinHorizontal(lipgloss.Left,
			btn("[↑↓] scroll repos"), btn("[/] search repos"), btn("[esc] clear search"), btn("[s] focus sidebar"),
		)
	}
}
