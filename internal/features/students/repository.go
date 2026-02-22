package students

import (
	"context"
	"database/sql"
	"fmt"
)

// Helper function to convert string to *string (NULL if empty)
func nullString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Lookup
func (r *Repository) GetGenders(ctx context.Context) ([]Gender, error) {
	query := `SELECT id, gender_name FROM genders ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get genders: %w", err)
	}
	defer rows.Close()

	var genders []Gender
	for rows.Next() {
		var g Gender
		if err := rows.Scan(&g.ID, &g.GenderName); err != nil {
			return nil, err
		}
		genders = append(genders, g)
	}

	return genders, nil
}

func (r *Repository) GetParentalStatusTypes(ctx context.Context) ([]ParentalStatusType, error) {
	query := `SELECT id, status_name FROM parental_status_types ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get parental status types: %w", err)
	}
	defer rows.Close()

	var statuses []ParentalStatusType
	for rows.Next() {
		var s ParentalStatusType
		if err := rows.Scan(&s.ID, &s.StatusName); err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}

	return statuses, nil
}

func (r *Repository) GetEnrollmentReasons(ctx context.Context) ([]EnrollmentReason, error) {
	query := `SELECT id, reason_text FROM enrollment_reasons ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment reasons: %w", err)
	}
	defer rows.Close()

	var reasons []EnrollmentReason
	for rows.Next() {
		var er EnrollmentReason
		if err := rows.Scan(&er.ID, &er.Text); err != nil {
			return nil, err
		}
		reasons = append(reasons, er)
	}

	return reasons, nil
}

func (r *Repository) GetIncomeRanges(ctx context.Context) ([]IncomeRange, error) {
	query := `SELECT id, range_text FROM income_ranges ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get income ranges: %w", err)
	}
	defer rows.Close()

	var ranges []IncomeRange
	for rows.Next() {
		var ir IncomeRange
		if err := rows.Scan(&ir.ID, &ir.RangeText); err != nil {
			return nil, err
		}
		ranges = append(ranges, ir)
	}

	return ranges, nil
}

func (r *Repository) GetStudentSupportTypes(ctx context.Context) ([]StudentSupportType, error) {
	query := `SELECT id, support_type_name FROM student_support_types ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get student support types: %w", err)
	}
	defer rows.Close()

	var supportTypes []StudentSupportType
	for rows.Next() {
		var s StudentSupportType
		if err := rows.Scan(&s.ID, &s.SupportTypeName); err != nil {
			return nil, err
		}
		supportTypes = append(supportTypes, s)
	}

	return supportTypes, nil
}

func (r *Repository) GetSiblingSupportTypes(ctx context.Context) ([]SibilingSupportType, error) {
	query := `SELECT id, name FROM sibling_support_types ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get sibling support types: %w", err)
	}
	defer rows.Close()

	var supportTypes []SibilingSupportType
	for rows.Next() {
		var s SibilingSupportType
		if err := rows.Scan(&s.ID, &s.SupportName); err != nil {
			return nil, err
		}
		supportTypes = append(supportTypes, s)
	}

	return supportTypes, nil
}

func (r *Repository) GetEducationalLevels(ctx context.Context) ([]EducationalLevel, error) {
	query := `SELECT id, level_name FROM educational_levels ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get educational levels: %w", err)
	}
	defer rows.Close()

	var levels []EducationalLevel
	for rows.Next() {
		var el EducationalLevel
		if err := rows.Scan(&el.ID, &el.LevelName); err != nil {
			return nil, err
		}
		levels = append(levels, el)
	}

	return levels, nil
}

func (r *Repository) GetCourses(ctx context.Context) ([]Course, error) {
	query := `SELECT id, code, course_name FROM courses ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}
	defer rows.Close()

	var courses []Course
	for rows.Next() {
		var c Course
		if err := rows.Scan(&c.ID, &c.Code, &c.CourseName); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}

	return courses, nil
}

