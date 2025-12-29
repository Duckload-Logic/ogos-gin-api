package students

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/olazo-johnalbert/duckload-api/internal/database"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// ======================================
// |                                   |
// |         RETRIEVE FUNCTIONS        |
// |                                   |
// =====================================

// GetStudentRecordByStudentID
func (r *Repository) GetStudentRecordByStudentID(
	ctx context.Context, userID int,
) (*StudentRecord, error) {
	studentRec := &StudentRecord{}
	query := `
		SELECT 
			student_record_id, user_id, 
			gender_id, civil_status_type_id,
			religion_type_id, height_cm, weight_kg,
			student_number, course, year_level, section,
			good_moral_status, has_derogatory_record,
			place_of_birth, birth_date, mobile_no
		FROM student_records
		WHERE user_id = ?
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&studentRec.ID, &studentRec.UserID,
		&studentRec.GenderID, &studentRec.CivilStatusTypeID,
		&studentRec.ReligionTypeID,
		&studentRec.HeightCm, &studentRec.WeightKg,
		&studentRec.StudentNumber, &studentRec.Course,
		&studentRec.YearLevel, &studentRec.Section,
		&studentRec.GoodMoralStatus, &studentRec.HasDerogatoryRecord,
		&studentRec.PlaceOfBirth, &studentRec.BirthDate,
		&studentRec.MobileNo,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get student record: %w", err)
	}

	return studentRec, nil
}

// GetGuardians
func (r *Repository) GetGuardians(
	ctx context.Context, studentRecordID int,
) ([]Guardian, error) {
	query := `
		SELECT
			g.guardian_id, g.educational_level_id, g.birth_date,
			g.last_name, g.first_name, g.middle_name,
			g.occupation, g.maiden_name, g.company_name,
			g.contact_number,
			sg.relationship_type_id, sg.is_primary_contact
		FROM guardians g
		INNER JOIN student_guardians sg ON g.guardian_id = sg.guardian_id
		WHERE sg.student_record_id = ?
		ORDER BY sg.is_primary_contact DESC
	`
	rows, err := r.db.QueryContext(ctx, query, studentRecordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var guardians []Guardian
	for rows.Next() {
		var g Guardian
		err := rows.Scan(
			&g.ID, &g.EducationalLevelID,
			&g.BirthDate, &g.LastName,
			&g.FirstName, &g.MiddleName,
			&g.Occupation, &g.MaidenName,
			&g.CompanyName, &g.ContactNumber,
			&g.RelationshipTypeID, &g.IsPrimaryContact,
		)
		if err != nil {
			return nil, err
		}

		guardians = append(guardians, g)
	}

	return guardians, nil
}

// GetPrimaryGuardian
func (r *Repository) GetPrimaryGuardian(
	ctx context.Context, studentRecordID int,
) (*Guardian, error) {
	query := `
		SELECT
			g.guardian_id, g.educational_level_id, g.birth_date,
			g.last_name, g.first_name, g.middle_name,
			g.occupation, g.maiden_name, g.company_name,
			g.contact_number,
			sg.relationship_type_id, sg.is_primary_contact
		FROM guardians g
		INNER JOIN student_guardians sg ON g.guardian_id = sg.guardian_id
		WHERE sg.student_record_id = ? AND sg.is_primary_contact = TRUE
		LIMIT 1
	`

	var g Guardian
	err := r.db.QueryRowContext(ctx, query, studentRecordID).Scan(
		&g.ID, &g.EducationalLevelID,
		&g.BirthDate, &g.LastName,
		&g.FirstName, &g.MiddleName,
		&g.Occupation, &g.MaidenName,
		&g.CompanyName, &g.ContactNumber,
		&g.RelationshipTypeID, &g.IsPrimaryContact,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &g, nil
}

// GetFamily
func (r *Repository) GetFamily(
	ctx context.Context, studentRecordID int,
) (*FamilyBackground, error) {
	familyBg := &FamilyBackground{}
	query := `
		SELECT 
			family_background_id, student_record_id,
			parental_status_id, parental_status_details,
			siblings_brothers, sibling_sisters,
			monthly_family_income
		FROM family_backgrounds
		WHERE student_record_id = ?
	`
	err := r.db.QueryRowContext(ctx, query, studentRecordID).Scan(
		&familyBg.ID, &familyBg.StudentRecordID,
		&familyBg.ParentalStatusID, &familyBg.ParentalStatusDetails,
		&familyBg.SiblingsBrothers, &familyBg.SiblingSisters,
		&familyBg.MonthlyFamilyIncome,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return familyBg, nil
}

// GetEducationalBackgrounds
func (r *Repository) GetEducationalBackgrounds(
	ctx context.Context, studentRecordID int,
) ([]EducationalBackground, error) {
	query := `
		SELECT
			educational_background_id, student_record_id,
			educational_level_id, school_name,
			location, school_type,
			year_completed, awards
		FROM educational_backgrounds
		WHERE student_record_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, studentRecordID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var backgrounds []EducationalBackground
	for rows.Next() {
		var bg EducationalBackground
		err := rows.Scan(
			&bg.ID, &bg.StudentRecordID,
			&bg.EducationalLevelID, &bg.SchoolName,
			&bg.Location, &bg.SchoolType,
			&bg.YearCompleted, &bg.Awards,
		)
		if err != nil {
			return nil, err
		}
		backgrounds = append(backgrounds, bg)
	}

	return backgrounds, nil
}

