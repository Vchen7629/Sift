package components

import (
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type FooterModel struct {
	height, width int
	Ctx *context.App
}

func (m *FooterModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m FooterModel) Init() tea.Cmd {
	return nil
}

func (m FooterModel) Update(msg tea.Msg) tea.Cmd {
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

func (m FooterModel) View() tea.View {
	appName := lipgloss.NewStyle().
		PaddingLeft(2).PaddingRight(1).
		Background(styles.Footer).Foreground(styles.Cold.AccentMid).      
		Render("Sift")    

	background := lipgloss.NewStyle().Background(styles.Footer).Width(m.width)

	content := background.Render(lipgloss.JoinHorizontal(lipgloss.Left, appName, m.navButtons()))

	return tea.NewView(content)
}

func (m FooterModel) navButtons() string {
	navBtnStyle := lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).Background(styles.Footer)
	selectedStyle := lipgloss.NewStyle().Foreground(styles.Warm.AccentBright).Bold(true)

	buttons := []struct {
		label string
		page  context.Page
	}{
		{"[1] repos", context.UserReposPage},
		{"[2] search issue", context.QueryPage},
		{"[3] theme", context.ThemePage},
	}

	rendered := make([]string, len(buttons))
	for i, btn := range buttons {
		label := btn.label
		if m.Ctx.CurrentPage == btn.page {
			label = selectedStyle.Render(label)
		}
		rendered[i] = navBtnStyle.Render(label)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, rendered...)
}