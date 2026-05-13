package styles

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

var (
	NavBtnStyle     = lipgloss.NewStyle().PaddingLeft(2)
	NavBtnTextStyle = lipgloss.NewStyle().Foreground(Divider).Bold(true)
	ActionBarBorder = lipgloss.NewStyle().BorderBottom(true).BorderStyle(lipgloss.ThickBorder()).BorderBottomForeground(Divider)
)

func FocusColor(theme Theme, idx, focusedIdx int) color.Color {
	if idx == focusedIdx {
		return theme.AccentBright
	}
	return theme.AccentMid
}