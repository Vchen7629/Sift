package styles

import (
	"charm.land/lipgloss/v2"
)

var (
	NavBtnStyle     = lipgloss.NewStyle().PaddingLeft(2)
	NavBtnTextStyle = lipgloss.NewStyle().Foreground(Divider).Bold(true)
	ActionBarBorder = lipgloss.NewStyle().BorderBottom(true).BorderStyle(lipgloss.ThickBorder()).BorderBottomForeground(Divider)
)
