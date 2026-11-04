package webfetch

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type testTransport func(*http.Request) (*http.Response, error)

func (f testTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type trackedBody struct {
	io.Reader
	closed bool
}

func (b *trackedBody) Close() error { b.closed = true; return nil }

func TestJSONClosesBodiesOnAllResponses(t *testing.T) {
	old := Client
	t.Cleanup(func() { Client = old })
	for _, test := range []struct {
		status    int
		body      string
		wantError bool
	}{
		{200, `{"name":"Pasta"}`, false},
		{200, `invalid json`, true},
		{404, `not found`, true},
	} {
		body := &trackedBody{Reader: strings.NewReader(test.body)}
		Client = &http.Client{Transport: testTransport(func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: test.status, Status: http.StatusText(test.status), Body: body}, nil
		})}
		var value map[string]string
		err := JSON("https://example.com/menu", &value)
		if test.wantError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, "Pasta", value["name"])
		}
		require.True(t, body.closed)
	}
}
