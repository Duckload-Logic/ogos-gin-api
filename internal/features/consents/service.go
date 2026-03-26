package consents

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/features/logs"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/storage"
)

const (
	DocPathFormat = "legal/%s/%s.md"
	VersionPrefix = "_v"
	FirstVersion  = 1
)

type Service struct {
	repo       RepositoryInterface
	logService *logs.Service
	storage    storage.FileStorage
}

func NewService(
	repo RepositoryInterface,
	logService *logs.Service,
	fileStorage storage.FileStorage,
) *Service {
	return &Service{repo: repo, logService: logService, storage: fileStorage}
}

func (s *Service) GetLatestDocument(
	ctx context.Context,
	docType string,
) (*LegalDocument, error) {
	dbDocType := ""
	switch docType {
	case "terms":
		dbDocType = "TERMS_OF_SERVICE"
	case "privacy":
		dbDocType = "PRIVACY_POLICY"
	}

	return s.repo.GetLatestDocument(ctx, dbDocType)
}

func (s *Service) GetLatestDocumentLocked(
	ctx context.Context,
	tx datastore.DB,
	docType string,
) (*LegalDocument, error) {
	dbDocType := ""
	switch docType {
	case "terms":
		dbDocType = "TERMS_OF_SERVICE"
	case "privacy":
		dbDocType = "PRIVACY_POLICY"
	}

	return s.repo.GetLatestDocumentLocked(ctx, tx, dbDocType)
}

// In consents/service.go
func (s *Service) GetDocumentContent(
	ctx context.Context,
	docType string,
) ([]byte, string, error) {
	// 1. Fetch latest document metadata
	doc, err := s.GetLatestDocument(ctx, docType)
	if err != nil {
		return nil, "", err
	}

	// 2. Download the blob content using your storage implementation
	//    We'll use a bytes.Buffer to capture the output.
	var buf bytes.Buffer
	err = s.storage.Download(ctx, doc.FileURL, &buf)
	if err != nil {
		return nil, "", err
	}

	// 3. Return the content and content type (you stored it as "text/markdown")
	return buf.Bytes(), "text/markdown; charset=utf-8", nil
}

func (s *Service) HasUserAccepted(
	ctx context.Context,
	userID string,
	docID int,
) (bool, error) {
	return s.repo.HasUserAccepted(ctx, userID, docID)
}

func (s *Service) SaveConsent(
	ctx context.Context,
	userID string,
	docID int,
) error {
	_, ipAddress, userAgent, userEmail := audit.ExtractMeta(ctx)

	err := s.repo.SaveConsent(ctx, userID, docID, ipAddress)
	if err == nil {
		s.logService.Record(ctx, s.repo.GetDB(), logs.LogEntry{
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
func (s *Service) ListUserConsentHistory(
	ctx context.Context,
	userID string,
) ([]UserConsent, error) {
	return s.repo.ListUserConsentHistory(ctx, userID)
}

func (s *Service) UploadNewDocument(
	ctx context.Context,
	docType string,
	file io.ReadSeeker,
	contentType string,
) (err error) { // Use named return for the defer check
	now := time.Now()
	today := strings.ReplaceAll(now.Format("2006.01.02"), "-", ".")

	return datastore.RunInTransaction(
		ctx,
		s.repo.GetDB(),
		func(tx datastore.DB) error {
			lastDoc, err := s.repo.GetLatestDocumentLocked(ctx, tx, docType)
			if err != nil {
				log.Printf(
					"[PostUploadNewDocument] Database Query Locked: %v",
					err,
				)
				return fmt.Errorf("failed to get latest document: %w", err)
			}

			versionCtr := FirstVersion
			if lastDoc != nil && strings.HasPrefix(lastDoc.Version, today) {
				parts := strings.Split(lastDoc.Version, VersionPrefix)
				if len(parts) == 2 {
					if val, e := strconv.Atoi(parts[1]); e == nil {
						versionCtr = val + 1
					}
				}
			}

			version := fmt.Sprintf("%s%s%d", today, VersionPrefix, versionCtr)
			filePath := fmt.Sprintf(DocPathFormat, docType, version)

			if err = s.storage.Upload(ctx, filePath, file, contentType); err != nil {
				log.Printf("[PostUploadNewDocument] Storage Upload: %v", err)
				return fmt.Errorf("failed to upload file: %w", err)
			}

			doc := LegalDocument{
				DocType: docType,
				Version: version,
				FileURL: filePath,
			}

			if err = s.repo.CreateNewVersion(ctx, tx, doc); err != nil {
				log.Printf("[PostUploadNewDocument] Database Insert: %v", err)
				return fmt.Errorf("failed to save document record: %w", err)
			}

			return nil
		},
	)
}
