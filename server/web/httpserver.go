package web

import (
	"net"
	"net/http"
)

func HTTPServe(l net.Listener) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Welcome to the new campus API."))
	})

	s := &http.Server{Handler: mux}
	return s.Serve(l)
}
