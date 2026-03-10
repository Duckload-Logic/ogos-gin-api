package consents

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// GetLatestDocument fetches the current active version of a policy from Azure/DB
func (r *Repository) GetLatestDocument(ctx context.Context, docType string) (*LegalDocument, error) {
	var doc LegalDocument
	query := fmt.Sprintf(
		`SELECT %s
		FROM legal_documents
        WHERE doc_type = ? AND
		is_active = TRUE
		LIMIT 1`,
		database.GetColumns(LegalDocument{}),
	)

	log.Printf("Query: %s, DocType: %s", query, docType)

	err := r.db.GetContext(ctx, &doc, query, docType)
	return &doc, err
}

// HasUserAccepted checks if the specific user has already signed this specific document ID
func (r *Repository) HasUserAccepted(ctx context.Context, email string, docID int) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM user_consents
			WHERE user_email = ? AND
			document_id = ?
		)`

	err := r.db.GetContext(ctx, &exists, query, email, docID)
	return exists, err
}

// SaveConsent records the "State" of the user's agreement
func (r *Repository) SaveConsent(ctx context.Context, email string, docID int, ip string) error {
	query := `
		INSERT INTO user_consents (user_email, document_id, ip_address)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE accepted_at = CURRENT_TIMESTAMP`

	_, err := r.db.ExecContext(ctx, query, email, docID, ip)
	return err
}

// ListUserConsentHistory for an admin "Compliance Dashboard"
func (r *Repository) ListUserConsentHistory(ctx context.Context, email string) ([]UserConsent, error) {
	var history []UserConsent
	query := fmt.Sprintf(`
		SELECT %s FROM user_consents uc
		JOIN legal_documents ld ON uc.document_id = ld.id
		WHERE uc.user_email = ? ORDER BY uc.accepted_at DESC
	`, database.GetColumns(UserConsent{}))

	err := r.db.SelectContext(ctx, &history, query, email)
	return history, err
}

// CreateNewVersion inserts a new document and optionally deactivates the old one
func (r *Repository) CreateNewVersion(ctx context.Context, doc LegalDocument) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// 1. Deactivate current active version of this type
	_, err = tx.ExecContext(
		ctx,
		"UPDATE legal_documents SET is_active = FALSE WHERE doc_type = ?",
		doc.DocType,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 2. Insert the new version as the active one
	query := `INSERT INTO legal_documents (doc_type, version, file_url, is_active)
              VALUES (?, ?, ?, TRUE)`
	_, err = tx.ExecContext(ctx, query, doc.DocType, doc.Version, doc.FileURL)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
