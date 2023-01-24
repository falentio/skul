package domain

import (
	"context"
	"mime/multipart"
	"net/http"
)

type FileServiceRead interface {
	GetFile(ctx context.Context, slug string) ([]byte, string, error)
}

type FileServiceWrite interface {
	DeleteFile(ctx context.Context, slug string) error
	CreateFile(ctx context.Context, f http.File, h multipart.FileHeader) (string, error)
}

type FileService interface {
	FileServiceRead
	FileServiceWrite
}
