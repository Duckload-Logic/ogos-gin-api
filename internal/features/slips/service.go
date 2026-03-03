package slips

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/builders"
	"github.com/olazo-johnalbert/duckload-api/internal/core/hash"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

const MaxFileSize = 5 * 1024 * 1024 // 5MB limit

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
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
				ID:         slips[s].UserID,
				FirstName:  slips[s].UserFirstName,
				MiddleName: structs.FromSqlNull(slips[s].UserMiddleName),
				LastName:   slips[s].UserLastName,
				Email:      slips[s].UserEmail,
			},
			Reason:        slips[s].Reason,
			DateOfAbsence: slips[s].DateOfAbsence,
			DateNeeded:    slips[s].DateNeeded,
			AdminNotes:    structs.FromSqlNull(slips[s].AdminNotes),
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

func (s *Service) GetSlipStats(ctx context.Context, userID *int, req *ListSlipRequest) ([]SlipStatusCount, error) {
	stats, err := s.repo.GetSlipStats(ctx, userID, req)
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
				ID:         slips[s].UserID,
				FirstName:  slips[s].UserFirstName,
				MiddleName: structs.FromSqlNull(slips[s].UserMiddleName),
				LastName:   slips[s].UserLastName,
				Email:      slips[s].UserEmail,
			},
			Reason:        slips[s].Reason,
			DateOfAbsence: slips[s].DateOfAbsence,
			DateNeeded:    slips[s].DateNeeded,
			AdminNotes:    structs.FromSqlNull(slips[s].AdminNotes),
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

func (s *Service) GetExcuseSlipsByUserID(ctx context.Context, userID int, req ListSlipRequest) (*ListSlipsDTO, error) {
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

	slips, err := s.repo.GetByUserID(ctx, userID, &req)
	if err != nil {
		return nil, err
	}

	var slipDTOs []SlipDTO
	for s := range slips {
		slipDTOs = append(slipDTOs, SlipDTO{
			ID: slips[s].ID,
			User: users.GetUserResponse{
				ID:         slips[s].UserID,
				FirstName:  slips[s].UserFirstName,
				MiddleName: structs.FromSqlNull(slips[s].UserMiddleName),
				LastName:   slips[s].UserLastName,
				Email:      slips[s].UserEmail,
			},
			Reason:        slips[s].Reason,
			DateOfAbsence: slips[s].DateOfAbsence,
			DateNeeded:    slips[s].DateNeeded,
			AdminNotes:    structs.FromSqlNull(slips[s].AdminNotes),
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

// SubmitExcuseSlip
func (s *Service) SubmitExcuseSlip(ctx context.Context, userID int, req CreateSlipRequest, files []*multipart.FileHeader) (*Slip, error) {

	// Check File Size
	if files[0].Size > MaxFileSize {
		return nil, fmt.Errorf("file too large: maximum allowed size is 5MB")
	}

	// Check File Type
	ext := strings.ToLower(filepath.Ext(files[0].Filename))
	allowedTypes := map[string]bool{
		".pdf": true, ".jpg": true, ".jpeg": true, ".png": true,
	}
	if !allowedTypes[ext] {
		return nil, fmt.Errorf("invalid file type: only PDF and images (JPG, PNG) are allowed")
	}

	parsedDate, err := time.Parse("2006-01-02", req.DateOfAbsence)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: expected YYYY-MM-DD")
	}

	if parsedDate.After(time.Now()) {
		return nil, fmt.Errorf("absence date cannot be in the future")
	}

	folderHash := hash.GetSHA256Hash(
		fmt.Sprintf(
			"%d%s%s%d",
			userID,
			req.DateOfAbsence,
			req.DateNeeded,
			time.Now().UnixNano()),
		8,
	)
	uploadDir := builders.BuildFileURL("slips", folderHash)

	// Create Directory If Missing
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("internal server error: unable to initialize file storage")
	}

	var savedPaths []string
	var fileURLs []string

	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		fileHash := hash.GetSHA256Hash(fmt.Sprintf("%s%d", file.Filename, time.Now().UnixNano()), 16)
		uniqueFileName := fileHash + ext

		finalPath := filepath.Join(uploadDir, uniqueFileName)

		if err := s.saveFileToDisk(file, finalPath); err != nil {
			os.Remove(uploadDir)
			return nil, fmt.Errorf("failed to save %s", file.Filename)
		}

		savedPaths = append(savedPaths, finalPath)
		fileURLs = append(fileURLs, fmt.Sprintf("/slips/%s/%s", folderHash, uniqueFileName))
	}

	slip := &Slip{
		UserID:        userID,
		Reason:        req.Reason,
		DateOfAbsence: req.DateOfAbsence,
		DateNeeded:    req.DateNeeded,
		CategoryID:    req.CategoryID,
		StatusID:      1,
	}

	slipID, err := s.repo.CreateSlip(ctx, slip) // Gets the slip.ID
	if err != nil {
		os.Remove(uploadDir)
		return nil, err
	}

	// 2. Loop to create attachment records
	for i, url := range fileURLs {
		attachment := &SlipAttachment{
			SlipID:   *slipID,
			FileName: files[i].Filename,
			FileURL:  url,
		}
		if err := s.repo.SaveSlipAttachment(ctx, attachment); err != nil {
			os.Remove(uploadDir)
			return nil, err
		}
	}

	return slip, nil
}

func (s *Service) saveFileToDisk(fileHeader *multipart.FileHeader, destPath string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
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

	err := s.repo.UpdateStatus(ctx, id, newStatus, adminNotes)
	if err != nil {
		return err
	}

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
