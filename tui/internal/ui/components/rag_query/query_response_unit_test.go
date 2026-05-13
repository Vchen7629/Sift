//go:build unit

package rag_query

import (
	"testing"
	"tui/internal/api"
	"tui/internal/ui/context"
	"tui/internal/ui/styles"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
)

func queryResponseCtx() *context.App {
	return &context.App{
		SelectedTheme: styles.Warm,
	}
}

func TestQueryResponse_SidebarFocusedGuard(t *testing.T) {
	tt := []struct {
		name             string
		isSidebarFocused bool
		key              rune
		wantIdx          int
	}{
		{"down blocked when sidebar focused", true, tea.KeyDown, 0},
		{"up blocked when sidebar focused", true, tea.KeyUp, 0},
		{"down works when sidebar not focused", false, tea.KeyDown, 1},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := NewRagQueryResponse(queryResponseCtx())
			m.queryRes = api.SearchRes{IssueSources: []api.IssueSource{{}, {}}}
			m.Update(tea.KeyPressMsg{Code: tc.key}, tc.isSidebarFocused)

			assert.Equal(t, tc.wantIdx, m.focusedIdx)
		})
	}
}

func TestQueryResponse_WindowSizeMsg(t *testing.T) {
	tt := []struct {
		name           string
		width          int
		height         int
		expectedWidth  int
		expectedHeight int
	}{
		{"resizing window properly", 100, 50, 98, 4},
		{"doesnt resize width negative", 1, 50, 0, 4},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := queryResponseCtx()
			ctx.MainWidth = tc.width

			m := NewRagQueryResponse(ctx)
			m.Update(tea.WindowSizeMsg{Width: 33, Height: 22}, false)

			assert.Equal(t, tc.expectedWidth, m.viewport.Width())
			assert.Equal(t, tc.expectedHeight, m.viewport.Height())
		})
	}
}

func TestQueryResponse_NewSearchQueryMsg(t *testing.T) {
	res := api.SearchRes{
		RepoName:     "my-repo",
		IssueSources: []api.IssueSource{{Title: "bug"}},
	}

	m := NewRagQueryResponse(queryResponseCtx())
	m.Update(NewSearchQueryMsg{Res: res}, false)

	assert.Equal(t, res, m.queryRes)
}
