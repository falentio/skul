package examination

import (
	"context"
	"errors"

	"github.com/falentio/raid-go"
	"gorm.io/gorm"

	"github.com/falentio/skul/internal/domain"
)

var _ domain.ExaminationRepository = new(ExaminationRepositoryGorm)

type ExaminationRepositoryGorm struct {
	DB *gorm.DB
}

func (r *ExaminationRepositoryGorm) CreateExamination(ctx context.Context, examination *domain.Examination) error {
	return r.DB.
		WithContext(ctx).
		Omit("Admin").
		Omit("EnteranceTokens").
		Omit("ExamineQuestions").
		Create(examination).
		Error
}

func (r *ExaminationRepositoryGorm) GetExamination(ctx context.Context, examinationID raid.Raid) (*domain.Examination, error) {
	ex := &domain.Examination{}
	ex.ID = examinationID
	err := r.DB.
		WithContext(ctx).
		Preload("Admin").
		Preload("EnteranceTokens").
		Preload("ExamineQuestions").
		Preload("ExamineQuestions.ExamineAnswers").
		Preload("ExamineQuestions.ExamineAttatchment").
		First(ex, "id = ?", ex.ID.String()).
		Error
	if err != nil {
		ex = nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = domain.ErrExaminationNotFound
	}
	return ex, err
}

func (r *ExaminationRepositoryGorm) ListExamination(ctx context.Context, o *domain.ListExaminationOptions) ([]*domain.Examination, error) {
	exs := make([]*domain.Examination, 0)
	err := r.DB.
		WithContext(ctx).
		Model(&domain.Examination{}).
		Where(o).
		Limit(o.Count).
		Offset(o.Offset).
		Find(exs).
		Error
	return exs, err
}

func (r *ExaminationRepositoryGorm) DeleteExamination(ctx context.Context, examinationID raid.Raid) error {
	ex := &domain.Examination{}
	ex.ID = examinationID
	return r.DB.WithContext(ctx).Delete(ex).Error
}

func (r *ExaminationRepositoryGorm) UpdateExamination(ctx context.Context, examination *domain.Examination) error {
	return r.DB.
		WithContext(ctx).
		Updates(examination).
		Error
}