func (r *Repository) GetCivilStatusTypes(ctx context.Context) ([]CivilStatusType, error) {
	query := `SELECT id, status_name FROM civil_status_types ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get civil status types: %w", err)
	}
	defer rows.Close()

	var statuses []CivilStatusType
	for rows.Next() {
		var s CivilStatusType
		if err := rows.Scan(&s.ID, &s.StatusName); err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}

	return statuses, nil
}

func (r *Repository) GetReligions(ctx context.Context) ([]Religion, error) {
	query := `SELECT id, religion_name FROM religions ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get religions: %w", err)
	}
	defer rows.Close()

	var religions []Religion
	for rows.Next() {
		var r Religion
		if err := rows.Scan(&r.ID, &r.ReligionName); err != nil {
			return nil, err
		}
		religions = append(religions, r)
	}

	return religions, nil
}

func (r *Repository) GetStudentRelationshipTypes(ctx context.Context) ([]StudentRelationshipType, error) {
	query := `SELECT id, relationship_name FROM student_relationship_types ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get student relationship types: %w", err)
	}
	defer rows.Close()

	var relationships []StudentRelationshipType
	for rows.Next() {
		var s StudentRelationshipType
		if err := rows.Scan(&s.ID, &s.RelationshipName); err != nil {
			return nil, err
		}
		relationships = append(relationships, s)
	}

	return relationships, nil
}

func (r *Repository) GetNatureOfResidenceTypes(ctx context.Context) ([]NatureOfResidenceType, error) {
	query := `SELECT id, residence_type_name FROM nature_of_residence_types ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get nature of residence types: %w", err)
	}
	defer rows.Close()

	var residences []NatureOfResidenceType
	for rows.Next() {
		var n NatureOfResidenceType
		if err := rows.Scan(&n.ID, &n.ResidenceTypeName); err != nil {
			return nil, err
		}
		residences = append(residences, n)
	}

	return residences, nil
}

