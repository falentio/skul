package seeder

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/falentio/skul/internal/domain"
	"github.com/falentio/skul/internal/service/admin"
)

func AdminRepository(r domain.AdminRepository) error {
	password, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	for _, a := range []*domain.Admin{
		{Username: "admin"},
		{Username: "foo"},
	} {
		a.ID = admin.AdminIDFactory.WithTimestampNow().WithRandom()
		a.Name = a.Username
		a.PasswordHash = string(password)

		if err := r.CreateAdmin(context.Background(), a); err != nil {
			return err
		}
	}
	return nil
}
