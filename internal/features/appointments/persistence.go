package appointments

import (
	"database/sql"
	"time"
)

type TimeSlotDB struct {
	ID   int    `db:"id"`
	Time string `db:"time"`
}

type AppointmentStatusDB struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	ColorKey string `db:"color_key"`
}

type StatusCountDB struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Count int    `db:"count"`
}

type AppointmentCategoryDB struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type AvailableTimeSlotViewDB struct {
	ID          int    `db:"time_slot_id"`
	Time        string `db:"time"`
	IsAvailable bool   `db:"is_available"`
}

type AppointmentDB struct {
	ID                    string         `db:"id"`
	IIRID                 string         `db:"iir_id"`
	Reason                sql.NullString `db:"reason"`
	AdminNotes            sql.NullString `db:"admin_notes"`
	WhenDate              string         `db:"when_date"`
	TimeSlotID            int            `db:"time_slot_id"`
	AppointmentCategoryID int            `db:"appointment_category_id"`
	StatusID              int            `db:"status_id"`
	UrgencyLevel          string         `db:"urgency_level"`
	UrgencyScore          float64        `db:"urgency_score"`
	CreatedAt             time.Time      `db:"created_at"`
	UpdatedAt             time.Time      `db:"updated_at"`
}

type AppointmentWithDetailsViewDB struct {
	ID             string         `db:"id"`
	IIRID          string         `db:"iir_id"`
	StudentNumber  string         `db:"student_number"`
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
	UrgencyLevel   string         `db:"urgency_level"`
	UrgencyScore   float64        `db:"urgency_score"`
}

type DailyStatusCountDB struct {
	Date             string `db:"date"`
	PendingCount     int    `db:"pending_count"`
	ScheduledCount   int    `db:"scheduled_count"`
	RescheduledCount int    `db:"rescheduled_count"`
}
