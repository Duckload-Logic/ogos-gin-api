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

func (r *Repository) ListStudents(
	ctx context.Context, offset int, limit int,
	course string, genderID int,
) ([]StudentProfileView, error) {
	query := `
        SELECT 
            sr.student_record_id,
            u.first_name,
            u.middle_name,
            u.last_name,       
            u.email,
            sp.course
        FROM student_records sr
        JOIN users u ON sr.user_id = u.user_id
        JOIN student_profiles sp ON sr.student_record_id = sp.student_record_id
        WHERE sr.is_submitted = TRUE
    `

	var args []interface{}

	if course != "" {
		query += " AND sp.course = ?"
		args = append(args, course)
	}

	if genderID != 0 {
		query += " AND sp.gender_id = ?"
		args = append(args, genderID)
	}

	query += `
        ORDER BY sr.student_record_id DESC
        LIMIT ? OFFSET ?
    `
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list students: %w", err)
	}
	defer rows.Close()

	var students []StudentProfileView
	for rows.Next() {
		var student StudentProfileView
		if err := rows.Scan(
			&student.StudentRecordID,
			&student.FirstName,
			&student.MiddleName,
			&student.LastName,
			&student.Email,
			&student.Course,
		); err != nil {
			return nil, err
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
			student_record_id, user_id, 
			is_submitted, created_at, updated_at
		FROM student_records
		WHERE user_id = ?
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&studentRec.ID,
		&studentRec.UserID,
		&studentRec.IsSubmitted,
		&studentRec.CreatedAt,
		&studentRec.UpdatedAt,
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
			religion, height_ft, weight_kg,
			student_number, course, high_school_gwa,
			place_of_birth, birth_date, contact_no
		FROM student_profiles
		WHERE student_record_id = ?
	`
	err := r.db.QueryRowContext(ctx, query, studentRecordID).Scan(
		&profile.ID, &profile.StudentRecordID,
		&profile.GenderID, &profile.CivilStatusTypeID,
		&profile.Religion,
		&profile.HeightFt, &profile.WeightKg,
		&profile.StudentNumber, &profile.Course, &profile.HighSchoolGWA,
		&profile.PlaceOfBirth, &profile.BirthDate,
		&profile.ContactNo,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get student profile: %w", err)
	}

	return profile, nil
}

// GetEmergencyContact
func (r *Repository) GetEmergencyContact(
	ctx context.Context, studentRecordID int,
) (*StudentEmergencyContact, error) {
	emergencyContact := &StudentEmergencyContact{}
	query := `
		SELECT
			emergency_contact_id, student_record_id,
			parent_id, emergency_contact_name,
			emergency_contact_phone, emergency_contact_relationship
		FROM student_emergency_contacts
		WHERE student_record_id = ?
	`
	err := r.db.QueryRowContext(ctx, query, studentRecordID).Scan(
		&emergencyContact.ID,
		&emergencyContact.StudentRecordID,
		&emergencyContact.ParentID,
		&emergencyContact.EmergencyContactName,
		&emergencyContact.EmergencyContactPhone,
		&emergencyContact.EmergencyContactRelationship,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get emergency contact: %w", err)
	}

	return emergencyContact, nil
}

// GetParents
func (r *Repository) GetParents(
	ctx context.Context, studentRecordID int,
) ([]ParentInfoView, error) {
	query := `
		SELECT
			g.parent_id, g.educational_level,
			g.birth_date, g.last_name, g.first_name,
			g.middle_name, g.occupation,
			g.company_name,
			sg.relationship
		FROM parents_info_view g
		JOIN student_parents sg ON g.parent_id = sg.parent_id
		WHERE sg.student_record_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, studentRecordID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var parents []ParentInfoView
	for rows.Next() {
		var g ParentInfoView
		err := rows.Scan(
			&g.ID, &g.EducationalLevel,
			&g.BirthDate, &g.LastName,
			&g.FirstName, &g.MiddleName,
			&g.Occupation,
			&g.CompanyName,
			&g.Relationship,
		)
		if err != nil {
			return nil, err
		}

		parents = append(parents, g)
	}

	return parents, nil
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
			monthly_family_income, guardian_name,
			guardian_address
		FROM family_backgrounds
		WHERE student_record_id = ?
	`
	err := r.db.QueryRowContext(ctx, query, studentRecordID).Scan(
		&familyBg.ID, &familyBg.StudentRecordID,
		&familyBg.ParentalStatusID, &familyBg.ParentalStatusDetails,
		&familyBg.SiblingsBrothers, &familyBg.SiblingSisters,
		&familyBg.MonthlyFamilyIncome, &familyBg.GuardianName,
		&familyBg.GuardianAddress,
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
			educational_level, school_name,
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
			&bg.EducationalLevel, &bg.SchoolName,
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
			employed_family_members_count, supports_studies_count,
			supports_family_count, financial_support,
			weekly_allowance
		FROM student_finances
		WHERE student_record_id = ?
	`
	err := r.db.QueryRowContext(ctx, query, studentRecordID).Scan(
		&finance.ID, &finance.StudentRecordID,
		&finance.EmployedFamilyMembersCount, &finance.SupportsStudiesCount,
		&finance.SupportsFamilyCount, &finance.FinancialSupport,
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
	course string, genderID int,
) (int, error) {
	query := `
        SELECT COUNT(*)
        FROM student_records sr
        JOIN users u ON sr.user_id = u.user_id
        JOIN student_profiles sp ON sr.student_record_id = sp.student_record_id
        WHERE (? = '' OR sp.course = ?)
        AND (? = 0 OR sp.gender_id = ?)
    `

	var total int
	err := r.db.QueryRowContext(
		ctx, query,
		course, course,
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
				religion, height_ft, 
				weight_kg, student_number, 
				course, high_school_gwa,
				place_of_birth, birth_date, contact_no
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				gender_id = VALUES(gender_id),
				civil_status_type_id = VALUES(civil_status_type_id),
				religion = VALUES(religion),
				height_ft = VALUES(height_ft), weight_kg = VALUES(weight_kg),
				course = VALUES(course), 
				place_of_birth = VALUES(place_of_birth),
				birth_date = VALUES(birth_date),
				contact_no = VALUES(contact_no)
		`

		result, err := tx.ExecContext(
			ctx, upsertQuery,
			profile.StudentRecordID, profile.GenderID,
			profile.CivilStatusTypeID, profile.Religion,
			profile.HeightFt, profile.WeightKg,
			profile.StudentNumber, profile.Course, profile.HighSchoolGWA,
			profile.PlaceOfBirth, profile.BirthDate,
			profile.ContactNo,
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

// SaveEmergencyContact
func (r *Repository) SaveEmergencyContact(
	ctx context.Context, emergencyContact *StudentEmergencyContact,
) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sql.Tx) error {
		query := `
			INSERT INTO student_emergency_contacts (
				student_record_id, parent_id,
				emergency_contact_name, emergency_contact_phone,
				emergency_contact_relationship
			) VALUES (?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				parent_id = VALUES(parent_id),
				emergency_contact_name = VALUES(emergency_contact_name),
				emergency_contact_phone = VALUES(emergency_contact_phone),
				emergency_contact_relationship = VALUES(emergency_contact_relationship)
		`

		_, err := tx.ExecContext(ctx, query,
			emergencyContact.StudentRecordID,
			emergencyContact.ParentID,
			emergencyContact.EmergencyContactName,
			emergencyContact.EmergencyContactPhone,
			emergencyContact.EmergencyContactRelationship,
		)

		return err
	})
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
			sibling_sisters, monthly_family_income,
			guardian_name, guardian_address
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			parental_status_id = VALUES(parental_status_id),
			parental_status_details = VALUES(parental_status_details),
			siblings_brothers = VALUES(siblings_brothers),
			sibling_sisters = VALUES(sibling_sisters),
			monthly_family_income = VALUES(monthly_family_income),
			guardian_name = VALUES(guardian_name),
			guardian_address = VALUES(guardian_address)
		
	`

		_, err := tx.ExecContext(ctx, query,
			family.StudentRecordID, family.ParentalStatusID,
			family.ParentalStatusDetails, family.SiblingsBrothers,
			family.SiblingSisters, family.MonthlyFamilyIncome,
			family.GuardianName, family.GuardianAddress,
		)

		return err
	})
}

// SaveParentsInfo
func (r *Repository) SaveParentsInfo(
	ctx context.Context, studentRecordID int,
	parents []Parent,
	links []StudentParent,
) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sql.Tx) error {
		// Cleanup old links for this student
		cleanupQuery := `
			DELETE FROM student_parents WHERE student_record_id = ?
		`
		if _, err := tx.ExecContext(
			ctx, cleanupQuery, studentRecordID,
		); err != nil {
			return err
		}

		// Prepare Queries
		parentQuery := `
            INSERT INTO parents (
                educational_level, birth_date, last_name, first_name,
                middle_name, occupation, company_name
            ) VALUES (?, ?, ?, ?, ?, ?, ?)`

		linkQuery := `
            INSERT INTO student_parents (
                student_record_id, parent_id, 
                relationship
            ) VALUES (?, ?, ?)`

		// Loop through and process each parent
		for i, g := range parents {
			result, err := tx.ExecContext(ctx, parentQuery,
				g.EducationalLevel, g.BirthDate, g.LastName, g.FirstName,
				g.MiddleName, g.Occupation, g.CompanyName,
			)
			if err != nil {
				return fmt.Errorf("failed to insert parent: %w", err)
			}

			parentID, err := result.LastInsertId()
			if err != nil {
				return err
			}

			_, err = tx.ExecContext(ctx, linkQuery,
				studentRecordID,
				parentID,
				links[i].Relationship,
			)
			if err != nil {
				return fmt.Errorf("failed to link parent to student: %w", err)
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
				student_record_id, educational_level, 
				school_name, location, 
				school_type, year_completed, 
				awards
			) VALUES (?, ?, ?, ?, ?, ?, ?)
		`
		for _, edu := range educations {
			_, err = tx.ExecContext(ctx, query,
				studentRecordID, edu.EducationalLevel,
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
				student_record_id, employed_family_members_count,
				supports_studies_count, supports_family_count,
				financial_support, weekly_allowance
			) VALUES (?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE
				employed_family_members_count = VALUES(employed_family_members_count),
				supports_studies_count = VALUES(supports_studies_count),
				supports_family_count = VALUES(supports_family_count),
				financial_support = VALUES(financial_support),
				weekly_allowance = VALUES(weekly_allowance)
		`

		_, err := tx.ExecContext(ctx, query,
			finance.StudentRecordID, finance.EmployedFamilyMembersCount,
			finance.SupportsStudiesCount, finance.SupportsFamilyCount,
			finance.FinancialSupport, finance.WeeklyAllowance,
		)

		return err
	})
}

func (r *Repository) MarkOnboardingComplete(
	ctx context.Context, studentRecordID int,
) error {
	query := `
		UPDATE student_records
		SET is_submitted = TRUE
		WHERE student_record_id = ?
	`

	_, err := r.db.ExecContext(ctx, query, studentRecordID)
	if err != nil {
		return fmt.Errorf("failed to mark onboarding complete: %w", err)
	}

	return nil
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
