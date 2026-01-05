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

// ListStudents
func (r *Repository) ListStudents(
	ctx context.Context, offset int, limit int,
	course string, yearLevel int, genderID int,
) ([]StudentProfileView, error) {
	query := `
		SELECT 
			sr.student_record_id,
			u.first_name,
			u.middle_name,
			u.last_name,       
			u.email,
			sp.course,
			sp.year_level,
			sp.section
		FROM student_records sr
		JOIN users u ON sr.user_id = u.user_id
		JOIN student_profiles sp ON sr.student_record_id = sp.student_record_id
		WHERE (? = '' OR sp.course = ?)
		AND (? = 0 OR sp.year_level = ?)
		AND (? = 0 OR sp.gender_id = ?)
		ORDER BY sr.student_record_id DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(
		ctx, query,
		course, course,
		yearLevel, yearLevel,
		genderID, genderID,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list students: %w", err)
	}

	defer rows.Close()

	var students []StudentProfileView
	for rows.Next() {
		var student StudentProfileView
		err := rows.Scan(
			&student.StudentRecordID,
			&student.FirstName,
			&student.MiddleName,
			&student.LastName,
			&student.Email,
			&student.Course,
			&student.YearLevel,
			&student.Section,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan student record: %w", err)
		}
		students = append(students, student)
	}

	return students, nil
}

// GetStudentEnrollmentReasons
func (r *Repository) GetStudentEnrollmentReasons(
	ctx context.Context, studentRecordID int,
) ([]StudentSelectedReason, error) {
	query := `
		SELECT student_record_id, reason_id, other_reason_text
		FROM student_selected_reasons
		WHERE student_record_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, studentRecordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []StudentSelectedReason{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	var reasons []StudentSelectedReason
	for rows.Next() {
		var reason StudentSelectedReason
		err := rows.Scan(
			&reason.StudentRecordID,
			&reason.ReasonID,
			&reason.OtherReasonText,
		)
		if err != nil {
			return nil, err
		}
		reasons = append(reasons, reason)
	}

	return reasons, nil
}

// GetStudentRecordByStudentID
func (r *Repository) GetStudentRecordByStudentID(
	ctx context.Context, userID int,
) (*StudentRecord, error) {
	studentRec := &StudentRecord{}
	query := `
		SELECT 
			student_record_id, user_id
		FROM student_records
		WHERE user_id = ?
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&studentRec.ID,
		&studentRec.UserID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get student record: %w", err)
	}

	return studentRec, nil
}

// GetStudentProfileByStudentRecordID
func (r *Repository) GetStudentProfileByStudentRecordID(
	ctx context.Context, studentRecordID int,
) (*StudentProfile, error) {
	profile := &StudentProfile{}
	query := `
		SELECT 
			student_profile_id, student_record_id,
			gender_id, civil_status_type_id,
			religion_type_id, height_cm, weight_kg,
			student_number, course, year_level, section,
			good_moral_status, has_derogatory_record,
			place_of_birth, birth_date, mobile_no
		FROM student_profiles
		WHERE student_record_id = ?
	`
	err := r.db.QueryRowContext(ctx, query, studentRecordID).Scan(
		&profile.ID, &profile.StudentRecordID,
		&profile.GenderID, &profile.CivilStatusTypeID,
		&profile.ReligionTypeID,
		&profile.HeightCm, &profile.WeightKg,
		&profile.StudentNumber, &profile.Course,
		&profile.YearLevel, &profile.Section,
		&profile.GoodMoralStatus, &profile.HasDerogatoryRecord,
		&profile.PlaceOfBirth, &profile.BirthDate,
		&profile.MobileNo,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get student profile: %w", err)
	}

	return profile, nil
}

// GetGuardians
func (r *Repository) GetGuardians(
	ctx context.Context, studentRecordID int,
) ([]GuardianInfoView, error) {
	query := `
		SELECT
			g.guardian_id, g.educational_level_id,
			g.birth_date, g.last_name, g.first_name,
			g.middle_name, g.occupation, g.maiden_name,
			g.company_name, g.contact_number,
			sg.relationship_type_id, sg.is_primary_contact
		FROM guardians_info_view g
		JOIN student_guardians sg ON g.guardian_id = sg.guardian_id
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

	var guardians []GuardianInfoView
	for rows.Next() {
		var g GuardianInfoView
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
			address_type, region_name,
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
			&addr.AddressType, &addr.RegionName,
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
			vision_remark, hearing_remark, 
			mobility_remark, speech_remark,
			general_health_remark, consulted_professional,
			consultation_reason, date_started, 
			num_sessions, date_concluded
		FROM student_health_records
		WHERE student_record_id = ?
	`
	err := r.db.QueryRowContext(ctx, query, studentRecordID).Scan(
		&healthRec.ID, &healthRec.StudentRecordID,
		&healthRec.VisionRemark, &healthRec.HearingRemark,
		&healthRec.MobilityRemark, &healthRec.SpeechRemark,
		&healthRec.GeneralHealthRemark, &healthRec.ConsultedProfessional,
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

// Add to repository.go
func (r *Repository) GetFinance(ctx context.Context, studentRecordID int) (*StudentFinance, error) {
	finance := &StudentFinance{}
	query := `
		SELECT
			finance_id, student_record_id,
			is_employed, supports_studies,
			supports_family, financial_support_type_id,
			weekly_allowance
		FROM student_finances
		WHERE student_record_id = ?
	`
	err := r.db.QueryRowContext(ctx, query, studentRecordID).Scan(
		&finance.ID, &finance.StudentRecordID,
		&finance.IsEmployed, &finance.SupportsStudies,
		&finance.SupportsFamily, &finance.FinancialSupportTypeID,
		&finance.WeeklyAllowance,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return finance, nil
}

func (r *Repository) GetTotalStudentsCount(
	ctx context.Context,
	course string, yearLevel int, genderID int,
) (int, error) {
	query := `
        SELECT COUNT(*)
        FROM student_records sr
        JOIN users u ON sr.user_id = u.user_id
        JOIN student_profiles sp ON sr.student_record_id = sp.student_record_id
        WHERE (? = '' OR sp.course = ?)
        AND (? = 0 OR sp.year_level = ?)
        AND (? = 0 OR sp.gender_id = ?)
    `

	var total int
	err := r.db.QueryRowContext(
		ctx, query,
		course, course,
		yearLevel, yearLevel,
		genderID, genderID,
	).Scan(&total)

	if err != nil {
		return 0, fmt.Errorf("failed to get total students count: %w", err)
	}

	return total, nil
}

// =====================================
// |                                   |
// |         UPSERT FUNCTIONS          |
// |                                   |
// =====================================

// CreateStudentRecord - Creates a basic student record
func (r *Repository) CreateStudentRecord(
	ctx context.Context, userID int,
) (int, error) {
	query := `
		INSERT INTO student_records (user_id) 
		VALUES (?)
	`

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to create student record: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// SaveEnrollmentReason
func (r *Repository) SaveEnrollmentReason(
	ctx context.Context, studentRecordID int,
	reasonID int, otherReasonText sql.NullString,
) error {
	query := `
		INSERT INTO student_selected_reasons 
		(student_record_id, reason_id, other_reason_text) 
		VALUES (?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query, studentRecordID, reasonID, otherReasonText)
	if err != nil {
		return fmt.Errorf("failed to save enrollment reason: %w", err)
	}

	return nil
}

// SaveStudentProfile
func (r *Repository) SaveStudentProfile(
	ctx context.Context, profile *StudentProfile,
) (int, error) {
	var studentProfileID int

	err := database.RunInTransaction(ctx, r.db, func(tx *sql.Tx) error {
		upsertQuery := `
			INSERT INTO student_profiles (
				student_record_id, gender_id, civil_status_type_id, 
				religion_type_id, height_cm, 
				weight_kg, student_number, 
				course, year_level, 
				section, good_moral_status,
				has_derogatory_record,
				place_of_birth, birth_date, mobile_no
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				gender_id = VALUES(gender_id),
				civil_status_type_id = VALUES(civil_status_type_id),
				religion_type_id = VALUES(religion_type_id),
				height_cm = VALUES(height_cm), weight_kg = VALUES(weight_kg),
				course = VALUES(course), year_level = VALUES(year_level), 
				section = VALUES(section),
				good_moral_status = VALUES(good_moral_status),
				has_derogatory_record = VALUES(has_derogatory_record),
				place_of_birth = VALUES(place_of_birth),
				birth_date = VALUES(birth_date),
				mobile_no = VALUES(mobile_no)
		`

		result, err := tx.ExecContext(
			ctx, upsertQuery,
			profile.StudentRecordID, profile.GenderID,
			profile.CivilStatusTypeID, profile.ReligionTypeID,
			profile.HeightCm, profile.WeightKg,
			profile.StudentNumber, profile.Course,
			profile.YearLevel, profile.Section,
			profile.GoodMoralStatus, profile.HasDerogatoryRecord,
			profile.PlaceOfBirth, profile.BirthDate,
			profile.MobileNo,
		)
		if err != nil {
			return err
		}

		// Get the student profile ID
		id, err := result.LastInsertId()
		if err != nil {
			// If no last insert id (in case of update), get the existing ID
			getIDQuery := `
				SELECT student_profile_id FROM student_profiles 
				WHERE student_record_id = ?
			`
			err = tx.QueryRowContext(
				ctx, getIDQuery,
				profile.StudentRecordID,
			).Scan(&studentProfileID)
			return err
		}

		studentProfileID = int(id)
		return nil
	})

	if err != nil {
		return 0, err
	}

	return studentProfileID, nil
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
				student_record_id, address_type,
				region_name, province_name,
				city_name, barangay_name,
				street_lot_blk, unit_no,
				building_name
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		for _, addr := range addresses {
			_, err = tx.ExecContext(ctx, query,
				studentRecordID, addr.AddressType,
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
				student_record_id, vision_remark,
				hearing_remark, mobility_remark,
				speech_remark, general_health_remark,
				consulted_professional, consultation_reason,
				date_started, num_sessions, date_concluded
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				vision_remark = VALUES(vision_remark),
				hearing_remark = VALUES(hearing_remark),
				mobility_remark = VALUES(mobility_remark),
				speech_remark = VALUES(speech_remark),
				general_health_remark = VALUES(general_health_remark),
				consulted_professional = VALUES(consulted_professional),
				consultation_reason = VALUES(consultation_reason),
				date_started = VALUES(date_started),
				num_sessions = VALUES(num_sessions),
				date_concluded = VALUES(date_concluded)
		`

		_, err := tx.ExecContext(ctx, query,
			health.StudentRecordID, health.VisionRemark,
			health.HearingRemark, health.MobilityRemark,
			health.SpeechRemark, health.GeneralHealthRemark,
			health.ConsultedProfessional, health.ConsultationReason,
			health.DateStarted, health.NumberOfSessions,
			health.DateConcluded,
		)

		return err
	})
}

// SaveFinanceInfo
func (r *Repository) SaveFinanceInfo(
	ctx context.Context, finance *StudentFinance,
) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO student_finances (
				student_record_id, is_employed,
				supports_studies, supports_family,
				financial_support_type_id, weekly_allowance
			) VALUES (?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				is_employed = VALUES(is_employed),
				supports_studies = VALUES(supports_studies),
				supports_family = VALUES(supports_family),
				financial_support_type_id = VALUES(financial_support_type_id),
				weekly_allowance = VALUES(weekly_allowance)
		`

		_, err := tx.ExecContext(ctx, query,
			finance.StudentRecordID, finance.IsEmployed,
			finance.SupportsStudies, finance.SupportsFamily,
			finance.FinancialSupportTypeID, finance.WeeklyAllowance,
		)

		return err
	})
}

// =====================================
// |                                   |
// |        DELETE FUNCTIONS          |
// |                                   |
// =====================================

// DeleteEnrollmentReasons
func (r *Repository) DeleteEnrollmentReasons(
	ctx context.Context, studentRecordID int,
) error {
	query := `DELETE FROM student_selected_reasons WHERE student_record_id = ?`

	_, err := r.db.ExecContext(ctx, query, studentRecordID)
	if err != nil {
		return fmt.Errorf("failed to delete enrollment reasons: %w", err)
	}

	return nil
}

// DeleteFinance
func (r *Repository) DeleteFinance(
	ctx context.Context, studentRecordID int,
) error {
	query := `DELETE FROM student_finances WHERE student_record_id = ?`

	_, err := r.db.ExecContext(ctx, query, studentRecordID)
	if err != nil {
		return fmt.Errorf("failed to delete finance info: %w", err)
	}

	return nil
}
