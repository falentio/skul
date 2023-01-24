package domain

import (
	"context"
	"errors"

	"github.com/falentio/raid-go"
)

const ExamineAttatchmentIDPrefix = "xat"

var (
	ErrExamineAttatchmentNotFound = errors.New("ExamineAttatchment: can not find examine attatchment")
)

type ExamineAttatchment struct {
	Model

	ExamineQuestionID raid.Raid `json:"examineQuestionID" gorm:"type:varchar(32);not null"`

	Type string `json:"type"`
	Slug string `json:"slug"`

	ExamineQuestion *ExamineQuestion `json:"examineQuestion"`
}

type ListExamineAttatchmentOptions struct {
	PaginateOptions

	ExamineQuestionID raid.Raid `json:"examineQuestionID"`
}

type ExamineAttatchmentRepositoryWrite interface {
	CreateExamineAttathcment(ctx context.Context, a *ExamineAttatchment) error
	DeleteExamineAttathcment(ctx context.Context, id raid.Raid) error
	UpdateExamineAttathcment(ctx context.Context, a *ExamineAttatchment) error
}

type ExamineAttatchmentRepositoryRead interface {
	GetExamineAttathcment(ctx context.Context, id raid.Raid) (*ExamineAttatchment, error)
	ListExamineAnswer(ctx context.Context, o *ListExaminationOptions) ([]*ExamineAttatchment, error)
}

type ExamineAttatchmentRepository interface {
	ExamineAttatchmentRepositoryRead
	ExamineAttatchmentRepositoryWrite
}
