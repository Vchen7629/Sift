package common

import (
	"strings"
	"tui/internal/ui/styles"

	"charm.land/lipgloss/v2"
)

// divider in between the main panel and sidebar
func VerticalDivider(height int) string {
	line := strings.Repeat("│\n", height-1) + "│"
	return lipgloss.NewStyle().Foreground(styles.Divider).Render(line)
}

func SpaceBetween(totalWidth, leftWidth, rightWidth, padding int) string {
	w := totalWidth - leftWidth - rightWidth - padding
	if w < 0 {
		w = 0
	}

	return lipgloss.NewStyle().Width(w).Render("")
}