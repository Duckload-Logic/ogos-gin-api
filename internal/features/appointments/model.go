package appointments

import (
	"database/sql"
	"time"
)

type TimeSlot struct {
	ID   int    `db:"id" json:"id"`
	Time string `db:"time" json:"time,omitempty"`
}

type AppointmentStatus struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name,omitempty"`
	ColorKey string `db:"color_key" json:"colorKey,omitempty"`
}

// StatusCount holds the count of appointments per status
type StatusCount struct {
	ID    int    `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Count int    `db:"count" json:"count"`
}

type AppointmentCategory struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name,omitempty"`
}

type AvailableTimeSlotView struct {
	ID          int    `db:"time_slot_id" json:"id"`
	Time        string `db:"time" json:"time,omitempty"`
	IsAvailable bool   `db:"is_available" json:"isAvailable"`
}

type Appointment struct {
	ID                    int            `db:"id" json:"id,omitempty"`
	UserID                int            `db:"user_id" json:"userId,omitempty"`
	Reason                sql.NullString `db:"reason" json:"reason,omitempty"`
	AdminNotes            sql.NullString `db:"admin_notes" json:"adminNotes,omitempty"`
	WhenDate              string         `db:"when_date" json:"whenDate"`
	TimeSlotID            int            `db:"time_slot_id" json:"scheduledTime"`
	AppointmentCategoryID int            `db:"appointment_category_id" json:"appointmentCategoryId"`
	StatusID              int            `db:"status_id" json:"statusId"`
	CreatedAt             time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt             time.Time      `db:"updated_at" json:"updatedAt"`
}

// AppointmentWithDetailsView holds appointment data with joined time slot, category, and status
type AppointmentWithDetailsView struct {
	ID             int            `db:"id"`
	UserFirstName  string         `db:"user_first_name"`
	UserMiddleName sql.NullString `db:"user_middle_name"`
	UserLastName   string         `db:"user_last_name"`
	UserEmail      string         `db:"user_email"`
	Reason         sql.NullString `db:"reason"`
	AdminNotes     sql.NullString `db:"admin_notes"`
	WhenDate       string         `db:"when_date"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
	TimeSlotID     int            `db:"time_slot_id"`
	TimeSlotTime   string         `db:"time_slot_time"`
	CategoryID     int            `db:"category_id"`
	CategoryName   string         `db:"category_name"`
	StatusID       int            `db:"status_id"`
	StatusName     string         `db:"status_name"`
	StatusColorKey string         `db:"status_color_key"`
}

type DailyStatusCount struct {
	Date             string `db:"date"`
	PendingCount     int    `db:"pending_count"`
	ScheduledCount   int    `db:"scheduled_count"`
	RescheduledCount int    `db:"rescheduled_count"`
}
