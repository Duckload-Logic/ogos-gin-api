package students

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Lookup
func (r *Repository) GetGenders(ctx context.Context) ([]Gender, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM genders ORDER BY id
	`, database.GetColumns(Gender{}))

	var genders []Gender
	err := r.db.SelectContext(ctx, &genders, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get genders: %w", err)
	}

	return genders, nil
}

func (r *Repository) GetParentalStatusTypes(ctx context.Context) ([]ParentalStatusType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM parental_status_types ORDER BY id
	`, database.GetColumns(ParentalStatusType{}))
	var statuses []ParentalStatusType
	err := r.db.SelectContext(ctx, &statuses, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get parental status types: %w", err)
	}

	return statuses, nil
}

func (r *Repository) GetEnrollmentReasons(ctx context.Context) ([]EnrollmentReason, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM enrollment_reasons ORDER BY id
	`, database.GetColumns(EnrollmentReason{}))
	var reasons []EnrollmentReason
	err := r.db.SelectContext(ctx, &reasons, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment reasons: %w", err)
	}

	return reasons, nil
}

func (r *Repository) GetIncomeRanges(ctx context.Context) ([]IncomeRange, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM income_ranges ORDER BY id
	`, database.GetColumns(IncomeRange{}))
	var ranges []IncomeRange
	err := r.db.SelectContext(ctx, &ranges, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get income ranges: %w", err)
	}

	return ranges, nil
}

func (r *Repository) GetStudentSupportTypes(ctx context.Context) ([]StudentSupportType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_support_types ORDER BY id
	`, database.GetColumns(StudentSupportType{}))
	var supportTypes []StudentSupportType
	err := r.db.SelectContext(ctx, &supportTypes, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get student support types: %w", err)
	}

	return supportTypes, nil
}

func (r *Repository) GetSiblingSupportTypes(ctx context.Context) ([]SibilingSupportType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM sibling_support_types ORDER BY id
	`, database.GetColumns(SibilingSupportType{}))

	var supportTypes []SibilingSupportType
	err := r.db.SelectContext(ctx, &supportTypes, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get sibling support types: %w", err)
	}

	return supportTypes, nil
}

func (r *Repository) GetEducationalLevels(ctx context.Context) ([]EducationalLevel, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM educational_levels ORDER BY id
	`, database.GetColumns(EducationalLevel{}))

	var levels []EducationalLevel
	err := r.db.SelectContext(ctx, &levels, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get educational levels: %w", err)
	}

	return levels, nil
}

func (r *Repository) GetCourses(ctx context.Context) ([]Course, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM courses ORDER BY id
	`, database.GetColumns(Course{}))

	var courses []Course
	err := r.db.SelectContext(ctx, &courses, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}

	return courses, nil
}

func (r *Repository) GetCivilStatusTypes(ctx context.Context) ([]CivilStatusType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM civil_status_types ORDER BY id
	`, database.GetColumns(CivilStatusType{}))

	var statuses []CivilStatusType
	err := r.db.SelectContext(ctx, &statuses, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get civil status types: %w", err)
	}

	return statuses, nil
}

func (r *Repository) GetReligions(ctx context.Context) ([]Religion, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM religions ORDER BY id
	`, database.GetColumns(Religion{}))

	var religions []Religion
	err := r.db.SelectContext(ctx, &religions, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get religions: %w", err)
	}

	return religions, nil
}

func (r *Repository) GetStudentRelationshipTypes(ctx context.Context) ([]StudentRelationshipType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_relationship_types ORDER BY id
	`, database.GetColumns(StudentRelationshipType{}))

	var relationships []StudentRelationshipType
	err := r.db.SelectContext(ctx, &relationships, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get student relationship types: %w", err)
	}

	return relationships, nil
}

func (r *Repository) GetNatureOfResidenceTypes(ctx context.Context) ([]NatureOfResidenceType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM nature_of_residence_types ORDER BY id
	`, database.GetColumns(NatureOfResidenceType{}))

	var residences []NatureOfResidenceType
	err := r.db.SelectContext(ctx, &residences, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get nature of residence types: %w", err)
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
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_personal_info
		WHERE iir_id = ?
	`, database.GetColumns(StudentPersonalInfo{}))

	var info StudentPersonalInfo
	err := r.db.GetContext(ctx, &info, query, iirID)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (r *Repository) GetEmergencyContactByIIRID(ctx context.Context, iirID int) (*EmergencyContact, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM emergency_contacts ec
		JOIN student_personal_info spi ON ec.id = spi.emergency_contact_id
		WHERE spi.iir_id = ?
	`, database.GetColumns(EmergencyContact{}))

	var ec EmergencyContact
	err := r.db.GetContext(ctx, &ec, query, iirID)
	if err != nil {
		return nil, err
	}

	return &ec, nil
}

