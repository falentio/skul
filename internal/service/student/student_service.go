package student

import (
	"context"
	"errors"

	"github.com/falentio/raid-go"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"

	"github.com/falentio/skul/internal/domain"
	"github.com/falentio/skul/internal/pkg/auth"
	"github.com/falentio/skul/internal/pkg/response"
	"github.com/falentio/skul/internal/pkg/validator"
	"github.com/falentio/skul/internal/pkg/xrand"
)

var StudentIDFactory = raid.NewRaid().WithPrefix("stu")
var _ domain.StudentService = new(StudentService)

type StudentService struct {
	StudentRepository domain.StudentRepository
	Auth              *auth.Auth
	Logger            zerolog.Logger
}

func (s *StudentService) CreateStudent(ctx context.Context, student *domain.Student) (res response.Response, err error) {
	adminID, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return
	}

	student.ID = StudentIDFactory.WithRandom().WithTimestampNow()
	student.AdminID = adminID
	student.Admin = nil
	student.EnteranceTokens = nil
	student.ExamineAnswer = nil

	err = validator.Struct(student)
	if err != nil {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(student.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	student.PasswordHash = string(hash)

	err = s.StudentRepository.CreateStudent(ctx, student)
	if err == domain.ErrStudentConflict {
		err = response.NewConflict(nil, "can not create student, data conflicting with other existing student")
	}
	if err != nil {
		return
	}

	res = response.NewCreated(student)
	return res, nil
}

func (s *StudentService) BatchCreateStudent(ctx context.Context, students []*domain.Student) (res response.Response, err error) {
	adminID, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix)
	if err != nil {
		return
	}
	for _, student := range students {
		student.ID = StudentIDFactory.WithRandom().WithTimestampNow()
		student.AdminID = adminID
		student.Admin = nil
		student.EnteranceTokens = nil
		student.ExamineAnswer = nil

		err = validator.Struct(student)
		if err != nil {
			return
		}
		if student.Password == "" {
			student.Password = xrand.Smol.GeneratePassword(10)
		}

		var hash []byte
		hash, err = bcrypt.GenerateFromPassword([]byte(student.Password), bcrypt.DefaultCost)
		if err != nil {
			return
		}
		student.PasswordHash = string(hash)
	}

	err = s.StudentRepository.BatchCreateStudent(ctx, students)
	if err == domain.ErrStudentConflict {
		err = response.NewConflict(nil, "can not create student, data conflicting with other existing student")
	}
	if err != nil {
		return
	}

	res = response.NewCreated(students)
	return res, nil
}

func (s *StudentService) GetStudent(ctx context.Context, studentID raid.Raid) (res response.Response, err error) {
	student, err := s.StudentRepository.GetStudent(ctx, studentID)
	if err != nil {
		return
	}

	res = response.NewOK(student)
	return res, nil
}

func (s *StudentService) LoginStudent(ctx context.Context, student *domain.Student) (res response.Response, err error) {
	stored, err := s.StudentRepository.GetStudentByUsername(ctx, student.Username)
	if errors.Is(err, domain.ErrStudentNotFound) {
		err = response.NewBadRequest(nil, "can not find student with username %q", student.Username)
	}
	if err != nil {
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(stored.PasswordHash), []byte(student.Password)); err != nil {
		err = response.NewBadRequest(nil, "password invalid")
		return
	}

	c, err := s.Auth.Sign(jwt.RegisteredClaims{
		Subject: stored.ID.String(),
	})

	r := response.NewNoContent()
	r.Cookies = append(r.Cookies, c)
	return r, nil
}

func (s *StudentService) ListStudent(ctx context.Context, opts *domain.ListStudentOptions) (res response.Response, err error) {
	students, err := s.StudentRepository.ListStudent(ctx, opts)
	if err != nil {
		return
	}

	res = response.NewPaginate(students, response.Page{
		Count:  opts.Count,
		Offset: opts.Offset,
	})
	return nil, nil
}

func (s *StudentService) DeleteStudent(ctx context.Context, studentID raid.Raid) (res response.Response, err error) {
	err = s.StudentRepository.DeleteStudent(ctx, studentID)
	if err != nil {
		return
	}

	res = response.NewNoContent()
	return res, nil
}

func (s *StudentService) UpdateStudent(ctx context.Context, student *domain.Student) (res response.Response, err error) {
	student.EnteranceTokens = nil
	student.ExamineAnswer = nil
	student.Admin = nil

	if err = validator.Struct(student); err != nil {
		return
	}

	if student.Password != "" {
		var hash []byte
		hash, err = bcrypt.GenerateFromPassword([]byte(student.Password), bcrypt.DefaultCost)
		if err != nil {
			return
		}
		student.Password = string(hash)
	}

	err = s.StudentRepository.UpdateStudent(ctx, student)
	if err != nil {
		return
	}

	res = response.NewOK(student)
	return res, nil
}
