package components

import (
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type StatusBarModel struct {
	height, width int
}

func (m *StatusBarModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m StatusBarModel) Init() tea.Cmd {
	return nil
}

func (m StatusBarModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (m StatusBarModel) View() tea.View {
	sideWidth := 12
	sidePadding := 3

	appName := lipgloss.NewStyle().
		Width(sideWidth).
		Background(styles.Surface).
		PaddingLeft(sidePadding).
		Render("Sift")

	navBtnStyle := lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(2)

	navBtnTextStyle := lipgloss.NewStyle().
		Foreground(styles.Warm.AccentBright).
		Bold(true)

	btn1 := navBtnStyle.Render(navBtnTextStyle.Render("[1] repos"))
	btn2 := navBtnStyle.Render(navBtnTextStyle.Render("[2] search"))
	btn3 := navBtnStyle.Render(navBtnTextStyle.Render("[3] theme"))

	buttons := lipgloss.JoinHorizontal(lipgloss.Left, btn1, btn2, btn3)

	middleWidth := m.width - sideWidth - sideWidth
	centered := lipgloss.NewStyle().
		Width(middleWidth).
		Align(lipgloss.Center).
		Render(buttons)
	
	authStatus := lipgloss.NewStyle().
		Width(sideWidth).
		Align(lipgloss.Right).
		PaddingRight(sidePadding).
		Render("username")

	return tea.NewView(lipgloss.JoinHorizontal(lipgloss.Left, appName, centered, authStatus))
}