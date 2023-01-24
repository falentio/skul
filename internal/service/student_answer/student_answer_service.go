package studentanswer

import (
	"context"
	"errors"

	"github.com/falentio/raid-go"

	"github.com/falentio/skul/internal/domain"
	"github.com/falentio/skul/internal/pkg/auth"
	"github.com/falentio/skul/internal/pkg/response"
)

type StudentAnswerService struct {
	StudentAnswerRepository domain.StudentAnswerRepository
	Auth                    *auth.Auth
}

func (s *StudentAnswerService) ListStudentAnswer(ctx context.Context, o *domain.ListStudentAnswerOptions) (response.Response, error) {
	id, err := s.Auth.GetSubjectRaid(ctx, "")
	if err != nil {
		return nil, err
	}

	if id.Prefix() != domain.AdminIDPrefix {
		o.StudentID = id
	}

	a, err := s.StudentAnswerRepository.ListStudentAnswer(ctx, o)
	if err != nil {
		return nil, err
	}

	return response.NewOK(a), nil
}

func (s *StudentAnswerService) CreateStudentAnswer(ctx context.Context, a *domain.StudentAnswer) (response.Response, error) {
	id, err := s.Auth.GetSubjectRaid(ctx, "")
	if err != nil {
		return nil, err
	}

	if id.Prefix() != domain.AdminIDPrefix {
		a.StudentID = id
	}

	a.ID = raid.NewRaid().WithPrefix(domain.StudentAnswerIDPrefix)
	if err := s.StudentAnswerRepository.CreateStudentAnswer(ctx, a); err != nil {
		return nil, err
	}

	return response.NewOK(a), nil
}

func (s *StudentAnswerService) DeleteStudentAnswer(ctx context.Context, answerID raid.Raid) (response.Response, error) {
	id, err := s.Auth.GetSubjectRaid(ctx, "")
	if err != nil {
		return nil, err
	}

	a, err := s.StudentAnswerRepository.GetStudentAnswer(ctx, answerID)
	if errors.Is(err, domain.ErrStudentAnswerNotFound) {
		return response.NewNoContent(), nil
	}
	if err != nil {
		return nil, err
	}

	if a.StudentID.String() != id.String() {
		return nil, response.NewForbidden(nil, "can not delete others student answer")
	}

	if err := s.StudentAnswerRepository.DeleteStudentAnswer(ctx, answerID); err != nil {
		return nil, err
	}

	return response.NewNoContent(), nil
}
