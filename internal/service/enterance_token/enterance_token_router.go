package enterancetoken

import (
	"encoding/json"
	"net/http"

	"github.com/falentio/raid-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/falentio/skul/internal/domain"
	"github.com/falentio/skul/internal/pkg/auth"
	"github.com/falentio/skul/internal/pkg/response"
)

type EnteranceTokenRouter struct {
	EnteranceTokenService *EnteranceTokenService
	Auth                  *auth.Auth
}

func (e *EnteranceTokenRouter) Route(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(e.Auth.VerifyMiddleware)
		r.Use(middleware.NoCache)
		r.Get("/{enteranceTokenID}", e.GetEnteranceToken)
		r.Get("/{enteranceTokenID}/examination", e.GetEnteranceToken)
		r.Post("/create", e.CreateEnteranceToken)
		r.Put("/{enteranceTokenID}", e.UpdateEnteranceToken)
		r.Delete("/{enteranceTokenID}", e.DeleteEnteranceToken)
	})
}

func (e *EnteranceTokenRouter) GetEnteranceToken(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "enteranceTokenID")
	id, err := raid.RaidFromString(idStr)
	if err != nil {
		err = response.NewBadRequest(nil, "invalid enteranceTokenID received: %q", idStr)
		response.HandleError(w, r, err)
		return
	}

	res, err := e.EnteranceTokenService.GetEnteranceToken(r.Context(), id)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (e *EnteranceTokenRouter) GetExamination(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "enteranceTokenID")
	id, err := raid.RaidFromString(idStr)
	if err != nil {
		err = response.NewBadRequest(nil, "invalid enteranceTokenID received: %q", idStr)
		response.HandleError(w, r, err)
		return
	}

	res, err := e.EnteranceTokenService.GetExamination(r.Context(), id)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (e *EnteranceTokenRouter) DeleteEnteranceToken(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "enteranceTokenID")
	id, err := raid.RaidFromString(idStr)
	if err != nil {
		err = response.NewBadRequest(nil, "invalid enteranceTokenID received: %q", idStr)
		response.HandleError(w, r, err)
		return
	}

	res, err := e.EnteranceTokenService.DeleteEnteranceToken(r.Context(), id)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (e *EnteranceTokenRouter) CreateEnteranceToken(w http.ResponseWriter, r *http.Request) {
	token := &domain.EnteranceToken{}

	if err := json.NewDecoder(r.Body).Decode(token); err != nil {
		err = response.NewBadRequest(nil, "failed to decode body")
		response.HandleError(w, r, err)
		return
	}

	res, err := e.EnteranceTokenService.CreateEnteranceToken(r.Context(), token)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (e *EnteranceTokenRouter) UpdateEnteranceToken(w http.ResponseWriter, r *http.Request) {
	token := &domain.EnteranceToken{}

	if err := json.NewDecoder(r.Body).Decode(token); err != nil {
		err = response.NewBadRequest(nil, "failed to decode body")
		response.HandleError(w, r, err)
		return
	}

	idStr := chi.URLParam(r, "enteranceTokenID")
	id, err := raid.RaidFromString(idStr)
	if err != nil {
		err = response.NewBadRequest(nil, "invalid enteranceTokenID received: %q", idStr)
		response.HandleError(w, r, err)
		return
	}

	token.ID = id
	res, err := e.EnteranceTokenService.UpdateEnteranceToken(r.Context(), token)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (e *EnteranceTokenRouter) ListEnteranceToken(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	o := &domain.ListEnteranceTokenOptions{}
	o.PageFromQuery(q)

	if id, err := raid.RaidFromString(q.Get("examinationID")); q.Has("examinationID") && err == nil {
		o.ExaminationID = id
	} else if err != nil {
		err = response.NewBadRequest(nil, "invalid value for examinationID in url query, received %q", q.Get("examinationID"))
		response.HandleError(w, r, err)
		return
	}

	res, err := e.EnteranceTokenService.ListEnteranceToken(r.Context(), o)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}
