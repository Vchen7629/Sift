package footer

import (
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type BaseModel struct {
	height, width int
	Ctx 		  *context.App
	NavButtons    *NavButtonsModel
	ThemeSelector *ThemeSelectorModel
}

func (m *BaseModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m BaseModel) Init() tea.Cmd {
	return nil
}

func (m BaseModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "1":
			m.Ctx.ThemeSelectorOpen = false
			m.Ctx.CurrentPage = context.UserReposPage
			return nil
		
		case "2": 
			m.Ctx.ThemeSelectorOpen = false
			m.Ctx.CurrentPage = context.QueryPage
			return nil

		case "3":
			m.Ctx.ThemeSelectorOpen = true
			return nil
		}
	}

	return nil
}

func (m BaseModel) View() tea.View {
	appName := lipgloss.NewStyle().
		PaddingLeft(2).PaddingRight(1).
		Background(styles.Footer).Foreground(styles.Cold.AccentMid).      
		Render("Sift")    

	background := lipgloss.NewStyle().Background(styles.Footer).Width(m.width)

	var content string
	if m.Ctx.ThemeSelectorOpen {
		content = background.Render(lipgloss.JoinHorizontal(
			lipgloss.Left, appName, m.NavButtons.View(), m.ThemeSelector.View().Content,
		))
	} else {
		content = background.Render(lipgloss.JoinHorizontal(lipgloss.Left, appName, m.NavButtons.View()))
	}

	return tea.NewView(content)
}