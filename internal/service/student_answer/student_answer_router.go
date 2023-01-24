package studentanswer

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/falentio/raid-go"
	"github.com/falentio/skul/internal/domain"
	"github.com/falentio/skul/internal/pkg/auth"
	"github.com/falentio/skul/internal/pkg/response"
	"github.com/falentio/skul/internal/pkg/validator"
)

type StudentAnswerRouter struct {
	StudentAnswerService *StudentAnswerService
	Auth                 *auth.Auth
}

func (s *StudentAnswerRouter) Route(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Use(s.Auth.VerifyMiddleware)
		r.Get("/list", s.ListStudentAnswer)
		r.Post("/create", s.CreateStudentAnswer)
		r.Delete("/{studentAnswerID}", s.DeleteStudentAnswer)
	})
}

func (s *StudentAnswerRouter) ListStudentAnswer(w http.ResponseWriter, r *http.Request) {
	o := &domain.ListStudentAnswerOptions{}
	q := r.URL.Query()
	if q.Has("studentID") {
		id, err := raid.RaidFromString(q.Get("studentID"))
		if err != nil {
			err = response.NewBadRequest(nil, "invalid value for query studentID")
			response.HandleError(w, r, err)
			return
		}
		o.StudentID = id
	}

	if q.Has("enteranceTokenID") {
		id, err := raid.RaidFromString(q.Get("enteranceTokenID"))
		if err != nil {
			err = response.NewBadRequest(nil, "invalid value for query enteranceTokenID")
			response.HandleError(w, r, err)
			return
		}
		o.EnteranceTokenID = id
	}

	if q.Has("examineAnswerID") {
		id, err := raid.RaidFromString(q.Get("examineAnswerID"))
		if err != nil {
			err = response.NewBadRequest(nil, "invalid value for query examineAnswerID")
			response.HandleError(w, r, err)
			return
		}
		o.ExamineAnswerID = id
	}

	res, err := s.StudentAnswerService.ListStudentAnswer(r.Context(), o)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (s *StudentAnswerRouter) DeleteStudentAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "studentAnswerID")
	id, err := raid.RaidFromString(idStr)
	if err != nil {
		err = response.NewBadRequest(nil, "invalid studentAnswerID received: %q", idStr)
		response.HandleError(w, r, err)
		return
	}

	res, err := s.StudentAnswerService.DeleteStudentAnswer(r.Context(), id)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (s *StudentAnswerRouter) CreateStudentAnswer(w http.ResponseWriter, r *http.Request) {
	a := &domain.StudentAnswer{}
	if err := json.NewDecoder(r.Body).Decode(a); err != nil {
		err = response.NewBadRequest(nil, "failed to decode body")
		response.HandleError(w, r, err)
		return
	}

	if err := validator.Struct(a); err != nil {
		response.HandleError(w, r, err)
		return
	}

	res, err := s.StudentAnswerService.CreateStudentAnswer(r.Context(), a)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}