func (r *Repository) GetGenderByID(ctx context.Context, genderID int) (*Gender, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM genders WHERE id = ?
	`, database.GetColumns(Gender{}))

	var gender Gender
	err := r.db.GetContext(ctx, &gender, query, genderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gender by ID: %w", err)
	}

	return &gender, nil
}

func (r *Repository) GetCivilStatusByID(ctx context.Context, statusID int) (*CivilStatusType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM civil_status_types WHERE id = ?
	`, database.GetColumns(CivilStatusType{}))

	var status CivilStatusType
	err := r.db.GetContext(ctx, &status, query, statusID)
	if err != nil {
		return nil, fmt.Errorf("failed to get civil status by ID: %w", err)
	}

	return &status, nil
}

func (r *Repository) GetReligionByID(ctx context.Context, religionID int) (*Religion, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM religions WHERE id = ?
	`, database.GetColumns(Religion{}))

	var religion Religion
	err := r.db.GetContext(ctx, &religion, query, religionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get religion by ID: %w", err)
	}

	return &religion, nil
}

func (r *Repository) GetCourseByID(ctx context.Context, courseID int) (*Course, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM courses WHERE id = ?
	`, database.GetColumns(Course{}))

	var course Course
	err := r.db.GetContext(ctx, &course, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course by ID: %w", err)
	}

	return &course, nil
}
func (r *Repository) GetStudentAddresses(ctx context.Context, iirID int) ([]StudentAddress, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_addresses
		WHERE iir_id = ?
	`, database.GetColumns(StudentAddress{}))

	var addresses []StudentAddress
	err := r.db.SelectContext(ctx, &addresses, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student addresses: %w", err)
	}

	return addresses, nil
}

func (r *Repository) GetAddressByID(ctx context.Context, addressID int) (*Address, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM addresses
		WHERE id = ?
	`, database.GetColumns(Address{}))

	var addr Address
	err := r.db.GetContext(ctx, &addr, query, addressID)
	if err != nil {
		return nil, fmt.Errorf("failed to get address by ID: %w", err)
	}

	return &addr, nil
}

func (r *Repository) GetStudentEducationalBackground(ctx context.Context, iirID int) (*EducationalBackground, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM educational_backgrounds
		WHERE iir_id = ?
	`, database.GetColumns(EducationalBackground{}))

	var eb EducationalBackground
	err := r.db.GetContext(ctx, &eb, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student educational background: %w", err)
	}

	return &eb, nil
}

func (r *Repository) GetSchoolDetailsByEBID(ctx context.Context, ebID int) ([]SchoolDetails, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM school_details
		WHERE eb_id = ?
		ORDER BY educational_level_id ASC
	`, database.GetColumns(SchoolDetails{}))

	var schoolDetails []SchoolDetails
	err := r.db.SelectContext(ctx, &schoolDetails, query, ebID)
	if err != nil {
		return nil, fmt.Errorf("failed to get school details by EB ID: %w", err)
	}

	return schoolDetails, nil
}

func (r *Repository) GetEducationalLevelByID(ctx context.Context, levelID int) (*EducationalLevel, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM educational_levels WHERE id = ?
	`, database.GetColumns(EducationalLevel{}))
	var el EducationalLevel
	err := r.db.GetContext(ctx, &el, query, levelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get educational level by ID: %w", err)
	}

	return &el, nil
}

func (r *Repository) GetStudentRelatedPersons(
	ctx context.Context, iirID int,
) ([]StudentRelatedPerson, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_related_persons WHERE iir_id = ?
	`, database.GetColumns(StudentRelatedPerson{}))

	var persons []StudentRelatedPerson

	err := r.db.SelectContext(ctx, &persons, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student related persons: %w", err)
	}

	return persons, nil
}

