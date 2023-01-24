package examineattatchment

import (
	"encoding/json"
	"net/http"

	"github.com/falentio/raid-go"
	"github.com/falentio/skul/internal/domain"
	"github.com/falentio/skul/internal/pkg/auth"
	"github.com/falentio/skul/internal/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ExamineAttatchmentRouter struct {
	ExamineAttatchmentService *ExamineAttatchmentService
	Auth                      *auth.Auth
}

func (a *ExamineAttatchmentRouter) Route(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(a.Auth.VerifyMiddleware)
		r.Use(middleware.NoCache)
		r.Get("/{examineAttatchmentID}", a.GetExamineAttatchment)
		r.Delete("/{examineAttatchmentID}", a.DeleteExamineAttatchment)
		r.Post("/create", a.CreateExamineAttatchment)
		r.Put("/{examineAttatchmentID}", a.UpdateExamineAttatchment)
	})
}

func (a *ExamineAttatchmentRouter) CreateExamineAttatchment(w http.ResponseWriter, r *http.Request) {
	at := &domain.ExamineAttatchment{}
	if err := json.NewDecoder(r.Body).Decode(at); err != nil {
		err = response.NewBadRequest(nil, "failed to decode body")
		response.HandleError(w, r, err)
		return
	}

	res, err := a.ExamineAttatchmentService.CreateExamineAttatchment(r.Context(), at)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (a *ExamineAttatchmentRouter) UpdateExamineAttatchment(w http.ResponseWriter, r *http.Request) {
	at := &domain.ExamineAttatchment{}
	if err := json.NewDecoder(r.Body).Decode(at); err != nil {
		err = response.NewBadRequest(nil, "failed to decode body")
		response.HandleError(w, r, err)
		return
	}

	idStr := chi.URLParam(r, "examineAttatchmentID")
	id, err := raid.RaidFromString(idStr)
	if err != nil {
		err = response.NewBadRequest(nil, "invalid examineAttatchmentID received: %q", idStr)
		response.HandleError(w, r, err)
		return
	}

	at.ID = id
	res, err := a.ExamineAttatchmentService.UpdateExamineAttatchment(r.Context(), at)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (a *ExamineAttatchmentRouter) GetExamineAttatchment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "examineAttatchmentID")
	id, err := raid.RaidFromString(idStr)
	if err != nil {
		err = response.NewBadRequest(nil, "invalid examineAttatchmentID received: %q", idStr)
		response.HandleError(w, r, err)
		return
	}

	res, err := a.ExamineAttatchmentService.GetExamineAttatchment(r.Context(), id)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (a *ExamineAttatchmentRouter) DeleteExamineAttatchment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "examineAttatchmentID")
	id, err := raid.RaidFromString(idStr)
	if err != nil {
		err = response.NewBadRequest(nil, "invalid examineAttatchmentID received: %q", idStr)
		response.HandleError(w, r, err)
		return
	}

	res, err := a.ExamineAttatchmentService.DeleteExamineAttatchment(r.Context(), id)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}
