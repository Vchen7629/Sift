//go:build unit

package service

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshallRequestBody(t *testing.T) {
	tt := []struct {
		name     string
		username string
		repoName string
	}{
		{"valid inputs", "vchen7629", "Sift"},
		{"empty username", "", "Sift"},
		{"empty repoName", "vchen7629", ""},
		{"both empty", "", ""},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			data, err := MarshalRequestBody(tc.username, tc.repoName)

			require.NoError(t, err)

			var got RequestBody
			require.NoError(t, json.Unmarshal(data, &got))

			assert.Equal(t, tc.username, got.UserId)
			assert.Equal(t, tc.repoName, got.RepoName)
		})
	}
}
