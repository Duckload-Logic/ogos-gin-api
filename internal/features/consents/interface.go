package consents

import (
	"context"
	"io"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

// ServiceInterface defines the business logic for legal document consents.
type ServiceInterface interface {
	GetLatestDocument(
		ctx context.Context,
		docType string,
	) (*LegalDocument, error)
	GetLatestDocumentLocked(
		ctx context.Context,
		tx datastore.DB,
		docType string,
	) (*LegalDocument, error)
	GetDocumentContent(
		ctx context.Context,
		docType string,
	) ([]byte, string, error)
	HasUserAccepted(ctx context.Context, userID string, docID int) (bool, error)
	SaveConsent(ctx context.Context, userID string, docID int) error
	ListUserConsentHistory(
		ctx context.Context,
		userID string,
	) ([]UserConsent, error)
	UploadNewDocument(
		ctx context.Context,
		docType string,
		file io.ReadSeeker,
		contentType string,
	) error
}

// RepositoryInterface defines the data access layer for legal document
// consents.
type RepositoryInterface interface {
	GetDB() *sqlx.DB
	BeginTx(ctx context.Context) (datastore.DB, error)
	GetLatestDocument(
		ctx context.Context,
		docType string,
	) (*LegalDocument, error)
	GetLatestDocumentLocked(
		ctx context.Context,
		tx datastore.DB,
		docType string,
	) (*LegalDocument, error)
	HasUserAccepted(ctx context.Context, userID string, docID int) (bool, error)
	SaveConsent(ctx context.Context, userID string, docID int, ip string) error
	ListUserConsentHistory(
		ctx context.Context,
		userID string,
	) ([]UserConsent, error)
	CreateNewVersion(
		ctx context.Context,
		tx datastore.DB,
		doc LegalDocument,
	) error
}