func (r *Repository) GetRelatedPersonByID(
	ctx context.Context, personID int,
) (*RelatedPerson, error) {
	query := fmt.Sprintf(`
		SELECT
			%s
		FROM related_persons
		WHERE id = ?
	`, database.GetColumns(RelatedPerson{}))

	var person RelatedPerson
	err := r.db.GetContext(ctx, &person, query, personID)
	if err != nil {
		return nil, fmt.Errorf("failed to get related person by ID: %w", err)
	}

	return &person, nil
}

func (r *Repository) GetStudentRelationshipByID(
	ctx context.Context, relationshipID int,
) (*StudentRelationshipType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_relationship_types WHERE id = ?
	`, database.GetColumns(StudentRelationshipType{}))

	var srt StudentRelationshipType
	err := r.db.GetContext(ctx, &srt, query, relationshipID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student relationship by ID: %w", err)
	}

	return &srt, nil
}

func (r *Repository) GetStudentFamilyBackground(ctx context.Context, iirID int) (*FamilyBackground, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM family_backgrounds
		WHERE iir_id = ?
	`, database.GetColumns(FamilyBackground{}))

	var fb FamilyBackground
	err := r.db.GetContext(ctx, &fb, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student family background: %w", err)
	}

	return &fb, nil
}

func (r *Repository) GetParentalStatusByID(ctx context.Context, statusID int) (*ParentalStatusType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM parental_status_types WHERE id = ?
	`, database.GetColumns(ParentalStatusType{}))

	var ps ParentalStatusType
	err := r.db.GetContext(ctx, &ps, query, statusID)
	if err != nil {
		return nil, fmt.Errorf("failed to get parental status by ID: %w", err)
	}

	return &ps, nil
}

func (r *Repository) GetNatureOfResidenceByID(ctx context.Context, residenceID int) (*NatureOfResidenceType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM nature_of_residence_types WHERE id = ?
	`, database.GetColumns(NatureOfResidenceType{}))

	var nr NatureOfResidenceType
	err := r.db.GetContext(ctx, &nr, query, residenceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nature of residence by ID: %w", err)
	}

	return &nr, nil
}

func (r *Repository) GetStudentSiblingSupport(ctx context.Context, fbID int) ([]StudentSiblingSupport, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_sibling_supports
		WHERE family_background_id = ?
	`, database.GetColumns(StudentSiblingSupport{}))

	var sss []StudentSiblingSupport
	err := r.db.SelectContext(ctx, &sss, query, fbID)
	if err != nil {
		return nil, fmt.Errorf("failed to query student sibling supports: %w", err)
	}

	return sss, nil
}

func (r *Repository) GetSiblingSupportTypeByID(ctx context.Context, supportID int) (*SibilingSupportType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM sibling_support_types WHERE id = ?
	`, database.GetColumns(SibilingSupportType{}))

	var sst SibilingSupportType
	err := r.db.GetContext(ctx, &sst, query, supportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sibling support type by ID: %w", err)
	}

	return &sst, nil
}

func (r *Repository) GetStudentFinancialInfo(ctx context.Context, iirID int) (*StudentFinance, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_finances
		WHERE iir_id = ?
	`, database.GetColumns(StudentFinance{}))

	var fi StudentFinance
	err := r.db.GetContext(ctx, &fi, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student financial info: %w", err)
	}

	return &fi, nil
}

func (r *Repository) GetFinancialSupportTypeByFinanceID(ctx context.Context, financeID int) ([]StudentFinancialSupport, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_financial_supports
		WHERE sf_id = ?
	`, database.GetColumns(StudentFinancialSupport{}))

	var sfs []StudentFinancialSupport
	err := r.db.SelectContext(ctx, &sfs, query, financeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query student financial supports: %w", err)
	}

	return sfs, nil
}

func (r *Repository) GetIncomeRangeByID(ctx context.Context, rangeID int) (*IncomeRange, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM income_ranges WHERE id = ?
	`, database.GetColumns(IncomeRange{}))

	var ir IncomeRange
	err := r.db.GetContext(ctx, &ir, query, rangeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get income range by ID: %w", err)
	}

	return &ir, nil
}

func (r *Repository) GetStudentSupportByID(ctx context.Context, supportID int) (*StudentSupportType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_support_types WHERE id = ?
	`, database.GetColumns(StudentSupportType{}))

	var sst StudentSupportType
	err := r.db.GetContext(ctx, &sst, query, supportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student support by ID: %w", err)
	}

	return &sst, nil
}

