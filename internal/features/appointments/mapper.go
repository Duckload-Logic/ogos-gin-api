package appointments

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// TimeSlot mappers
func (m TimeSlotDB) ToDomain() TimeSlot {
	return TimeSlot(m)
}

func (d TimeSlot) ToPersistence() TimeSlotDB {
	return TimeSlotDB(d)
}

// Status mappers
func (m AppointmentStatusDB) ToDomain() AppointmentStatus {
	return AppointmentStatus(m)
}

// Category mappers
func (m AppointmentCategoryDB) ToDomain() AppointmentCategory {
	return AppointmentCategory(m)
}

// AvailableTimeSlot mappers
func (m AvailableTimeSlotViewDB) ToDomain() AvailableTimeSlotView {
	return AvailableTimeSlotView(m)
}

// Appointment mappers
func (m AppointmentDB) ToDomain() Appointment {
	return Appointment{
		ID:                    m.ID,
		IIRID:                 m.IIRID,
		Reason:                structs.FromSqlNull(m.Reason),
		AdminNotes:            structs.FromSqlNull(m.AdminNotes),
		WhenDate:              m.WhenDate,
		TimeSlotID:            m.TimeSlotID,
		AppointmentCategoryID: m.AppointmentCategoryID,
		StatusID:              m.StatusID,
		UrgencyLevel:          m.UrgencyLevel,
		UrgencyScore:          m.UrgencyScore,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
	}
}

func (d Appointment) ToPersistence() AppointmentDB {
	return AppointmentDB{
		ID:                    d.ID,
		IIRID:                 d.IIRID,
		Reason:                structs.ToSqlNull(d.Reason),
		AdminNotes:            structs.ToSqlNull(d.AdminNotes),
		WhenDate:              d.WhenDate,
		TimeSlotID:            d.TimeSlotID,
		AppointmentCategoryID: d.AppointmentCategoryID,
		StatusID:              d.StatusID,
		UrgencyLevel:          d.UrgencyLevel,
		UrgencyScore:          d.UrgencyScore,
		CreatedAt:             d.CreatedAt,
		UpdatedAt:             d.UpdatedAt,
	}
}

// Details View mapper
func (m AppointmentWithDetailsViewDB) ToDomain() AppointmentWithDetailsView {
	return AppointmentWithDetailsView{
		ID:             m.ID,
		IIRID:          m.IIRID,
		StudentNumber:  m.StudentNumber,
		UserFirstName:  m.UserFirstName,
		UserMiddleName: structs.FromSqlNull(m.UserMiddleName),
		UserLastName:   m.UserLastName,
		UserEmail:      m.UserEmail,
		Reason:         structs.FromSqlNull(m.Reason),
		AdminNotes:     structs.FromSqlNull(m.AdminNotes),
		WhenDate:       m.WhenDate,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		TimeSlotID:     m.TimeSlotID,
		TimeSlotTime:   m.TimeSlotTime,
		CategoryID:     m.CategoryID,
		CategoryName:   m.CategoryName,
		StatusID:       m.StatusID,
		StatusName:     m.StatusName,
		StatusColorKey: m.StatusColorKey,
		UrgencyLevel:   m.UrgencyLevel,
		UrgencyScore:   m.UrgencyScore,
	}
}

// Status Count mappers
func (m StatusCountDB) ToDomain() StatusCount {
	return StatusCount(m)
}

func (m DailyStatusCountDB) ToDomain() DailyStatusCount {
	return DailyStatusCount(m)
}
