package slips

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// Category mappers
func (m SlipCategoryDB) ToDomain() SlipCategory {
	return SlipCategory(m)
}

func (d SlipCategory) ToPersistence() SlipCategoryDB {
	return SlipCategoryDB(d)
}

// Status mappers
func (m SlipStatusDB) ToDomain() SlipStatus {
	return SlipStatus(m)
}

func (d SlipStatus) ToPersistence() SlipStatusDB {
	return SlipStatusDB(d)
}

// Details View mappers
func (m SlipWithDetailsViewDB) ToDomain() SlipWithDetailsView {
	return SlipWithDetailsView{
		ID:             m.ID,
		IIRID:          m.IIRID,
		UserID:         m.UserID,
		UserFirstName:  m.UserFirstName,
		UserMiddleName: structs.FromSqlNull(m.UserMiddleName),
		UserLastName:   m.UserLastName,
		UserEmail:      m.UserEmail,
		StudentNumber:  m.StudentNumber,
		Reason:         m.Reason,
		DateOfAbsence:  m.DateOfAbsence,
		DateNeeded:     m.DateNeeded,
		AdminNotes:     structs.FromSqlNull(m.AdminNotes),
		CategoryID:     m.CategoryID,
		CategoryName:   m.CategoryName,
		StatusID:       m.StatusID,
		StatusName:     m.StatusName,
		UrgencyScore:   m.UrgencyScore,
		StatusColorKey: m.StatusColorKey,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

// Status Count mappers
func (m SlipStatusCountDB) ToDomain() SlipStatusCount {
	return SlipStatusCount(m)
}

func (d SlipStatusCount) ToPersistence() SlipStatusCountDB {
	return SlipStatusCountDB(d)
}

// Slip mappers
func (m SlipDB) ToDomain() Slip {
	return Slip{
		ID:            m.ID,
		IIRID:         m.IIRID,
		Reason:        m.Reason,
		DateOfAbsence: m.DateOfAbsence,
		DateNeeded:    m.DateNeeded,
		AdminNotes:    structs.FromSqlNull(m.AdminNotes),
		CategoryID:    m.CategoryID,
		StatusID:      m.StatusID,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func (d Slip) ToPersistence() SlipDB {
	return SlipDB{
		ID:            d.ID,
		IIRID:         d.IIRID,
		Reason:        d.Reason,
		DateOfAbsence: d.DateOfAbsence,
		DateNeeded:    d.DateNeeded,
		AdminNotes:    structs.ToSqlNull(d.AdminNotes),
		CategoryID:    d.CategoryID,
		StatusID:      d.StatusID,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}

// Attachment mappers
func (m SlipAttachmentDB) ToDomain() SlipAttachment {
	return SlipAttachment{
		FileID:         m.FileID,
		SlipID:         structs.FromSqlNull(m.SlipID),
		AttachmentType: m.AttachmentType,
		FileName:       m.FileName,
		FileURL:        m.FileURL,
	}
}

func (d SlipAttachment) ToPersistence() SlipAttachmentDB {
	return SlipAttachmentDB{
		FileID:         d.FileID,
		SlipID:         structs.ToSqlNull(d.SlipID),
		AttachmentType: d.AttachmentType,
		FileName:       d.FileName,
		FileURL:        d.FileURL,
	}
}
