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
		case "3":
			m.Ctx.ThemeSelectorOpen = !m.Ctx.ThemeSelectorOpen
			return nil
		}
	}

	return tea.Batch(m.NavButtons.Update(msg), m.ThemeSelector.Update(msg))
}

func (m BaseModel) View() tea.View {
	appName := lipgloss.NewStyle().
		PaddingLeft(2).PaddingRight(1).
		Background(styles.Footer).Foreground(m.Ctx.SelectedTheme.AccentMid).      
		Render("Sift")    

	background := lipgloss.NewStyle().Background(styles.Footer).Width(m.width)

	var content string
	if m.Ctx.ThemeSelectorOpen {
		content = background.Render(lipgloss.JoinHorizontal(
			lipgloss.Left, appName, m.NavButtons.View(), m.themeBtn(), m.ThemeSelector.View().Content,
		))
	} else {
		content = background.Render(lipgloss.JoinHorizontal(lipgloss.Left, appName, m.NavButtons.View(), m.themeBtn()))
	}

	return tea.NewView(content)
}

func (m BaseModel) themeBtn() string {
	themeLabel := "[3] theme"
	if m.Ctx.ThemeSelectorOpen {
		themeLabel = lipgloss.NewStyle().
			PaddingRight(2).
			Background(styles.Footer).Foreground(m.Ctx.SelectedTheme.AccentBright).Bold(true).
			Render(themeLabel)
	}

	return lipgloss.NewStyle().
		PaddingLeft(1).PaddingRight(1).
		Background(styles.Footer).
		Render(themeLabel)
}