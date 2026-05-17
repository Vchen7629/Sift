//go:build unit

package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexRepo_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)

		cookie, err := r.Cookie("JSESSIONID")
		require.NoError(t, err, "expected JSESSIONID cookie")
		assert.Equal(t, "mytoken", cookie.Value)

		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	orig := indexBaseUrl
	indexBaseUrl = ts.URL

	t.Cleanup(func() { indexBaseUrl = orig })

	err := IndexRepo("mytoken", "some-repo")

	require.NoError(t, err)
}

func TestGetJobStatus_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)

		cookie, err := r.Cookie("JSESSIONID")
		require.NoError(t, err, "expected JSESSIONID cookie")
		assert.Equal(t, "mytoken", cookie.Value)

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	orig := indexBaseUrl
	indexBaseUrl = ts.URL

	t.Cleanup(func() { indexBaseUrl = orig })

	_, err := GetJobStatus("mytoken", "some-repo")

	require.NoError(t, err)
}

func TestIndexRepo_InvalidResponse(t *testing.T) {
	tt := []struct {
		name        string
		statusCode  int
		expectedErr string
	}{
		{"Non 202 response (401) returns error", 401, "unexpected error sending req: 401"},
		{"403 Forbidden returns ErrUnauthorized", 403, "unauthorized"},
		{"Unexpected status returns error", 500, "unexpected error sending req: 500"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)

				cookie, err := r.Cookie("JSESSIONID")
				require.NoError(t, err, "expected JSESSIONID cookie")
				assert.Equal(t, "mytoken", cookie.Value)

				w.WriteHeader(tc.statusCode)
			}))
			defer ts.Close()

			orig := indexBaseUrl
			indexBaseUrl = ts.URL
			t.Cleanup(func() { indexBaseUrl = orig })

			err := IndexRepo("mytoken", "some-repo")

			require.EqualError(t, err, tc.expectedErr)
		})
	}
}

func TestGetJobStatus_InvalidResponse(t *testing.T) {
	tt := []struct {
		name        string
		statusCode  int
		expectedErr error
	}{
		{"404 not found returns nil", 404, nil},
		{"403 Forbidden returns ErrUnauthorized", 403, errors.New("unauthorized")},
		{"Unexpected status returns error", 500, errors.New("unexpected error sending req: 500")},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)

				cookie, err := r.Cookie("JSESSIONID")
				require.NoError(t, err, "expected JSESSIONID cookie")
				assert.Equal(t, "mytoken", cookie.Value)

				w.WriteHeader(tc.statusCode)
			}))
			defer ts.Close()

			orig := indexBaseUrl
			indexBaseUrl = ts.URL
			t.Cleanup(func() { indexBaseUrl = orig })

			_, err := GetJobStatus("mytoken", "some-repo")

			require.Equal(t, err, tc.expectedErr)
		})
	}
}

func TestIndexRepo_NetworkError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(0)
	}))

	orig := indexBaseUrl
	indexBaseUrl = ts.URL
	t.Cleanup(func() { indexBaseUrl = orig })

	ts.Close()

	err := IndexRepo("mytoken", "some-repo")

	require.Error(t, err)
}

func TestGetJobStatus_NetworkError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(0)
	}))

	orig := indexBaseUrl
	indexBaseUrl = ts.URL
	t.Cleanup(func() { indexBaseUrl = orig })

	ts.Close()

	_, err := GetJobStatus("mytoken", "some-repo")

	require.Error(t, err)
}
