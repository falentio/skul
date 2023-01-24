package file

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"

	"github.com/falentio/skul/internal/pkg/auth"
	"github.com/falentio/skul/internal/pkg/response"
)

type FileRouter struct {
	FileService *FileService
	Auth        *auth.Auth
	Logger      zerolog.Logger
}

func (f *FileRouter) Route(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Use(f.Auth.VerifyMiddleware)
		r.Get("/{slug}", f.GetFile)
		r.Delete("/{slug}", f.DeleteFile)
		r.Post("/create", f.CreateFile)
	})
}

func (f *FileRouter) CreateFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 6); err != nil {
		response.HandleError(w, r, err)
		return
	}

	file, h, err := r.FormFile("file")
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res, err := f.FileService.CreateFile(r.Context(), file, h)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (f *FileRouter) DeleteFile(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	res, err := f.FileService.DeleteFile(r.Context(), slug)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (f *FileRouter) GetFile(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	b, h, err := f.FileService.GetFile(r.Context(), slug)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	w.Header().Set("content-type", h)
	if _, err := w.Write(b); err != nil {
		response.HandleError(w, r, err)
		return
	}
}