func (r *Repository) GetStudentHealthRecord(ctx context.Context, iirID int) (*StudentHealthRecord, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_health_records
		WHERE iir_id = ?
	`, database.GetColumns(StudentHealthRecord{}))

	var hr StudentHealthRecord
	err := r.db.GetContext(ctx, &hr, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student health record: %w", err)
	}

	return &hr, nil
}

func (r *Repository) GetStudentConsultations(ctx context.Context, iirID int) ([]StudentConsultation, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_consultations
		WHERE iir_id = ?
	`, database.GetColumns(StudentConsultation{}))

	var consultations []StudentConsultation
	err := r.db.SelectContext(ctx, &consultations, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student consultations: %w", err)
	}

	return consultations, nil
}

func (r *Repository) GetStudentActivities(ctx context.Context, iirID int) ([]StudentActivity, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_activities
		WHERE iir_id = ?
	`, database.GetColumns(StudentActivity{}))

	var activities []StudentActivity
	err := r.db.SelectContext(ctx, &activities, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student activities: %w", err)
	}

	return activities, nil
}

func (r *Repository) GetActivityOptionByID(ctx context.Context, optionID int) (*ActivityOption, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM activity_options WHERE id = ?
	`, database.GetColumns(ActivityOption{}))

	var ao ActivityOption
	err := r.db.GetContext(ctx, &ao, query, optionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity option by ID: %w", err)
	}

	return &ao, nil
}

func (r *Repository) GetStudentSubjectPreferences(ctx context.Context, iirID int) ([]StudentSubjectPreference, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_subject_preferences
		WHERE iir_id = ?
	`, database.GetColumns(StudentSubjectPreference{}))

	var preferences []StudentSubjectPreference
	err := r.db.SelectContext(ctx, &preferences, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student subject preferences: %w", err)
	}

	return preferences, nil
}

func (r *Repository) GetStudentHobbies(ctx context.Context, iirID int) ([]StudentHobby, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_hobbies
		WHERE iir_id = ?
	`, database.GetColumns(StudentHobby{}))

	var hobbies []StudentHobby
	err := r.db.SelectContext(ctx, &hobbies, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student hobbies: %w", err)
	}

	return hobbies, nil
}

func (r *Repository) GetStudentTestResults(ctx context.Context, iirID int) ([]TestResult, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM test_results
		WHERE iir_id = ?
	`, database.GetColumns(TestResult{}))

	var results []TestResult
	err := r.db.SelectContext(ctx, &results, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student test results: %w", err)
	}

	return results, nil
}

func (r *Repository) GetStudentSignificantNotes(ctx context.Context, iirID int) ([]SignificantNote, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM significant_notes
		WHERE iir_id = ?
	`, database.GetColumns(SignificantNote{}))

	var notes []SignificantNote
	err := r.db.SelectContext(ctx, &notes, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student significant notes: %w", err)
	}

	return notes, nil
}

// Save and Upsert
func (r *Repository) CreateIIRRecord(ctx context.Context, iir *IIRRecord) (int, error) {
	var id int
	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(IIRRecord{}, []string{"created_at", "updated_at"})

		query := fmt.Sprintf(`
			INSERT INTO iir_records (%s)
			VALUES (%s)
		`, cols, vals)
		result, err := tx.NamedExecContext(ctx, query, iir)
		if err != nil {
			return fmt.Errorf("failed to create IIR record: %w", err)
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID for IIR record: %w", err)
		}

		id = int(lastID)
		return nil
	})
	return id, err
}

func (r *Repository) UpsertStudentPersonalInfo(ctx context.Context, info *StudentPersonalInfo) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(StudentPersonalInfo{}, []string{"created_at", "updated_at", "address_id"})
		updateCols := database.GetOnDuplicateKeyUpdateStatement(StudentPersonalInfo{}, []string{"created_at", "updated_at", "iir_id"})

		query := fmt.Sprintf(`
			INSERT INTO student_personal_info (%s)
			VALUES (%s)
			ON DUPLICATE KEY UPDATE %s
		`, cols, vals, updateCols)

		_, err := tx.NamedExecContext(ctx, query, info)
		return err
	})
}

