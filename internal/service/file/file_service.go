package file

import (
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"

	"github.com/falentio/raid-go"
	"github.com/gofiber/storage"
	"github.com/rs/zerolog"

	"github.com/falentio/skul/internal/domain"
	"github.com/falentio/skul/internal/pkg/auth"
	"github.com/falentio/skul/internal/pkg/response"
)

type FileService struct {
	Storage storage.Storage
	Auth    *auth.Auth
	Logger  zerolog.Logger
}

func (s *FileService) CreateFile(ctx context.Context, f multipart.File, h *multipart.FileHeader) (response.Response, error) {
	if _, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix); err != nil {
		return nil, err
	}

	var slug = fmt.Sprintf(
		"%s.%s",
		raid.NewRaid().WithPrefix(domain.FileSlugPrefix),
		filepath.Ext(h.Filename),
	)

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if err := s.Storage.Set(slug, b, 0); err != nil {
		return nil, err
	}

	return response.NewOK(map[string]string{"slug": slug}), nil
}

func (s *FileService) GetFile(ctx context.Context, slug string) ([]byte, string, error) {
	if _, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix); err != nil {
		return nil, "", err
	}

	m := mime.TypeByExtension(filepath.Ext(slug))

	b, err := s.Storage.Get(slug)
	if err != nil {
		return nil, m, nil
	}

	return b, m, nil
}

func (s *FileService) DeleteFile(ctx context.Context, slug string) (response.Response, error) {
	if _, err := s.Auth.GetSubjectRaid(ctx, domain.AdminIDPrefix); err != nil {
		return nil, err
	}

	return response.NewNoContent(), s.Storage.Delete(slug)
}
