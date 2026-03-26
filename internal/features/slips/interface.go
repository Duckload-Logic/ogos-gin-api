package slips

import (
	"context"
	"io"
	"mime/multipart"
)

// ServiceInterface defines the business logic for managing excuse slips.
type ServiceInterface interface {
	GetSlipStatuses(ctx context.Context) ([]SlipStatus, error)
	GetSlipCategories(ctx context.Context) ([]SlipCategory, error)
	GetUrgentSlips(ctx context.Context, req *ListSlipRequest) (*ListSlipsDTO, error)
	GetSlipStats(ctx context.Context, iirID *string, req *ListSlipRequest) ([]SlipStatusCount, error)
	GetAllExcuseSlips(ctx context.Context, req ListSlipRequest) (*ListSlipsDTO, error)
	GetExcuseSlipsByIIRID(ctx context.Context, iirID string, req ListSlipRequest) (*ListSlipsDTO, error)
	GetSlipAttachments(ctx context.Context, slipID string) ([]AttachmentDTO, error)
	GetAttachmentFile(ctx context.Context, attachmentID string) (*SlipAttachment, error)
	SubmitExcuseSlip(ctx context.Context, iirID string, req CreateSlipRequest, files []*multipart.FileHeader) (*Slip, error)
	DownloadAttachment(ctx context.Context, attachmentID string, writer io.Writer) (*SlipAttachment, error)
	UpdateExcuseSlipStatus(ctx context.Context, id string, newStatus string, adminNotes string) error
}

// RepositoryInterface defines the data access layer for managing excuse slips.
type RepositoryInterface interface {
	CreateSlip(ctx context.Context, slip *Slip) (*string, error)
	SaveSlipAttachment(ctx context.Context, attachment *SlipAttachment) error
	CheckStudentExistence(ctx context.Context, studentID int) (bool, error)
	GetSlipStatuses(ctx context.Context) ([]SlipStatus, error)
	GetSlipCategories(ctx context.Context) ([]SlipCategory, error)
	GetSlipStats(ctx context.Context, iirID *string, req *ListSlipRequest) ([]SlipStatusCount, error)
	GetTotalSlipsCount(ctx context.Context, req *ListSlipRequest) (int, error)
	GetTotalUrgentSlipsCount(ctx context.Context, req *ListSlipRequest) (int, error)
	GetUrgentSlips(ctx context.Context, req *ListSlipRequest) ([]SlipWithDetailsView, error)
	GetAll(ctx context.Context, req *ListSlipRequest) ([]SlipWithDetailsView, error)
	GetByUserID(ctx context.Context, userID string, req *ListSlipRequest) ([]SlipWithDetailsView, error)
	GetByIIRID(ctx context.Context, iirID string, req *ListSlipRequest) ([]SlipWithDetailsView, error)
	GetSlipByID(ctx context.Context, id string) (*Slip, error)
	GetSlipAttachments(ctx context.Context, slipID string) ([]SlipAttachment, error)
	GetAttachmentByID(ctx context.Context, attachmentID string) (*SlipAttachment, error)
	UpdateStatus(ctx context.Context, id string, statusName string, adminNotes string) error
	Delete(ctx context.Context, id string) error
}
