package excuseslips

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Define constants for security limits
const MaxFileSize = 5 * 1024 * 1024 // 5MB limit

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ==========================================
// |                                        |
// |      EXCUSE SLIP SERVICE FUNCTIONS     |
// |                                        |
// ==========================================

// SubmitExcuseSlip
func (s *Service) SubmitExcuseSlip(ctx context.Context, req CreateExcuseSlipRequest, file *multipart.FileHeader) (*ExcuseSlip, error) {

	// Check File Size (Prevent Disk Filling Attacks)
	if file.Size > MaxFileSize {
		return nil, fmt.Errorf("file too large: maximum allowed size is 5MB")
	}

	// Check File Type (Prevent Malicious Scripts)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedTypes := map[string]bool{
		".pdf": true, ".jpg": true, ".jpeg": true, ".png": true,
	}
	if !allowedTypes[ext] {
		return nil, fmt.Errorf("invalid file type: only PDF and images (JPG, PNG) are allowed")
	}

	// Ensure Student Record Exists (Database Check)
	exists, err := s.repo.CheckStudentExistence(ctx, req.StudentRecordID)
	if err != nil {
		log.Printf("[Error] Failed to validate student ID %d: %v", req.StudentRecordID, err)
		return nil, fmt.Errorf("internal server error: validation failed")
	}
	if !exists {
		return nil, fmt.Errorf("unauthorized: invalid student record")
	}

	// Parse and Validate Date (Prevent Future Dates)
	parsedDate, err := time.Parse("2006-01-02", req.AbsenceDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format (use YYYY-MM-DD)")
	}
	if parsedDate.After(time.Now()) {
		return nil, fmt.Errorf("absence date cannot be in the future")
	}

	// Determine Upload Directory (Env Var > Default)
	baseDir := os.Getenv("UPLOAD_DIR")
	if baseDir == "" {
		baseDir = "uploads"
	}
	uploadDir := filepath.Join(baseDir, "excuse_slips")

	// Create Directory If Missing
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("[Error] Failed to create upload directory '%s': %v", uploadDir, err)
		return nil, fmt.Errorf("internal server error: unable to initialize file storage")
	}

	// Generate Unique File Name & Sanitize Path
	sanitizedFilename := filepath.Base(file.Filename)
	uniqueFileName := fmt.Sprintf("%d_%s", time.Now().Unix(), sanitizedFilename)
	filePath := filepath.Join(uploadDir, uniqueFileName)

	// Save File to Disk
	if err := s.saveFileToDisk(file, filePath); err != nil {
		log.Printf("[Error] Failed to write file to '%s': %v", filePath, err)
		return nil, fmt.Errorf("internal server error: unable to save attached file")
	}

	fileURL := fmt.Sprintf("/uploads/excuse_slips/%s", uniqueFileName)

	slip := &ExcuseSlip{
		StudentRecordID: req.StudentRecordID,
		Reason:          req.Reason,
		Date_of_absence: parsedDate,
		FilePath:        filePath, 
		FileURL:         fileURL,  
		Status:          "Pending",
	}

	// Save to Database
	if err := s.repo.Create(ctx, slip); err != nil {
		if removeErr := os.Remove(filePath); removeErr != nil {
			log.Printf("[Warning] Failed to remove orphaned file '%s': %v", filePath, removeErr)
		}

		log.Printf("[Error] Database insert failed for StudentID %d: %v", req.StudentRecordID, err)
		return nil, fmt.Errorf("internal server error: unable to submit excuse slip")
	}

	return slip, nil
}

// ==========================================
// |                                        |
// |            HELPER FUNCTIONS            |
// |                                        |
// ==========================================

// saveFileToDisk handles the low-level IO operations
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