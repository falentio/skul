package examination

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/falentio/raid-go"
	"github.com/falentio/skul/internal/domain"
	"github.com/falentio/skul/internal/pkg/auth"
	"github.com/falentio/skul/internal/pkg/response"
)

type ExaminationRouter struct {
	ExaminationService domain.ExaminationService
	Auth               *auth.Auth
}

func (e *ExaminationRouter) Route(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Use(e.Auth.VerifyMiddleware)
		r.Get("/list", e.ListExamination)
		r.Get("/{examinationID}", e.GetExamination)
		r.Delete("/{examinationID}", e.DeleteExamination)
		r.Post("/create", e.CreateExamination)
		r.Put("/update", e.UpdateExamination)
	})
}

func (e *ExaminationRouter) ListExamination(w http.ResponseWriter, r *http.Request) {
	o := &domain.ListExaminationOptions{}

	if err := o.PageFromQuery(r.URL.Query()); err != nil {
		err = response.NewBadRequest(nil, "invalid pagination query value")
		response.HandleError(w, r, err)
		return
	}

	res, err := e.ExaminationService.ListExamination(r.Context(), o)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (e *ExaminationRouter) GetExamination(w http.ResponseWriter, r *http.Request) {
	examinationIDStr := chi.URLParam(r, "examinationID")
	examinationID, err := raid.RaidFromString(examinationIDStr)
	if err != nil {
		err = response.NewBadRequest(nil, "invalid examination id, received: %q", examinationIDStr)
		response.HandleError(w, r, err)
		return
	}

	res, err := e.ExaminationService.GetExamination(r.Context(), examinationID)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (e *ExaminationRouter) DeleteExamination(w http.ResponseWriter, r *http.Request) {
	examinationIDStr := chi.URLParam(r, "examinationID")
	examinationID, err := raid.RaidFromString(examinationIDStr)
	if err != nil {
		err = response.NewBadRequest(nil, "invalid examination id, received: %q", examinationIDStr)
		response.HandleError(w, r, err)
		return
	}

	res, err := e.ExaminationService.DeleteExamination(r.Context(), examinationID)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (e *ExaminationRouter) CreateExamination(w http.ResponseWriter, r *http.Request) {
	ex := &domain.Examination{}

	if err := json.NewDecoder(r.Body).Decode(ex); err != nil {
		err = response.NewBadRequest(nil, "failed to decode body")
		response.HandleError(w, r, err)
		return
	}

	res, err := e.ExaminationService.CreateExamination(r.Context(), ex)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (e *ExaminationRouter) UpdateExamination(w http.ResponseWriter, r *http.Request) {
	ex := &domain.Examination{}

	if err := json.NewDecoder(r.Body).Decode(ex); err != nil {
		err = response.NewBadRequest(nil, "failed to decode body")
		response.HandleError(w, r, err)
		return
	}

	res, err := e.ExaminationService.UpdateExamination(r.Context(), ex)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}
