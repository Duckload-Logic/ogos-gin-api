package notes

import "time"

type SignificantNoteDTO struct {
	ID        int       `json:"id,omitempty"`
	Note      string    `json:"note" binding:"required"`
	Remarks   string    `json:"remarks" binding:"required"`
	CreatedAt time.Time `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}
