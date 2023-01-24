package domain

import (
	"context"
	"errors"

	"github.com/falentio/raid-go"

	"github.com/falentio/skul/internal/pkg/response"
)

const ExaminationIDPrefix = "exa"

var (
	ErrExaminationNotFound = errors.New("examination: can not find examination")
)

type Examination struct {
	Model

	AdminID raid.Raid `json:"adminID" gorm:"type:varchar(32);not null"`

	Name            string `json:"name"`
	DurationMinutes uint   `json:"durationMinutes"`
	QuestionCount   int    `json:"questionCount"`

	Admin            *Admin             `json:"admin"`
	EnteranceTokens  []*EnteranceToken  `json:"enteranceTokens"`
	ExamineQuestions []*ExamineQuestion `json:"examineQuestions"`
}

type ListExaminationOptions struct {
	PaginateOptions
}

type ExaminationRepositoryRead interface {
	GetExamination(ctx context.Context, examinationID raid.Raid) (*Examination, error)
	ListExamination(ctx context.Context, o *ListExaminationOptions) ([]*Examination, error)
}

type ExaminationRepositoryWrite interface {
	CreateExamination(ctx context.Context, examination *Examination) error
	DeleteExamination(ctx context.Context, examinationID raid.Raid) error
	UpdateExamination(ctx context.Context, examination *Examination) error
}

type ExaminationRepository interface {
	ExaminationRepositoryRead
	ExaminationRepositoryWrite
}

type ExaminationServiceRead interface {
	GetExamination(ctx context.Context, examinationID raid.Raid) (response.Response, error)
	ListExamination(ctx context.Context, o *ListExaminationOptions) (response.Response, error)
}

type ExaminationServiceWrite interface {
	CreateExamination(ctx context.Context, examination *Examination) (response.Response, error)
	DeleteExamination(ctx context.Context, examinationID raid.Raid) (response.Response, error)
	UpdateExamination(ctx context.Context, examination *Examination) (response.Response, error)
}

type ExaminationService interface {
	ExaminationServiceRead
	ExaminationServiceWrite
}
