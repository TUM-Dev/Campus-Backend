package cron

import (
	"io"
	"net/http"
	"strings"
	"testing"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/internal/webfetch"
	"github.com/stretchr/testify/require"
)

type httpTestTransport func(*http.Request) (*http.Response, error)

func (f httpTestTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestFailedFileDownloadReturnsError(t *testing.T) {
	old := webfetch.Client
	t.Cleanup(func() { webfetch.Client = old })
	webfetch.Client = &http.Client{Transport: httpTestTransport(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 404, Status: "404 Not Found", Body: io.NopCloser(strings.NewReader("not found"))}, nil
	})}
	require.Error(t, downloadFile("https://example.com/missing.jpg", t.TempDir()+"/cover.jpg"))
}

func TestStudentClubFetchFailureDoesNotPanic(t *testing.T) {
	old := webfetch.Client
	t.Cleanup(func() { webfetch.Client = old })
	for _, code := range []int{404, 503} {
		webfetch.Client = &http.Client{Transport: httpTestTransport(func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: code, Status: http.StatusText(code), Body: io.NopCloser(strings.NewReader("unavailable"))}, nil
		})}
		require.NotPanics(t, func() { require.Error(t, (&CronService{}).studentClubCron(pb.Language_German)) })
	}
}