func (r *Repository) UpsertEmergencyContact(ctx context.Context, ec *EmergencyContact) (int, error) {
	var id int
	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(EmergencyContact{}, []string{"created_at", "updated_at"})
		updateCols := database.GetOnDuplicateKeyUpdateStatement(EmergencyContact{}, []string{"created_at", "updated_at", "iir_id"})

		query := fmt.Sprintf(`
			INSERT INTO emergency_contacts (%s)
			VALUES (%s)
			ON DUPLICATE KEY UPDATE %s
		`, cols, vals, updateCols)
		result, err := tx.NamedExecContext(ctx, query, ec)
		if err != nil {
			return fmt.Errorf("failed to upsert emergency contact: %w", err)
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID for emergency contact: %w", err)
		}

		id = int(lastID)
		return nil
	})
	return id, err
}

func (r *Repository) UpsertStudentAddress(ctx context.Context, sa *StudentAddress) (int, error) {
	var id int
	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(StudentAddress{}, []string{"created_at", "updated_at"})
		updateCols := database.GetOnDuplicateKeyUpdateStatement(StudentAddress{}, []string{"created_at", "updated_at", "iir_id"})

		query := fmt.Sprintf(`
			INSERT INTO student_addresses (%s)
			VALUES (%s)
			ON DUPLICATE KEY UPDATE %s
		`, cols, vals, updateCols)
		result, err := tx.NamedExecContext(ctx, query, sa)
		if err != nil {
			return fmt.Errorf("failed to upsert student address: %w", err)
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID for student address: %w", err)
		}

		id = int(lastID)
		return nil
	})
	return id, err
}

func (r *Repository) UpsertAddress(ctx context.Context, addr *Address) (int, error) {
	var id int
	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(Address{}, []string{"created_at", "updated_at"})
		updateCols := database.GetOnDuplicateKeyUpdateStatement(Address{}, []string{"created_at", "updated_at"})

		query := fmt.Sprintf(`
			INSERT INTO addresses (%s)
			VALUES (%s)
			ON DUPLICATE KEY UPDATE %s
		`, cols, vals, updateCols)
		result, err := tx.NamedExecContext(ctx, query, addr)
		if err != nil {
			return fmt.Errorf("failed to upsert address: %w", err)
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID for address: %w", err)
		}

		id = int(lastID)
		return nil
	})
	return id, err
}

func (r *Repository) CreateStudentSelectedReason(ctx context.Context, ssr *StudentSelectedReason) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		query := `
			INSERT INTO student_selected_reasons (iir_id, reason_id, other_reason_text)
			VALUES (:iir_id, :reason_id, :other_reason_text)
		`
		_, err := tx.NamedExecContext(ctx, query, ssr)
		if err != nil {
			return fmt.Errorf("failed to create student selected reason: %w", err)
		}
		return nil
	})
}

func (r *Repository) DeleteStudentSelectedReasons(ctx context.Context, iirID int) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		query := `DELETE FROM student_selected_reasons WHERE iir_id = ?`
		_, err := tx.ExecContext(ctx, query, iirID)
		if err != nil {
			return fmt.Errorf("failed to delete student selected reasons: %w", err)
		}
		return nil
	})
}

func (r *Repository) UpsertRelatedPerson(ctx context.Context, rp *RelatedPerson) (int, error) {
	var id int
	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(RelatedPerson{}, []string{"created_at", "updated_at"})
		updateCols := database.GetOnDuplicateKeyUpdateStatement(RelatedPerson{}, []string{"created_at", "updated_at"})

		query := fmt.Sprintf(`
			INSERT INTO related_persons (%s)
			VALUES (%s)
			ON DUPLICATE KEY UPDATE %s
		`, cols, vals, updateCols)
		result, err := tx.NamedExecContext(ctx, query, rp)
		if err != nil {
			return fmt.Errorf("failed to upsert related person: %w", err)
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID for related person: %w", err)
		}

		id = int(lastID)
		return nil
	})
	return id, err
}

