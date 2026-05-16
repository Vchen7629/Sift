package rag_query

import (
	"tui/internal/api"
	"tui/internal/ui/common"
	"tui/internal/ui/context"

	tea "charm.land/bubbletea/v2"
)

type SearchBarModel struct {
	*common.SearchBar
}

func NewSearchBar(ctx *context.App, placeholderText string) *SearchBarModel {
	return &SearchBarModel{SearchBar: common.NewSearchBar(ctx, placeholderText)}
}

type searchQueryLoadingMsg struct{ name string }

func (m *SearchBarModel) Update(msg tea.Msg, selectedRepo string) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "/":
			return m.ToggleFocus()
		case "esc":
			if m.IsFocused {
				m.TextInput.Reset()
			}
			return nil

		case "enter":
			if !m.IsFocused || selectedRepo == "" {
				return nil
			}
			cmd := m.newSearchQuery(selectedRepo)

			return tea.Batch(cmd, func() tea.Msg { return searchQueryLoadingMsg{name: selectedRepo} })
		}
	case tea.WindowSizeMsg:
		m.TextInput.SetWidth(m.Ctx.MainWidth - 10)
	}

	return m.UpdateInput(msg)
}

type NewSearchQueryMsg struct{ Res api.SearchRes }
type NewSearchQueryErr struct{ RepoName, Err string }

func (m *SearchBarModel) newSearchQuery(repoName string) tea.Cmd {
	return func() tea.Msg {
		searchRes, err := api.Search(m.Ctx.SessionToken, m.Ctx.Username, repoName, m.TextInput.Value())
		if err != nil {
			return NewSearchQueryErr{RepoName: repoName, Err: err.Error()}
		}

		return NewSearchQueryMsg{Res: searchRes}
	}
}
