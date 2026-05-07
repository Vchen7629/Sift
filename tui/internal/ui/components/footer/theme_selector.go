package footer

import (
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ThemeSelectorModel struct {
	ctx *context.App
}

func NewThemeSelector(ctx *context.App) *ThemeSelectorModel {
	return &ThemeSelectorModel{ctx: ctx}
}

func (m *ThemeSelectorModel) Init() tea.Cmd {
	return nil
}

func (m *ThemeSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *ThemeSelectorModel) View() tea.View {
	label := lipgloss.NewStyle().PaddingRight(2).Background(styles.Footer)

	Mono := label.PaddingLeft(1).Foreground(styles.Mono.AccentBright).Render("Mono")
	Warm := label.Foreground(styles.Warm.AccentBright).Render("Warm")
	Cold := label.Foreground(styles.Cold.AccentBright).Render("Cold")

	return tea.NewView(lipgloss.JoinHorizontal(lipgloss.Left, Mono, Warm, Cold))
}
