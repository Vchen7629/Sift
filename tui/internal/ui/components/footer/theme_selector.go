package footer

import (
	"image/color"
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

var nextTheme = map[styles.Theme]styles.Theme{
	styles.Mono: styles.Warm,
	styles.Warm: styles.Cold,
	styles.Cold: styles.Mono,
}

var prevTheme = map[styles.Theme]styles.Theme{
	styles.Mono: styles.Cold,
	styles.Warm: styles.Mono,
	styles.Cold: styles.Warm,
}

func (m *ThemeSelectorModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if !m.ctx.ThemeSelectorOpen {
			return nil
		}
		switch msg.String() {
		case "right":
			m.ctx.SelectedTheme = nextTheme[m.ctx.SelectedTheme]
		case "left":
			m.ctx.SelectedTheme = prevTheme[m.ctx.SelectedTheme]
		}
	}

	return nil
}

func (m *ThemeSelectorModel) View() tea.View {
	label := lipgloss.NewStyle().PaddingRight(2).Foreground(styles.TextDim).Background(styles.Footer)

	themes := []struct {
		label         string
		theme         styles.Theme
		selectedColor color.Color
	}{
		{"Mono", styles.Mono, styles.Mono.AccentBright},
		{"Warm", styles.Warm, styles.Warm.AccentBright},
		{"Cold", styles.Cold, styles.Cold.AccentBright},
	}

	themeBtns := make([]string, 3)
	for i, theme := range themes {
		name := theme.label
		selectedLabel := lipgloss.NewStyle().Inherit(label).Foreground(theme.selectedColor)
		
		if m.ctx.SelectedTheme == theme.theme {
			name = selectedLabel.Render(name)
		}

		themeBtns[i] = label.Render(name)
	}

	navBtn := lipgloss.NewStyle().Background(styles.Footer).Foreground(m.ctx.SelectedTheme.AccentMid).Render("[◄ ►] navigate")

	return tea.NewView(lipgloss.JoinHorizontal(lipgloss.Left, append(themeBtns, navBtn)...))
}
