package student

import (
	"context"
	"errors"
	"strings"

	"github.com/falentio/raid-go"
	"gorm.io/gorm"

	"github.com/falentio/skul/internal/domain"
)

var _ domain.StudentRepository = new(StudentRepositoryGorm)

type StudentRepositoryGorm struct {
	DB *gorm.DB
}

func (r *StudentRepositoryGorm) CreateStudent(ctx context.Context, student *domain.Student) error {
	err := r.DB.WithContext(ctx).Create(student).Error
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "unique") {
		return domain.ErrStudentConflict
	}
	return err
}

func (r *StudentRepositoryGorm) BatchCreateStudent(ctx context.Context, students []*domain.Student) error {
	err := r.DB.WithContext(ctx).Create(students).Error
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "unique") {
		return domain.ErrStudentConflict
	}
	return err
}

func (r *StudentRepositoryGorm) GetStudent(ctx context.Context, studentID raid.Raid) (*domain.Student, error) {
	student := &domain.Student{Model: domain.Model{ID: studentID}}
	err := r.DB.
		WithContext(ctx).
		Preload("Admin").
		Preload("ExamineAnswer", func(db *gorm.DB) *gorm.DB {
			return db.Limit(10)
		}).
		Preload("EnteranceTokens", func(db *gorm.DB) *gorm.DB {
			return db.Limit(10)
		}).
		First(student, "id = ?", studentID.String()).
		Error
	if err != nil {
		student = nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = domain.ErrStudentNotFound
	}
	return student, err
}

func (r *StudentRepositoryGorm) GetStudentByUsername(ctx context.Context, username string) (*domain.Student, error) {
	student := &domain.Student{Username: username}
	err := r.DB.
		WithContext(ctx).
		Preload("Admin").
		Preload("ExamineAnswer", func(db *gorm.DB) *gorm.DB {
			return db.Limit(10)
		}).
		Preload("EnteranceTokens", func(db *gorm.DB) *gorm.DB {
			return db.Limit(10)
		}).
		First(student, "username = ?", username).
		Error
	if err != nil {
		student = nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = domain.ErrStudentNotFound
	}
	return student, err
}

func (r *StudentRepositoryGorm) ListStudent(ctx context.Context, opts *domain.ListStudentOptions) ([]*domain.Student, error) {
	students := make([]*domain.Student, 0)
	err := r.DB.WithContext(ctx).
		Model(&domain.Student{}).
		Where(opts).
		Order("ID").
		Offset(opts.Offset).
		Limit(opts.Count).
		Find(students).
		Error
	return students, err
}

func (r *StudentRepositoryGorm) DeleteStudent(ctx context.Context, studentID raid.Raid) error {
	err := r.DB.WithContext(ctx).Delete(&domain.Student{Model: domain.Model{ID: studentID}}).Error
	return err
}

func (r *StudentRepositoryGorm) UpdateStudent(ctx context.Context, student *domain.Student) error {
	err := r.DB.WithContext(ctx).Updates(student).Error
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "unique") {
		return domain.ErrStudentConflict
	}
	return err
}
