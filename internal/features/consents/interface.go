package consents

import (
	"context"
	"io"

	"github.com/jmoiron/sqlx"
)

// ServiceInterface defines the business logic for legal document consents.
type ServiceInterface interface {
	GetLatestDocument(ctx context.Context, docType string) (*LegalDocument, error)
	GetLatestDocumentLocked(ctx context.Context, tx *sqlx.Tx, docType string) (*LegalDocument, error)
	GetDocumentContent(ctx context.Context, docType string) ([]byte, string, error)
	HasUserAccepted(ctx context.Context, userID string, docID int) (bool, error)
	SaveConsent(ctx context.Context, userID string, docID int) error
	ListUserConsentHistory(ctx context.Context, userID string) ([]UserConsent, error)
	UploadNewDocument(ctx context.Context, docType string, file io.ReadSeeker, contentType string) error
}

// RepositoryInterface defines the data access layer for legal document consents.
type RepositoryInterface interface {
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	GetLatestDocument(ctx context.Context, docType string) (*LegalDocument, error)
	GetLatestDocumentLocked(ctx context.Context, tx *sqlx.Tx, docType string) (*LegalDocument, error)
	HasUserAccepted(ctx context.Context, userID string, docID int) (bool, error)
	SaveConsent(ctx context.Context, userID string, docID int, ip string) error
	ListUserConsentHistory(ctx context.Context, userID string) ([]UserConsent, error)
	CreateNewVersion(ctx context.Context, tx *sqlx.Tx, doc LegalDocument) error
}
