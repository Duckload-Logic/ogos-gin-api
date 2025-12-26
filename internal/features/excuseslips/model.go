package excuseslips

import "time"

type ExcuseSlip struct {
	ID              int       `json:"id" db:"excuse_slip_id"`
	StudentRecordID int       `json:"studentRecordId" db:"student_record_id"`
	Reason          string    `json:"reason" db:"reason"`
	Date_of_absence     time.Time `json:"date_of_absence" db:"date_of_absence"`
	FilePath        string    `json:"-"`
	FileURL         string    `json:"fileUrl"`
	Status          string    `json:"status" db:"status"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateExcuseSlipRequest struct {
	StudentRecordID int    `form:"studentRecordId" binding:"required"`
	Reason          string `form:"reason" binding:"required"`
	AbsenceDate     string `form:"absenceDate" binding:"required"` 
}
