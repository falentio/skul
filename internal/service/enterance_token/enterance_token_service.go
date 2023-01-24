package enterancetoken

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/falentio/raid-go"

	"github.com/falentio/skul/internal/domain"
	"github.com/falentio/skul/internal/pkg/auth"
	"github.com/falentio/skul/internal/pkg/response"
	"github.com/falentio/skul/internal/pkg/validator"
	"github.com/falentio/skul/internal/pkg/xrand"
)

var EnteranceTokenIDFactory = raid.NewRaid().WithPrefix(domain.EnteranceTokenIDPrefix)

type EnteranceTokenService struct {
	EnteranceTokenRepository domain.EnteranceTokenRepository
	ExaminationRepository    domain.ExaminationRepositoryRead
	ExamineStudentRepository domain.ExamineStudentRepositoryRead
	Auth                     *auth.Auth
}

func (s *EnteranceTokenService) GetEnteranceToken(ctx context.Context, tokenID raid.Raid) (response.Response, error) {
	id, err := s.Auth.GetSubjectRaid(ctx, "")
	if err != nil {
		return nil, err
	}

	if id.Prefix() != domain.AdminIDPrefix {
		stds, err := s.ExamineStudentRepository.ListExamineStudent(ctx, &domain.ListExamineStudentOptions{
			StudentID:        id,
			EnteranceTokenID: tokenID,
		})
		if err != nil {
			return nil, err
		}
		if len(stds) == 0 {
			return nil, response.NewForbidden(nil, "student with id %q does has permission to use enterance token with id %q", id, tokenID)
		}
	}

	token, err := s.EnteranceTokenRepository.GetEnteranceToken(ctx, tokenID)
	if errors.Is(err, domain.ErrEnteranceTokenNotFound) {
		return nil, response.NewNotFound(nil, "can not find enterance token with id %q", tokenID.String())
	}
	if err != nil {
		return nil, err
	}
	if id.Prefix() != domain.AdminIDPrefix {
		token.Students = nil
	}

	return response.NewOK(token), nil
}

func (s *EnteranceTokenService) GetExamination(ctx context.Context, tokenID raid.Raid) (response.Response, error) {
	id, err := s.Auth.GetSubjectRaid(ctx, "")
	if err != nil {
		return nil, err
	}

	if id.Prefix() != domain.AdminIDPrefix {
		stds, err := s.ExamineStudentRepository.ListExamineStudent(ctx, &domain.ListExamineStudentOptions{
			StudentID:        id,
			EnteranceTokenID: tokenID,
		})
		if err != nil {
			return nil, err
		}
		if len(stds) == 0 {
			return nil, response.NewForbidden(nil, "student with id %q does has permission to use enterance token with id %q", id, tokenID)
		}
	}

	token, err := s.EnteranceTokenRepository.GetEnteranceToken(ctx, tokenID)
	if errors.Is(err, domain.ErrEnteranceTokenNotFound) {
		return nil, response.NewNotFound(nil, "can not find enterance token with id %q", tokenID.String())
	}
	if err != nil {
		return nil, err
	}

	exa, err := s.ExaminationRepository.GetExamination(ctx, token.ExaminationID)
	if err != nil {
		return nil, err
	}
	if len(exa.ExamineQuestions) < int(exa.QuestionCount) {
		return nil, response.NewBadRequest(nil, "invalid examination, question count is less then desired ammount")
	}

	if id.Prefix() != domain.AdminIDPrefix {
		rng := rand.New(rand.NewSource(0))
		seed := fmt.Sprintf("%s%s%s", tokenID, exa.ID, id)
		xrand.Seed(rng, seed)

		rng.Shuffle(len(exa.ExamineQuestions), func(i, j int) {
			exa.ExamineQuestions[i], exa.ExamineQuestions[j] = exa.ExamineQuestions[j], exa.ExamineQuestions[i]
		})
		exa.ExamineQuestions = exa.ExamineQuestions[:exa.QuestionCount]

		for _, q := range exa.ExamineQuestions {
			rng.Shuffle(len(q.ExamineAnswers), func(i, j int) {
				q.ExamineAnswers[i], q.ExamineAnswers[j] = q.ExamineAnswers[j], q.ExamineAnswers[i]
			})

			answers := make([]*domain.ExamineAnswer, 0, q.AnswerCount)
			for _, a := range q.ExamineAnswers {
				if a.Correct {
					answers = append(answers, a)
					break
				}
			}

			if len(answers) == 0 {
				return nil, response.NewBadRequest(nil, "invalid examination, question with id %q does not has correct answer", q.ID)
			}

			for i, a := range q.ExamineAnswers {
				if len(answers) >= q.AnswerCount {
					break
				}
				if !a.Correct {
					answers = append(answers, q.ExamineAnswers[i])
					break
				}
			}

			if len(answers) < int(q.AnswerCount) {
				return nil, response.NewBadRequest(nil, "invalid examination, answers count is less then desired ammount")
			}

			q.ExamineAnswers = answers

			rng.Shuffle(len(q.ExamineAnswers), func(i, j int) {
				q.ExamineAnswers[i], q.ExamineAnswers[j] = q.ExamineAnswers[j], q.ExamineAnswers[i]
			})
		}
	}

	return response.NewOK(exa), nil
}

func (s *EnteranceTokenService) ListEnteranceToken(ctx context.Context, o *domain.ListEnteranceTokenOptions) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	ts, err := s.EnteranceTokenRepository.ListEnteranceToken(ctx, o)
	if err != nil {
		return nil, err
	}

	return response.NewOK(ts), nil
}

func (s *EnteranceTokenService) CreateEnteranceToken(ctx context.Context, token *domain.EnteranceToken) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	if err := validator.Struct(token); err != nil {
		return nil, err
	}

	if err := s.EnteranceTokenRepository.CreateEnteranceToken(ctx, token); err != nil {
		return nil, err
	}

	return response.NewOK(token), nil
}

func (s *EnteranceTokenService) UpdateEnteranceToken(ctx context.Context, token *domain.EnteranceToken) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	if err := validator.Struct(token); err != nil {
		return nil, err
	}

	if err := s.EnteranceTokenRepository.UpdateEnteranceToken(ctx, token); err != nil {
		return nil, err
	}

	return response.NewOK(token), nil
}

func (s *EnteranceTokenService) BatchCreateEnteranceToken(ctx context.Context, tokens []*domain.EnteranceToken) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	for _, token := range tokens {
		if err := validator.Struct(token); err != nil {
			return nil, err
		}
	}

	if err := s.EnteranceTokenRepository.BatchCreateEnteranceToken(ctx, tokens); err != nil {
		return nil, err
	}

	return response.NewOK(tokens), nil
}

func (s *EnteranceTokenService) DeleteEnteranceToken(ctx context.Context, tokenID raid.Raid) (response.Response, error) {
	_, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return nil, err
	}

	if err := s.EnteranceTokenRepository.DeleteEnteranceToken(ctx, tokenID); err != nil {
		return nil, err
	}

	return nil, nil
}