// GetAddresses
func (r *Repository) GetAddresses(
	ctx context.Context, studentRecordID int,
) ([]StudentAddress, error) {
	query := `
		SELECT
			student_address_id, student_record_id,
			address_type_id, region_name,
			province_name, city_name,
			barangay_name, street_lot_blk,
			unit_no, building_name
		FROM student_addresses
		WHERE student_record_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, studentRecordID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var addresses []StudentAddress
	for rows.Next() {
		var addr StudentAddress
		err := rows.Scan(
			&addr.ID, &addr.StudentRecordID,
			&addr.AddressTypeID, &addr.RegionName,
			&addr.ProvinceName, &addr.CityName,
			&addr.BarangayName, &addr.StreetLotBlk,
			&addr.UnitNo, &addr.BuildingName,
		)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, addr)
	}

	return addresses, nil
}

// GetHealthRecord
func (r *Repository) GetHealthRecord(
	ctx context.Context, studentRecordID int,
) (*StudentHealthRecord, error) {
	healthRec := &StudentHealthRecord{}
	query := `
		SELECT
			health_id, student_record_id,
			vision_remark_id, hearing_remark_id, 
			mobility_remark_id, speech_remark_id,
			general_health_remark_id, consulted_professional,
			consultation_reason, date_started, 
			num_sessions, date_concluded
		FROM student_health_records
		WHERE student_record_id = ?
	`
	err := r.db.QueryRowContext(ctx, query, studentRecordID).Scan(
		&healthRec.ID, &healthRec.StudentRecordID,
		&healthRec.VisionRemarkID, &healthRec.HearingRemarkID,
		&healthRec.MobilityRemarkID, &healthRec.SpeechRemarkID,
		&healthRec.GeneralHealthRemarkID, &healthRec.ConsultedProfessional,
		&healthRec.ConsultationReason, &healthRec.DateStarted,
		&healthRec.NumberOfSessions, &healthRec.DateConcluded,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return healthRec, nil
}

// =====================================
// |                                   |
// |         UPSERT FUNCTIONS          |
// |                                   |
// =====================================

// SaveBaseProfileInfo
func (r *Repository) SaveBaseProfileInfo(
	ctx context.Context, record *StudentRecord,
) (int, error) {
	var studentRecordID int

	err := database.RunInTransaction(ctx, r.db, func(tx *sql.Tx) error {
		upsertQuery := `
			INSERT INTO student_records (
				user_id, civil_status_type_id, 
				religion_type_id, height_cm, 
				weight_kg, student_number, 
				course, year_level, 
				section, good_moral_status,
				has_derogatory_record, gender_id,
				place_of_birth, birth_date, mobile_no
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				civil_status_type_id = VALUES(civil_status_type_id),
				religion_type_id = VALUES(religion_type_id),
				height_cm = VALUES(height_cm), weight_kg = VALUES(weight_kg),
				course = VALUES(course), year_level = VALUES(year_level), 
				section = VALUES(section),
				good_moral_status = VALUES(good_moral_status),
				has_derogatory_record = VALUES(has_derogatory_record),
				gender_id = VALUES(gender_id),
				place_of_birth = VALUES(place_of_birth),
				birth_date = VALUES(birth_date),
				mobile_no = VALUES(mobile_no)
		`

		_, err := tx.ExecContext(
			ctx, upsertQuery,
			record.UserID, record.CivilStatusTypeID,
			record.ReligionTypeID, record.HeightCm,
			record.WeightKg, record.StudentNumber,
			record.Course, record.YearLevel,
			record.Section, record.GoodMoralStatus,
			record.HasDerogatoryRecord,
			record.GenderID,
			record.PlaceOfBirth,
			record.BirthDate,
			record.MobileNo,
		)
		if err != nil {
			return err
		}

		getRecIDQuery := `
			SELECT student_record_id FROM student_records 
			WHERE user_id = ?
		`
		err = tx.QueryRowContext(
			ctx, getRecIDQuery,
			record.UserID,
		).Scan(&studentRecordID)

		return err
	})

	if err != nil {
		return 0, err
	}

	return studentRecordID, nil
}

// SaveFamilyInfo
func (r *Repository) SaveFamilyInfo(
	ctx context.Context, family *FamilyBackground,
) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sql.Tx) error {
		query := `
		INSERT INTO family_backgrounds (
			student_record_id, parental_status_id,
			parental_status_details, siblings_brothers,
			sibling_sisters, monthly_family_income
		) VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			parental_status_id = VALUES(parental_status_id),
			parental_status_details = VALUES(parental_status_details),
			siblings_brothers = VALUES(siblings_brothers),
			sibling_sisters = VALUES(sibling_sisters),
			monthly_family_income = VALUES(monthly_family_income)
	`

		_, err := tx.ExecContext(ctx, query,
			family.StudentRecordID, family.ParentalStatusID,
			family.ParentalStatusDetails, family.SiblingsBrothers,
			family.SiblingSisters, family.MonthlyFamilyIncome,
		)

		return err
	})
}

// SaveGuardiansInfo
func (r *Repository) SaveGuardiansInfo(
	ctx context.Context, studentRecordID int,
	guardians []Guardian,
	links []StudentGuardian,
) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sql.Tx) error {
		// Cleanup old links for this student
		cleanupQuery := `
			DELETE FROM student_guardians WHERE student_record_id = ?
		`
		if _, err := tx.ExecContext(
			ctx, cleanupQuery, studentRecordID,
		); err != nil {
			return err
		}

		// Prepare Queries
		guardianQuery := `
            INSERT INTO guardians (
                educational_level_id, birth_date, last_name, first_name,
                middle_name, occupation, maiden_name, company_name,
				contact_number
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

		linkQuery := `
            INSERT INTO student_guardians (
                student_record_id, guardian_id, 
                relationship_type_id, is_primary_contact
            ) VALUES (?, ?, ?, ?)`

		// Loop through and process each guardian
		for i, g := range guardians {
			result, err := tx.ExecContext(ctx, guardianQuery,
				g.EducationalLevelID, g.BirthDate, g.LastName, g.FirstName,
				g.MiddleName, g.Occupation, g.MaidenName, g.CompanyName,
				g.ContactNumber,
			)
			if err != nil {
				return fmt.Errorf("failed to insert guardian: %w", err)
			}

			guardianID, err := result.LastInsertId()
			if err != nil {
				return err
			}

			// 4. Create the link in student_guardians
			_, err = tx.ExecContext(ctx, linkQuery,
				studentRecordID,
				guardianID,
				links[i].RelationshipTypeID,
				links[i].IsPrimaryContact,
			)
			if err != nil {
				return fmt.Errorf("failed to link guardian to student: %w", err)
			}
		}
		return nil
	})
}

