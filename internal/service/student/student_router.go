package student

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/falentio/raid-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"

	"github.com/falentio/skul/internal/domain"
	"github.com/falentio/skul/internal/pkg/auth"
	"github.com/falentio/skul/internal/pkg/response"
)

type StudentRouter struct {
	StudentService domain.StudentService
	Auth           *auth.Auth
}

func (s *StudentRouter) Route(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(5, time.Minute))
		r.Post("/login", s.LoginStudent)
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Use(s.Auth.VerifyMiddleware)
		r.Post("/create", s.CreateStudent)
		r.Post("/create-batch", s.BatchCreateStudent)
		r.Get("/list", s.ListSutdent)
		r.Get("/{studentID}", s.GetStudent)
		r.Delete("/{studentID}", s.DeleteStudent)
		r.Put("/info", s.UpdateStudent)
		r.Put("/info", s.UpdateStudent)
	})
}

func (s *StudentRouter) CreateStudent(w http.ResponseWriter, r *http.Request) {
	adminID, err := s.Auth.GetSubjectRaid(r.Context(), domain.AdminIDPrefix)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	student := &domain.Student{AdminID: adminID}

	if err := json.NewDecoder(r.Body).Decode(student); err != nil {
		err = response.NewBadRequest(nil, "failed to decode body")
		response.HandleError(w, r, err)
		return
	}

	res, err := s.StudentService.CreateStudent(r.Context(), student)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (s *StudentRouter) BatchCreateStudent(w http.ResponseWriter, r *http.Request) {
	_, err := s.Auth.GetSubjectRaid(r.Context(), domain.AdminIDPrefix)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	students := make([]*domain.Student, 0)

	if err := json.NewDecoder(r.Body).Decode(&students); err != nil {
		err = response.NewBadRequest(nil, "failed to decode body")
		response.HandleError(w, r, err)
		return
	}

	res, err := s.StudentService.BatchCreateStudent(r.Context(), students)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (s *StudentRouter) LoginStudent(w http.ResponseWriter, r *http.Request) {
	student := &domain.Student{}

	if err := json.NewDecoder(r.Body).Decode(student); err != nil {
		err = response.NewBadRequest(nil, "failed to decode body")
		response.HandleError(w, r, err)
		return
	}

	res, err := s.StudentService.LoginStudent(r.Context(), student)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (s *StudentRouter) GetStudent(w http.ResponseWriter, r *http.Request) {
	_, err := s.Auth.GetSubjectRaid(r.Context(), domain.AdminIDPrefix)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	studentIDStr := chi.URLParam(r, "studentID")
	studentID, err := raid.RaidFromString(studentIDStr)
	if err != nil {
		err = response.NewBadRequest(nil, "invalid student id, received: %q", studentIDStr)
		response.HandleError(w, r, err)
		return
	}

	res, err := s.StudentService.GetStudent(r.Context(), studentID)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (s *StudentRouter) GetStudentByID(w http.ResponseWriter, r *http.Request) {
	id, err := s.Auth.GetSubjectRaid(r.Context(), "")
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res, err := s.StudentService.GetStudent(r.Context(), id)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (s *StudentRouter) ListSutdent(w http.ResponseWriter, r *http.Request) {
	_, err := s.Auth.GetSubjectRaid(r.Context(), domain.AdminIDPrefix)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	l := &domain.ListStudentOptions{}
	q := r.URL.Query()
	adminID := q.Get("adminID")
	presenceNumber := q.Get("presenceNumber")
	if presenceNumber == "" {
		presenceNumber = "0"
	}
	count := q.Get("count")
	if count == "" {
		count = "10"
	}
	offset := q.Get("offset")
	if offset == "" {
		offset = "0"
	}

	if l.PresenceNumber, err = strconv.Atoi(presenceNumber); err != nil {
		err = response.NewBadRequest(map[string]string{"presenceNumber": "invalid"}, "invalid query param presenceNumber, expected int but got %q", presenceNumber)
		response.HandleError(w, r, err)
		return
	}
	if l.Count, err = strconv.Atoi(count); err != nil {
		err = response.NewBadRequest(map[string]string{"count": "invalid"}, "invalid query param count, expected int but got %q", count)
		response.HandleError(w, r, err)
		return
	}
	if l.Offset, err = strconv.Atoi(offset); err != nil {
		err = response.NewBadRequest(map[string]string{"offset": "invalid"}, "invalid query param offset, expected int but got %q", offset)
		response.HandleError(w, r, err)
		return
	}
	if l.AdminID, err = raid.RaidFromString(adminID); err != nil {
		err = response.NewBadRequest(map[string]string{"adminID": "invalid"}, "invalid query param adminID, got %q", adminID)
		response.HandleError(w, r, err)
		return
	}
	l.Class = q.Get("class")
	l.Grade = q.Get("grade")
	l.Name = q.Get("name")

	res, err := s.StudentService.ListStudent(r.Context(), l)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (s *StudentRouter) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	_, err := s.Auth.GetSubjectRaid(r.Context(), domain.AdminIDPrefix)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	studentIDStr := chi.URLParam(r, "studentID")
	studentID, err := raid.RaidFromString(studentIDStr)
	if err != nil {
		err = response.NewBadRequest(nil, "invalid student id, received: %q", studentIDStr)
		response.HandleError(w, r, err)
		return
	}

	res, err := s.StudentService.DeleteStudent(r.Context(), studentID)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (s *StudentRouter) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	_, err := s.Auth.GetSubjectRaid(r.Context(), domain.AdminIDPrefix)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	student := &domain.Student{}
	if err := json.NewDecoder(r.Body).Decode(student); err != nil {
		err = response.NewBadRequest(nil, "failed to decode body")
		response.HandleError(w, r, err)
		return
	}

	res, err := s.StudentService.UpdateStudent(r.Context(), student)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}
