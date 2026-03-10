package consents

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/storage"
	"github.com/olazo-johnalbert/duckload-api/internal/features/logs"
)

type Service struct {
	repo       *Repository
	logService *logs.Service
	storage    storage.FileStorage
}

func NewService(repo *Repository, logService *logs.Service, fileStorage storage.FileStorage) *Service {
	return &Service{repo: repo, logService: logService, storage: fileStorage}
}

func (s *Service) GetLatestDocument(ctx context.Context, docType string) (*LegalDocument, error) {
	dbDocType := ""
	switch docType {
	case "terms":
		dbDocType = "TERMS_OF_SERVICE"
	case "privacy":
		dbDocType = "PRIVACY_POLICY"
	}

	return s.repo.GetLatestDocument(ctx, dbDocType)
}

func (s *Service) HasUserAccepted(ctx context.Context, userID int, docID int) (bool, error) {
	return s.repo.HasUserAccepted(ctx, userID, docID)
}

func (s *Service) SaveConsent(ctx context.Context, userID int, docID int) error {
	_, ipAddress, userAgent, userEmail := audit.ExtractMeta(ctx)

	err := s.repo.SaveConsent(ctx, userID, docID, ipAddress)
	if err == nil {
		s.logService.Record(ctx, logs.LogEntry{
			Category:  logs.CategoryConsent,
			Action:    logs.ActionTermsAccepted,
			Message:   fmt.Sprintf("%s accepted legal document", userEmail),
			UserID:    userID,
			UserEmail: userEmail,
			IPAddress: ipAddress,
			UserAgent: userAgent,
		})
	}

	return err
}

// ListUserConsentHistory (Optional) for an admin "Compliance Dashboard"
func (s *Service) ListUserConsentHistory(ctx context.Context, userID int) ([]UserConsent, error) {
	return s.repo.ListUserConsentHistory(ctx, userID)
}

func (s *Service) UploadNewDocument(ctx context.Context, docType string, file io.ReadSeeker, contentType string) error {
	today := time.Now().Format("2006-01-02")
	version := strings.Join(strings.Split(today, "-"), ".")

	filePath := fmt.Sprintf("legal/%s/%s_%d.md", docType, version, time.Now().Unix())

	if err := s.storage.Upload(ctx, filePath, file, contentType); err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	doc := LegalDocument{
		DocType: docType,
		Version: version,
		FileURL: filePath,
	}

	return s.repo.CreateNewVersion(ctx, doc)
}
