//go:build unit

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSession_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "Bearer mytoken", r.Header.Get("Authorization"))

		http.SetCookie(w, &http.Cookie{Name: "JSESSIONID", Value: "abc123"})
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	orig := authBaseUrl
	authBaseUrl = ts.URL + "/auth"

	t.Cleanup(func() { authBaseUrl = orig })

	token, err := NewSession("mytoken")
	require.NoError(t, err)
	assert.Equal(t, "abc123", token)
}

func TestNewSession_InvalidRequest(t *testing.T) {
	tt := []struct {
		name        string
		statusCode  int
		hasCookie   bool
		expectedErr string
	}{
		{"Non 200 response (401) returns error", 401, true, "unexpected error sending req: 401"},
		{"200 response but no cookie returns error", 200, false, "no session token in response"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "Bearer mytoken", r.Header.Get("Authorization"))

				if tc.hasCookie {
					http.SetCookie(w, &http.Cookie{Name: "JSESSIONID", Value: "abc123"})
				}
				w.WriteHeader(tc.statusCode)
			}))
			defer ts.Close()

			orig := authBaseUrl
			authBaseUrl = ts.URL + "/auth"
			t.Cleanup(func() { authBaseUrl = orig })

			_, err := NewSession("mytoken")

			assert.Equal(t, tc.expectedErr, err.Error())
		})
	}
}

func TestNewSession_NetworkError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "Bearer mytoken", r.Header.Get("Authorization"))

		w.WriteHeader(0)
	}))

	orig := authBaseUrl
	authBaseUrl = ts.URL + "/auth"
	t.Cleanup(func() { authBaseUrl = orig })

	ts.Close()

	_, err := NewSession("mytoken")

	require.Error(t, err)
}
