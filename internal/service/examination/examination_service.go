package examination

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

var _ domain.ExaminationService = new(ExaminationService)
var ExaminationIDFactory = raid.NewRaid().WithPrefix(domain.ExaminationIDPrefix)

type ExaminationService struct {
	ExaminationRepository domain.ExaminationRepository
	Auth                  *auth.Auth
	Logger                zerolog.Logger
}

func (s *ExaminationService) GetExamination(ctx context.Context, examinationID raid.Raid) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	ex, err := s.ExaminationRepository.GetExamination(ctx, examinationID)
	if errors.Is(err, domain.ErrExaminationNotFound) {
		err = response.NewNotFound(nil, "can not find examination with id %q", examinationID)
	}
	if err != nil {
		return nil, err
	}

	return response.NewOK(ex), nil
}

func (s *ExaminationService) ListExamination(ctx context.Context, o *domain.ListExaminationOptions) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	exs, err := s.ExaminationRepository.ListExamination(ctx, o)
	if err != nil {
		return nil, err
	}

	return response.NewPaginate(exs, response.Page{
		Count:  o.Count,
		Offset: o.Offset,
	}), nil
}

func (s *ExaminationService) CreateExamination(ctx context.Context, examination *domain.Examination) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	examination.ID = ExaminationIDFactory.WithTimestampNow().WithRandom()

	if err := validator.Struct(examination); err != nil {
		return nil, err
	}

	if err := s.ExaminationRepository.CreateExamination(ctx, examination); err != nil {
		return nil, err
	}

	return response.NewOK(examination), nil
}

func (s *ExaminationService) DeleteExamination(ctx context.Context, examinationID raid.Raid) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	if err := s.ExaminationRepository.DeleteExamination(ctx, examinationID); err != nil {
		return nil, err
	}

	return response.NewNoContent(), nil
}

func (s *ExaminationService) UpdateExamination(ctx context.Context, examination *domain.Examination) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	if err := validator.Struct(examination); err != nil {
		return nil, err
	}

	if err := s.ExaminationRepository.UpdateExamination(ctx, examination); err != nil {
		return nil, err
	}

	return nil, nil
}
