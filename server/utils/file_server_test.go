package utils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

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
