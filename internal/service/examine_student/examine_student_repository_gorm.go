package examinestudent

import (
	"context"

	"github.com/falentio/skul/internal/domain"
	"gorm.io/gorm"
)

type ExamineStudetnRepositoryGorm struct {
	DB *gorm.DB
}

func (r *ExamineStudetnRepositoryGorm) ListExamineStudent(ctx context.Context, o *domain.ListExamineStudentOptions) ([]*domain.ExamineStudent, error) {
	ess := make([]*domain.ExamineStudent, 0)
	err := r.DB.WithContext(ctx).
		Model(&domain.ExamineStudent{}).
		Where(o).
		Order("ID").
		Offset(o.Offset).
		Where("studentID = ?", o.StudentID.String()).
		Where("enteranceTokenID = ?", o.EnteranceTokenID.String()).
		Limit(o.Count).
		Find(ess).
		Error
	return ess, err
}