func (r *Repository) UpsertStudentRelatedPerson(ctx context.Context, srp *StudentRelatedPerson) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(StudentRelatedPerson{}, []string{"created_at", "updated_at"})
		updateCols := database.GetOnDuplicateKeyUpdateStatement(StudentRelatedPerson{}, []string{"created_at", "updated_at"})

		query := fmt.Sprintf(`
			INSERT INTO student_related_persons (%s)
			VALUES (%s)
			ON DUPLICATE KEY UPDATE %s
		`, cols, vals, updateCols)

		_, err := tx.NamedExecContext(ctx, query, srp)
		if err != nil {
			return fmt.Errorf("failed to upsert student related person: %w", err)
		}
		return nil
	})
}

func (r *Repository) DeleteStudentRelatedPersons(ctx context.Context, iirID int) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		query := `DELETE FROM student_related_persons WHERE iir_id = ?`
		_, err := tx.ExecContext(ctx, query, iirID)
		if err != nil {
			return fmt.Errorf("failed to delete student related persons: %w", err)
		}
		return nil
	})
}

func (r *Repository) UpsertFamilyBackground(ctx context.Context, fb *FamilyBackground) (int, error) {
	var id int
	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(FamilyBackground{}, []string{"created_at", "updated_at"})
		updateCols := database.GetOnDuplicateKeyUpdateStatement(FamilyBackground{}, []string{"created_at", "updated_at", "iir_id"})

		query := fmt.Sprintf(`
			INSERT INTO family_backgrounds (%s)
			VALUES (%s)
			ON DUPLICATE KEY UPDATE %s
		`, cols, vals, updateCols)
		result, err := tx.NamedExecContext(ctx, query, fb)
		if err != nil {
			return fmt.Errorf("failed to upsert family background: %w", err)
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID for family background: %w", err)
		}

		id = int(lastID)
		return nil
	})
	return id, err
}

func (r *Repository) CreateStudentSiblingSupport(ctx context.Context, sss *StudentSiblingSupport) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		query := `
			INSERT INTO student_sibling_supports (family_background_id, support_type_id)
			VALUES (:family_background_id, :support_type_id)
		`
		_, err := tx.NamedExecContext(ctx, query, sss)
		if err != nil {
			return fmt.Errorf("failed to create student sibling support: %w", err)
		}
		return nil
	})
}

func (r *Repository) DeleteStudentSiblingSupportsByFamilyID(ctx context.Context, familyBackgroundID int) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		query := `DELETE FROM student_sibling_supports WHERE family_background_id = ?`
		_, err := tx.ExecContext(ctx, query, familyBackgroundID)
		if err != nil {
			return fmt.Errorf("failed to delete student sibling supports: %w", err)
		}
		return nil
	})
}

func (r *Repository) UpsertEducationalBackground(ctx context.Context, eb *EducationalBackground) (int, error) {
	var id int
	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(EducationalBackground{}, []string{"created_at", "updated_at"})
		updateCols := database.GetOnDuplicateKeyUpdateStatement(EducationalBackground{}, []string{"created_at", "updated_at", "iir_id"})

		query := fmt.Sprintf(`
			INSERT INTO educational_backgrounds (%s)
			VALUES (%s)
			ON DUPLICATE KEY UPDATE %s
		`, cols, vals, updateCols)
		result, err := tx.NamedExecContext(ctx, query, eb)
		if err != nil {
			return fmt.Errorf("failed to upsert educational background: %w", err)
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID for educational background: %w", err)
		}

		id = int(lastID)
		return nil
	})
	return id, err
}

func (r *Repository) UpsertSchoolDetails(ctx context.Context, sd *SchoolDetails) (int, error) {
	var id int
	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(SchoolDetails{}, []string{"created_at", "updated_at"})
		updateCols := database.GetOnDuplicateKeyUpdateStatement(SchoolDetails{}, []string{"created_at", "updated_at"})

		query := fmt.Sprintf(`
			INSERT INTO school_details (%s)
			VALUES (%s)
			ON DUPLICATE KEY UPDATE %s
		`, cols, vals, updateCols)
		result, err := tx.NamedExecContext(ctx, query, sd)
		if err != nil {
			return fmt.Errorf("failed to upsert school details: %w", err)
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID for school details: %w", err)
		}

		id = int(lastID)
		return nil
	})
	return id, err
}