// Retrieve - Count
func (r *Repository) GetTotalStudentsCount(
	ctx context.Context,
	search string,
	courseID int,
	genderID int,
	yearLevel int,
) (int, error) {
	// 1. Base query
	sql := `SELECT COUNT(iir.id) FROM iir_records iir
            JOIN users u ON iir.user_id = u.id
            JOIN student_personal_info spi ON iir.id = spi.iir_id
            WHERE iir.is_submitted = TRUE`

	var args []interface{}

	if courseID > 0 {
		sql += " AND spi.course_id = ?"
		args = append(args, courseID)
	}

	if genderID > 0 {
		sql += " AND spi.gender_id = ?"
		args = append(args, genderID)
	}

	if yearLevel > 0 {
		sql += " AND spi.year_level = ?"
		args = append(args, yearLevel)
	}

	if search != "" {
		sql += ` AND (u.first_name LIKE ?
                 OR u.last_name LIKE ?
                 OR u.email LIKE ?
                 OR spi.student_number LIKE ?)`

		pattern := "%" + search + "%"
		args = append(args, pattern, pattern, pattern, pattern)
	}

	var total int
	err := r.db.QueryRowContext(ctx, sql, args...).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// Retrieve - List
func (r *Repository) ListStudents(
	ctx context.Context, search string, offset int, limit int, orderBy string,
	courseID int, genderID int, yearLevel int,
) ([]StudentProfileView, error) {
	query := `
        SELECT
			iir.id as iir_id,
			usr.id as user_id,
			usr.first_name,
			usr.middle_name,
			usr.last_name,
			spi.gender_id,
			usr.email,
			spi.student_number,
			spi.course_id,
			spi.section,
			spi.year_level
		FROM iir_records iir
		JOIN users usr ON iir.user_id = usr.id
		JOIN student_personal_info spi ON iir.id = spi.iir_id
		WHERE iir.is_submitted = TRUE
		AND (usr.first_name LIKE ? OR usr.middle_name LIKE ? OR usr.last_name LIKE ? OR spi.student_number LIKE ? OR usr.email LIKE ?)
    `

	var args []interface{}
	args = append(args, "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")

	if courseID != 0 {
		query += " AND spi.course_id = ?"
		args = append(args, courseID)
	}

	if genderID != 0 {
		query += " AND spi.gender_id = ?"
		args = append(args, genderID)
	}

	if yearLevel != 0 {
		query += " AND spi.year_level = ?"
		args = append(args, yearLevel)
	}

	allowedSortColumns := map[string]string{
		"last_name":      "usr.last_name",
		"first_name":     "usr.first_name",
		"year_level":     "spi.year_level",
		"course_id":      "spi.course_id",
		"created_at":     "iir.created_at",
		"updated_at":     "iir.updated_at",
		"student_number": "spi.student_number",
		"iir_id":         "iir.id",
	}

	sortColumn, ok := allowedSortColumns[orderBy]
	if !ok {
		sortColumn = allowedSortColumns["last_name"]
	}

	query += fmt.Sprintf(" ORDER BY %s ASC, iir.id ASC", sortColumn)

	query += `
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
			&student.IIRID,
			&student.UserID,
			&student.FirstName,
			&student.MiddleName,
			&student.LastName,
			&student.GenderID,
			&student.Email,
			&student.StudentNumber,
			&student.CourseID,
			&student.Section,
			&student.YearLevel,
		); err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return students, nil
}

func (r *Repository) GetStudentBasicInfo(ctx context.Context, iirID int) (*StudentBasicInfoView, error) {
	query := `
		SELECT u.id, u.first_name, u.middle_name, u.last_name, u.email
		FROM users u
		JOIN iir_records iir ON u.id = iir.user_id
		WHERE iir.id = ?
	`

	var info StudentBasicInfoView
	err := r.db.QueryRowContext(ctx, query, iirID).Scan(
		&info.ID,
		&info.FirstName,
		&info.MiddleName,
		&info.LastName,
		&info.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get student basic info: %w", err)
	}

	return &info, nil
}

func (r *Repository) GetStudentIIRByUserID(ctx context.Context, userID int) (*IIRRecord, error) {
	query := `
		SELECT id, user_id, is_submitted, created_at, updated_at
		FROM iir_records
		WHERE user_id = ?
		LIMIT 1
	`

	var iir IIRRecord
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&iir.ID,
		&iir.UserID,
		&iir.IsSubmitted,
		&iir.CreatedAt,
		&iir.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &iir, nil
}

func (r *Repository) GetStudentIIR(ctx context.Context, iirID int) (*IIRRecord, error) {
	query := `
		SELECT id, user_id, is_submitted, created_at, updated_at
		FROM iir_records
		WHERE id = ?
	`

	var iir IIRRecord
	err := r.db.QueryRowContext(ctx, query, iirID).Scan(
		&iir.ID,
		&iir.UserID,
		&iir.IsSubmitted,
		&iir.CreatedAt,
		&iir.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &iir, nil
}

func (r *Repository) GetStudentEnrollmentReasons(ctx context.Context, iirID int) ([]StudentSelectedReason, error) {
	query := `
		SELECT iir_id, reason_id, other_reason_text
		FROM student_selected_reasons
		WHERE iir_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student enrollment reasons: %w", err)
	}
	defer rows.Close()

	var reasons []StudentSelectedReason
	for rows.Next() {
		var sr StudentSelectedReason
		if err := rows.Scan(&sr.IIRID, &sr.ReasonID, &sr.OtherReasonText); err != nil {
			return nil, err
		}
		reasons = append(reasons, sr)
	}

	return reasons, nil
}

func (r *Repository) GetEnrollmentReasonByID(ctx context.Context, reasonID int) (*EnrollmentReason, error) {
	query := `SELECT id, reason_text FROM enrollment_reasons WHERE id = ?`
	var er EnrollmentReason
	err := r.db.QueryRowContext(ctx, query, reasonID).Scan(&er.ID, &er.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment reason by ID: %w", err)
	}

	return &er, nil
}

