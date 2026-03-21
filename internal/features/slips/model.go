package slips

import (
	"database/sql"
	"time"
)

type SlipCategory struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type SlipStatus struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	ColorKey string `json:"colorKey" db:"color_key"`
}

type SlipWithDetailsView struct {
	ID             string         `json:"id" db:"id"`
	IIRID          string         `json:"iirId" db:"iir_id"`
	UserFirstName  string         `json:"userFirstName" db:"user_first_name"`
	UserMiddleName sql.NullString `json:"userMiddleName" db:"user_middle_name"`
	UserLastName   string         `json:"userLastName" db:"user_last_name"`
	UserEmail      string         `json:"userEmail" db:"user_email"`
	Reason         string         `json:"reason" db:"reason"`
	DateOfAbsence  string         `json:"dateOfAbsence" db:"date_of_absence"`
	DateNeeded     string         `json:"dateNeeded" db:"date_needed"`
	AdminNotes     sql.NullString `json:"adminNotes,omitempty" db:"admin_notes"`
	CategoryID     int            `json:"categoryId" db:"category_id"`
	CategoryName   string         `json:"categoryName" db:"category_name"`
	StatusID       int            `json:"statusId" db:"status_id"`
	StatusName     string         `json:"statusName" db:"status_name"`
	UrgencyScore   int            `json:"urgencyScore" db:"urgency_score"`
	StatusColorKey string         `json:"statusColorKey" db:"status_color_key"`
	CreatedAt      time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time      `json:"updatedAt" db:"updated_at"`
}

type SlipStatusCount struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	ColorKey string `db:"color_key" json:"colorKey"`
	Count    int    `db:"count" json:"count"`
}

type Slip struct {
	ID            string         `json:"id" db:"id"`
	IIRID         string         `json:"iirId" db:"iir_id"`
	Reason        string         `json:"reason" db:"reason"`
	DateOfAbsence string         `json:"dateOfAbsence" db:"date_of_absence"`
	DateNeeded    string         `json:"dateNeeded" db:"date_needed"`
	AdminNotes    sql.NullString `json:"adminNotes,omitempty" db:"admin_notes"`
	CategoryID    int            `json:"categoryId" db:"category_id"`
	StatusID      int            `json:"statusId" db:"status_id"`
	CreatedAt     time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time      `json:"updatedAt" db:"updated_at"`
}

type SlipAttachment struct {
	ID       string  `json:"id" db:"id"`
	SlipID   *string `json:"slipId" db:"admission_slip_id"`
	FileName string  `json:"fileName" db:"file_name"`
	FileURL  string  `json:"fileUrl" db:"file_url"`
}
