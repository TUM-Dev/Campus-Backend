package utils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPublicFileServer(t *testing.T) {
	root := t.TempDir()
	write := func(rel, content string) {
		p := filepath.Join(root, rel)
		require.NoError(t, os.MkdirAll(filepath.Dir(p), 0o755))
		require.NoError(t, os.WriteFile(p, []byte(content), 0o644))
	}
	write("news/newspread/a.jpg", "news image")
	write("feedback/0b5c-uuid/0.png", "private")
	write("rating/cafeterias/1/abc.jpeg", "private")
	write("tmp/x", "private")
	write("index.html", "root index must not be served in place of a listing")

	srv := httptest.NewServer(http.StripPrefix("/files/", PublicFileServer(root)))
	defer srv.Close()

	cases := []struct {
		path string
		code int
	}{
		{"/files/news/newspread/a.jpg", http.StatusOK},
		{"/files/", http.StatusNotFound},
		{"/files/news/", http.StatusNotFound},
		{"/files/news", http.StatusNotFound},
		{"/files/feedback/", http.StatusNotFound},
		{"/files/feedback/0b5c-uuid/0.png", http.StatusNotFound},
		{"/files/Feedback/0b5c-uuid/0.png", http.StatusNotFound},
		{"/files/news/../feedback/0b5c-uuid/0.png", http.StatusNotFound},
		{"/files/./rating/cafeterias/1/abc.jpeg", http.StatusNotFound},
		{"/files/tmp/x", http.StatusNotFound},
		{"/files/does-not-exist.png", http.StatusNotFound},
	}
	client := &http.Client{} // follows redirects, so the final target is what counts
	for _, c := range cases {
		t.Run(c.path, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, srv.URL+c.path, nil)
			require.NoError(t, err)
			// keep "..", "." and the like as sent, like a hostile client would
			req.URL.Opaque = c.path
			res, err := client.Do(req)
			require.NoError(t, err)
			_ = res.Body.Close()
			require.Equal(t, c.code, res.StatusCode)
		})
	}
}

func TestPublicFileServerSandboxesDocuments(t *testing.T) {
	root := t.TempDir()
	svg := `<svg xmlns="http://www.w3.org/2000/svg"><style>rect { fill: red; }</style><script>window.scriptRan = true</script><rect width="10" height="10"/></svg>`
	html := `<html><script>window.scriptRan = true</script></html>`
	require.NoError(t, os.WriteFile(filepath.Join(root, "logo.svg"), []byte(svg), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "page.html"), []byte(html), 0o644))
	handler := http.StripPrefix("/files/", PublicFileServer(root))

	cases := []struct {
		name, method, path, header, value, body, contentType string
		code                                                 int
	}{
		{name: "svg", method: http.MethodGet, path: "/files/logo.svg", code: http.StatusOK, body: svg, contentType: "image/svg+xml"},
		{name: "html", method: http.MethodGet, path: "/files/page.html", code: http.StatusOK, body: html, contentType: "text/html"},
		{name: "head", method: http.MethodHead, path: "/files/logo.svg", code: http.StatusOK, contentType: "image/svg+xml"},
		{name: "range", method: http.MethodGet, path: "/files/logo.svg", header: "Range", value: "bytes=0-4", code: http.StatusPartialContent, body: svg[:5], contentType: "image/svg+xml"},
		{name: "not modified", method: http.MethodGet, path: "/files/logo.svg", header: "If-Modified-Since", value: time.Now().Add(time.Hour).UTC().Format(http.TimeFormat), code: http.StatusNotModified},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			request := httptest.NewRequest(c.method, c.path, nil)
			if c.header != "" {
				request.Header.Set(c.header, c.value)
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)

			require.Equal(t, c.code, response.Code)
			require.Equal(t, c.body, response.Body.String())
			if c.contentType != "" {
				require.Contains(t, response.Header().Get("Content-Type"), c.contentType)
			}
			policy := response.Header().Get("Content-Security-Policy")
			// A bare sandbox directive grants neither scripts nor same-origin access.
			require.Contains(t, strings.Split(policy, ";"), "sandbox")
			require.Contains(t, policy, "default-src 'none'")
			require.Equal(t, "nosniff", response.Header().Get("X-Content-Type-Options"))
		})
	}
}
