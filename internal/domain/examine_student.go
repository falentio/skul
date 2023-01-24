package domain

import (
	"context"
	"time"

	"github.com/falentio/raid-go"
	"gorm.io/gorm"
)

type ExamineStudent struct {
	gorm.Model

	EnteranceTokenID raid.Raid `json:"enteranceTokenID" gorm:"type:varchar(32);primaryKey;not null"`
	StudentID        raid.Raid `json:"graderID" gorm:"type:varchar(32);primaryKey;not null"`

	DueDate time.Time `json:"dueDate"`

	EnteranceToken *EnteranceToken `json:"enteranceToken"`
	Student        *Student        `json:"student"`
}

type ListExamineStudentOptions struct {
	PaginateOptions

	StudentID        raid.Raid
	EnteranceTokenID raid.Raid
}

type ExamineStudentRepositoryRead interface {
	ListExamineStudent(ctx context.Context, o *ListExamineStudentOptions) ([]*ExamineStudent, error)
}

type ExamineStudentRepositoryWrite interface {
	CreateExamineStudent(ctx context.Context, examineStudent *ExamineStudent) error
	DeleteExamineStudent(ctx context.Context, examineStudentID raid.Raid) error
}

type ExamineStudentRepository interface {
	ExamineStudentRepositoryRead
	ExamineStudentRepositoryWrite
}
