package examineanswer

import (
	"context"

	"github.com/falentio/raid-go"
	"github.com/falentio/skul/internal/domain"
	"github.com/falentio/skul/internal/pkg/auth"
	"github.com/falentio/skul/internal/pkg/response"
	"github.com/falentio/skul/internal/pkg/validator"
)

type ExamineAnswerService struct {
	ExamineAnswerRepository domain.ExamineAnswerRepository
	Auth                    *auth.Auth
}

func (s *ExamineAnswerService) CreateExamineAnswer(ctx context.Context, as *domain.ExamineAnswer) (response.Response, error) {
	if _, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix); err != nil {
		return nil, err
	}

	if err := validator.Struct(as); err != nil {
		return nil, err
	}

	err := s.ExamineAnswerRepository.CreateExamineAnswer(ctx, as)
	if err != nil {
		return nil, err
	}

	return response.NewOK(as), nil
}

func (s *ExamineAnswerService) UpdateExamineAnswer(ctx context.Context, as *domain.ExamineAnswer) (response.Response, error) {
	if _, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix); err != nil {
		return nil, err
	}

	if err := validator.Struct(as); err != nil {
		return nil, err
	}

	err := s.ExamineAnswerRepository.CreateExamineAnswer(ctx, as)
	if err != nil {
		return nil, err
	}

	return response.NewOK(as), nil
}

func (s *ExamineAnswerService) GetExamineAnswer(ctx context.Context, id raid.Raid) (response.Response, error) {
	if _, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix); err != nil {
		return nil, err
	}

	as, err := s.ExamineAnswerRepository.GetExamineAnswer(ctx, id)
	if err != nil {
		return nil, err
	}

	return response.NewOK(as), nil
}

func (s *ExamineAnswerService) DeleteExamineAnswer(ctx context.Context, id raid.Raid) (response.Response, error) {
	if _, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix); err != nil {
		return nil, err
	}

	err := s.ExamineAnswerRepository.DeleteExamineAnswer(ctx, id)
	if err != nil {
		return nil, err
	}

	return response.NewNoContent(), nil
}
