package admin

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"

	"github.com/falentio/skul/internal/domain"
	"github.com/falentio/skul/internal/pkg/auth"
	"github.com/falentio/skul/internal/pkg/response"
)

type AdminRouter struct {
	AdminService domain.AdminService
	Auth         *auth.Auth
}

func (ar *AdminRouter) Route(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(5, time.Minute))
		r.Post("/", ar.CreateAdmin)
		r.Post("/login", ar.LoginAdmin)
	})
	r.Group(func(r chi.Router) {
		r.Use(ar.Auth.VerifyMiddleware)
		r.Use(middleware.NoCache)
		r.Get("/info", ar.GetAdminByID)
		r.Put("/info", ar.UpdateAdmin)
		r.Delete("/info", ar.DeleteAdmin)
	})
}

func (ar *AdminRouter) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	admin := &domain.Admin{}

	if err := json.NewDecoder(r.Body).Decode(admin); err != nil {
		err = response.NewBadRequest(nil, "failed to decode body")
		response.HandleError(w, r, err)
		return
	}

	res, err := ar.AdminService.CreateAdmin(r.Context(), admin)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (ar *AdminRouter) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	admin := &domain.Admin{}
	if err := json.NewDecoder(r.Body).Decode(admin); err != nil {
		err = response.NewBadRequest(nil, "failed to decode body")
		response.HandleError(w, r, err)
		return
	}

	res, err := ar.AdminService.LoginAdmin(r.Context(), admin)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (ar *AdminRouter) GetAdminByID(w http.ResponseWriter, r *http.Request) {
	adminID, err := ar.Auth.GetSubjectRaid(r.Context(), domain.AdminIDPrefix)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res, err := ar.AdminService.GetAdminByID(r.Context(), adminID)
	if err != nil {
		response.HandleError(w, r, auth.ErrInvalidToken)
		return
	}

	res.ServeHTTP(w, r)
}

func (ar *AdminRouter) UpdateAdmin(w http.ResponseWriter, r *http.Request) {
	adminID, err := ar.Auth.GetSubjectRaid(r.Context(), domain.AdminIDPrefix)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	admin := &domain.Admin{}
	if err := json.NewDecoder(r.Body).Decode(admin); err != nil {
		err = response.NewBadRequest(nil, "failed to decode body")
		response.HandleError(w, r, err)
		return
	}
	admin.ID = adminID

	res, err := ar.AdminService.UpdateAdmin(r.Context(), admin)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}

func (ar *AdminRouter) DeleteAdmin(w http.ResponseWriter, r *http.Request) {
	adminID, err := ar.Auth.GetSubjectRaid(r.Context(), domain.AdminIDPrefix)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res, err := ar.AdminService.DeleteAdmin(r.Context(), adminID)
	if err != nil {
		response.HandleError(w, r, err)
		return
	}

	res.ServeHTTP(w, r)
}
