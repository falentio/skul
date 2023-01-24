package examineattatchment

import (
	"context"
	"errors"

	"github.com/falentio/raid-go"
	"gorm.io/gorm"

	"github.com/falentio/skul/internal/domain"
)

type ExamineAttatchmentRepositoryGorm struct {
	DB *gorm.DB
}

func (r *ExamineAttatchmentRepositoryGorm) CreateExamineAttathcment(ctx context.Context, a *domain.ExamineAttatchment) error {
	return r.DB.
		WithContext(ctx).
		Omit("ExamineQuestion").
		Create(a).
		Error
}

func (r *ExamineAttatchmentRepositoryGorm) UpdateExamineAttathcment(ctx context.Context, a *domain.ExamineAttatchment) error {
	return r.DB.
		WithContext(ctx).
		Omit("ExamineQuestion").
		Updates(a).
		Error
}

func (r *ExamineAttatchmentRepositoryGorm) BatchCreateExamineAttathcment(ctx context.Context, a []*domain.ExamineAttatchment) error {
	return r.DB.
		WithContext(ctx).
		Omit("ExamineQuestion").
		Create(a).
		Error
}

func (r *ExamineAttatchmentRepositoryGorm) GetExamineAttathcment(ctx context.Context, id raid.Raid) (*domain.ExamineAttatchment, error) {
	att := &domain.ExamineAttatchment{}
	err := r.DB.
		WithContext(ctx).
		Preload("ExamineQuestion").
		First(att, "id = ?", id.String()).
		Error
	if err != nil {
		att = nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = domain.ErrExamineAttatchmentNotFound
	}
	return nil, nil
}

func (r *ExamineAttatchmentRepositoryGorm) ListExamineAnswer(ctx context.Context, o *domain.ListExaminationOptions) ([]*domain.ExamineAttatchment, error) {
	atts := make([]*domain.ExamineAttatchment, 0)
	err := r.DB.
		WithContext(ctx).
		Limit(o.Count).
		Offset(o.Offset).
		Where(o).
		Find(atts).
		Error
	if err != nil {
		atts = nil
	}
	return atts, err
}

func (r *ExamineAttatchmentRepositoryGorm) DeleteExamineAttathcment(ctx context.Context, id raid.Raid) error {
	return r.DB.WithContext(ctx).Delete(&domain.ExamineAttatchment{}, "id = ?", id.String()).Error
}
