//go:build unit

package rag_query

import (
	"testing"
	"tui/internal/api"
	"tui/internal/ui/context"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
)

func queryResponseCtx() *context.App {
	return &context.App{}
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
			m.Update(tea.KeyPressMsg{Code: tc.key}, tc.isSidebarFocused, false)

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
			m.Update(tea.WindowSizeMsg{Width: 33, Height: 22}, false, false)

			assert.Equal(t, tc.expectedWidth, m.viewport.Width())
			assert.Equal(t, tc.expectedHeight, m.viewport.Height())
		})
	}
}

func TestQueryResponse_NewSearchQueryMsg(t *testing.T) {
	tt := []struct {
		name          string
		isReauthed    bool
		newSessToken  string
		expectedToken string
	}{
		{"sets query res properly", false, "new-token", "old-token"},
		{"updates the token on reauth", true, "new-token", "new-token"},
		{"does not overwrite token when not reauthed", false, "old-token-2", "old-token"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			res := api.SearchRes{
				RepoName:     "my-repo",
				IssueSources: []api.IssueSource{{Title: "bug"}},
			}

			ctx := queryResponseCtx()
			m := NewRagQueryResponse(ctx)
			m.ctx.SessionToken = "old-token"
			msg := NewSearchQueryMsg{Res: res, NewSessionToken: tc.newSessToken, isReauthed: tc.isReauthed}

			m.Update(msg, false, false)

			assert.Equal(t, res, m.queryRes)
			assert.Equal(t, tc.expectedToken, ctx.SessionToken)
		})
	}
}

func TestQueryResponse_UpdateMessages(t *testing.T) {
	tt := []struct {
		name         string
		msg          tea.Msg
		setupLoading bool
		wantLoading  bool
		wantCmd      bool
		wantQueryRes *api.SearchRes
	}{
		{
			name:        "searchQueryLoadingMsg sets loading and starts spinner",
			msg:         searchQueryLoadingMsg{},
			wantLoading: true,
			wantCmd:     true,
		},
		{
			name:    "spinner.TickMsg returns tick cmd",
			msg:     spinner.TickMsg{},
			wantCmd: true,
		},
		{
			name:         "NewSearchQueryMsg clears loading and sets result",
			msg:          NewSearchQueryMsg{Res: api.SearchRes{RepoName: "repo"}},
			setupLoading: true,
			wantLoading:  false,
			wantQueryRes: &api.SearchRes{RepoName: "repo"},
		},
		{
			name:         "NewSearchQueryErr clears loading and sets error",
			msg:          NewSearchQueryErr{RepoName: "repo", Err: "failed"},
			setupLoading: true,
			wantLoading:  false,
			wantQueryRes: &api.SearchRes{RepoName: "repo", Summary: "failed"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := NewRagQueryResponse(queryResponseCtx())
			m.loadingSearchQuery = tc.setupLoading
			cmd := m.Update(tc.msg, false, false)

			assert.Equal(t, tc.wantLoading, m.loadingSearchQuery)
			assert.Equal(t, tc.wantCmd, cmd != nil)
			if tc.wantQueryRes != nil {
				assert.Equal(t, *tc.wantQueryRes, m.queryRes)
			}
		})
	}
}
