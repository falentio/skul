package domain

import (
	"context"
	"errors"

	"github.com/falentio/raid-go"
)

const StudentAnswerIDPrefix = "sta"

var (
	ErrStudentAnswerNotFound = errors.New("StudentAnswer: can not find student answer")
)

type StudentAnswer struct {
	Model

	ExamineAnswerID  raid.Raid `json:"examineAnswerID" gorm:"type:varchar(32);not null"`
	StudentID        raid.Raid `json:"studentID" gorm:"type:varchar(32);not null"`
	EnteranceTokenID raid.Raid `json:"enteranceTokenID" gorm:"type:varchar(32);not null"`

	Student        *Student        `json:"student"`
	ExamineAnswer  *ExamineAnswer  `json:"examineAnswer"`
	EnteranceToken *EnteranceToken `json:"enteranceToken"`
}

type ListStudentAnswerOptions struct {
	ExamineAnswerID  raid.Raid
	StudentID        raid.Raid
	EnteranceTokenID raid.Raid
}

type StudentAnswerRepositoryRead interface {
	ListStudentAnswer(ctx context.Context, o *ListStudentAnswerOptions) ([]*StudentAnswer, error)
	GetStudentAnswer(ctx context.Context, id raid.Raid) (*StudentAnswer, error)
}

type StudentAnswerRepositoryWrite interface {
	CreateStudentAnswer(ctx context.Context, a *StudentAnswer) error
	DeleteStudentAnswer(ctx context.Context, id raid.Raid) error
}

type StudentAnswerRepository interface {
	StudentAnswerRepositoryRead
	StudentAnswerRepositoryWrite
}