func (r *Repository) GetStudentPersonalInfo(ctx context.Context, iirID int) (*StudentPersonalInfo, error) {
	query := `
		SELECT
			id, iir_id, student_number, gender_id, civil_status_id,
			religion_id, height_ft, weight_kg, complexion,
			high_school_gwa, course_id, year_level, section,
			place_of_birth, date_of_birth, is_employed,
			employer_name, employer_address, mobile_number, telephone_number,
			emergency_contact_name, emergency_contact_number, emergency_contact_relationship_id,
			emergency_contact_address_id,
			created_at, updated_at
		FROM student_personal_info
		WHERE iir_id = ?
	`

	var info StudentPersonalInfo
	err := r.db.QueryRowContext(ctx, query, iirID).Scan(
		&info.ID,
		&info.IIRID,
		&info.StudentNumber,
		&info.GenderID,
		&info.CivilStatusID,
		&info.ReligionID,
		&info.HeightFt,
		&info.WeightKg,
		&info.Complexion,
		&info.HighSchoolGWA,
		&info.CourseID,
		&info.YearLevel,
		&info.Section,
		&info.PlaceOfBirth,
		&info.DateOfBirth,
		&info.IsEmployed,
		&info.EmployerName,
		&info.EmployerAddress,
		&info.MobileNumber,
		&info.TelephoneNumber,
		&info.EmergencyContactName,
		&info.EmergencyContactNumber,
		&info.EmergencyContactRelationshipID,
		&info.EmergencyContactAddressID,
		&info.CreatedAt,
		&info.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (r *Repository) GetGenderByID(ctx context.Context, genderID int) (*Gender, error) {
	query := `SELECT id, gender_name FROM genders WHERE id = ?`
	var gender Gender
	err := r.db.QueryRowContext(ctx, query, genderID).Scan(&gender.ID, &gender.GenderName)
	if err != nil {
		return nil, fmt.Errorf("failed to get gender by ID: %w", err)
	}

	return &gender, nil
}

func (r *Repository) GetCivilStatusByID(ctx context.Context, statusID int) (*CivilStatusType, error) {
	query := `SELECT id, status_name FROM civil_status_types WHERE id = ?`
	var status CivilStatusType
	err := r.db.QueryRowContext(ctx, query, statusID).Scan(&status.ID, &status.StatusName)
	if err != nil {
		return nil, fmt.Errorf("failed to get civil status by ID: %w", err)
	}

	return &status, nil
}

func (r *Repository) GetReligionByID(ctx context.Context, religionID int) (*Religion, error) {
	query := `SELECT id, religion_name FROM religions WHERE id = ?`
	var religion Religion
	err := r.db.QueryRowContext(ctx, query, religionID).Scan(&religion.ID, &religion.ReligionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get religion by ID: %w", err)
	}

	return &religion, nil
}

func (r *Repository) GetCourseByID(ctx context.Context, courseID int) (*Course, error) {
	query := `SELECT id, code, course_name FROM courses WHERE id = ?`
	var course Course
	err := r.db.QueryRowContext(ctx, query, courseID).Scan(&course.ID, &course.Code, &course.CourseName)
	if err != nil {
		return nil, fmt.Errorf("failed to get course by ID: %w", err)
	}

	return &course, nil
}

func (r *Repository) GetStudentAddresses(ctx context.Context, iirID int) ([]StudentAddress, error) {
	query := `
		SELECT id, iir_id, address_id, address_type, created_at, updated_at
		FROM student_addresses
		WHERE iir_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student addresses: %w", err)
	}
	defer rows.Close()

	var addresses []StudentAddress
	for rows.Next() {
		var sa StudentAddress
		if err := rows.Scan(
			&sa.ID,
			&sa.IIRID,
			&sa.AddressID,
			&sa.AddressType,
			&sa.CreatedAt,
			&sa.UpdatedAt); err != nil {
			return nil, err
		}
		addresses = append(addresses, sa)
	}

	return addresses, nil
}

func (r *Repository) GetAddressByID(ctx context.Context, addressID int) (*Address, error) {
	query := `
		SELECT
			id, region, city, barangay, street_detail,
			created_at, updated_at
		FROM addresses
		WHERE id = ?
	`
	var addr Address
	err := r.db.QueryRowContext(ctx, query, addressID).Scan(
		&addr.ID,
		&addr.Region,
		&addr.City,
		&addr.Barangay,
		&addr.StreetDetail,
		&addr.CreatedAt,
		&addr.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get address by ID: %w", err)
	}

	return &addr, nil
}

func (r *Repository) GetStudentEducationalBackground(ctx context.Context, iirID int) (*EducationalBackground, error) {
	query := `
		SELECT
			id, iir_id, nature_of_schooling, interrupted_details,
			created_at, updated_at
		FROM educational_backgrounds
		WHERE iir_id = ?
	`
	var eb EducationalBackground
	err := r.db.QueryRowContext(ctx, query, iirID).Scan(
		&eb.ID,
		&eb.IIRID,
		&eb.NatureOfSchooling,
		&eb.InterruptedDetails,
		&eb.CreatedAt,
		&eb.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get student educational background: %w", err)
	}

	return &eb, nil
}

func (r *Repository) GetSchoolDetailsByEBID(ctx context.Context, ebID int) ([]SchoolDetails, error) {
	query := `
		SELECT
			id, eb_id, educational_level_id, school_name,
			school_address, school_type, year_started,
			year_completed, awards, created_at, updated_at
		FROM school_details
		WHERE eb_id = ?
		ORDER BY educational_level_id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, ebID)
	if err != nil {
		return nil, fmt.Errorf("failed to get school details by EB ID: %w", err)
	}
	defer rows.Close()

	var schoolDetails []SchoolDetails
	for rows.Next() {
		var sd SchoolDetails
		if err := rows.Scan(
			&sd.ID,
			&sd.EBID,
			&sd.EducationalLevelID,
			&sd.SchoolName,
			&sd.SchoolAddress,
			&sd.SchoolType,
			&sd.YearStarted,
			&sd.YearCompleted,
			&sd.Awards,
			&sd.CreatedAt,
			&sd.UpdatedAt); err != nil {
			return nil, err
		}
		schoolDetails = append(schoolDetails, sd)
	}

	return schoolDetails, nil
}

func (r *Repository) GetEducationalLevelByID(ctx context.Context, levelID int) (*EducationalLevel, error) {
	query := `SELECT id, level_name FROM educational_levels WHERE id = ?`
	var el EducationalLevel
	err := r.db.QueryRowContext(ctx, query, levelID).Scan(&el.ID, &el.LevelName)
	if err != nil {
		return nil, fmt.Errorf("failed to get educational level by ID: %w", err)
	}

	return &el, nil
}

func (r *Repository) GetStudentRelatedPersons(
	ctx context.Context, iirID int,
) ([]StudentRelatedPerson, error) {
	query := `
		SELECT iir_id, related_person_id, relationship_id,
			is_parent, is_guardian, is_living,
			created_at, updated_at
		FROM student_related_persons
		WHERE iir_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student related persons: %w", err)
	}
	defer rows.Close()

	var persons []StudentRelatedPerson
	for rows.Next() {
		var srp StudentRelatedPerson
		if err := rows.Scan(
			&srp.IIRID,
			&srp.RelatedPersonID,
			&srp.RelationshipID,
			&srp.IsParent,
			&srp.IsGuardian,
			&srp.IsLiving,
			&srp.CreatedAt,
			&srp.UpdatedAt); err != nil {
			return nil, err
		}
		persons = append(persons, srp)
	}

	return persons, nil
}

func (r *Repository) GetRelatedPersonByID(
	ctx context.Context, personID int,
) (*RelatedPerson, error) {
	query := `
		SELECT
			id, first_name, middle_name, last_name, address_id,
			educational_level, date_of_birth, occupation,
			employer_name, employer_address, contact_number,
			created_at, updated_at
		FROM related_persons
		WHERE id = ?
	`
	var person RelatedPerson
	err := r.db.QueryRowContext(ctx, query, personID).Scan(
		&person.ID,
		&person.FirstName,
		&person.MiddleName,
		&person.LastName,
		&person.AddressID,
		&person.EducationalLevel,
		&person.DateOfBirth,
		&person.Occupation,
		&person.EmployerName,
		&person.EmployerAddress,
		&person.ContactNumber,
		&person.CreatedAt,
		&person.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get related person by ID: %w", err)
	}

	return &person, nil
}

func (r *Repository) GetStudentRelationshipByID(
	ctx context.Context, relationshipID int,
) (*StudentRelationshipType, error) {
	query := `SELECT id, relationship_name FROM student_relationship_types WHERE id = ?`
	var srt StudentRelationshipType
	err := r.db.QueryRowContext(ctx, query, relationshipID).Scan(&srt.ID, &srt.RelationshipName)
	if err != nil {
		return nil, fmt.Errorf("failed to get student relationship by ID: %w", err)
	}

	return &srt, nil
}

func (r *Repository) GetStudentFamilyBackground(ctx context.Context, iirID int) (*FamilyBackground, error) {
	query := `
		SELECT
			id, iir_id, parental_status_id, parental_status_details, brothers,
			sisters, employed_siblings, ordinal_position,
			have_quiet_place_to_study, is_sharing_room,
			nature_of_residence_id, created_at, updated_at
		FROM family_backgrounds
		WHERE iir_id = ?
	`

	var fb FamilyBackground
	err := r.db.QueryRowContext(ctx, query, iirID).Scan(
		&fb.ID,
		&fb.IIRID,
		&fb.ParentalStatusID,
		&fb.ParentalStatusDetails,
		&fb.Brothers,
		&fb.Sisters,
		&fb.EmployedSiblings,
		&fb.OrdinalPosition,
		&fb.HaveQuietPlaceToStudy,
		&fb.IsSharingRoom,
		&fb.NatureOfResidenceId,
		&fb.CreatedAt,
		&fb.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get student family background: %w", err)
	}

	return &fb, nil
}

func (r *Repository) GetParentalStatusByID(ctx context.Context, statusID int) (*ParentalStatusType, error) {
	query := `SELECT id, status_name FROM parental_status_types WHERE id = ?`
	var ps ParentalStatusType
	err := r.db.QueryRowContext(ctx, query, statusID).Scan(&ps.ID, &ps.StatusName)
	if err != nil {
		return nil, fmt.Errorf("failed to get parental status by ID: %w", err)
	}

	return &ps, nil
}

func (r *Repository) GetNatureOfResidenceByID(ctx context.Context, residenceID int) (*NatureOfResidenceType, error) {
	query := `SELECT id, residence_type_name FROM nature_of_residence_types WHERE id = ?`
	var nr NatureOfResidenceType
	err := r.db.QueryRowContext(ctx, query, residenceID).Scan(&nr.ID, &nr.ResidenceTypeName)
	if err != nil {
		return nil, fmt.Errorf("failed to get nature of residence by ID: %w", err)
	}

	return &nr, nil
}

func (r *Repository) GetStudentSiblingSupport(ctx context.Context, fbID int) ([]StudentSiblingSupport, error) {
	query := `
		SELECT
			family_background_id, support_type_id
		FROM student_sibling_supports
		WHERE family_background_id = ?
	`
	var sss []StudentSiblingSupport
	rows, err := r.db.QueryContext(ctx, query, fbID)
	if err != nil {
		return nil, fmt.Errorf("failed to query student sibling supports: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ss StudentSiblingSupport
		err := rows.Scan(
			&ss.FamilyBackgroundID,
			&ss.SupportTypeID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan student sibling support: %w", err)
		}
		sss = append(sss, ss)
	}

	return sss, nil
}

func (r *Repository) GetSiblingSupportTypeByID(ctx context.Context, supportID int) (*SibilingSupportType, error) {
	query := `SELECT id, name FROM sibling_support_types WHERE id = ?`
	var sst SibilingSupportType
	err := r.db.QueryRowContext(ctx, query, supportID).Scan(&sst.ID, &sst.SupportName)
	if err != nil {
		return nil, fmt.Errorf("failed to get sibling support type by ID: %w", err)
	}

	return &sst, nil
}

func (r *Repository) GetStudentFinancialInfo(ctx context.Context, iirID int) (*StudentFinance, error) {
	query := `
		SELECT
			id, iir_id, monthly_family_income_range_id, other_income_details,
			weekly_allowance, created_at, updated_at
		FROM student_finances
		WHERE iir_id = ?
	`
	var fi StudentFinance
	err := r.db.QueryRowContext(ctx, query, iirID).Scan(
		&fi.ID,
		&fi.IIRID,
		&fi.MonthlyFamilyIncomeRangeID,
		&fi.OtherIncomeDetails,
		&fi.WeeklyAllowance,
		&fi.CreatedAt,
		&fi.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get student financial info: %w", err)
	}

	return &fi, nil
}

func (r *Repository) GetFinancialSupportTypeByFinanceID(ctx context.Context, financeID int) ([]StudentFinancialSupport, error) {
	query := `
		SELECT
			sf_id, support_type_id, created_at, updated_at
		FROM student_financial_supports
		WHERE sf_id = ?
	`
	var sfs []StudentFinancialSupport
	rows, err := r.db.QueryContext(ctx, query, financeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query student financial supports: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var sf StudentFinancialSupport
		err := rows.Scan(
			&sf.StudentFinanceID,
			&sf.SupportTypeID,
			&sf.CreatedAt,
			&sf.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan student financial support: %w", err)
		}
		sfs = append(sfs, sf)
	}

	return sfs, nil
}

func (r *Repository) GetIncomeRangeByID(ctx context.Context, rangeID int) (*IncomeRange, error) {
	query := `SELECT id, range_text FROM income_ranges WHERE id = ?`
	var ir IncomeRange
	err := r.db.QueryRowContext(ctx, query, rangeID).Scan(&ir.ID, &ir.RangeText)
	if err != nil {
		return nil, fmt.Errorf("failed to get income range by ID: %w", err)
	}

	return &ir, nil
}

func (r *Repository) GetStudentSupportByID(ctx context.Context, supportID int) (*StudentSupportType, error) {
	query := `SELECT id, support_type_name FROM student_support_types WHERE id = ?`
	var sst StudentSupportType
	err := r.db.QueryRowContext(ctx, query, supportID).Scan(&sst.ID, &sst.SupportTypeName)
	if err != nil {
		return nil, fmt.Errorf("failed to get student support by ID: %w", err)
	}

	return &sst, nil
}

func (r *Repository) GetStudentHealthRecord(ctx context.Context, iirID int) (*StudentHealthRecord, error) {
	query := `
		SELECT
			id, iir_id, vision_has_problem, vision_details, hearing_has_problem,
			hearing_details, speech_has_problem, speech_details,
			general_health_has_problem, general_health_details,
			created_at, updated_at
		FROM student_health_records
		WHERE iir_id = ?
	`
	var hr StudentHealthRecord
	err := r.db.QueryRowContext(ctx, query, iirID).Scan(
		&hr.ID,
		&hr.IIRID,
		&hr.VisionHasProblem,
		&hr.VisionDetails,
		&hr.HearingHasProblem,
		&hr.HearingDetails,
		&hr.SpeechHasProblem,
		&hr.SpeechDetails,
		&hr.GeneralHealthHasProblem,
		&hr.GeneralHealthDetails,
		&hr.CreatedAt,
		&hr.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get student health record: %w", err)
	}

	return &hr, nil
}

func (r *Repository) GetStudentConsultations(ctx context.Context, iirID int) ([]StudentConsultation, error) {
	query := `
		SELECT
			id, iir_id, professional_type,
			has_consulted, when_date, for_what,
			created_at, updated_at
		FROM student_consultations
		WHERE iir_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student consultations: %w", err)
	}
	defer rows.Close()

	var consultations []StudentConsultation
	for rows.Next() {
		var sc StudentConsultation
		if err := rows.Scan(
			&sc.ID,
			&sc.IIRID,
			&sc.ProfessionalType,
			&sc.HasConsulted,
			&sc.WhenDate,
			&sc.ForWhat,
			&sc.CreatedAt,
			&sc.UpdatedAt,
		); err != nil {
			return nil, err
		}
		consultations = append(consultations, sc)
	}

	return consultations, nil
}

