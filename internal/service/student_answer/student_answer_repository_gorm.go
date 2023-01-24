package studentanswer

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/falentio/raid-go"
	"github.com/falentio/skul/internal/domain"
)

type StudentAnswerRepositoryGorm struct {
	DB *gorm.DB
}

func (r *StudentAnswerRepositoryGorm) ListStudentAnswer(ctx context.Context, o *domain.ListStudentAnswerOptions) ([]*domain.StudentAnswer, error) {
	a := make([]*domain.StudentAnswer, 0)
	err := r.DB.
		WithContext(ctx).
		Model(&domain.StudentAnswer{}).
		Preload("Examination").
		Preload("EnteranceToken").
		Preload("Student").
		Where(o).
		Find(a).
		Error
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *StudentAnswerRepositoryGorm) GetStudentAnswer(ctx context.Context, id raid.Raid) (*domain.StudentAnswer, error) {
	a := &domain.StudentAnswer{}
	err := r.DB.
		WithContext(ctx).
		First(a, "id = ?", id.String()).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = domain.ErrStudentAnswerNotFound
	}
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *StudentAnswerRepositoryGorm) CreateStudentAnswer(ctx context.Context, a *domain.StudentAnswer) error {
	return r.DB.
		WithContext(ctx).
		Omit("Examination").
		Omit("EnteranceToken").
		Omit("Student").
		Create(a).
		Error
}

func (r *StudentAnswerRepositoryGorm) DeleteStudentAnswer(ctx context.Context, id raid.Raid) error {
	return r.DB.
		WithContext(ctx).
		Delete(&domain.StudentAnswer{}, "id = ?", id.String()).
		Error
}
