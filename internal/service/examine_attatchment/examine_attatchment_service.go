package examineattatchment

import (
	"context"
	"errors"

	"github.com/rs/zerolog"

	"github.com/falentio/raid-go"
	"github.com/falentio/skul/internal/domain"
	"github.com/falentio/skul/internal/pkg/auth"
	"github.com/falentio/skul/internal/pkg/response"
	"github.com/falentio/skul/internal/pkg/validator"
)

type ExamineAttatchmentService struct {
	ExamineAttatchmentRepository domain.ExamineAttatchmentRepository
	Auth                         *auth.Auth
	Logger                       zerolog.Logger
}

func (s *ExamineAttatchmentService) CreateExamineAttatchment(ctx context.Context, a *domain.ExamineAttatchment) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	if err := validator.Struct(a); err != nil {
		return nil, err
	}

	a.ID = raid.NewRaid().WithPrefix(domain.ExamineAttatchmentIDPrefix)
	if err := s.ExamineAttatchmentRepository.CreateExamineAttathcment(ctx, a); err != nil {
		return nil, err
	}

	return response.NewOK(a), nil
}

func (s *ExamineAttatchmentService) UpdateExamineAttatchment(ctx context.Context, a *domain.ExamineAttatchment) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	if err := validator.Struct(a); err != nil {
		return nil, err
	}

	if err := s.ExamineAttatchmentRepository.UpdateExamineAttathcment(ctx, a); err != nil {
		return nil, err
	}

	return response.NewOK(a), nil
}

func (s *ExamineAttatchmentService) GetExamineAttatchment(ctx context.Context, id raid.Raid) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	a, err := s.ExamineAttatchmentRepository.GetExamineAttathcment(ctx, id)
	if errors.Is(err, domain.ErrExamineAttatchmentNotFound) {
		return nil, response.NewNotFound(nil, "can not find examine attatchment with id %q", id)
	}
	if err != nil {
		return nil, err
	}

	return response.NewOK(a), nil
}

func (s *ExamineAttatchmentService) DeleteExamineAttatchment(ctx context.Context, id raid.Raid) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	err = s.ExamineAttatchmentRepository.DeleteExamineAttathcment(ctx, id)
	if err != nil {
		return nil, err
	}

	return response.NewNoContent(), nil
}
