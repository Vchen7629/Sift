//go:build unit

package rag_query

import (
	"testing"
	"tui/internal/types"
	"tui/internal/ui/common"
	"tui/internal/ui/context"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
)

func testCtx() *context.App {
	return &context.App{}
}

var fakeIndexedRepos = []types.IndexedRepo{{Name: "Coco"}, {Name: "Cottom"}, {Name: "Maple"}}

func TestSidebarUpdate_SidebarFocusedGuard(t *testing.T) {
	tt := []struct {
		name             string
		isSidebarFocused bool
		isSearching      bool
		key              rune
		wantIdx          int
	}{
		{"down blocked when not focused", false, false, tea.KeyDown, 0},
		{"up blocked when not focused", false, false, tea.KeyUp, 0},
		{"down blocked when searching", true, true, tea.KeyDown, 0},
		{"up blocked when searching", true, true, tea.KeyUp, 0},
		{"down works when focused and not searching", true, false, tea.KeyDown, 1},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := NewSidebar(testCtx())

			m.Update(common.FetchIndexedRepoMsg{IndexedRepos: fakeIndexedRepos}, false, tc.isSearching)
			m.Update(tea.KeyPressMsg{Code: tc.key}, tc.isSidebarFocused, tc.isSearching)

			assert.Equal(t, m.focusedIdx, tc.wantIdx)
		})
	}
}

func TestSidebarUpdate_KeyPress_Enter(t *testing.T) {
	tc := []struct {
		name    string
		repos   []types.IndexedRepo
		wantCmd bool
	}{
		{"no-op when repos empty", nil, false},
		{"returns cmd when repos present", fakeIndexedRepos, true},
	}

	for _, tc := range tc {
		t.Run(tc.name, func(t *testing.T) {
			m := NewSidebar(testCtx())
			m.Update(common.FetchIndexedRepoMsg{IndexedRepos: tc.repos}, false, false)
			cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter}, true, false)

			if tc.wantCmd {
				assert.NotNil(t, cmd)
			} else {
				assert.Nil(t, cmd)
			}
		})
	}
}

func TestSidebarUpdate_FetchIndexedRepoMsg(t *testing.T) {
	tt := []struct {
		name           string
		indexedRepos   []types.IndexedRepo
		expectedTeaMsg func() tea.Msg
	}{
		{
			name:           "returns nil if indexed repos is empty",
			indexedRepos:   []types.IndexedRepo{},
			expectedTeaMsg: nil,
		},
		{
			name: "returns first repo selected msg if indexed repo non empty",
			indexedRepos: []types.IndexedRepo{
				{TotalDependencies: 2, Name: "idk", LastIndexed: "yesterday", Dependencies: nil},
				{TotalDependencies: 4, Name: "idk3", LastIndexed: "yesterday=2", Dependencies: nil},
			},
			expectedTeaMsg: func() tea.Msg { return SelectRepoMsg{RepoName: "idk"} },
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := testCtx()
			m := NewSidebar(ctx)
			cmd := m.Update(common.FetchIndexedRepoMsg{
				IndexedRepos: tc.indexedRepos, NewSessionToken: "new-session", IsReauthed: false,
			}, true, false)

			if tc.expectedTeaMsg == nil {
				assert.Nil(t, cmd)
			} else {
				assert.Equal(t, tc.expectedTeaMsg(), cmd())
			}
		})
	}
}