func (r *Repository) DeleteSchoolDetailsByEBID(ctx context.Context, ebID int) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		query := `DELETE FROM school_details WHERE eb_id = ?`
		_, err := tx.ExecContext(ctx, query, ebID)
		if err != nil {
			return fmt.Errorf("failed to delete school details: %w", err)
		}
		return nil
	})
}

func (r *Repository) UpsertStudentHealthRecord(ctx context.Context, hr *StudentHealthRecord) (int, error) {
	var id int
	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(StudentHealthRecord{}, []string{"created_at", "updated_at"})
		updateCols := database.GetOnDuplicateKeyUpdateStatement(StudentHealthRecord{}, []string{"created_at", "updated_at", "iir_id"})

		query := fmt.Sprintf(`
			INSERT INTO student_health_records (%s)
			VALUES (%s)
			ON DUPLICATE KEY UPDATE %s
		`, cols, vals, updateCols)
		result, err := tx.NamedExecContext(ctx, query, hr)
		if err != nil {
			return fmt.Errorf("failed to upsert student health record: %w", err)
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID for student health record: %w", err)
		}

		id = int(lastID)
		return nil
	})
	return id, err
}

func (r *Repository) UpsertStudentConsultation(ctx context.Context, sc *StudentConsultation) (int, error) {
	var id int
	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(StudentConsultation{}, []string{"created_at", "updated_at"})
		updateCols := database.GetOnDuplicateKeyUpdateStatement(StudentConsultation{}, []string{"created_at", "updated_at", "iir_id"})

		query := fmt.Sprintf(`
			INSERT INTO student_consultations (%s)
			VALUES (%s)
			ON DUPLICATE KEY UPDATE %s
		`, cols, vals, updateCols)
		result, err := tx.NamedExecContext(ctx, query, sc)
		if err != nil {
			return fmt.Errorf("failed to upsert student consultation: %w", err)
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID for student consultation: %w", err)
		}

		id = int(lastID)
		return nil
	})
	return id, err
}

func (r *Repository) UpsertStudentFinance(ctx context.Context, sf *StudentFinance) (int, error) {
	var id int
	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(StudentFinance{}, []string{"created_at", "updated_at"})
		updateCols := database.GetOnDuplicateKeyUpdateStatement(StudentFinance{}, []string{"created_at", "updated_at", "iir_id"})

		query := fmt.Sprintf(`
			INSERT INTO student_finances (%s)
			VALUES (%s)
			ON DUPLICATE KEY UPDATE %s
		`, cols, vals, updateCols)
		result, err := tx.NamedExecContext(ctx, query, sf)
		if err != nil {
			return fmt.Errorf("failed to upsert student finance: %w", err)
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID for student finance: %w", err)
		}

		id = int(lastID)
		return nil
	})
	return id, err
}

func (r *Repository) CreateStudentFinancialSupport(ctx context.Context, sfs *StudentFinancialSupport) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		query := `
			INSERT INTO student_financial_supports (sf_id, support_type_id, created_at, updated_at)
			VALUES (:sf_id, :support_type_id, :created_at, :updated_at)
		`
		_, err := tx.NamedExecContext(ctx, query, sfs)
		if err != nil {
			return fmt.Errorf("failed to create student financial support: %w", err)
		}
		return nil
	})
}

func (r *Repository) DeleteStudentFinancialSupportsByFinanceID(ctx context.Context, financeID int) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		query := `DELETE FROM student_financial_supports WHERE sf_id = ?`
		_, err := tx.ExecContext(ctx, query, financeID)
		if err != nil {
			return fmt.Errorf("failed to delete student financial supports: %w", err)
		}
		return nil
	})
}

func (r *Repository) CreateStudentActivity(ctx context.Context, sa *StudentActivity) (int, error) {
	var id int
	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(StudentActivity{}, []string{"created_at", "updated_at"})

		query := fmt.Sprintf(`
			INSERT INTO student_activities (%s)
			VALUES (%s)
		`, cols, vals)
		result, err := tx.NamedExecContext(ctx, query, sa)
		if err != nil {
			return fmt.Errorf("failed to create student activity: %w", err)
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID for student activity: %w", err)
		}

		id = int(lastID)
		return nil
	})
	return id, err
}

