package enterancetoken

import (
	"context"
	"errors"

	"github.com/falentio/raid-go"
	"gorm.io/gorm"

	"github.com/falentio/skul/internal/domain"
)

var _ domain.EnteranceTokenRepository = new(EnteranceTokenRepositoryGorm)

type EnteranceTokenRepositoryGorm struct {
	DB *gorm.DB
}

func (r *EnteranceTokenRepositoryGorm) CreateEnteranceToken(ctx context.Context, token *domain.EnteranceToken) error {
	return r.DB.
		WithContext(ctx).
		Omit("Examination").
		Omit("Students").
		Create(token).
		Error
}

func (r *EnteranceTokenRepositoryGorm) UpdateEnteranceToken(ctx context.Context, token *domain.EnteranceToken) error {
	return r.DB.
		WithContext(ctx).
		Omit("Examination").
		Omit("Students").
		Updates(token).
		Error
}

func (r *EnteranceTokenRepositoryGorm) BatchCreateEnteranceToken(ctx context.Context, tokens []*domain.EnteranceToken) error {
	return r.DB.
		WithContext(ctx).
		Omit("Examination").
		Omit("Students").
		Create(tokens).
		Error
}

func (r *EnteranceTokenRepositoryGorm) GetEnteranceToken(ctx context.Context, tokenID raid.Raid) (*domain.EnteranceToken, error) {
	token := &domain.EnteranceToken{}

	err := r.DB.
		WithContext(ctx).
		Preload("Examination").
		First(token, "id = ?", tokenID.String()).
		Error
	if err != nil {
		token = nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrEnteranceTokenNotFound
	}

	return token, err
}

func (r *EnteranceTokenRepositoryGorm) ListEnteranceToken(ctx context.Context, o *domain.ListEnteranceTokenOptions) ([]*domain.EnteranceToken, error) {
	tokens := make([]*domain.EnteranceToken, 0)

	err := r.DB.
		WithContext(ctx).
		Find(tokens).
		Where(o).
		Order("ID").
		Limit(o.Count).
		Offset(o.Offset).
		Error
	if err != nil {
		tokens = nil
	}

	return tokens, err
}

func (r *EnteranceTokenRepositoryGorm) DeleteEnteranceToken(ctx context.Context, tokenID raid.Raid) error {
	return r.DB.
		WithContext(ctx).
		Delete(&domain.EnteranceToken{}, "id = ?", tokenID.String()).
		Error
}
