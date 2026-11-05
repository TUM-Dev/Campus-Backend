package utils

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// privateStorageDirs are directories below the storage root that hold user uploads
// or scratch data and must never be served over HTTP.
// Feedback attachments are only read from disk, and rating images are returned as bytes by the API.
var privateStorageDirs = []string{"feedback", "rating", "tmp", "storage", "lost+found"}

// PublicFileServer serves the files below root, without directory listings and
// without the private directories in privateStorageDirs.
// It expects the URL prefix to be stripped already (see http.StripPrefix).
func PublicFileServer(root string) http.Handler {
	files := http.FileServer(noDirFS{http.Dir(root)})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPrivateStoragePath(r.URL.Path) {
			http.NotFound(w, r)
			return
		}
		files.ServeHTTP(w, r)
	})
}

func isPrivateStoragePath(p string) bool {
	clean := strings.TrimPrefix(path.Clean("/"+p), "/")
	top, _, _ := strings.Cut(clean, "/")
	for _, dir := range privateStorageDirs {
		if strings.EqualFold(top, dir) {
			return true
		}
	}
	return false
}

// noDirFS hides directories, so that http.FileServer answers 404 instead of a listing.
type noDirFS struct{ fs http.FileSystem }

func (n noDirFS) Open(name string) (http.File, error) {
	f, err := n.fs.Open(name)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	if stat.IsDir() {
		_ = f.Close()
		return nil, fs.ErrNotExist
	}
	return f, nil
}
