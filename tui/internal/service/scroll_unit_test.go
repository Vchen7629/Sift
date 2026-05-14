//go:build unit

package service

import (
	"strings"
	"testing"

	"charm.land/bubbles/v2/viewport"
	"github.com/stretchr/testify/assert"
)

func makeViewport(height, yOffset int) *viewport.Model {
	vp := viewport.New(viewport.WithHeight(height))
	vp.SetContent(strings.Repeat("line\n", height+yOffset+20))
	vp.SetYOffset(yOffset)

	return &vp
}

func TestScrollToFocused(t *testing.T) {
	tt := []struct{
		name		  string
		height		  int
		initialOffset int
		focusedIndex  int
		cardHeight    int
		wantOffset    int
	}{
		{
			name: "item fully visible, no scroll",
			height: 5, initialOffset: 0,
			focusedIndex: 1, cardHeight: 2,
			wantOffset: 0,
		},
		{
			name: "item below visible area, scrolls down just enough",
			height: 3, initialOffset: 0,
			focusedIndex: 2, cardHeight: 2,
			wantOffset: 3,
		},
		{
			name: "item above visible area, scrolls up just enough",
			height: 3, initialOffset: 6,
			focusedIndex: 1, cardHeight: 2,
			wantOffset: 2,
		},
		{
			name: "item exactly at bottom boundary, no scroll",
			height: 4, initialOffset: 0,
			focusedIndex: 1, cardHeight: 2,
			wantOffset: 0,
		},
		{
			name: "item exactly at top boundary, no scroll",
			height: 4, initialOffset: 2,
			focusedIndex: 1, cardHeight: 2,
			wantOffset: 2,
		},
		{
			name: "scroll back to top when first item focused",
			height: 3, initialOffset: 5,
			focusedIndex: 0, cardHeight: 3,
			wantOffset: 0,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			vp := makeViewport(tc.height, tc.initialOffset)
			ScrollToFocused(vp, tc.focusedIndex, tc.cardHeight)

			assert.Equal(t, tc.wantOffset, vp.YOffset())
		})
	}
}

func TestNavigateUp(t *testing.T) {
	tt := []struct {
		name    string
		start   int
		wantIdx int
	}{
		{"moves up normally", 2, 1},
		{"moves up from middle", 1, 0},
		{"does not go negative", 0, 0},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			vp := viewport.New()
			idx := tc.start
			NavigateUp(&idx, &vp, 1)
			assert.Equal(t, tc.wantIdx, idx, "should be this index after navigating")
		})
	}
}

func TestNavigateDown(t *testing.T) {
	tt := []struct {
		name    string
		start   int
		listLen int
		wantIdx int
	}{
		{"moves down normally", 0, 3, 1},
		{"moves down from middle", 1, 3, 2},
		{"does not move past end", 2, 3, 2},
		{"does not move on single item list", 0, 1, 0},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			vp := viewport.New()
			idx := tc.start
			NavigateDown(&idx, tc.listLen, &vp, 1)
			assert.Equal(t, tc.wantIdx, idx, "should be this index after navigating")
		})
	}
}