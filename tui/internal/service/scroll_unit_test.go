//go:build unit

package service

import (
	"testing"

	"charm.land/bubbles/v2/viewport"
	"github.com/stretchr/testify/assert"
)

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
