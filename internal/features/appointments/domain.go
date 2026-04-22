package appointments

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type TimeSlot struct {
	ID   int
	Time string
}

type AppointmentStatus struct {
	ID       int
	Name     string
	ColorKey string
}

type StatusCount struct {
	ID    int
	Name  string
	Count int
}

type AppointmentCategory struct {
	ID   int
	Name string
}

type AvailableTimeSlotView struct {
	ID          int
	Time        string
	IsAvailable bool
}

type Appointment struct {
	ID                    string
	IIRID                 string
	Reason                structs.NullableString
	AdminNotes            structs.NullableString
	WhenDate              string
	TimeSlotID            int
	AppointmentCategoryID int
	StatusID              int
	UrgencyLevel          string
	UrgencyScore          float64
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type AppointmentWithDetailsView struct {
	ID             string
	IIRID          string
	StudentNumber  string
	UserFirstName  string
	UserMiddleName structs.NullableString
	UserLastName   string
	UserEmail      string
	Reason         structs.NullableString
	AdminNotes     structs.NullableString
	WhenDate       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	TimeSlotID     int
	TimeSlotTime   string
	CategoryID     int
	CategoryName   string
	StatusID       int
	StatusName     string
	StatusColorKey string
	UrgencyLevel   string
	UrgencyScore   float64
}

type DailyStatusCount struct {
	Date             string
	PendingCount     int
	ScheduledCount   int
	RescheduledCount int
}
