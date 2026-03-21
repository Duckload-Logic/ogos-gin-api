package slips

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/hash"
	"github.com/olazo-johnalbert/duckload-api/internal/core/storage"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/logs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

const MaxFileSize = 5 * 1024 * 1024 // 5MB limit

type Service struct {
	repo       *Repository
	logService *logs.Service
	storage    storage.FileStorage
}

func NewService(repo *Repository, logService *logs.Service, storage storage.FileStorage) *Service {
	return &Service{repo: repo, logService: logService, storage: storage}
}

func (s *Service) GetSlipStatuses(ctx context.Context) ([]SlipStatus, error) {
	statuses, err := s.repo.GetSlipStatuses(ctx)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

func (s *Service) GetSlipCategories(ctx context.Context) ([]SlipCategory, error) {
	categories, err := s.repo.GetSlipCategories(ctx)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *Service) GetUrgentSlips(ctx context.Context, req *ListSlipRequest) (*ListSlipsDTO, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Page > 100 {
		req.PageSize = 100
	}
	if req.PageSize <= 0 {
		req.PageSize = 12
	}

	slips, err := s.repo.GetUrgentSlips(ctx, req)
	if err != nil {
		return nil, err
	}

	var slipDTOs []SlipDTO
	for s := range slips {
		slipDTOs = append(slipDTOs, SlipDTO{
			ID: slips[s].ID,
			User: users.GetUserResponse{
				ID:        "",
				FirstName: slips[s].UserFirstName,
				MiddleName: structs.FromSqlNull(
					slips[s].UserMiddleName,
				),
				LastName: slips[s].UserLastName,
				Email:    slips[s].UserEmail,
			},
			Reason:        slips[s].Reason,
			DateOfAbsence: slips[s].DateOfAbsence,
			DateNeeded:    slips[s].DateNeeded,
			AdminNotes: structs.FromSqlNull(
				slips[s].AdminNotes,
			),
			Category: SlipCategory{
				ID:   slips[s].CategoryID,
				Name: slips[s].CategoryName,
			},
			Status: SlipStatus{
				ID:       slips[s].StatusID,
				Name:     slips[s].StatusName,
				ColorKey: slips[s].StatusColorKey,
			},
			CreatedAt: slips[s].CreatedAt,
			UpdatedAt: slips[s].UpdatedAt,
		})
	}

	req.StatusID = 1
	total, err := s.repo.GetTotalUrgentSlipsCount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get slips count: %w", err)
	}

	listSlipDTO := ListSlipsDTO{
		Slips:      slipDTOs,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: (total + req.PageSize - 1) / req.PageSize,
	}

	return &listSlipDTO, nil
}

func (s *Service) GetSlipStats(
	ctx context.Context,
	iirID *string,
	req *ListSlipRequest,
) ([]SlipStatusCount, error) {
	stats, err := s.repo.GetSlipStats(ctx, iirID, req)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (s *Service) GetAllExcuseSlips(ctx context.Context, req ListSlipRequest) (*ListSlipsDTO, error) {
	// 1. Get raw data from repository
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	slips, err := s.repo.GetAll(ctx, &req)
	if err != nil {
		return nil, err
	}

	var slipDTOs []SlipDTO
	for s := range slips {
		slipDTOs = append(slipDTOs, SlipDTO{
			ID: slips[s].ID,
			User: users.GetUserResponse{
				ID:        "",
				FirstName: slips[s].UserFirstName,
				MiddleName: structs.FromSqlNull(
					slips[s].UserMiddleName,
				),
				LastName: slips[s].UserLastName,
				Email:    slips[s].UserEmail,
			},
			Reason:        slips[s].Reason,
			DateOfAbsence: slips[s].DateOfAbsence,
			DateNeeded:    slips[s].DateNeeded,
			AdminNotes: structs.FromSqlNull(
				slips[s].AdminNotes,
			),
			Category: SlipCategory{
				ID:   slips[s].CategoryID,
				Name: slips[s].CategoryName,
			},
			Status: SlipStatus{
				ID:       slips[s].StatusID,
				Name:     slips[s].StatusName,
				ColorKey: slips[s].StatusColorKey,
			},
			CreatedAt: slips[s].CreatedAt,
			UpdatedAt: slips[s].UpdatedAt,
		})
	}

	total, err := s.repo.GetTotalSlipsCount(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to get slips count: %w", err)
	}

	listSlipsDTO := ListSlipsDTO{
		Slips:      slipDTOs,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: (total + req.PageSize - 1) / req.PageSize,
	}

	return &listSlipsDTO, nil
}

func (s *Service) GetExcuseSlipsByIIRID(
	ctx context.Context,
	iirID string,
	req ListSlipRequest,
) (*ListSlipsDTO, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	slips, err := s.repo.GetByIIRID(ctx, iirID, &req)
	if err != nil {
		return nil, err
	}

	var slipDTOs []SlipDTO
	for s := range slips {
		slipDTOs = append(slipDTOs, SlipDTO{
			ID: slips[s].ID,
			User: users.GetUserResponse{
				ID:        "",
				FirstName: slips[s].UserFirstName,
				MiddleName: structs.FromSqlNull(
					slips[s].UserMiddleName,
				),
				LastName: slips[s].UserLastName,
				Email:    slips[s].UserEmail,
			},
			Reason:        slips[s].Reason,
			DateOfAbsence: slips[s].DateOfAbsence,
			DateNeeded:    slips[s].DateNeeded,
			AdminNotes: structs.FromSqlNull(
				slips[s].AdminNotes,
			),
			Category: SlipCategory{
				ID:   slips[s].CategoryID,
				Name: slips[s].CategoryName,
			},
			Status: SlipStatus{
				ID:       slips[s].StatusID,
				Name:     slips[s].StatusName,
				ColorKey: slips[s].StatusColorKey,
			},
			CreatedAt: slips[s].CreatedAt,
			UpdatedAt: slips[s].UpdatedAt,
		})
	}

	total, err := s.repo.GetTotalSlipsCount(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get slips count: %w",
			err,
		)
	}

	var listSlipsDTO ListSlipsDTO
	listSlipsDTO = ListSlipsDTO{
		Slips:      slipDTOs,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: (total + req.PageSize - 1) / req.PageSize,
	}

	return &listSlipsDTO, nil
}

