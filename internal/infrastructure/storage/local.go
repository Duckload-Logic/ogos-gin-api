package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// DiskStorage implements FileStorage using the local filesystem.
type DiskStorage struct {
	baseDir string
}

func NewDiskStorage(baseDir string) *DiskStorage {
	return &DiskStorage{baseDir: baseDir}
}

func (d *DiskStorage) Upload(
	ctx context.Context,
	path string,
	reader io.ReadSeeker,
	contentType string,
) error {
	fullPath := filepath.Join(d.baseDir, filepath.FromSlash(path))

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (d *DiskStorage) Download(
	ctx context.Context,
	path string,
	writer io.Writer,
) error {
	fullPath := filepath.Join(d.baseDir, filepath.FromSlash(path))

	f, err := os.Open(fullPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(writer, f); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return nil
}
