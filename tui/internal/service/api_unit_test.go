//go:build unit

package service

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshallRequestBody(t *testing.T) {
	tc := []struct{
		name 	 string
		username string
		repoName string
	}{
		{"valid inputs", "vchen7629", "Sift"},
		{"empty username", "", "Sift"},
		{"empty repoName", "vchen7629", ""},
		{"both empty", "", ""},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			data, err := MarshalRequestBody(tt.username, tt.repoName)

			require.NoError(t, err)

			var got RequestBody
			require.NoError(t, json.Unmarshal(data, &got))

			assert.Equal(t, tt.username, got.UserId)
			assert.Equal(t, tt.repoName, got.RepoName)
		})
	}
}