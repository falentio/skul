package examinequestion

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/falentio/raid-go"
	"github.com/falentio/skul/internal/domain"
)

var _ domain.ExamineQuestionRepository = new(ExamineQuestionRepositoryGorm)

type ExamineQuestionRepositoryGorm struct {
	DB *gorm.DB
}

func (r *ExamineQuestionRepositoryGorm) CreateExamineQuestion(ctx context.Context, examineQuestion *domain.ExamineQuestion) error {
	return r.DB.
		WithContext(ctx).
		Omit("Examination").
		Omit("ExamineAnswer").
		Omit("ExamineAttatchments").
		Create(examineQuestion).
		Error
}

func (r *ExamineQuestionRepositoryGorm) GetExamineQuestion(ctx context.Context, examineQuestionID raid.Raid) (*domain.ExamineQuestion, error) {
	q := &domain.ExamineQuestion{}
	err := r.DB.
		WithContext(ctx).
		Preload("Examination").
		Preload("ExamineAnswer").
		Preload("ExamineAttatchments").
		First(q, "id = ?", examineQuestionID.String()).
		Error
	if err != nil {
		q = nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = domain.ErrExamineQuestionNotFound
	}

	return q, err
}

func (r *ExamineQuestionRepositoryGorm) ListExamineQuestion(ctx context.Context, o *domain.ListExamineQuestionOptions) ([]*domain.ExamineQuestion, error) {
	qs := make([]*domain.ExamineQuestion, 0)

	err := r.DB.
		WithContext(ctx).
		Find(qs).
		Where(o).
		Limit(o.Count).
		Offset(o.Offset).
		Error
	if err != nil {
		qs = nil
	}

	return qs, err
}

func (r *ExamineQuestionRepositoryGorm) DeleteExamineQuestion(ctx context.Context, examineQuestionID raid.Raid) error {
	return r.DB.
		WithContext(ctx).
		Delete(&domain.ExamineQuestion{}, "id = ?", examineQuestionID.String()).
		Error
}

func (r *ExamineQuestionRepositoryGorm) UpdateExamineQuestion(ctx context.Context, examineQuestion *domain.ExamineQuestion) error {
	return r.DB.
		WithContext(ctx).
		Updates(examineQuestion).
		Error
}
