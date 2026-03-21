package consents

import "time"

type LegalDocument struct {
	ID        int    `db:"id" json:"id"`
	DocType   string `db:"doc_type" json:"docType"`   // e.g., "privacy_policy", "terms_of_service"
	Version   string `db:"version" json:"version"`    // e.g., "v1.0", "2024-01-01"
	FileURL   string `db:"file_url" json:"fileUrl"`   // URL to the document file (PDF, HTML, etc.)
	IsActive  bool   `db:"is_active" json:"isActive"` // Indicates if this is the current active version
	CreatedAt string `db:"created_at" json:"createdAt"`
}

type UserConsent struct {
	ID         int       `db:"id" json:"id"`
	UserID     string    `db:"user_id" json:"userId"`
	DocumentID int       `db:"document_id" json:"documentId"`
	AcceptedAt time.Time `db:"accepted_at" json:"acceptedAt"`
	IPAddress  string    `db:"ip_address" json:"ipAddress"`
}