func (r *Repository) DeleteStudentActivitiesByIIRID(ctx context.Context, iirID int) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		query := `DELETE FROM student_activities WHERE iir_id = ?`
		_, err := tx.ExecContext(ctx, query, iirID)
		if err != nil {
			return fmt.Errorf("failed to delete student activities: %w", err)
		}
		return nil
	})
}

func (r *Repository) CreateStudentSubjectPreference(ctx context.Context, ssp *StudentSubjectPreference) (int, error) {
	var id int
	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(StudentSubjectPreference{}, []string{"created_at", "updated_at"})

		query := fmt.Sprintf(`
			INSERT INTO student_subject_preferences (%s)
			VALUES (%s)
		`, cols, vals)
		result, err := tx.NamedExecContext(ctx, query, ssp)
		if err != nil {
			return fmt.Errorf("failed to create student subject preference: %w", err)
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID for student subject preference: %w", err)
		}

		id = int(lastID)
		return nil
	})
	return id, err
}

func (r *Repository) DeleteStudentSubjectPreferencesByIIRID(ctx context.Context, iirID int) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		query := `DELETE FROM student_subject_preferences WHERE iir_id = ?`
		_, err := tx.ExecContext(ctx, query, iirID)
		if err != nil {
			return fmt.Errorf("failed to delete student subject preferences: %w", err)
		}
		return nil
	})
}

func (r *Repository) CreateStudentHobby(ctx context.Context, sh *StudentHobby) (int, error) {
	var id int
	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(StudentHobby{}, []string{"created_at", "updated_at"})

		query := fmt.Sprintf(`
			INSERT INTO student_hobbies (%s)
			VALUES (%s)
		`, cols, vals)
		result, err := tx.NamedExecContext(ctx, query, sh)
		if err != nil {
			return fmt.Errorf("failed to create student hobby: %w", err)
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID for student hobby: %w", err)
		}

		id = int(lastID)
		return nil
	})
	return id, err
}

func (r *Repository) DeleteStudentHobbiesByIIRID(ctx context.Context, iirID int) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		query := `DELETE FROM student_hobbies WHERE iir_id = ?`
		_, err := tx.ExecContext(ctx, query, iirID)
		if err != nil {
			return fmt.Errorf("failed to delete student hobbies: %w", err)
		}
		return nil
	})
}

func (r *Repository) CreateTestResult(ctx context.Context, tr *TestResult) (int, error) {
	var id int
	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(TestResult{}, []string{"created_at", "updated_at"})

		query := fmt.Sprintf(`
			INSERT INTO test_results (%s)
			VALUES (%s)
		`, cols, vals)
		result, err := tx.NamedExecContext(ctx, query, tr)
		if err != nil {
			return fmt.Errorf("failed to create test result: %w", err)
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID for test result: %w", err)
		}

		id = int(lastID)
		return nil
	})
	return id, err
}

func (r *Repository) DeleteTestResultsByIIRID(ctx context.Context, iirID int) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		query := `DELETE FROM test_results WHERE iir_id = ?`
		_, err := tx.ExecContext(ctx, query, iirID)
		if err != nil {
			return fmt.Errorf("failed to delete test results: %w", err)
		}
		return nil
	})
}

func (r *Repository) CreateSignificantNote(ctx context.Context, sn *SignificantNote) (int, error) {
	var id int
	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(SignificantNote{}, []string{"created_at", "updated_at"})

		query := fmt.Sprintf(`
			INSERT INTO significant_notes (%s)
			VALUES (%s)
		`, cols, vals)
		result, err := tx.NamedExecContext(ctx, query, sn)
		if err != nil {
			return fmt.Errorf("failed to create significant note: %w", err)
		}

		lastID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID for significant note: %w", err)
		}

		id = int(lastID)
		return nil
	})
	return id, err
}

func (r *Repository) DeleteSignificantNotesByIIRID(ctx context.Context, iirID int) error {
	return database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		query := `DELETE FROM significant_notes WHERE iir_id = ?`
		_, err := tx.ExecContext(ctx, query, iirID)
		if err != nil {
			return fmt.Errorf("failed to delete significant notes: %w", err)
		}
		return nil
	})
}
