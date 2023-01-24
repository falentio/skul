package admin

import (
	"context"
	"errors"

	"github.com/falentio/raid-go"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	"github.com/falentio/skul/internal/domain"
	"github.com/falentio/skul/internal/pkg/auth"
	"github.com/falentio/skul/internal/pkg/response"
	"github.com/falentio/skul/internal/pkg/validator"
)

var AdminIDFactory = raid.NewRaid().WithPrefix(domain.AdminIDPrefix)
var _ domain.AdminService = new(AdminService)

type AdminService struct {
	AdminRepository domain.AdminRepository
	Auth            *auth.Auth
	Logger          zerolog.Logger
}

func (s *AdminService) CreateAdmin(ctx context.Context, admin *domain.Admin) (res response.Response, err error) {
	admin.ID = AdminIDFactory.WithRandom().WithTimestampNow()
	admin.Students = nil
	admin.Examinations = nil

	err = validator.Struct(admin)
	if err != nil {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	admin.PasswordHash = string(hash)

	err = s.AdminRepository.CreateAdmin(ctx, admin)
	if errors.Is(err, domain.ErrAdminConflict) {
		err = response.NewConflict(nil, "can not create admin, data conflicted with other existing admin")
	}
	if err != nil {
		return
	}

	res = response.NewCreated(admin)
	return
}

func (s *AdminService) LoginAdmin(ctx context.Context, admin *domain.Admin) (res response.Response, err error) {
	a, err := s.AdminRepository.GetAdminByUsername(ctx, admin.Username)
	if errors.Is(err, domain.ErrAdminNotFound) {
		err = response.NewBadRequest(nil, "can not find admin with username %q", admin.Username)
	}
	if err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(admin.Password))
	if err != nil {
		err = response.NewBadRequest(map[string]string{"password": "not match"}, "failed to login to user with username %q", admin.Username)
		return
	}

	c, err := s.Auth.Sign(jwt.RegisteredClaims{
		Subject: a.ID.String(),
	})
	r := response.NewNoContent()
	r.Cookies = append(r.Cookies, c)

	return r, nil
}

func (s *AdminService) GetAdminByID(ctx context.Context, adminID raid.Raid) (res response.Response, err error) {
	admin, err := s.AdminRepository.GetAdminByID(ctx, adminID)
	if err == domain.ErrAdminNotFound {
		err = response.NewNotFound(nil, "can not find admin with id %q", adminID.String())
	}
	if err != nil {
		return
	}
	admin.Password = ""

	res = response.NewOK(admin)
	return
}

func (s *AdminService) GetAdminByUsername(ctx context.Context, username string) (res response.Response, err error) {
	admin, err := s.AdminRepository.GetAdminByUsername(ctx, username)
	if err == domain.ErrAdminNotFound {
		err = response.NewNotFound(nil, "can not find admin with username %q", username)
	}
	if err != nil {
		return
	}
	admin.Password = ""

	res = response.NewOK(admin)
	return
}

func (s *AdminService) UpdateAdmin(ctx context.Context, admin *domain.Admin) (res response.Response, err error) {
	admin.Students = nil
	admin.Examinations = nil

	err = validator.Struct(admin)
	if err != nil {
		return
	}

	if admin.Password != "" {
		var hash []byte
		hash, err = bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return
		}
		admin.Password = string(hash)
	}

	err = s.AdminRepository.UpdateAdmin(ctx, admin)
	if errors.Is(err, domain.ErrAdminConflict) {
		err = response.NewUnprocessableEntity(nil, "can not create admin, data conflicted with other existing admin")
	}
	if err != nil {
		return
	}

	res = response.NewOK(admin)
	return
}

func (s *AdminService) DeleteAdmin(ctx context.Context, adminID raid.Raid) (res response.Response, err error) {
	err = s.AdminRepository.DeleteAdmin(ctx, adminID)
	if err != nil {
		return
	}

	res = response.NewNoContent()
	return
}
