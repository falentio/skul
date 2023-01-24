package admin

import (
	"context"
	"errors"
	"strings"

	"github.com/falentio/raid-go"
	"gorm.io/gorm"

	"github.com/falentio/skul/internal/domain"
)

var _ domain.AdminRepository = new(AdminRepositoryGorm)

type AdminRepositoryGorm struct {
	DB *gorm.DB
}

func (r *AdminRepositoryGorm) CreateAdmin(ctx context.Context, admin *domain.Admin) error {
	err := r.DB.
		WithContext(ctx).
		Omit("Examinations").
		Omit("Students").
		Create(admin).
		Error
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "unique") {
		return domain.ErrAdminConflict
	}
	return err
}

func (r *AdminRepositoryGorm) GetAdminByID(ctx context.Context, adminID raid.Raid) (*domain.Admin, error) {
	admin := &domain.Admin{Model: domain.Model{ID: adminID}}
	err := r.DB.
		WithContext(ctx).
		Preload("Examinations", func(db *gorm.DB) *gorm.DB {
			return db.Limit(10)
		}).
		Preload("Students", func(db *gorm.DB) *gorm.DB {
			return db.Limit(10)
		}).
		First(admin, "id = ?", adminID.String()).
		Error
	if err != nil {
		admin = nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = domain.ErrAdminNotFound
	}
	return admin, err
}

func (r *AdminRepositoryGorm) GetAdminByUsername(ctx context.Context, username string) (*domain.Admin, error) {
	admin := &domain.Admin{}
	err := r.DB.
		WithContext(ctx).
		Preload("Examinations", func(db *gorm.DB) *gorm.DB {
			return db.Limit(10)
		}).
		Preload("Students", func(db *gorm.DB) *gorm.DB {
			return db.Limit(10)
		}).
		First(admin, "username = ?", username).
		Error
	if err != nil {
		admin = nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = domain.ErrAdminNotFound
	}
	return admin, err
}

func (r *AdminRepositoryGorm) DeleteAdmin(ctx context.Context, adminID raid.Raid) error {
	err := r.DB.
		WithContext(ctx).
		Delete(&domain.Admin{Model: domain.Model{ID: adminID}}).
		Error
	return err
}

func (r *AdminRepositoryGorm) UpdateAdmin(ctx context.Context, admin *domain.Admin) error {
	err := r.DB.
		WithContext(ctx).
		Updates(admin).
		Error
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "unique") {
		return domain.ErrAdminConflict
	}
	return err
}
