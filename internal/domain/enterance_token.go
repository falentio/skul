package domain

import (
	"context"
	"errors"
	"time"

	"github.com/falentio/raid-go"
	"github.com/falentio/skul/internal/pkg/response"
)

const EnteranceTokenIDPrefix = "eto"

var (
	ErrEnteranceTokenNotFound = errors.New("EnteranceToken: can not find enterance token")
)

type EnteranceToken struct {
	Model

	ExaminationID raid.Raid `json:"examinationID" gorm:"type:varchar(32);not null"`

	EnteranceFrom  time.Time `json:"enteranceFrom"`
	EnteranceUntil time.Time `json:"enteranceUntil"`

	Examination *Examination `json:"examination"`
	Students    []*Student   `json:"student" gorm:"many2many:examine_student"`
}

type ListEnteranceTokenOptions struct {
	PaginateOptions

	ExaminationID raid.Raid `json:"examinationID"`
}

type EnteranceTokenRepositoryRead interface {
	GetEnteranceToken(ctx context.Context, tokenID raid.Raid) (*EnteranceToken, error)
	ListEnteranceToken(ctx context.Context, o *ListEnteranceTokenOptions) ([]*EnteranceToken, error)
}

type EnteranceTokenRepositoryWrite interface {
	CreateEnteranceToken(ctx context.Context, token *EnteranceToken) error
	UpdateEnteranceToken(ctx context.Context, token *EnteranceToken) error
	BatchCreateEnteranceToken(ctx context.Context, tokens []*EnteranceToken) error
	DeleteEnteranceToken(ctx context.Context, tokenID raid.Raid) error
}

type EnteranceTokenRepository interface {
	EnteranceTokenRepositoryRead
	EnteranceTokenRepositoryWrite
}

type EnteranceTokenServiceRead interface {
	GetEnteranceToken(ctx context.Context, tokenID raid.Raid) (response.Response, error)
	ListEnteranceToken(ctx context.Context, o *ListEnteranceTokenOptions) (response.Response, error)
}

type EnteranceTokenServiceWrite interface {
	CreateEnteranceToken(ctx context.Context, token *EnteranceToken) (response.Response, error)
	UpdateEnteranceToken(ctx context.Context, token *EnteranceToken) (response.Response, error)
	BatchCreateEnteranceToken(ctx context.Context, tokens []*EnteranceToken) (response.Response, error)
	DeleteEnteranceToken(ctx context.Context, tokenID raid.Raid) (response.Response, error)
}

type EnteranceTokenService interface {
	EnteranceTokenServiceRead
	EnteranceTokenServiceWrite
}
