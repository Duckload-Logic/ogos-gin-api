package slips

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/request"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

type ListSlipRequest struct {
	request.PaginationParams
	Search    string `form:"search,omitempty"`
	StatusID  int    `form:"status_id"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	OrderBy   string `form:"order_by" binding:"omitempty,oneof=created_at"`
}

type ListSlipsDTO struct {
	Slips      []SlipDTO `json:"slips"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	PageSize   int       `json:"pageSize"`
	TotalPages int       `json:"totalPages"`
}

type SlipDTO struct {
	ID            int                    `json:"id,omitempty"`
	User          users.GetUserResponse  `json:"user,omitempty"`
	Reason        string                 `json:"reason" form:"reason" binding:"required"`
	DateOfAbsence string                 `json:"dateOfAbsence" form:"dateOfAbsence" binding:"required"`
	DateNeeded    string                 `json:"dateNeeded" form:"dateNeeded" binding:"required"`
	AdminNotes    structs.NullableString `json:"adminNotes,omitempty"`
	Category      SlipCategory           `json:"category" form:"categoryId" binding:"required"`
	Status        SlipStatus             `json:"status,omitempty"`
	CreatedAt     time.Time              `json:"createdAt,omitempty"`
	UpdatedAt     time.Time              `json:"updatedAt,omitempty"`
}

type AttachmentDTO struct {
	ID       int    `json:"id"`
	FileName string `json:"fileName"`
	FileURL  string `json:"fileUrl"`
}

type CreateSlipRequest struct {
	Reason        string `json:"reason" form:"reason" binding:"required"`
	DateOfAbsence string `json:"dateOfAbsence" form:"dateOfAbsence" binding:"required"`
	DateNeeded    string `json:"dateNeeded" form:"dateNeeded" binding:"required"`
	CategoryID    int    `json:"categoryId" form:"categoryId" binding:"required"`
}

type UpdateStatusRequest struct {
	Status     string `json:"status" binding:"required"`
	AdminNotes string `json:"adminNotes"`
}
