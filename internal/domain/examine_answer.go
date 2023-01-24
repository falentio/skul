package domain

import (
	"context"
	"errors"

	"github.com/falentio/raid-go"
)

const ExamineAnswerIDPrefix = "xan"

var (
	ErrExamineAnswerNotFound = errors.New("ExamineAnswer: can not find examine answer")
)

type ExamineAnswer struct {
	Model

	ExaminationID     raid.Raid `json:"examinationID" gorm:"type:varchar(32);not null"`
	ExamineQuestionID raid.Raid `json:"examineQuestionID" gorm:"type:varchar(32);not null"`

	Correct bool   `json:"correct"`
	Answer  string `json:"answer"`

	Examination     *Examination     `json:"examination"`
	ExamineQuestion *ExamineQuestion `json:"examineQuestion"`
}

type ListExamineAnswerOptions struct {
	PaginateOptions

	ExaminationID raid.Raid `json:"examinationID"`
}

type ExamineAnswerRepositoryRead interface {
	GetExamineAnswer(ctx context.Context, examineAnswerID raid.Raid) (*ExamineAnswer, error)
	ListExamineAnswer(ctx context.Context, o *ListExamineAnswerOptions) ([]*ExamineAnswer, error)
}

type ExamineAnswerRepositoryWrite interface {
	CreateExamineAnswer(ctx context.Context, examineAnswer *ExamineAnswer) error
	BatchCreateExamineAnswer(ctx context.Context, examineAnswers []*ExamineAnswer) error
	DeleteExamineAnswer(ctx context.Context, examineAnswerID raid.Raid) error
	UpdateExamineAnswer(ctx context.Context, examineAnswer *ExamineAnswer) error
}

type ExamineAnswerRepository interface {
	ExamineAnswerRepositoryRead
	ExamineAnswerRepositoryWrite
}
