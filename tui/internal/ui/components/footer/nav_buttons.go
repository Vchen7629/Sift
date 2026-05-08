package footer

import (
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type NavButtonsModel struct {
	ctx *context.App
}

func NewNavButtons(ctx *context.App) *NavButtonsModel {
	return &NavButtonsModel{ctx: ctx}
}

func (m NavButtonsModel) Init() tea.Cmd {
	return nil
}

func (m NavButtonsModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "1":
			m.ctx.CurrentPage = context.UserReposPage
			m.ctx.ThemeSelectorOpen = false
			return nil
		
		case "2":
			m.ctx.CurrentPage = context.QueryPage
			m.ctx.ThemeSelectorOpen = false
			return nil
		}
	}

	return nil
}

func (m NavButtonsModel) View() string {
	navBtnStyle := lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).Background(styles.Footer)
	selectedStyle := lipgloss.NewStyle().Foreground(m.ctx.SelectedTheme.AccentBright).Bold(true)

	buttons := []struct {
		label string
		page  context.Page
	}{
		{"[1] repo index status", context.UserReposPage},
		{"[2] search issue", context.QueryPage},
	}

	rendered := make([]string, len(buttons)+1)
	for i, btn := range buttons {
		label := btn.label
		if m.ctx.CurrentPage == btn.page {
			label = selectedStyle.Render(label)
		}
		rendered[i] = navBtnStyle.Render(label)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, rendered...)
}