// SaveEducationInfo
func (r *Repository) SaveEducationInfo(
	ctx context.Context, studentRecordID int,
	educations []EducationalBackground,
) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sql.Tx) error {
		cleanupQuery := `
			DELETE FROM educational_backgrounds WHERE student_record_id = ?
		`
		_, err := tx.ExecContext(ctx, cleanupQuery, studentRecordID)
		if err != nil {
			return err
		}

		query := `
			INSERT INTO educational_backgrounds (
				student_record_id, educational_level_id, 
				school_name, location, 
				school_type, year_completed, 
				awards
			) VALUES (?, ?, ?, ?, ?, ?, ?)
		`
		for _, edu := range educations {
			_, err = tx.ExecContext(ctx, query,
				studentRecordID, edu.EducationalLevelID,
				edu.SchoolName, edu.Location,
				edu.SchoolType, edu.YearCompleted,
				edu.Awards,
			)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// SaveAddressInfo
func (r *Repository) SaveAddressInfo(
	ctx context.Context, studentRecordID int,
	addresses []StudentAddress,
) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sql.Tx) error {
		cleanupQuery := `
			DELETE FROM student_addresses WHERE student_record_id = ?
		`
		_, err := tx.ExecContext(ctx, cleanupQuery, studentRecordID)
		if err != nil {
			return err
		}

		query := `
			INSERT INTO student_addresses (
				student_record_id, address_type_id,
				region_name, province_name,
				city_name, barangay_name,
				street_lot_blk, unit_no,
				building_name
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		for _, addr := range addresses {
			_, err = tx.ExecContext(ctx, query,
				studentRecordID, addr.AddressTypeID,
				addr.RegionName, addr.ProvinceName,
				addr.CityName, addr.BarangayName,
				addr.StreetLotBlk, addr.UnitNo,
				addr.BuildingName,
			)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// SaveHealthRecord
func (r *Repository) SaveHealthRecord(
	ctx context.Context, health *StudentHealthRecord,
) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO student_health_records (
				student_record_id, vision_remark_id,
				hearing_remark_id, mobility_remark_id,
				speech_remark_id, general_health_remark_id,
				consulted_professional, consultation_reason,
				date_started, num_sessions, date_concluded
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				vision_remark_id = VALUES(vision_remark_id),
				hearing_remark_id = VALUES(hearing_remark_id),
				mobility_remark_id = VALUES(mobility_remark_id),
				speech_remark_id = VALUES(speech_remark_id),
				general_health_remark_id = VALUES(general_health_remark_id),
				consulted_professional = VALUES(consulted_professional),
				consultation_reason = VALUES(consultation_reason),
				date_started = VALUES(date_started),
				num_sessions = VALUES(num_sessions),
				date_concluded = VALUES(date_concluded)
		`

		_, err := tx.ExecContext(ctx, query,
			health.StudentRecordID, health.VisionRemarkID,
			health.HearingRemarkID, health.MobilityRemarkID,
			health.SpeechRemarkID, health.GeneralHealthRemarkID,
			health.ConsultedProfessional, health.ConsultationReason,
			health.DateStarted, health.NumberOfSessions,
			health.DateConcluded,
		)

		return err
	})
}