func (s *Service) GetSlipAttachments(ctx context.Context, slipID int) ([]AttachmentDTO, error) {
	attachments, err := s.repo.GetSlipAttachments(ctx, slipID)
	if err != nil {
		return nil, err
	}

	var attachmentDTOs []AttachmentDTO
	for a := range attachments {
		// Keep FileURL as the URL path (e.g., /slips/{hash}/{filename})
		// Don't convert it to filesystem path - the frontend needs the URL path
		attachmentDTOs = append(attachmentDTOs, AttachmentDTO{
			ID:       attachments[a].ID,
			FileName: attachments[a].FileName,
			FileURL:  attachments[a].FileURL,
		})
	}

	return attachmentDTOs, nil
}

func (s *Service) GetAttachmentFile(ctx context.Context, attachmentID int) (*SlipAttachment, error) {
	attachment, err := s.repo.GetAttachmentByID(ctx, attachmentID)
	if err != nil {
		return nil, err
	}

	if attachment == nil {
		return nil, fmt.Errorf("attachment not found")
	}

	return attachment, nil
}

// SubmitExcuseSlip creates a new slip with attachments.
func (s *Service) SubmitExcuseSlip(
	ctx context.Context,
	iirID string,
	req CreateSlipRequest,
	files []*multipart.FileHeader,
) (*Slip, error) {
	// Check File Size
	if files[0].Size > MaxFileSize {
		return nil, fmt.Errorf(
			"file too large: maximum 5MB allowed",
		)
	}

	// Check File Type
	ext := strings.ToLower(filepath.Ext(files[0].Filename))
	allowedTypes := map[string]bool{
		".pdf":  true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	if !allowedTypes[ext] {
		return nil, fmt.Errorf(
			"invalid file type: PDF and images only",
		)
	}

	parsedDate, err := time.Parse("2006-01-02", req.DateOfAbsence)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid date format: YYYY-MM-DD",
		)
	}

	if parsedDate.After(time.Now()) {
		return nil, fmt.Errorf(
			"absence date cannot be in future",
		)
	}

	folderHash := hash.GetSHA256Hash(
		fmt.Sprintf(
			"%s%s%s%d",
			iirID,
			req.DateOfAbsence,
			req.DateNeeded,
			time.Now().UnixNano(),
		),
		8,
	)

	var fileURLs []string

	for _, file := range files {
		ext := strings.ToLower(
			filepath.Ext(file.Filename),
		)
		fileHash := hash.GetSHA256Hash(
			fmt.Sprintf(
				"%s%d",
				file.Filename,
				time.Now().UnixNano(),
			),
			16,
		)
		uniqueFileName := fileHash + ext

		blobPath := fmt.Sprintf(
			"slips/%s/%s",
			folderHash,
			uniqueFileName,
		)

		if err := s.uploadToBlob(
			ctx,
			file,
			blobPath,
		); err != nil {
			return nil, fmt.Errorf(
				"failed to upload %s: %w",
				file.Filename,
				err,
			)
		}

		fileURLs = append(
			fileURLs,
			fmt.Sprintf(
				"/slips/%s/%s",
				folderHash,
				uniqueFileName,
			),
		)
	}

	slip := &Slip{
		ID:            uuid.New().String(),
		IIRID:         iirID,
		Reason:        req.Reason,
		DateOfAbsence: req.DateOfAbsence,
		DateNeeded:    req.DateNeeded,
		CategoryID:    req.CategoryID,
		StatusID:      1,
	}

	slipID, err := s.repo.CreateSlip(ctx, slip)
	if err != nil {
		return nil, err
	}

	// Loop to create attachment records
	for i, url := range fileURLs {
		attachment := &SlipAttachment{
			ID:       uuid.New().String(),
			SlipID:   slipID,
			FileName: files[i].Filename,
			FileURL:  url,
		}
		if err := s.repo.SaveSlipAttachment(
			ctx,
			attachment,
		); err != nil {
			return nil, err
		}
	}

	auditUserID, ipAddress, userAgent, auditUserEmail := audit.ExtractMeta(ctx)
	s.logService.Record(ctx, logs.LogEntry{
		Category: logs.CategoryAudit,
		Action:   logs.ActionSlipCreated,
		Message: fmt.Sprintf(
			"Excuse slip #%d submitted by %s",
			slip.ID,
			auditUserEmail,
		),
		UserID:    auditUserID,
		UserEmail: auditUserEmail,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Metadata: map[string]interface{}{
			"entity_type": "slip",
			"entity_id":   slip.ID,
			"new_values":  req,
		},
	})

	return slip, nil
}

