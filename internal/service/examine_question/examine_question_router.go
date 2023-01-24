package examinequestion

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

type ExamineQuestionRouter struct {
	Auth                   *auth.Auth
	ExamineQuestionService *ExamineQuestionService
}

func (q *ExamineQuestionRouter) Route(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Use(q.Auth.VerifyMiddleware)
	})
}

func (q *ExamineQuestionRouter) CreateExamineQuestion(w http.ResponseWriter, r *http.Request) {
	question := &domain.ExamineQuestion{}

	if err := json.NewDecoder(r.Body).Decode(question); err != nil {
		err = response.NewBadRequest(nil, "failed to decode body")
		response.HandleError(w, r, err)
		return
	}

	res, err := q.ExamineQuestionService.CreateExamineQuestion(r.Context(), question)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (q *ExamineQuestionRouter) UpdateExamineQuestion(w http.ResponseWriter, r *http.Request) {
	question := &domain.ExamineQuestion{}

	if err := json.NewDecoder(r.Body).Decode(question); err != nil {
		err = response.NewBadRequest(nil, "failed to decode body")
		response.HandleError(w, r, err)
		return
	}

	res, err := q.ExamineQuestionService.UpdateExamineQuestion(r.Context(), question)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (q *ExamineQuestionRouter) DeleteExamineQuestion(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "examineQuestionID")
	id, err := raid.RaidFromString(idStr)
	if err != nil {
		err = response.NewBadRequest(nil, "invalid examineQuestionID received: %q", idStr)
		response.HandleError(w, r, err)
		return
	}

	res, err := q.ExamineQuestionService.DeleteExamineQuestion(r.Context(), id)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (q *ExamineQuestionRouter) GetExamineQuestion(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "examineQuestionID")
	id, err := raid.RaidFromString(idStr)
	if err != nil {
		err = response.NewBadRequest(nil, "invalid examineQuestionID received: %q", idStr)
		response.HandleError(w, r, err)
		return
	}

	res, err := q.ExamineQuestionService.GetExamineQuestion(r.Context(), id)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (q *ExamineQuestionRouter) ListExamineQuestion(w http.ResponseWriter, r *http.Request) {
	o := &domain.ListExamineQuestionOptions{}
	if err := o.PageFromQuery(r.URL.Query()); err != nil {
		response.HandleError(w, r, err)
		return
	}

	if r.URL.Query().Has("examineQuestionID") {
		idStr := r.URL.Query().Get("examineQuestionID")
		id, err := raid.RaidFromString(idStr)
		if err != nil {
			err = response.NewBadRequest(nil, "invalid examineQuestionID received: %q", idStr)
			response.HandleError(w, r, err)
			return
		}

		o.ExaminationID = id
	}

	res, err := q.ExamineQuestionService.ListExamineQuestion(r.Context(), o)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}
