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

	err := r.db.GetContext(ctx, &doc, query, docType)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil // No active document found, return nil without error
		}
		return nil, err
	}

	return &doc, nil
}

func (r *Repository) GetLatestDocumentLocked(ctx context.Context, tx *sqlx.Tx, docType string) (*LegalDocument, error) {
	var doc LegalDocument
	query := fmt.Sprintf(
		`SELECT %s
		FROM legal_documents
		WHERE doc_type = ? AND
		is_active = TRUE
		LIMIT 1
		FOR UPDATE`, // Lock the selected row for update
		database.GetColumns(LegalDocument{}),
	)

	err := tx.GetContext(ctx, &doc, query, docType)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil // No active document found, return nil without error
		}
		return nil, err
	}

	return &doc, nil
}

// HasUserAccepted checks if the specific user has already signed this specific document ID
func (r *Repository) HasUserAccepted(ctx context.Context, userID string, docID int) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM user_consents
			WHERE user_id = ? AND
			document_id = ?
		)`

	err := r.db.GetContext(ctx, &exists, query, userID, docID)
	return exists, err
}

// SaveConsent records the "State" of the user's agreement
func (r *Repository) SaveConsent(ctx context.Context, userID string, docID int, ip string) error {
	cols, vals := database.GetInsertStatement(UserConsent{}, []string{"created_at", "email", "accepted_at"})
	query := fmt.Sprintf(`
		INSERT INTO user_consents (%s)
		VALUES (%s)`, cols, vals)

	log.Printf("Query: %s, Vals: %s %d %s", query, userID, docID, ip)

	_, err := r.db.NamedExecContext(ctx, query, UserConsent{
		UserID:     userID,
		DocumentID: docID,
		IPAddress:  ip,
	})
	if err != nil {
		return err
	}

	return nil
}

// ListUserConsentHistory for an admin "Compliance Dashboard"
func (r *Repository) ListUserConsentHistory(ctx context.Context, userID string) ([]UserConsent, error) {
	var history []UserConsent
	query := fmt.Sprintf(`
		SELECT %s FROM user_consents uc
		JOIN legal_documents ld ON uc.document_id = ld.id
		WHERE uc.user_id = ? ORDER BY uc.accepted_at DESC
	`, database.GetColumns(UserConsent{}))

	err := r.db.SelectContext(ctx, &history, query, userID)
	return history, err
}

func (r *Repository) CreateNewVersion(ctx context.Context, tx *sqlx.Tx, doc LegalDocument) error {
	// Deactivate current active version
	_, err := tx.ExecContext(ctx, "UPDATE legal_documents SET is_active = FALSE WHERE doc_type = ?", doc.DocType)
	if err != nil {
		return err
	}
	// Insert the new version
	query := `INSERT INTO legal_documents (doc_type, version, file_url, is_active) VALUES (?, ?, ?, TRUE)`
	_, err = tx.ExecContext(ctx, query, doc.DocType, doc.Version, doc.FileURL)
	if err != nil {
		return err
	}
	return nil
}
