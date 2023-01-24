package domain

import (
	"context"
	"errors"

	"github.com/falentio/raid-go"
	"github.com/falentio/skul/internal/pkg/response"
)

var (
	ErrAdminNotFound = errors.New("admin: can not found admin")
	ErrAdminConflict = errors.New("admin: admin data already used")
)

const AdminIDPrefix = "adm"

type Admin struct {
	Model

	Name         string `json:"name" validate:"max=32" gorm:"type:varchar(32)"`
	Username     string `json:"username" validate:"max=32" gorm:"type:varchar(32);unique;index:,expression:LOWER(username)"`
	PasswordHash string `json:"-"`
	Password     string `json:"password,omitempty" validate:"min=8" gorm:"-"`

	Examinations []*Examination `json:"examinations"`
	Students     []*Student     `json:"students"`
}

type AdminRepositoryRead interface {
	GetAdminByID(ctx context.Context, adminID raid.Raid) (*Admin, error)
	GetAdminByUsername(ctx context.Context, username string) (*Admin, error)
}

type AdminRepositoryWrite interface {
	CreateAdmin(ctx context.Context, admin *Admin) error
	DeleteAdmin(ctx context.Context, adminID raid.Raid) error
	UpdateAdmin(ctx context.Context, admin *Admin) error
}

type AdminRepository interface {
	AdminRepositoryRead
	AdminRepositoryWrite
}

type AdminServiceRead interface {
	GetAdminByID(ctx context.Context, adminID raid.Raid) (response.Response, error)
	GetAdminByUsername(ctx context.Context, username string) (response.Response, error)
}

type AdminServiceWrite interface {
	LoginAdmin(ctx context.Context, admin *Admin) (response.Response, error)
	CreateAdmin(ctx context.Context, admin *Admin) (response.Response, error)
	DeleteAdmin(ctx context.Context, adminID raid.Raid) (response.Response, error)
	UpdateAdmin(ctx context.Context, admin *Admin) (response.Response, error)
}

type AdminService interface {
	AdminServiceRead
	AdminServiceWrite
}