func (r *Repository) GetStudentActivities(ctx context.Context, iirID int) ([]StudentActivity, error) {
	query := `
		SELECT
			id, iir_id, option_id, other_specification,
			role, role_specification, created_at, updated_at
		FROM student_activities
		WHERE iir_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student activities: %w", err)
	}
	defer rows.Close()

	var activities []StudentActivity
	for rows.Next() {
		var sa StudentActivity
		if err := rows.Scan(
			&sa.ID,
			&sa.IIRID,
			&sa.OptionID,
			&sa.OtherSpecification,
			&sa.Role,
			&sa.RoleSpecification,
			&sa.CreatedAt,
			&sa.UpdatedAt); err != nil {
			return nil, err
		}
		activities = append(activities, sa)
	}

	return activities, nil
}

func (r *Repository) GetActivityOptionByID(ctx context.Context, optionID int) (*ActivityOption, error) {
	query := `SELECT id, name, category, is_active FROM activity_options WHERE id = ?`
	var ao ActivityOption
	err := r.db.QueryRowContext(ctx, query, optionID).Scan(
		&ao.ID,
		&ao.Name,
		&ao.Category,
		&ao.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity option by ID: %w", err)
	}

	return &ao, nil
}

func (r *Repository) GetStudentSubjectPreferences(ctx context.Context, iirID int) ([]StudentSubjectPreference, error) {
	query := `
		SELECT
			id, iir_id, subject_name, is_favorite, created_at, updated_at
		FROM student_subject_preferences
		WHERE iir_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student subject preferences: %w", err)
	}
	defer rows.Close()

	var preferences []StudentSubjectPreference
	for rows.Next() {
		var sp StudentSubjectPreference
		if err := rows.Scan(
			&sp.ID,
			&sp.IIRID,
			&sp.SubjectName,
			&sp.IsFavorite,
			&sp.CreatedAt,
			&sp.UpdatedAt); err != nil {
			return nil, err
		}
		preferences = append(preferences, sp)
	}

	return preferences, nil
}

