package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
)

type BlobStorage struct {
	client        *azblob.Client
	containerName string
}

func NewBlobStorage(connectionString, containerName string) (*BlobStorage, error) {
	client, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure Blob client: %w", err)
	}

	return &BlobStorage{
		client:        client,
		containerName: containerName,
	}, nil
}

// Upload uploads data from a reader to the given blob path.
func (b *BlobStorage) Upload(ctx context.Context, blobPath string, reader io.ReadSeeker, contentType string) error {
	_, err := b.client.UploadStream(ctx, b.containerName, blobPath, reader, &azblob.UploadStreamOptions{
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType: &contentType,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to upload blob %q: %w", blobPath, err)
	}
	return nil
}

// Download streams the blob content into the provided writer.
func (b *BlobStorage) Download(ctx context.Context, blobPath string, writer io.Writer) error {
	resp, err := b.client.DownloadStream(ctx, b.containerName, blobPath, nil)
	if err != nil {
		return fmt.Errorf("failed to download blob %q: %w", blobPath, err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read blob %q: %w", blobPath, err)
	}
	return nil
}

// GetContentType returns the content type of the blob.
func (b *BlobStorage) GetContentType(ctx context.Context, blobPath string) (string, error) {
	resp, err := b.client.DownloadStream(ctx, b.containerName, blobPath, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get blob properties %q: %w", blobPath, err)
	}
	defer resp.Body.Close()

	if resp.ContentType != nil {
		return *resp.ContentType, nil
	}
	return "application/octet-stream", nil
}
