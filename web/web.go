package web

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed dist/*
var fileSystem embed.FS
var FileSystem fs.FS
var FileServer http.Handler

func init() {
	f, err := fs.Sub(fileSystem, "dist")
	if err != nil {
		panic(err)
	}
	FileSystem = f

	fileServer := http.FileServer(http.FS(f))
	FileServer = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w = &ResponseWriter{w, r, fileServer, false}
		w.Header().Set("Cache-Control", "public, max-age=3600")
		if strings.HasPrefix(r.URL.Path, "/_app/immutable/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		fileServer.ServeHTTP(w, r)
	})
}

type ResponseWriter struct {
	http.ResponseWriter
	r          *http.Request
	FileServer http.Handler
	close      bool
}

func (w *ResponseWriter) WriteHeader(code int) {
	if code == 404 {
		w.close = true
		w.Header().Set("Cache-Control", "private, max-age=0, no-store, no-cache")
		w.r.URL.Path = "/404.html"
		w.FileServer.ServeHTTP(w.ResponseWriter, w.r)
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	if w.close {
		return 0, nil
	}
	return w.ResponseWriter.Write(b)
}
