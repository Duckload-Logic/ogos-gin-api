package slips

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type SlipCategory struct {
	ID   int
	Name string
}

type SlipStatus struct {
	ID       int
	Name     string
	ColorKey string
}

type SlipWithDetailsView struct {
	ID             string
	IIRID          string
	UserID         string
	UserFirstName  string
	UserMiddleName structs.NullableString
	UserLastName   string
	UserEmail      string
	StudentNumber  string
	Reason         string
	DateOfAbsence  string
	DateNeeded     string
	AdminNotes     structs.NullableString
	CategoryID     int
	CategoryName   string
	StatusID       int
	StatusName     string
	UrgencyScore   int
	StatusColorKey string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type SlipStatusCount struct {
	ID       int
	Name     string
	ColorKey string
	Count    int
}

type Slip struct {
	ID            string
	IIRID         string
	Reason        string
	DateOfAbsence string
	DateNeeded    string
	AdminNotes    structs.NullableString
	CategoryID    int
	StatusID      int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type SlipAttachment struct {
	FileID         string
	SlipID         structs.NullableString
	AttachmentType string
	FileName       string
	FileURL        string
}
