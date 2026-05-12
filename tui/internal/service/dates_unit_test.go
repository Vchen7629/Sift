//go:build unit

package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func makeTs(d time.Duration) string {
	return time.Now().Add(-d).Format(time.RFC3339)
}

func TestFormatRelativeDate(t *testing.T) {
	tc := []struct{
		name     string
		input    string
		expected string
	}{
		// edge cases
		{"empty string return unknown", "", "unknown"},
		{"bad format returns unknown", "not-a-date", "unknown"},
		// minute cases
		{"1 minute ago", makeTs(61*time.Second), "1 minute ago"},
		{"59 minutes ago", makeTs(59*time.Minute), "59 minutes ago"},
		// hours cases
		{"1 hour ago", makeTs(61*time.Minute), "1 hour ago"},
		{"23 hours ago", makeTs(23*time.Hour), "23 hours ago"},
		// days cases
		{"1 day ago", makeTs(25*time.Hour), "1 day ago"},
		{"6 days ago", makeTs(6*24*time.Hour), "6 days ago"},
		// weeks cases
		{"1 week ago", makeTs(8*24*time.Hour), "1 week ago"},
		{"3 weeks ago", makeTs(22*24*time.Hour), "3 weeks ago"},
		// months cases
		{"1 month ago", makeTs(31*24*time.Hour), "1 month ago"},
		{"11 months ago", makeTs(11*30*24*time.Hour), "11 months ago"},
		// years cases
		{"1 year ago", makeTs(366*24*time.Hour), "1 year ago"},
		{"2 years ago", makeTs(2*366*24*time.Hour), "2 years ago"},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			fmtDate := FormatRelativeDate(tt.input)

			assert.Equal(t, tt.expected, fmtDate)
		})
	}
}