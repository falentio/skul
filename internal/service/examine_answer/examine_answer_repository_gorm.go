package examineanswer

import (
	"context"
	"errors"

	"github.com/falentio/raid-go"
	"gorm.io/gorm"

	"github.com/falentio/skul/internal/domain"
)

var _ domain.ExamineAnswerRepository = new(ExamineAnswerRepositoryGorm)

type ExamineAnswerRepositoryGorm struct {
	DB *gorm.DB
}

func (r *ExamineAnswerRepositoryGorm) CreateExamineAnswer(ctx context.Context, examineAnswer *domain.ExamineAnswer) error {
	return r.DB.
		WithContext(ctx).
		Omit("Examination").
		Omit("ExamineQuestion").
		Create(examineAnswer).
		Error
}

func (r *ExamineAnswerRepositoryGorm) UpdateExamineAnswer(ctx context.Context, examineAnswer *domain.ExamineAnswer) error {
	return r.DB.
		WithContext(ctx).
		Omit("Examination").
		Omit("ExamineQuestion").
		Updates(examineAnswer).
		Error
}

func (r *ExamineAnswerRepositoryGorm) BatchCreateExamineAnswer(ctx context.Context, examineAnswer []*domain.ExamineAnswer) error {
	return r.DB.
		WithContext(ctx).
		Omit("Examination").
		Omit("ExamineQuestion").
		Create(examineAnswer).
		Error
}

func (r *ExamineAnswerRepositoryGorm) GetExamineAnswer(ctx context.Context, examineAnswerID raid.Raid) (*domain.ExamineAnswer, error) {
	a := &domain.ExamineAnswer{}

	err := r.DB.
		WithContext(ctx).
		Preload("Examination").
		Preload("ExamineQuestion").
		First(a, "id = ?", examineAnswerID.String()).
		Error
	if err != nil {
		a = nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrExamineAnswerNotFound
	}

	return a, err
}

func (r *ExamineAnswerRepositoryGorm) ListExamineAnswer(ctx context.Context, o *domain.ListExamineAnswerOptions) ([]*domain.ExamineAnswer, error) {
	as := make([]*domain.ExamineAnswer, 0)

	err := r.DB.
		WithContext(ctx).
		Order("ID").
		Limit(o.Count).
		Offset(o.Offset).
		Find(as).
		Error
	if err != nil {
		as = nil
	}
	return as, err
}

func (r *ExamineAnswerRepositoryGorm) DeleteExamineAnswer(ctx context.Context, examineAnswerID raid.Raid) error {
	return r.DB.
		WithContext(ctx).
		Delete(&domain.ExamineAnswer{}, "id = ?", examineAnswerID.String()).
		Error
}
