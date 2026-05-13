package user_repo

import (
	"strings"
	"tui/internal/api"
	"tui/internal/types"
	"tui/internal/ui/common"
	"tui/internal/ui/context"

	tea "charm.land/bubbletea/v2"
)

type SearchBarModel struct {
	*common.SearchBar
	OriginalGHRepoList      []api.RepoApiRes
	OriginalIndexedRepoList map[string]*types.IndexedRepo
}

func NewSearchBar(ctx *context.App, placeholderText string) *SearchBarModel {
	return &SearchBarModel{
		SearchBar: common.NewSearchBar(ctx, placeholderText),
	}
}

type searchQueryMsg struct {
	filteredGHRepos      []api.RepoApiRes
	filteredIndexedRepos map[string]*types.IndexedRepo
}

func (m *SearchBarModel) Update(msg tea.Msg, isSidebarFocused bool) tea.Cmd {
	if isSidebarFocused {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "/":
			return m.ToggleFocus()
		case "esc":
			if m.IsFocused {
				m.TextInput.Reset()
			}
			return func() tea.Msg {
				return searchQueryMsg{
					filteredGHRepos: m.OriginalGHRepoList, filteredIndexedRepos: m.OriginalIndexedRepoList,
				}
			}
		// default case to detect when users enter a new key in searchbar
		default:
			if !m.IsFocused {
				return nil
			}

			prevValue := m.TextInput.Value()
			cmd := m.UpdateInput(msg)
			if m.TextInput.Value() != prevValue {
				filteredGHRepos, filteredIndexedRepos := m.filterNameMatch(m.TextInput.Value())

				return tea.Batch(cmd, func() tea.Msg {
					return searchQueryMsg{
						filteredGHRepos: filteredGHRepos, filteredIndexedRepos: filteredIndexedRepos,
					}
				})
			}

			return cmd
		}
	case tea.WindowSizeMsg:
		m.TextInput.SetWidth(m.Ctx.MainWidth - 10)
	}

	return m.UpdateInput(msg)
}

// this filters the github repo list and indexedrepo map for repos whose
// names contain the searchQuery
func (m *SearchBarModel) filterNameMatch(
	searchQuery string,
) ([]api.RepoApiRes, map[string]*types.IndexedRepo) {
	query := strings.ToLower(searchQuery)
	filteredMap := map[string]*types.IndexedRepo{}

	var exact, startsWith, contains []api.RepoApiRes

	for _, repo := range m.OriginalGHRepoList {
		name := strings.ToLower(repo.Name)
		switch {
		case name == query:
			exact = append(exact, repo)
		case strings.HasPrefix(name, query):
			startsWith = append(startsWith, repo)
		case strings.Contains(name, query):
			contains = append(contains, repo)
		}
	}

	filteredRepos := append(append(exact, startsWith...), contains...)

	for _, repo := range filteredRepos {
		if indexed, ok := m.OriginalIndexedRepoList[repo.Name]; ok {
			filteredMap[repo.Name] = indexed
		}
	}

	return filteredRepos, filteredMap
}
