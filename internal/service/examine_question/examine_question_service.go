package examinequestion

import (
	"context"
	"errors"

	"github.com/falentio/raid-go"
	"github.com/rs/zerolog"

	"github.com/falentio/skul/internal/domain"
	"github.com/falentio/skul/internal/pkg/auth"
	"github.com/falentio/skul/internal/pkg/response"
	"github.com/falentio/skul/internal/pkg/validator"
)

var ExamineQuestionIDFactory = raid.NewRaid().WithPrefix(domain.ExamineQuestionIDPrefix)

type ExamineQuestionService struct {
	ExamineQuestionRepository    domain.ExamineQuestionRepository
	ExamineAnswerRepository      domain.ExamineAnswerRepositoryWrite
	ExamineAttatchmentRepository domain.ExamineAttatchmentRepositoryWrite
	Auth                         *auth.Auth
	Logger                       zerolog.Logger
}

func (s *ExamineQuestionService) CreateExamineQuestion(ctx context.Context, q *domain.ExamineQuestion) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	if err := validator.Struct(q); err != nil {
		return nil, err
	}

	q.ID = ExamineQuestionIDFactory.WithRandom().WithTimestampNow()
	if err := s.ExamineQuestionRepository.CreateExamineQuestion(ctx, q); err != nil {
		return nil, err
	}

	q.ExamineAttatchment.ID = raid.NewRaid().WithPrefix(domain.ExamineAttatchmentIDPrefix)
	if err := s.ExamineAttatchmentRepository.CreateExamineAttathcment(ctx, q.ExamineAttatchment); err != nil {
		return nil, err
	}

	for _, a := range q.ExamineAnswers {
		a.ID = raid.NewRaid().WithPrefix(domain.ExamineAnswerIDPrefix)
	}
	if err := s.ExamineAnswerRepository.BatchCreateExamineAnswer(ctx, q.ExamineAnswers); err != nil {
		return nil, err
	}

	return response.NewOK(q), nil
}

func (s *ExamineQuestionService) UpdateExamineQuestion(ctx context.Context, q *domain.ExamineQuestion) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	if err := validator.Struct(q); err != nil {
		return nil, err
	}

	if err := s.ExamineQuestionRepository.UpdateExamineQuestion(ctx, q); err != nil {
		return nil, err
	}

	return response.NewOK(q), nil
}

func (s *ExamineQuestionService) GetExamineQuestion(ctx context.Context, id raid.Raid) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	q, err := s.ExamineQuestionRepository.GetExamineQuestion(ctx, id)
	if errors.Is(err, domain.ErrExamineQuestionNotFound) {
		err = response.NewBadRequest(nil, "can not find examine question with id %q", id)
	}
	if err != nil {
		return nil, err
	}

	return response.NewOK(q), nil
}

func (s *ExamineQuestionService) ListExamineQuestion(ctx context.Context, o *domain.ListExamineQuestionOptions) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	qs, err := s.ExamineQuestionRepository.ListExamineQuestion(ctx, o)
	if err != nil {
		return nil, err
	}

	return response.NewPaginate(qs, response.Page{
		Count:  o.Count,
		Offset: o.Offset,
		Page:   o.Page,
	}), nil
}

func (s *ExamineQuestionService) DeleteExamineQuestion(ctx context.Context, id raid.Raid) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	if err := s.ExamineQuestionRepository.DeleteExamineQuestion(ctx, id); err != nil {
		return nil, err
	}

	return response.NewNoContent(), nil
}
