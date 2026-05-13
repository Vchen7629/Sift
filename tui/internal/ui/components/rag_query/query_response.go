package rag_query

import (
	"fmt"
	"strconv"
	"tui/internal/api"
	"tui/internal/service"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type RagQueryResponseModel struct {
	ctx                *context.App
	queryRes           api.SearchRes
	focusedIdx         int
	viewport           viewport.Model
	loadingSearchQuery bool
	spinner            spinner.Model
}

func NewRagQueryResponse(ctx *context.App) *RagQueryResponseModel {
	s := spinner.New()
	s.Spinner = spinner.Points

	return &RagQueryResponseModel{
		ctx:        ctx,
		queryRes:   api.SearchRes{},
		focusedIdx: 0,
		spinner:    s,
	}
}

func (m *RagQueryResponseModel) Init() tea.Cmd {
	return nil
}

func (m *RagQueryResponseModel) Update(msg tea.Msg, isSidebarFocused bool) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if isSidebarFocused {
			break
		}

		switch msg.String() {
		case "up":
			service.NavigateUp(&m.focusedIdx, &m.viewport, 1)
		case "down":
			service.NavigateDown(&m.focusedIdx, len(m.queryRes.IssueSources), &m.viewport, 1)
		}

	case tea.WindowSizeMsg:
		m.viewport.SetWidth(max(0, m.ctx.MainWidth-2))
		m.viewport.SetHeight(4)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)

		return cmd
	case searchQueryLoadingMsg:
		m.loadingSearchQuery = true

		return m.spinner.Tick
	case NewSearchQueryMsg:
		m.loadingSearchQuery = false
		m.queryRes = msg.Res

		return nil
	case NewSearchQueryErr:
		m.loadingSearchQuery = false
		m.queryRes = api.SearchRes{
			RepoName:     msg.RepoName,
			NumSources:   0,
			IssueSources: nil,
			Summary:      msg.Err,
		}

		return nil
	}

	return nil
}

func (m *RagQueryResponseModel) View() tea.View {
	mainWidth := m.ctx.MainWidth - 2

	contentPos := lipgloss.NewStyle().Width(mainWidth).MarginLeft(2).MarginTop(1).MarginBottom(1)
	outerBorder := lipgloss.NewStyle().Width(mainWidth-2).Padding(0, 1).
		Border(lipgloss.DoubleBorder()).BorderForeground(m.ctx.SelectedTheme.BorderFocused)

	switch {
	case m.loadingSearchQuery:
		m.spinner.Style = lipgloss.NewStyle().Foreground(m.ctx.SelectedTheme.AccentBright)

		return tea.NewView(placeholderComponent(
			mainWidth, "Searching relevant issues "+m.spinner.View(), outerBorder, contentPos,
		))

	case len(m.queryRes.IssueSources) == 0 && m.queryRes.Summary != "No relevant issues were found for your query.":
		return tea.NewView(placeholderComponent(
			mainWidth, "Enter a query to search across your GitHub issues", outerBorder, contentPos,
		))		
	case m.queryRes.Summary == "No relevant issues were found for your query.":
		innerWidth := mainWidth - 4
		centeredText := lipgloss.Place(innerWidth, 6, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(m.queryRes.Summary),
		)
		content := lipgloss.JoinVertical(lipgloss.Top, m.header(), centeredText)

		return tea.NewView(contentPos.Render(outerBorder.Render(content)))
	default:
		answerText := lipgloss.NewStyle().
			MarginBottom(2).Width(mainWidth - 6).MaxHeight(6).
			BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(styles.TextDim).
			Render(m.queryRes.Summary)

		content := lipgloss.JoinVertical(lipgloss.Top, m.header(), answerText, m.source().Content)

		return tea.NewView(contentPos.Render(outerBorder.Render(content)))
	}
}

func (m *RagQueryResponseModel) header() string {
	cardTitle := lipgloss.NewStyle().Bold(true).MarginRight(3).Render("answer")
	name := lipgloss.NewStyle().Foreground(m.ctx.SelectedTheme.AccentBright).MarginRight(1).Render(m.queryRes.RepoName)
	numSources := lipgloss.NewStyle().Foreground(styles.TextMuted).Render(fmt.Sprintf("· %d sources", m.queryRes.NumSources))

	padding := lipgloss.NewStyle().PaddingBottom(1)

	return padding.Render(lipgloss.JoinHorizontal(lipgloss.Left, cardTitle, name, numSources))
}

func (m *RagQueryResponseModel) source() tea.View {
	var sources []string

	for i, dependency := range m.queryRes.IssueSources {
		sourceCard := m.sourceCard(i, dependency)

		sources = append(sources, sourceCard)
	}

	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, sources...))
	return tea.NewView(m.viewport.View())
}

func (m *RagQueryResponseModel) sourceCard(idx int, dependency api.IssueSource) string {
	textColor := styles.FocusColor(m.ctx.SelectedTheme, idx, m.focusedIdx)

	dependencyText := lipgloss.NewStyle().Foreground(textColor)

	id := lipgloss.NewStyle().PaddingRight(2).Render(strconv.Itoa(idx + 1))
	name := lipgloss.NewStyle().Render(dependency.Url)

	return dependencyText.Render(lipgloss.JoinHorizontal(lipgloss.Left, id, name))
}

// placeholder component when no search query/loading
func placeholderComponent(mainWidth int, text string, outerBorder, contentPos lipgloss.Style) string {
	innerWidth := mainWidth - 4
	placeholder := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(text)

	centered := lipgloss.Place(innerWidth, 10, lipgloss.Center, lipgloss.Center, placeholder)

	box := outerBorder.Height(10).Render(centered)
	return contentPos.Render(box)
}
