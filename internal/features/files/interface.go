package files

import (
	"context"
	"mime/multipart"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type ServiceInterface interface {
	UploadFile(
		ctx context.Context,
		file *multipart.FileHeader,
		prefix string,
	) (File, error)
	UploadFiles(
		ctx context.Context,
		files []*multipart.FileHeader,
		prefix string,
	) ([]File, error)
	GetFileByID(ctx context.Context, id string) (*File, error)
	DeleteFile(ctx context.Context, id string) error
	GetOCRResult(ctx context.Context, fileID string) (*OCRResult, error)
}

type RepositoryInterface interface {
	WithTransaction(ctx context.Context, fn func(datastore.DB) error) error
	GetDB() *sqlx.DB
	GetFileByID(ctx context.Context, id string) (*File, error)
	GetFilesByIDs(ctx context.Context, ids []string) ([]File, error)
	Create(ctx context.Context, tx datastore.DB, file File) (string, error)
	CreateBulk(
		ctx context.Context,
		tx datastore.DB,
		files []File,
	) ([]string, error)
	Delete(ctx context.Context, tx datastore.DB, id string) error
	SaveOCRResult(ctx context.Context, tx datastore.DB, result OCRResult) error
	GetOCRResultByFileID(ctx context.Context, fileID string) (*OCRResult, error)
}
