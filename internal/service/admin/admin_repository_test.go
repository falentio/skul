package admin

import (
	"context"
	"testing"

	"github.com/falentio/raid-go"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/falentio/skul/internal/domain"
)

func TestAdminRepository(t *testing.T) {
	t.Parallel()
	db, err := gorm.Open(sqlite.Open(":memory:?cache=shared"))
	if err != nil {
		t.Fatal(err.Error())
	}
	if err := db.AutoMigrate(&domain.Admin{}, &domain.Examination{}, &domain.Student{}); err != nil {
		t.Fatal(err.Error())
	}

	for _, tc := range []struct {
		Name string
		Repo domain.AdminRepository
	}{
		{"gorm", &AdminRepositoryGorm{db}},
	} {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			admin1 := &domain.Admin{Model: domain.Model{ID: raid.NewRaid()}, Username: "foo"}
			err := tc.Repo.CreateAdmin(context.Background(), admin1)
			if err != nil {
				t.Error("failed to create admin", err)
			}

			admin2 := &domain.Admin{Model: domain.Model{ID: raid.NewRaid()}, Username: "foo"}
			err = tc.Repo.CreateAdmin(context.Background(), admin2)
			if err == nil || err != domain.ErrAdminConflict {
				t.Errorf("%#+v", err)
				t.Error("CreateAdmin must return domain.ErrAdminConflict while creating with existing data")
			}
			admin2.Username = "bar"
			err = tc.Repo.CreateAdmin(context.Background(), admin2)
			if err != nil {
				t.Error("failed to create admin", err)
			}

			storedAdmin, err := tc.Repo.GetAdminByID(context.Background(), admin1.ID)
			if err != nil {
				t.Error("failed to get admin", err)
			}
			if storedAdmin.Username != admin1.Username {
				t.Errorf("invalid admin data returned, username changed from %s to %s", admin1.Username, storedAdmin.Username)
			}

			storedAdmin, err = tc.Repo.GetAdminByUsername(context.Background(), admin1.Username)
			if err != nil {
				t.Error("failed to get admin", err)
			}
			if storedAdmin.Username != admin1.Username {
				t.Errorf("invalid admin data returned, username changed from %s to %s", admin1.Username, storedAdmin.Username)
			}

			err = tc.Repo.UpdateAdmin(context.Background(), &domain.Admin{Model: domain.Model{ID: admin1.ID}, Username: "baz"})
			if err != nil {
				t.Error("failed to update admin", err)
			}

			storedAdmin, err = tc.Repo.GetAdminByID(context.Background(), admin1.ID)
			if err != nil {
				t.Error("failed to get admin", err)
			}
			if storedAdmin.Username != "baz" {
				t.Errorf("invalid admin data returned, username changed from %s to %s", admin1.Username, storedAdmin.Username)
			}

			err = tc.Repo.DeleteAdmin(context.Background(), admin1.ID)
			if err != nil {
				t.Error("failed to delete admin", err)
			}

			storedAdmin, err = tc.Repo.GetAdminByID(context.Background(), admin1.ID)
			if err == nil {
				t.Error("failed to delete admin")
			}
			if err != domain.ErrAdminNotFound {
				t.Error("GetAdminByID must return domain.ErrAdminNotFound while gettting unexists admin")
			}
			if storedAdmin != nil {
				t.Error("GetAdminByID must return nil admin if error")
			}

			storedAdmin, err = tc.Repo.GetAdminByID(context.Background(), raid.NilRaid)
			if err == nil {
				t.Error("error not returned while getting unexists admin ID")
			}
			if err != domain.ErrAdminNotFound {
				t.Error("GetAdminByID must return domain.ErrAdminNotFound while gettting unexists admin")
			}
			if storedAdmin != nil {
				t.Error("GetAdminByID must return nil admin if error")
			}
		})
	}
}
