package rag_query

import (
	"fmt"
	"strconv"
	"tui/internal/service"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type RagQueryResponseModel struct {
	ctx           *context.App
	mainWidth 	  int
	answer    	  answerModel
	Focused 	  source
	viewport  	  viewport.Model
}

type answerModel struct {
	dependencyName, text string
	numSources 	   		   int8
	sources			       []source
}

type source struct {
	id 					 int
	link, version, label string 
}

func NewRagQueryResponse(ctx *context.App) *RagQueryResponseModel {
	return &RagQueryResponseModel{
		ctx: ctx, 
		answer: answerModel{
			dependencyName: "Sift", 
			text: "RAG pipeline for semantic search over GitHub issues. Built with Java, Spring Boot, and OpenSearch.",
			numSources: 7,
			sources: []source{
				{id: 0, link: "github.com/idk/axios", version: "v1.7.2", label: "issue"},
				{id: 1, link: "github.com/idk/lodash", version: "v4.17.21", label: "changelog"},
				{id: 2, link: "github.com/idk/moment", version: "v2.29.4", label: "issue"},
				{id: 3, link: "github.com/idk/react-query", version: "v5.28.0", label: "issue"},
				{id: 4, link: "github.com/idk/classnames", version: "v2.3.2", label: "issue"},
				{id: 5, link: "github.com/idk/prop-types", version: "v15.8.1", label: "issue"},
				{id: 6, link: "github.com/idk/redux", version: "v4.2.1", label: "changelog"},
			},
		},
		Focused: source{id: 0, link: "github.com/idk/axios", version: "v1.7.2", label: "issue"},
	}
}

func (m RagQueryResponseModel) Init() tea.Cmd {
	return nil
}

func (m *RagQueryResponseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cardHeight := lipgloss.Height(m.sourceCard(m.Focused))

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "down":
			if m.Focused.id < len(m.answer.sources) - 1 {
				m.Focused.id++
				m.Focused = m.answer.sources[m.Focused.id]
				service.ScrollToFocused(&m.viewport, m.Focused.id, cardHeight)
			}
		case "up":
			if m.Focused.id > 0 {
				m.Focused.id--
				m.Focused = m.answer.sources[m.Focused.id]
				service.ScrollToFocused(&m.viewport, m.Focused.id, cardHeight)
			}
		}
	}

	return m, nil
}

func (m *RagQueryResponseModel) View() tea.View {
	mainWidth := m.ctx.MainWidth - 2

	contentPos := lipgloss.NewStyle().Width(mainWidth).MarginLeft(2).MarginTop(1).MarginBottom(1)
	outerBorder := lipgloss.NewStyle().Width(mainWidth - 2).Padding(0, 1).
		Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#747373"))

	answerText := lipgloss.NewStyle().
		PaddingBottom(1).Width(mainWidth - 6).MaxHeight(6).
		BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(styles.TextDim).
		Render(m.answer.text)

	content := lipgloss.JoinVertical(lipgloss.Top, m.header(), answerText, m.source().Content)

	return tea.NewView(contentPos.Render(outerBorder.Render(content)))
}

func (m RagQueryResponseModel) header() string {
	cardTitle := lipgloss.NewStyle().Bold(true).MarginRight(3).Render("answer")
	name := lipgloss.NewStyle().Foreground(styles.Warm.AccentBright).MarginRight(1).Render(m.answer.dependencyName)
	numSources := lipgloss.NewStyle().Foreground(styles.TextMuted).Render(fmt.Sprintf("· %d sources", m.answer.numSources))

	padding := lipgloss.NewStyle().PaddingBottom(1)

	return padding.Render(lipgloss.JoinHorizontal(lipgloss.Left, cardTitle, name, numSources))
}

func (m *RagQueryResponseModel) source() tea.View {
	var sources []string

	for _, dependency := range m.answer.sources {
		sourceCard := m.sourceCard(dependency)

		sources = append(sources, sourceCard)
	}

	m.viewport.SetWidth(m.ctx.MainWidth - 2)
	m.viewport.SetHeight(4)
	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, sources...))
	return tea.NewView(m.viewport.View())
}

func (m *RagQueryResponseModel) sourceCard(dependency source) string {
	textColor := styles.TextDim

	if m.Focused.id == dependency.id {
		textColor = styles.Warm.AccentBright
	}

	dependencyText := lipgloss.NewStyle().Foreground(textColor)

	id := lipgloss.NewStyle().PaddingRight(2).Render(strconv.Itoa(dependency.id + 1))
	name := lipgloss.NewStyle().Width(45).Render(dependency.link)
	version := lipgloss.NewStyle().Width(10).Render(dependency.version)
	label := lipgloss.NewStyle().Render(fmt.Sprintf("· %s", dependency.label))

	return lipgloss.NewStyle().Render(dependencyText.Render(lipgloss.JoinHorizontal(
		lipgloss.Left, id, name, version, label,
	)))
}

// Todo: sidebar to show a list of repos to switch