func (s *Service) uploadToBlob(ctx context.Context, fileHeader *multipart.FileHeader, blobPath string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	contentType := http.DetectContentType(data)
	reader := bytes.NewReader(data)

	return s.storage.Upload(ctx, blobPath, reader, contentType)
}

// DownloadAttachment streams the attachment from Azure Blob Storage.
func (s *Service) DownloadAttachment(ctx context.Context, attachmentID int, writer io.Writer) (*SlipAttachment, error) {
	attachment, err := s.repo.GetAttachmentByID(ctx, attachmentID)
	if err != nil {
		return nil, err
	}
	if attachment == nil {
		return nil, fmt.Errorf("attachment not found")
	}

	// Convert URL path "/slips/hash/file" to blob path "slips/hash/file"
	blobPath := strings.TrimPrefix(attachment.FileURL, "/")

	if err := s.storage.Download(ctx, blobPath, writer); err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return attachment, nil
}

func (s *Service) UpdateExcuseSlipStatus(ctx context.Context, id int, newStatus string, adminNotes string) error {
	validStatuses := map[string]bool{
		"Pending":      true,
		"Approved":     true,
		"Rejected":     true,
		"For Revision": true,
	}

	if !validStatuses[newStatus] {
		return fmt.Errorf("invalid status: must be 'Pending', 'Approved', 'Rejected', or 'For Revision'")
	}

	// Fetch old state for audit trail
	oldSlip, _ := s.repo.GetSlipByID(ctx, id)

	err := s.repo.UpdateStatus(ctx, id, newStatus, adminNotes)
	if err != nil {
		return err
	}

	auditUserID, ipAddress, userAgent, auditUserEmail := audit.ExtractMeta(ctx)
	s.logService.Record(ctx, logs.LogEntry{
		Category:  logs.CategoryAudit,
		Action:    logs.ActionSlipStatusUpdated,
		Message:   fmt.Sprintf("Excuse slip #%d status changed to '%s' by %s", id, newStatus, auditUserEmail),
		UserID:    auditUserID,
		UserEmail: auditUserEmail,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Metadata: map[string]interface{}{
			"entity_type": "slip",
			"entity_id":   id,
			"old_values":  oldSlip,
			"new_values": map[string]interface{}{
				"status":     newStatus,
				"adminNotes": adminNotes,
			},
		},
	})

	return nil
}

// func (s *Service) DeleteExcuseSlip(ctx context.Context, id int) error {
// 	slip, err := s.repo.GetByID(ctx, id)
// 	if err != nil {
// 		return err
// 	}
// 	if slip == nil {
// 		return fmt.Errorf("excuse slip not found")
// 	}

// 	err = s.repo.Delete(ctx, id)
// 	if err != nil {
// 		return err
// 	}

// 	if slip.FileURL != "" {
// 		err := os.Remove(slip.FileURL)
// 		if err != nil {
// 			log.Printf("[Warning] Failed to delete file '%s': %v", slip.FileURL, err)
// 		}
// 	}

// 	return nil
// }
