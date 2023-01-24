package domain

import (
	"context"
	"errors"

	"github.com/falentio/raid-go"
)

const ExamineQuestionIDPrefix = "xqs"

var ErrExamineQuestionNotFound = errors.New("ExamineQuestion: can not find examine question")

type ExamineQuestion struct {
	Model

	ExaminationID raid.Raid `json:"examinationID" gorm:"type:varchar(32);not null"`

	Question    string `json:"question"`
	AnswerCount int    `json:"answerCount" validate:"min=1"`

	Examination        *Examination        `json:"examination"`
	ExamineAnswers     []*ExamineAnswer    `json:"examineAnswers"`
	ExamineAttatchment *ExamineAttatchment `json:"examineAttatchment"`
}

type ListExamineQuestionOptions struct {
	PaginateOptions

	ExaminationID raid.Raid `json:"examinationID"`
}

type ExamineQuestionRepositoryRead interface {
	GetExamineQuestion(ctx context.Context, examineQuestionID raid.Raid) (*ExamineQuestion, error)
	ListExamineQuestion(ctx context.Context, o *ListExamineQuestionOptions) ([]*ExamineQuestion, error)
}

type ExamineQuestionRepositoryWrite interface {
	CreateExamineQuestion(ctx context.Context, examineQuestion *ExamineQuestion) error
	DeleteExamineQuestion(ctx context.Context, examineQuestionID raid.Raid) error
	UpdateExamineQuestion(ctx context.Context, examineQuestion *ExamineQuestion) error
}

type ExamineQuestionRepository interface {
	ExamineQuestionRepositoryRead
	ExamineQuestionRepositoryWrite
}
