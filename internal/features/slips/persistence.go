package slips

import (
	"database/sql"
	"time"
)

type SlipCategoryDB struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type SlipStatusDB struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	ColorKey string `db:"color_key"`
}

type SlipWithDetailsViewDB struct {
	ID             string         `db:"id"`
	IIRID          string         `db:"iir_id"`
	UserID         string         `db:"user_id"`
	UserFirstName  string         `db:"user_first_name"`
	UserMiddleName sql.NullString `db:"user_middle_name"`
	UserLastName   string         `db:"user_last_name"`
	UserEmail      string         `db:"user_email"`
	StudentNumber  string         `db:"student_number"`
	Reason         string         `db:"reason"`
	DateOfAbsence  string         `db:"date_of_absence"`
	DateNeeded     string         `db:"date_needed"`
	AdminNotes     sql.NullString `db:"admin_notes"`
	CategoryID     int            `db:"category_id"`
	CategoryName   string         `db:"category_name"`
	StatusID       int            `db:"status_id"`
	StatusName     string         `db:"status_name"`
	UrgencyScore   int            `db:"urgency_score"`
	StatusColorKey string         `db:"status_color_key"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
}

type SlipStatusCountDB struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	ColorKey string `db:"color_key"`
	Count    int    `db:"count"`
}

type SlipDB struct {
	ID            string         `db:"id"`
	IIRID         string         `db:"iir_id"`
	Reason        string         `db:"reason"`
	DateOfAbsence string         `db:"date_of_absence"`
	DateNeeded    string         `db:"date_needed"`
	AdminNotes    sql.NullString `db:"admin_notes"`
	CategoryID    int            `db:"category_id"`
	StatusID      int            `db:"status_id"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
}

type SlipAttachmentDB struct {
	FileID         string         `db:"file_id"`
	SlipID         sql.NullString `db:"admission_slip_id"`
	AttachmentType string         `db:"attachment_type"`
	FileName       string         `db:"file_name"`
	FileURL        string         `db:"file_url"`
}
