package components

import (
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type StatusBarModel struct {
	height, width int
	Ctx *context.App
}

func (m *StatusBarModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m StatusBarModel) Init() tea.Cmd {
	return nil
}

func (m StatusBarModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "1":
			m.Ctx.CurrentPage = context.UserReposPage
			return nil
		
		case "2": 
			m.Ctx.CurrentPage = context.QueryPage
			return nil
		}
	}

	return nil
}

func (m StatusBarModel) View() tea.View {
	appName := lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(1).
		BorderRight(true).                                                                                                                                                                          
		BorderStyle(lipgloss.MarkdownBorder()).
		Background(styles.Footer).
		Foreground(styles.Cold.AccentMid).      
		Render("Sift")    

	background := lipgloss.NewStyle().
		Background(styles.Footer).
		Width(m.width)

	content := background.Render(lipgloss.JoinHorizontal(lipgloss.Left, appName, m.navButtons()))

	return tea.NewView(content)
}

func (m StatusBarModel) navButtons() string {
	navBtnStyle := lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1).
		BorderRight(true).
		Background(styles.Footer).
		BorderStyle(lipgloss.MarkdownBorder())

	navBtnTextStyle := lipgloss.NewStyle().
		Foreground(styles.Warm.AccentBright).
		Background(styles.Footer). 
		Bold(true)

	btn1 := navBtnStyle.Render(navBtnTextStyle.Render("[1] repos"))
	btn2 := navBtnStyle.Render(navBtnTextStyle.Render("[2] search"))
	btn3 := navBtnStyle.Render(navBtnTextStyle.Render("[3] theme"))

	return lipgloss.JoinHorizontal(lipgloss.Left, btn1, btn2, btn3)
}