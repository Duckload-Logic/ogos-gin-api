package storage

import (
	"context"
	"io"
)

// FileStorage defines the interface for file upload/download operations.
// Implementations: azure.BlobStorage (prod), local.DiskStorage (dev).
type FileStorage interface {
	Upload(
		ctx context.Context,
		path string,
		reader io.ReadSeeker,
		contentType string,
	) error
	Download(ctx context.Context, path string, writer io.Writer) error
	Delete(ctx context.Context, path string) error
}
