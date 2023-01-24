package examineanswer

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

type ExamineAnswerRouter struct {
	ExamineAnswerService *ExamineAnswerService
	Auth                 *auth.Auth
}

func (e *ExamineAnswerRouter) Route(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(e.Auth.VerifyMiddleware)
		r.Use(middleware.NoCache)
		r.Get("/{examineAnswerID}", e.GetExamineAnswer)
		r.Delete("/{examineAnswerID}", e.DeleteExamineAnswer)
		r.Post("/create", e.CreateExamineAnswer)
		r.Put("/{examineAnswerID}", e.UpdateExamineAnswer)
	})
}

func (e *ExamineAnswerRouter) CreateExamineAnswer(w http.ResponseWriter, r *http.Request) {
	as := &domain.ExamineAnswer{}
	if err := json.NewDecoder(r.Body).Decode(as); err != nil {
		err = response.NewBadRequest(nil, "failed to decode body")
		response.HandleError(w, r, err)
		return
	}

	res, err := e.ExamineAnswerService.CreateExamineAnswer(r.Context(), as)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (e *ExamineAnswerRouter) UpdateExamineAnswer(w http.ResponseWriter, r *http.Request) {
	as := &domain.ExamineAnswer{}
	if err := json.NewDecoder(r.Body).Decode(as); err != nil {
		err = response.NewBadRequest(nil, "failed to decode body")
		response.HandleError(w, r, err)
		return
	}

	idStr := chi.URLParam(r, "examineAnswerID")
	id, err := raid.RaidFromString(idStr)
	if err != nil {
		err = response.NewBadRequest(nil, "invalid examineAnswerID received: %q", idStr)
		response.HandleError(w, r, err)
		return
	}

	as.ID = id
	res, err := e.ExamineAnswerService.UpdateExamineAnswer(r.Context(), as)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (e *ExamineAnswerRouter) GetExamineAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "examineAnswerID")
	id, err := raid.RaidFromString(idStr)
	if err != nil {
		err = response.NewBadRequest(nil, "invalid examineAnswerID received: %q", idStr)
		response.HandleError(w, r, err)
		return
	}

	res, err := e.ExamineAnswerService.GetExamineAnswer(r.Context(), id)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (e *ExamineAnswerRouter) DeleteExamineAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "examineAnswerID")
	id, err := raid.RaidFromString(idStr)
	if err != nil {
		err = response.NewBadRequest(nil, "invalid examineAnswerID received: %q", idStr)
		response.HandleError(w, r, err)
		return
	}

	res, err := e.ExamineAnswerService.DeleteExamineAnswer(r.Context(), id)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}
