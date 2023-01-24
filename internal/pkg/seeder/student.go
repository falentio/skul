package seeder

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/falentio/skul/internal/domain"
	"github.com/falentio/skul/internal/service/admin"
	"github.com/falentio/skul/internal/service/student"
)

func StudentRepository(s domain.StudentRepository) error {
	password, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	for _, std := range []*domain.Student{
		{Username: "foo"},
		{Username: "bar"},
	} {
		std.ID = student.StudentIDFactory.WithTimestampNow().WithRandom()
		std.AdminID = admin.AdminIDFactory.WithTimestampNow().WithRandom()
		std.Name = std.Username
		std.PasswordHash = string(password)

		if err := s.CreateStudent(context.Background(), std); err != nil {
			return err
		}
	}

	return nil
}