func (r *Repository) GetStudentHobbies(ctx context.Context, iirID int) ([]StudentHobby, error) {
	query := `
		SELECT
			id, iir_id, hobby_name, priority_rank, created_at, updated_at
		FROM student_hobbies
		WHERE iir_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student hobbies: %w", err)
	}
	defer rows.Close()

	var hobbies []StudentHobby
	for rows.Next() {
		var sh StudentHobby
		if err := rows.Scan(
			&sh.ID,
			&sh.IIRID,
			&sh.HobbyName,
			&sh.PriorityRank,
			&sh.CreatedAt,
			&sh.UpdatedAt); err != nil {
			return nil, err
		}
		hobbies = append(hobbies, sh)
	}

	return hobbies, nil
}

func (r *Repository) GetStudentTestResults(ctx context.Context, iirID int) ([]TestResult, error) {
	query := `
		SELECT
			id, iir_id, test_date, test_name,
			raw_score, percentile, description,
			created_at, updated_at
		FROM test_results
		WHERE iir_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student test results: %w", err)
	}
	defer rows.Close()

	var results []TestResult
	for rows.Next() {
		var tr TestResult
		if err := rows.Scan(
			&tr.ID,
			&tr.IIRID,
			&tr.TestDate,
			&tr.TestName,
			&tr.RawScore,
			&tr.Percentile,
			&tr.Description,
			&tr.CreatedAt,
			&tr.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, tr)
	}

	return results, nil
}

func (r *Repository) GetStudentSignificantNotes(ctx context.Context, iirID int) ([]SignificantNote, error) {
	query := `
		SELECT
			id, iir_id, note_date, incident_description, remarks,
			created_at, updated_at
		FROM significant_notes
		WHERE iir_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student significant notes: %w", err)
	}
	defer rows.Close()

	var notes []SignificantNote
	for rows.Next() {
		var sn SignificantNote
		if err := rows.Scan(
			&sn.ID,
			&sn.IIRID,
			&sn.NoteDate,
			&sn.IncidentDescription,
			&sn.Remarks,
			&sn.CreatedAt,
			&sn.UpdatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, sn)
	}

	return notes, nil
}
