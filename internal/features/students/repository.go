package students

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetDB() *sqlx.DB {
	return r.db
}

func (r *Repository) BeginTx(ctx context.Context) (datastore.DB, error) {
	return r.db.BeginTxx(ctx, nil)
}

// Lookup
func (r *Repository) GetGenders(ctx context.Context) ([]Gender, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM genders ORDER BY id
	`, datastore.GetColumns(Gender{}))

	var genders []Gender
	err := r.db.SelectContext(ctx, &genders, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get genders: %w", err)
	}

	return genders, nil
}

func (r *Repository) GetParentalStatusTypes(
	ctx context.Context,
) ([]ParentalStatusType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM parental_status_types ORDER BY id
	`, datastore.GetColumns(ParentalStatusType{}))
	var statuses []ParentalStatusType
	err := r.db.SelectContext(ctx, &statuses, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get parental status types: %w", err)
	}

	return statuses, nil
}

func (r *Repository) GetEnrollmentReasons(
	ctx context.Context,
) ([]EnrollmentReason, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM enrollment_reasons ORDER BY id
	`, datastore.GetColumns(EnrollmentReason{}))
	var reasons []EnrollmentReason
	err := r.db.SelectContext(ctx, &reasons, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment reasons: %w", err)
	}

	return reasons, nil
}

func (r *Repository) GetIncomeRanges(
	ctx context.Context,
) ([]IncomeRange, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM income_ranges ORDER BY id
	`, datastore.GetColumns(IncomeRange{}))
	var ranges []IncomeRange
	err := r.db.SelectContext(ctx, &ranges, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get income ranges: %w", err)
	}

	return ranges, nil
}

func (r *Repository) GetStudentSupportTypes(
	ctx context.Context,
) ([]StudentSupportType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_support_types ORDER BY id
	`, datastore.GetColumns(StudentSupportType{}))
	var supportTypes []StudentSupportType
	err := r.db.SelectContext(ctx, &supportTypes, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get student support types: %w", err)
	}

	return supportTypes, nil
}

func (r *Repository) GetSiblingSupportTypes(
	ctx context.Context,
) ([]SibilingSupportType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM sibling_support_types ORDER BY id
	`, datastore.GetColumns(SibilingSupportType{}))

	var supportTypes []SibilingSupportType
	err := r.db.SelectContext(ctx, &supportTypes, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get sibling support types: %w", err)
	}

	return supportTypes, nil
}

func (r *Repository) GetEducationalLevels(
	ctx context.Context,
) ([]EducationalLevel, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM educational_levels ORDER BY id
	`, datastore.GetColumns(EducationalLevel{}))

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
	`, datastore.GetColumns(Course{}))

	var courses []Course
	err := r.db.SelectContext(ctx, &courses, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}

	return courses, nil
}

func (r *Repository) GetCivilStatusTypes(
	ctx context.Context,
) ([]CivilStatusType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM civil_status_types ORDER BY id
	`, datastore.GetColumns(CivilStatusType{}))

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
	`, datastore.GetColumns(Religion{}))

	var religions []Religion
	err := r.db.SelectContext(ctx, &religions, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get religions: %w", err)
	}

	return religions, nil
}

func (r *Repository) GetStudentRelationshipTypes(
	ctx context.Context,
) ([]StudentRelationshipType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_relationship_types ORDER BY id
	`, datastore.GetColumns(StudentRelationshipType{}))

	var relationships []StudentRelationshipType
	err := r.db.SelectContext(ctx, &relationships, query)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student relationship types: %w",
			err,
		)
	}

	return relationships, nil
}

func (r *Repository) GetNatureOfResidenceTypes(
	ctx context.Context,
) ([]NatureOfResidenceType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM nature_of_residence_types ORDER BY id
	`, datastore.GetColumns(NatureOfResidenceType{}))

	var residences []NatureOfResidenceType
	err := r.db.SelectContext(ctx, &residences, query)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get nature of residence types: %w",
			err,
		)
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
	query, args := r.applyStudentFilters(
		`SELECT COUNT(iir.id) FROM iir_records iir
         JOIN users u ON iir.user_id = u.id
         JOIN student_personal_info spi ON iir.id = spi.iir_id
         WHERE iir.is_submitted = TRUE`,
		nil,
		search,
		courseID,
		genderID,
		yearLevel,
	)

	var total int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *Repository) applyStudentFilters(
	query string,
	args []interface{},
	search string,
	courseID, genderID, yearLevel int,
) (string, []interface{}) {
	if args == nil {
		args = []interface{}{}
	}

	if courseID > 0 {
		query += " AND spi.course_id = ?"
		args = append(args, courseID)
	}

	if genderID > 0 {
		query += " AND spi.gender_id = ?"
		args = append(args, genderID)
	}

	if yearLevel > 0 {
		query += " AND spi.year_level = ?"
		args = append(args, yearLevel)
	}

	if search != "" {
		query += ` AND (u.first_name LIKE ?
                 OR u.last_name LIKE ?
                 OR u.email LIKE ?
                 OR spi.student_number LIKE ?)`

		pattern := "%" + search + "%"
		args = append(args, pattern, pattern, pattern, pattern)
	}

	return query, args
}

// Retrieve - List
func (r *Repository) ListStudents(
	ctx context.Context, search string, offset int, limit int, orderBy string,
	courseID int, genderID int, yearLevel int,
) ([]StudentProfileView, error) {
	query := `
        SELECT
			iir.id as iir_id,
			u.id as user_id,
			u.first_name,
			u.middle_name,
			u.last_name,
			spi.suffix_name,
			spi.gender_id,
			u.email,
			spi.student_number,
			spi.course_id,
			spi.section,
			spi.year_level
		FROM iir_records iir
		JOIN users u ON iir.user_id = u.id
		JOIN student_personal_info spi ON iir.id = spi.iir_id
		WHERE iir.is_submitted = TRUE
    `

	query, args := r.applyStudentFilters(
		query,
		nil,
		search,
		courseID,
		genderID,
		yearLevel,
	)

	allowedSortColumns := map[string]string{
		"last_name":      "u.last_name",
		"first_name":     "u.first_name",
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
			&student.SuffixName,
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

func (r *Repository) GetStudentBasicInfo(
	ctx context.Context,
	iirID string,
) (*StudentBasicInfoView, error) {
	query := `
		SELECT u.id, u.email, u.first_name, u.middle_name, u.last_name
		FROM users u
		JOIN iir_records iir ON u.id = iir.user_id
		WHERE iir.id = ?
	`

	var info StudentBasicInfoView
	err := r.db.QueryRowContext(ctx, query, iirID).Scan(
		&info.UserID,
		&info.Email,
		&info.FirstName,
		&info.MiddleName,
		&info.LastName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get student basic info: %w", err)
	}

	return &info, nil
}

func (r *Repository) GetIIRDraftByUserID(
	ctx context.Context,
	userID string,
) (*IIRDraft, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM iir_drafts WHERE user_id = ? LIMIT 1
	`, datastore.GetColumns(IIRDraft{}))

	var draft IIRDraft
	err := r.db.GetContext(ctx, &draft, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &draft, nil
}

func (r *Repository) GetStudentIIRByUserID(
	ctx context.Context,
	userID string,
) (*IIRRecord, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM iir_records WHERE user_id = ? LIMIT 1
	`, datastore.GetColumns(IIRRecord{}))

	var iir IIRRecord
	err := r.db.GetContext(ctx, &iir, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &iir, nil
}

func (r *Repository) GetStudentIIR(
	ctx context.Context,
	iirID string,
) (*IIRRecord, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM iir_records
		WHERE id = ?
	`, datastore.GetColumns(IIRRecord{}))

	var iir IIRRecord
	err := r.db.GetContext(ctx, &iir, query, iirID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &iir, nil
}

func (r *Repository) GetStudentEnrollmentReasons(
	ctx context.Context,
	iirID string,
) ([]StudentSelectedReason, error) {
	query := `
		SELECT iir_id, reason_id, other_reason_text
		FROM student_selected_reasons
		WHERE iir_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, iirID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student enrollment reasons: %w",
			err,
		)
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

func (r *Repository) GetEnrollmentReasonByID(
	ctx context.Context,
	reasonID int,
) (*EnrollmentReason, error) {
	query := `SELECT id, reason_text FROM enrollment_reasons WHERE id = ?`
	var er EnrollmentReason
	err := r.db.QueryRowContext(ctx, query, reasonID).Scan(&er.ID, &er.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment reason by ID: %w", err)
	}

	return &er, nil
}

func (r *Repository) GetStudentPersonalInfo(
	ctx context.Context,
	iirID string,
) (*StudentPersonalInfo, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_personal_info
		WHERE iir_id = ?
	`, datastore.GetColumns(StudentPersonalInfo{}))

	var info StudentPersonalInfo
	err := r.db.GetContext(ctx, &info, query, iirID)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (r *Repository) GetEmergencyContactByIIRID(
	ctx context.Context,
	iirID string,
) (*EmergencyContact, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM emergency_contacts
		WHERE iir_id = ?
	`, datastore.GetColumns(EmergencyContact{}))

	var ec EmergencyContact
	err := r.db.GetContext(ctx, &ec, query, iirID)
	if err != nil {
		return nil, err
	}

	return &ec, nil
}

func (r *Repository) GetGenderByID(
	ctx context.Context,
	genderID int,
) (*Gender, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM genders WHERE id = ?
	`, datastore.GetColumns(Gender{}))

	var gender Gender
	err := r.db.GetContext(ctx, &gender, query, genderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gender by ID: %w", err)
	}

	return &gender, nil
}

func (r *Repository) GetCivilStatusByID(
	ctx context.Context,
	statusID int,
) (*CivilStatusType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM civil_status_types WHERE id = ?
	`, datastore.GetColumns(CivilStatusType{}))

	var status CivilStatusType
	err := r.db.GetContext(ctx, &status, query, statusID)
	if err != nil {
		return nil, fmt.Errorf("failed to get civil status by ID: %w", err)
	}

	return &status, nil
}

func (r *Repository) GetReligionByID(
	ctx context.Context,
	religionID int,
) (*Religion, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM religions WHERE id = ?
	`, datastore.GetColumns(Religion{}))

	var religion Religion
	err := r.db.GetContext(ctx, &religion, query, religionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get religion by ID: %w", err)
	}

	return &religion, nil
}

func (r *Repository) GetCourseByID(
	ctx context.Context,
	courseID int,
) (*Course, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM courses WHERE id = ?
	`, datastore.GetColumns(Course{}))

	var course Course
	err := r.db.GetContext(ctx, &course, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course by ID: %w", err)
	}

	return &course, nil
}

func (r *Repository) GetStudentAddresses(
	ctx context.Context,
	iirID string,
) ([]StudentAddress, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_addresses
		WHERE iir_id = ?
	`, datastore.GetColumns(StudentAddress{}))

	var addresses []StudentAddress
	err := r.db.SelectContext(ctx, &addresses, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student addresses: %w", err)
	}

	return addresses, nil
}

func (r *Repository) GetStudentEducationalBackground(
	ctx context.Context,
	iirID string,
) (*EducationalBackground, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM educational_backgrounds
		WHERE iir_id = ?
	`, datastore.GetColumns(EducationalBackground{}))

	var eb EducationalBackground
	err := r.db.GetContext(ctx, &eb, query, iirID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student educational background: %w",
			err,
		)
	}

	return &eb, nil
}

func (r *Repository) GetSchoolDetailsByEBID(
	ctx context.Context,
	ebID int,
) ([]SchoolDetails, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM school_details
		WHERE eb_id = ?
		ORDER BY educational_level_id ASC
	`, datastore.GetColumns(SchoolDetails{}))

	var schoolDetails []SchoolDetails
	err := r.db.SelectContext(ctx, &schoolDetails, query, ebID)
	if err != nil {
		return nil, fmt.Errorf("failed to get school details by EB ID: %w", err)
	}

	return schoolDetails, nil
}

func (r *Repository) GetEducationalLevelByID(
	ctx context.Context,
	levelID int,
) (*EducationalLevel, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM educational_levels WHERE id = ?
	`, datastore.GetColumns(EducationalLevel{}))
	var el EducationalLevel
	err := r.db.GetContext(ctx, &el, query, levelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get educational level by ID: %w", err)
	}

	return &el, nil
}

func (r *Repository) GetStudentRelatedPersons(
	ctx context.Context, iirID string,
) ([]StudentRelatedPerson, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_related_persons WHERE iir_id = ?
	`, datastore.GetColumns(StudentRelatedPerson{}))

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
	`, datastore.GetColumns(RelatedPerson{}))

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
	`, datastore.GetColumns(StudentRelationshipType{}))

	var srt StudentRelationshipType
	err := r.db.GetContext(ctx, &srt, query, relationshipID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student relationship by ID: %w",
			err,
		)
	}

	return &srt, nil
}

func (r *Repository) GetStudentFamilyBackground(
	ctx context.Context,
	iirID string,
) (*FamilyBackground, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM family_backgrounds
		WHERE iir_id = ?
	`, datastore.GetColumns(FamilyBackground{}))

	var fb FamilyBackground
	err := r.db.GetContext(ctx, &fb, query, iirID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student family background: %w",
			err,
		)
	}

	return &fb, nil
}

func (r *Repository) GetParentalStatusByID(
	ctx context.Context,
	statusID int,
) (*ParentalStatusType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM parental_status_types WHERE id = ?
	`, datastore.GetColumns(ParentalStatusType{}))

	var ps ParentalStatusType
	err := r.db.GetContext(ctx, &ps, query, statusID)
	if err != nil {
		return nil, fmt.Errorf("failed to get parental status by ID: %w", err)
	}

	return &ps, nil
}

func (r *Repository) GetNatureOfResidenceByID(
	ctx context.Context,
	residenceID int,
) (*NatureOfResidenceType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM nature_of_residence_types WHERE id = ?
	`, datastore.GetColumns(NatureOfResidenceType{}))

	var nr NatureOfResidenceType
	err := r.db.GetContext(ctx, &nr, query, residenceID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get nature of residence by ID: %w",
			err,
		)
	}

	return &nr, nil
}

func (r *Repository) GetStudentSiblingSupport(
	ctx context.Context,
	fbID int,
) ([]StudentSiblingSupport, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_sibling_supports
		WHERE family_background_id = ?
	`, datastore.GetColumns(StudentSiblingSupport{}))

	var sss []StudentSiblingSupport
	err := r.db.SelectContext(ctx, &sss, query, fbID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to query student sibling supports: %w",
			err,
		)
	}

	return sss, nil
}

func (r *Repository) GetSiblingSupportTypeByID(
	ctx context.Context,
	supportID int,
) (*SibilingSupportType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM sibling_support_types WHERE id = ?
	`, datastore.GetColumns(SibilingSupportType{}))

	var sst SibilingSupportType
	err := r.db.GetContext(ctx, &sst, query, supportID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get sibling support type by ID: %w",
			err,
		)
	}

	return &sst, nil
}

func (r *Repository) GetStudentFinancialInfo(
	ctx context.Context,
	iirID string,
) (*StudentFinance, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_finances
		WHERE iir_id = ?
	`, datastore.GetColumns(StudentFinance{}))

	var fi StudentFinance
	err := r.db.GetContext(ctx, &fi, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student financial info: %w", err)
	}

	return &fi, nil
}

func (r *Repository) GetFinancialSupportTypeByFinanceID(
	ctx context.Context,
	financeID int,
) ([]StudentFinancialSupport, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_financial_supports
		WHERE sf_id = ?
	`, datastore.GetColumns(StudentFinancialSupport{}))

	var sfs []StudentFinancialSupport
	err := r.db.SelectContext(ctx, &sfs, query, financeID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to query student financial supports: %w",
			err,
		)
	}

	return sfs, nil
}

func (r *Repository) GetIncomeRangeByID(
	ctx context.Context,
	rangeID int,
) (*IncomeRange, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM income_ranges WHERE id = ?
	`, datastore.GetColumns(IncomeRange{}))

	var ir IncomeRange
	err := r.db.GetContext(ctx, &ir, query, rangeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get income range by ID: %w", err)
	}

	return &ir, nil
}

func (r *Repository) GetStudentSupportByID(
	ctx context.Context,
	supportID int,
) (*StudentSupportType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_support_types WHERE id = ?
	`, datastore.GetColumns(StudentSupportType{}))

	var sst StudentSupportType
	err := r.db.GetContext(ctx, &sst, query, supportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student support by ID: %w", err)
	}

	return &sst, nil
}

func (r *Repository) GetStudentHealthRecord(
	ctx context.Context,
	iirID string,
) (*StudentHealthRecord, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_health_records
		WHERE iir_id = ?
	`, datastore.GetColumns(StudentHealthRecord{}))

	var hr StudentHealthRecord
	err := r.db.GetContext(ctx, &hr, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student health record: %w", err)
	}

	return &hr, nil
}

func (r *Repository) GetActivityOptions(
	ctx context.Context,
) ([]ActivityOption, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM activity_options ORDER BY id
	`, datastore.GetColumns(ActivityOption{}))

	var options []ActivityOption
	err := r.db.SelectContext(ctx, &options, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity options: %w", err)
	}

	return options, nil
}

func (r *Repository) GetStudentConsultations(
	ctx context.Context,
	iirID string,
) ([]StudentConsultation, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_consultations
		WHERE iir_id = ?
	`, datastore.GetColumns(StudentConsultation{}))

	var consultations []StudentConsultation
	err := r.db.SelectContext(ctx, &consultations, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student consultations: %w", err)
	}

	return consultations, nil
}

func (r *Repository) GetStudentActivities(
	ctx context.Context,
	iirID string,
) ([]StudentActivity, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_activities
		WHERE iir_id = ?
	`, datastore.GetColumns(StudentActivity{}))

	var activities []StudentActivity
	err := r.db.SelectContext(ctx, &activities, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student activities: %w", err)
	}

	return activities, nil
}

func (r *Repository) GetActivityOptionByID(
	ctx context.Context,
	optionID int,
) (*ActivityOption, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM activity_options WHERE id = ?
	`, datastore.GetColumns(ActivityOption{}))

	var ao ActivityOption
	err := r.db.GetContext(ctx, &ao, query, optionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity option by ID: %w", err)
	}

	return &ao, nil
}

func (r *Repository) GetStudentSubjectPreferences(
	ctx context.Context,
	iirID string,
) ([]StudentSubjectPreference, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_subject_preferences
		WHERE iir_id = ?
	`, datastore.GetColumns(StudentSubjectPreference{}))

	var preferences []StudentSubjectPreference
	err := r.db.SelectContext(ctx, &preferences, query, iirID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student subject preferences: %w",
			err,
		)
	}

	return preferences, nil
}

func (r *Repository) GetStudentHobbies(
	ctx context.Context,
	iirID string,
) ([]StudentHobby, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_hobbies
		WHERE iir_id = ?
	`, datastore.GetColumns(StudentHobby{}))

	var hobbies []StudentHobby
	err := r.db.SelectContext(ctx, &hobbies, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student hobbies: %w", err)
	}

	return hobbies, nil
}

func (r *Repository) GetStudentTestResults(
	ctx context.Context,
	iirID string,
) ([]TestResult, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM test_results
		WHERE iir_id = ?
	`, datastore.GetColumns(TestResult{}))

	var results []TestResult
	err := r.db.SelectContext(ctx, &results, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student test results: %w", err)
	}

	return results, nil
}

// Save and Upsert
func (r *Repository) UpsertIIRDraft(
	ctx context.Context,
	draft IIRDraft,
) (int, error) {
	var id int
	err := datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		var err error
		id, err = r.upsertIIRDraftTx(ctx, txn, draft)
		return err
	})
	return id, err
}

func (r *Repository) upsertIIRDraftTx(
	ctx context.Context,
	tx datastore.DB,
	draft IIRDraft,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(
		IIRDraft{},
		[]string{"created_at", "updated_at"},
	)
	onDuplicateKey := datastore.GetOnDuplicateKeyUpdateStatement(
		IIRDraft{},
		[]string{"created_at", "updated_at"},
	)
	query := fmt.Sprintf(`
		INSERT INTO iir_drafts (%s)
		VALUES (%s)
		ON DUPLICATE KEY UPDATE %s
	`, cols, vals, onDuplicateKey)
	result, err := tx.NamedExecContext(ctx, query, draft)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert IIR draft: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get last insert ID for IIR draft: %w",
			err,
		)
	}

	return int(lastID), nil
}

func (r *Repository) UpsertIIRRecord(
	ctx context.Context,
	tx datastore.DB,
	iir *IIRRecord,
) (string, error) {
	if tx != nil {
		return r.upsertIIRRecordTx(ctx, tx, iir)
	}

	var id string
	err := datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		var err error
		id, err = r.upsertIIRRecordTx(ctx, txn, iir)
		return err
	})
	return id, err
}

func (r *Repository) upsertIIRRecordTx(
	ctx context.Context,
	tx datastore.DB,
	iir *IIRRecord,
) (string, error) {
	if iir.ID == "" {
		existing, err := r.GetStudentIIRByUserID(ctx, iir.UserID)
		if err == nil && existing != nil {
			iir.ID = existing.ID
		} else {
			iir.ID = uuid.New().String()
		}
	}
	cols, vals := datastore.GetInsertStatement(
		IIRRecord{},
		[]string{"created_at", "updated_at"},
	)
	onDuplicateKey := datastore.GetOnDuplicateKeyUpdateStatement(
		IIRRecord{},
		[]string{"created_at", "updated_at"},
	)
	query := fmt.Sprintf(`
		INSERT INTO iir_records (id, %s)
		VALUES (:id, %s)
		ON DUPLICATE KEY UPDATE %s
	`, cols, vals, onDuplicateKey)
	_, err := tx.NamedExecContext(ctx, query, iir)
	if err != nil {
		return "", fmt.Errorf("failed to upsert IIR record: %w", err)
	}

	return iir.ID, nil
}

func (r *Repository) UpsertStudentPersonalInfo(
	ctx context.Context,
	tx datastore.DB,
	info *StudentPersonalInfo,
) error {
	if tx != nil {
		return r.upsertStudentPersonalInfoTx(ctx, tx, info)
	}

	return datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		return r.upsertStudentPersonalInfoTx(ctx, txn, info)
	})
}

func (r *Repository) upsertStudentPersonalInfoTx(
	ctx context.Context,
	tx datastore.DB,
	info *StudentPersonalInfo,
) error {
	cols, vals := datastore.GetInsertStatement(
		StudentPersonalInfo{},
		[]string{"created_at", "updated_at", "address_id"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentPersonalInfo{},
		[]string{"created_at", "updated_at", "iir_id"},
	)

	query := fmt.Sprintf(`
		INSERT INTO student_personal_info (%s)
		VALUES (%s)
		ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)

	_, err := tx.NamedExecContext(ctx, query, info)
	return err
}

func (r *Repository) UpsertEmergencyContact(
	ctx context.Context,
	tx datastore.DB,
	ec *EmergencyContact,
) (int, error) {
	if tx != nil {
		return r.upsertEmergencyContactTx(ctx, tx, ec)
	}

	var id int
	err := datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		var err error
		id, err = r.upsertEmergencyContactTx(ctx, txn, ec)
		return err
	})
	return id, err
}

func (r *Repository) upsertEmergencyContactTx(
	ctx context.Context,
	tx datastore.DB,
	ec *EmergencyContact,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(
		EmergencyContact{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		EmergencyContact{},
		[]string{"created_at", "updated_at", "iir_id"},
	)

	query := fmt.Sprintf(`
		INSERT INTO emergency_contacts (%s)
		VALUES (%s)
		ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)
	result, err := tx.NamedExecContext(ctx, query, ec)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert emergency contact: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get last insert ID for emergency contact: %w",
			err,
		)
	}

	return int(lastID), nil
}

func (r *Repository) UpsertStudentAddress(
	ctx context.Context,
	tx datastore.DB,
	sa *StudentAddress,
) (int, error) {
	if tx != nil {
		return r.upsertStudentAddressTx(ctx, tx, sa)
	}

	var id int
	err := datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		var err error
		id, err = r.upsertStudentAddressTx(ctx, txn, sa)
		return err
	})
	return id, err
}

func (r *Repository) upsertStudentAddressTx(
	ctx context.Context,
	tx datastore.DB,
	sa *StudentAddress,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(
		StudentAddress{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentAddress{},
		[]string{"created_at", "updated_at", "iir_id"},
	)

	query := fmt.Sprintf(`
		INSERT INTO student_addresses (%s)
		VALUES (%s)
		ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)
	result, err := tx.NamedExecContext(ctx, query, sa)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert student address: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get last insert ID for student address: %w",
			err,
		)
	}

	return int(lastID), nil
}

func (r *Repository) CreateStudentSelectedReason(
	ctx context.Context,
	tx datastore.DB,
	ssr *StudentSelectedReason,
) error {
	if tx != nil {
		return r.createStudentSelectedReasonTx(ctx, tx, ssr)
	}

	return datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		return r.createStudentSelectedReasonTx(ctx, txn, ssr)
	})
}

func (r *Repository) createStudentSelectedReasonTx(
	ctx context.Context,
	tx datastore.DB,
	ssr *StudentSelectedReason,
) error {
	query := `
		INSERT INTO student_selected_reasons (iir_id, reason_id, other_reason_text)
		VALUES (:iir_id, :reason_id, :other_reason_text)
	`
	_, err := tx.NamedExecContext(ctx, query, ssr)
	if err != nil {
		return fmt.Errorf("failed to create student selected reason: %w", err)
	}
	return nil
}

func (r *Repository) DeleteStudentSelectedReasons(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	if tx != nil {
		return r.deleteStudentSelectedReasonsTx(ctx, tx, iirID)
	}

	return datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		return r.deleteStudentSelectedReasonsTx(ctx, txn, iirID)
	})
}

func (r *Repository) deleteStudentSelectedReasonsTx(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	query := `DELETE FROM student_selected_reasons WHERE iir_id = ?`
	_, err := tx.ExecContext(ctx, query, iirID)
	if err != nil {
		return fmt.Errorf("failed to delete student selected reasons: %w", err)
	}
	return nil
}

func (r *Repository) UpsertRelatedPerson(
	ctx context.Context,
	tx datastore.DB,
	rp *RelatedPerson,
) (int, error) {
	if tx != nil {
		return r.upsertRelatedPersonTx(ctx, tx, rp)
	}

	var id int
	err := datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		var err error
		id, err = r.upsertRelatedPersonTx(ctx, txn, rp)
		return err
	})
	return id, err
}

func (r *Repository) upsertRelatedPersonTx(
	ctx context.Context,
	tx datastore.DB,
	rp *RelatedPerson,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(
		RelatedPerson{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		RelatedPerson{},
		[]string{"created_at", "updated_at"},
	)

	query := fmt.Sprintf(`
		INSERT INTO related_persons (%s)
		VALUES (%s)
		ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)
	result, err := tx.NamedExecContext(ctx, query, rp)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert related person: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get last insert ID for related person: %w",
			err,
		)
	}

	return int(lastID), nil
}

func (r *Repository) UpsertStudentRelatedPerson(
	ctx context.Context,
	tx datastore.DB,
	srp *StudentRelatedPerson,
) error {
	if tx != nil {
		return r.upsertStudentRelatedPersonTx(ctx, tx, srp)
	}

	return datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		return r.upsertStudentRelatedPersonTx(ctx, txn, srp)
	})
}

func (r *Repository) upsertStudentRelatedPersonTx(
	ctx context.Context,
	tx datastore.DB,
	srp *StudentRelatedPerson,
) error {
	cols, vals := datastore.GetInsertStatement(
		StudentRelatedPerson{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentRelatedPerson{},
		[]string{"created_at", "updated_at"},
	)

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
}

func (r *Repository) DeleteStudentRelatedPersons(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	if tx != nil {
		return r.deleteStudentRelatedPersonsTx(ctx, tx, iirID)
	}

	return datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		return r.deleteStudentRelatedPersonsTx(ctx, txn, iirID)
	})
}

func (r *Repository) deleteStudentRelatedPersonsTx(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	query := `DELETE FROM student_related_persons WHERE iir_id = ?`
	_, err := tx.ExecContext(ctx, query, iirID)
	if err != nil {
		return fmt.Errorf("failed to delete student related persons: %w", err)
	}
	return nil
}

func (r *Repository) UpsertFamilyBackground(
	ctx context.Context,
	tx datastore.DB,
	fb *FamilyBackground,
) (int, error) {
	if tx != nil {
		return r.upsertFamilyBackgroundTx(ctx, tx, fb)
	}

	var id int
	err := datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		var err error
		id, err = r.upsertFamilyBackgroundTx(ctx, txn, fb)
		return err
	})
	return id, err
}

func (r *Repository) upsertFamilyBackgroundTx(
	ctx context.Context,
	tx datastore.DB,
	fb *FamilyBackground,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(
		FamilyBackground{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		FamilyBackground{},
		[]string{"created_at", "updated_at", "iir_id"},
	)

	query := fmt.Sprintf(`
		INSERT INTO family_backgrounds (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)
	result, err := tx.NamedExecContext(ctx, query, fb)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert family background: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get last insert ID for family background: %w",
			err,
		)
	}

	return int(lastID), nil
}

func (r *Repository) CreateStudentSiblingSupport(
	ctx context.Context,
	tx datastore.DB,
	sss *StudentSiblingSupport,
) error {
	if tx != nil {
		return r.createStudentSiblingSupportTx(ctx, tx, sss)
	}

	return datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		return r.createStudentSiblingSupportTx(ctx, txn, sss)
	})
}

func (r *Repository) createStudentSiblingSupportTx(
	ctx context.Context,
	tx datastore.DB,
	sss *StudentSiblingSupport,
) error {
	cols, vals := datastore.GetInsertStatement(
		StudentSiblingSupport{},
		[]string{"created_at", "updated_at"},
	)
	query := fmt.Sprintf(`
		INSERT INTO student_sibling_supports (%s) VALUES (%s)
	`, cols, vals)
	_, err := tx.NamedExecContext(ctx, query, sss)
	if err != nil {
		return fmt.Errorf("failed to create student sibling support: %w", err)
	}
	return nil
}

func (r *Repository) DeleteStudentSiblingSupportsByFamilyID(
	ctx context.Context,
	tx datastore.DB,
	familyBackgroundID int,
) error {
	if tx != nil {
		return r.deleteStudentSiblingSupportsByFamilyIDTx(
			ctx,
			tx,
			familyBackgroundID,
		)
	}

	return datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		return r.deleteStudentSiblingSupportsByFamilyIDTx(
			ctx,
			txn,
			familyBackgroundID,
		)
	})
}

func (r *Repository) deleteStudentSiblingSupportsByFamilyIDTx(
	ctx context.Context,
	tx datastore.DB,
	familyBackgroundID int,
) error {
	query := `DELETE FROM student_sibling_supports WHERE family_background_id = ?`
	_, err := tx.ExecContext(ctx, query, familyBackgroundID)
	if err != nil {
		return fmt.Errorf("failed to delete student sibling supports: %w", err)
	}
	return nil
}

func (r *Repository) UpsertEducationalBackground(
	ctx context.Context,
	tx datastore.DB,
	eb *EducationalBackground,
) (int, error) {
	if tx != nil {
		return r.upsertEducationalBackgroundTx(ctx, tx, eb)
	}

	var id int
	err := datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		var err error
		id, err = r.upsertEducationalBackgroundTx(ctx, txn, eb)
		return err
	})
	return id, err
}

func (r *Repository) upsertEducationalBackgroundTx(
	ctx context.Context,
	tx datastore.DB,
	eb *EducationalBackground,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(
		EducationalBackground{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		EducationalBackground{},
		[]string{"created_at", "updated_at", "iir_id"},
	)

	query := fmt.Sprintf(`
		INSERT INTO educational_backgrounds (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)
	result, err := tx.NamedExecContext(ctx, query, eb)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert educational background: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get last insert ID for educational background: %w",
			err,
		)
	}

	return int(lastID), nil
}

func (r *Repository) UpsertSchoolDetails(
	ctx context.Context,
	tx datastore.DB,
	sd *SchoolDetails,
) (int, error) {
	if tx != nil {
		return r.upsertSchoolDetailsTx(ctx, tx, sd)
	}

	var id int
	err := datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		var err error
		id, err = r.upsertSchoolDetailsTx(ctx, txn, sd)
		return err
	})
	return id, err
}

func (r *Repository) upsertSchoolDetailsTx(
	ctx context.Context,
	tx datastore.DB,
	sd *SchoolDetails,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(
		SchoolDetails{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		SchoolDetails{},
		[]string{"created_at", "updated_at"},
	)

	query := fmt.Sprintf(`
		INSERT INTO school_details (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)
	result, err := tx.NamedExecContext(ctx, query, sd)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert school details: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get last insert ID for school details: %w",
			err,
		)
	}

	return int(lastID), nil
}

func (r *Repository) DeleteSchoolDetailsByEBID(
	ctx context.Context,
	tx datastore.DB,
	ebID int,
) error {
	if tx != nil {
		return r.deleteSchoolDetailsByEBIDTx(ctx, tx, ebID)
	}

	return datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		return r.deleteSchoolDetailsByEBIDTx(ctx, txn, ebID)
	})
}

func (r *Repository) deleteSchoolDetailsByEBIDTx(
	ctx context.Context,
	tx datastore.DB,
	ebID int,
) error {
	query := `DELETE FROM school_details WHERE eb_id = ?`
	_, err := tx.ExecContext(ctx, query, ebID)
	if err != nil {
		return fmt.Errorf("failed to delete school details: %w", err)
	}
	return nil
}

func (r *Repository) UpsertStudentHealthRecord(
	ctx context.Context,
	tx datastore.DB,
	hr *StudentHealthRecord,
) (int, error) {
	if tx != nil {
		return r.upsertStudentHealthRecordTx(ctx, tx, hr)
	}

	var id int
	err := datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		var err error
		id, err = r.upsertStudentHealthRecordTx(ctx, txn, hr)
		return err
	})
	return id, err
}

func (r *Repository) upsertStudentHealthRecordTx(
	ctx context.Context,
	tx datastore.DB,
	hr *StudentHealthRecord,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(
		StudentHealthRecord{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentHealthRecord{},
		[]string{"created_at", "updated_at", "iir_id"},
	)

	query := fmt.Sprintf(`
		INSERT INTO student_health_records (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)
	result, err := tx.NamedExecContext(ctx, query, hr)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert student health record: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get last insert ID for student health record: %w",
			err,
		)
	}

	return int(lastID), nil
}

func (r *Repository) UpsertStudentConsultation(
	ctx context.Context,
	tx datastore.DB,
	sc *StudentConsultation,
) (int, error) {
	if tx != nil {
		return r.upsertStudentConsultationTx(ctx, tx, sc)
	}

	var id int
	err := datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		var err error
		id, err = r.upsertStudentConsultationTx(ctx, txn, sc)
		return err
	})
	return id, err
}

func (r *Repository) upsertStudentConsultationTx(
	ctx context.Context,
	tx datastore.DB,
	sc *StudentConsultation,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(
		StudentConsultation{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentConsultation{},
		[]string{"created_at", "updated_at", "iir_id"},
	)

	query := fmt.Sprintf(`
		INSERT INTO student_consultations (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)
	result, err := tx.NamedExecContext(ctx, query, sc)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert student consultation: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get last insert ID for student consultation: %w",
			err,
		)
	}

	return int(lastID), nil
}

func (r *Repository) UpsertStudentFinance(
	ctx context.Context,
	tx datastore.DB,
	sf *StudentFinance,
) (int, error) {
	if tx != nil {
		return r.upsertStudentFinanceTx(ctx, tx, sf)
	}

	var id int
	err := datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		var err error
		id, err = r.upsertStudentFinanceTx(ctx, txn, sf)
		return err
	})
	return id, err
}

func (r *Repository) upsertStudentFinanceTx(
	ctx context.Context,
	tx datastore.DB,
	sf *StudentFinance,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(
		StudentFinance{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentFinance{},
		[]string{"created_at", "updated_at", "iir_id"},
	)

	query := fmt.Sprintf(`
		INSERT INTO student_finances (%s)
		VALUES (%s)
		ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)
	result, err := tx.NamedExecContext(ctx, query, sf)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert student finance: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get last insert ID for student finance: %w",
			err,
		)
	}

	return int(lastID), nil
}

func (r *Repository) CreateStudentFinancialSupport(
	ctx context.Context,
	tx datastore.DB,
	sfs *StudentFinancialSupport,
) error {
	if tx != nil {
		return r.createStudentFinancialSupportTx(ctx, tx, sfs)
	}

	return datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		return r.createStudentFinancialSupportTx(ctx, txn, sfs)
	})
}

func (r *Repository) createStudentFinancialSupportTx(
	ctx context.Context,
	tx datastore.DB,
	sfs *StudentFinancialSupport,
) error {
	cols, vals := datastore.GetInsertStatement(
		StudentFinancialSupport{},
		[]string{"created_at", "updated_at"},
	)
	query := fmt.Sprintf(`
		INSERT INTO student_financial_supports (%s)
		VALUES (%s)
	`, cols, vals)
	_, err := tx.NamedExecContext(ctx, query, sfs)
	if err != nil {
		return fmt.Errorf("failed to create student financial support: %w", err)
	}
	return nil
}

func (r *Repository) DeleteStudentFinancialSupportsByFinanceID(
	ctx context.Context,
	tx datastore.DB,
	financeID int,
) error {
	if tx != nil {
		return r.deleteStudentFinancialSupportsByFinanceIDTx(ctx, tx, financeID)
	}

	return datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		return r.deleteStudentFinancialSupportsByFinanceIDTx(
			ctx,
			txn,
			financeID,
		)
	})
}

func (r *Repository) deleteStudentFinancialSupportsByFinanceIDTx(
	ctx context.Context,
	tx datastore.DB,
	financeID int,
) error {
	query := `DELETE FROM student_financial_supports WHERE sf_id = ?`
	_, err := tx.ExecContext(ctx, query, financeID)
	if err != nil {
		return fmt.Errorf(
			"failed to delete student financial supports: %w",
			err,
		)
	}
	return nil
}

func (r *Repository) CreateStudentActivity(
	ctx context.Context,
	tx datastore.DB,
	sa *StudentActivity,
) (int, error) {
	if tx != nil {
		return r.createStudentActivityTx(ctx, tx, sa)
	}

	var id int
	err := datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		var err error
		id, err = r.createStudentActivityTx(ctx, txn, sa)
		return err
	})
	return id, err
}

func (r *Repository) createStudentActivityTx(
	ctx context.Context,
	tx datastore.DB,
	sa *StudentActivity,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(
		StudentActivity{},
		[]string{"created_at", "updated_at"},
	)

	query := fmt.Sprintf(`
		INSERT INTO student_activities (%s)
		VALUES (%s)
	`, cols, vals)
	result, err := tx.NamedExecContext(ctx, query, sa)
	if err != nil {
		return 0, fmt.Errorf("failed to create student activity: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get last insert ID for student activity: %w",
			err,
		)
	}

	return int(lastID), nil
}

func (r *Repository) DeleteStudentActivitiesByIIRID(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	if tx != nil {
		return r.deleteStudentActivitiesByIIRIDTx(ctx, tx, iirID)
	}

	return datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		return r.deleteStudentActivitiesByIIRIDTx(ctx, txn, iirID)
	})
}

func (r *Repository) deleteStudentActivitiesByIIRIDTx(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	query := `DELETE FROM student_activities WHERE iir_id = ?`
	_, err := tx.ExecContext(ctx, query, iirID)
	if err != nil {
		return fmt.Errorf("failed to delete student activities: %w", err)
	}
	return nil
}

func (r *Repository) CreateStudentSubjectPreference(
	ctx context.Context,
	tx datastore.DB,
	ssp *StudentSubjectPreference,
) (int, error) {
	if tx != nil {
		return r.createStudentSubjectPreferenceTx(ctx, tx, ssp)
	}

	var id int
	err := datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		var err error
		id, err = r.createStudentSubjectPreferenceTx(ctx, txn, ssp)
		return err
	})
	return id, err
}

func (r *Repository) createStudentSubjectPreferenceTx(
	ctx context.Context,
	tx datastore.DB,
	ssp *StudentSubjectPreference,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(
		StudentSubjectPreference{},
		[]string{"created_at", "updated_at"},
	)

	query := fmt.Sprintf(`
		INSERT INTO student_subject_preferences (%s)
		VALUES (%s)
	`, cols, vals)
	result, err := tx.NamedExecContext(ctx, query, ssp)
	if err != nil {
		return 0, fmt.Errorf(
			"failed to create student subject preference: %w",
			err,
		)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get last insert ID for student subject preference: %w",
			err,
		)
	}

	return int(lastID), nil
}

func (r *Repository) DeleteStudentSubjectPreferencesByIIRID(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	if tx != nil {
		return r.deleteStudentSubjectPreferencesByIIRIDTx(ctx, tx, iirID)
	}

	return datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		return r.deleteStudentSubjectPreferencesByIIRIDTx(ctx, txn, iirID)
	})
}

func (r *Repository) deleteStudentSubjectPreferencesByIIRIDTx(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	query := `DELETE FROM student_subject_preferences WHERE iir_id = ?`
	_, err := tx.ExecContext(ctx, query, iirID)
	if err != nil {
		return fmt.Errorf(
			"failed to delete student subject preferences: %w",
			err,
		)
	}
	return nil
}

func (r *Repository) CreateStudentHobby(
	ctx context.Context,
	tx datastore.DB,
	sh *StudentHobby,
) (int, error) {
	if tx != nil {
		return r.createStudentHobbyTx(ctx, tx, sh)
	}

	var id int
	err := datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		var err error
		id, err = r.createStudentHobbyTx(ctx, txn, sh)
		return err
	})
	return id, err
}

func (r *Repository) createStudentHobbyTx(
	ctx context.Context,
	tx datastore.DB,
	sh *StudentHobby,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(
		StudentHobby{},
		[]string{"created_at", "updated_at"},
	)

	query := fmt.Sprintf(`
		INSERT INTO student_hobbies (%s)
		VALUES (%s)
	`, cols, vals)
	result, err := tx.NamedExecContext(ctx, query, sh)
	if err != nil {
		return 0, fmt.Errorf("failed to create student hobby: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get last insert ID for student hobby: %w",
			err,
		)
	}

	return int(lastID), nil
}

func (r *Repository) DeleteStudentHobbiesByIIRID(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	if tx != nil {
		return r.deleteStudentHobbiesByIIRIDTx(ctx, tx, iirID)
	}

	return datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		return r.deleteStudentHobbiesByIIRIDTx(ctx, txn, iirID)
	})
}

func (r *Repository) deleteStudentHobbiesByIIRIDTx(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	query := `DELETE FROM student_hobbies WHERE iir_id = ?`
	_, err := tx.ExecContext(ctx, query, iirID)
	if err != nil {
		return fmt.Errorf("failed to delete student hobbies: %w", err)
	}
	return nil
}

func (r *Repository) CreateTestResult(
	ctx context.Context,
	tx datastore.DB,
	tr *TestResult,
) (int, error) {
	if tx != nil {
		return r.createTestResultTx(ctx, tx, tr)
	}

	var id int
	err := datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		var err error
		id, err = r.createTestResultTx(ctx, txn, tr)
		return err
	})
	return id, err
}

func (r *Repository) createTestResultTx(
	ctx context.Context,
	tx datastore.DB,
	tr *TestResult,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(
		TestResult{},
		[]string{"created_at", "updated_at"},
	)

	query := fmt.Sprintf(`
		INSERT INTO test_results (%s)
		VALUES (%s)
	`, cols, vals)
	result, err := tx.NamedExecContext(ctx, query, tr)
	if err != nil {
		return 0, fmt.Errorf("failed to create test result: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get last insert ID for test result: %w",
			err,
		)
	}

	return int(lastID), nil
}

func (r *Repository) DeleteTestResultsByIIRID(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	if tx != nil {
		return r.deleteTestResultsByIIRIDTx(ctx, tx, iirID)
	}

	return datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		return r.deleteTestResultsByIIRIDTx(ctx, txn, iirID)
	})
}

func (r *Repository) deleteTestResultsByIIRIDTx(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	query := `DELETE FROM test_results WHERE iir_id = ?`
	_, err := tx.ExecContext(ctx, query, iirID)
	if err != nil {
		return fmt.Errorf("failed to delete test results: %w", err)
	}
	return nil
}

func (r *Repository) DeleteSignificantNotesByIIRID(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	if tx != nil {
		return r.deleteSignificantNotesByIIRIDTx(ctx, tx, iirID)
	}

	return datastore.RunInTransaction(ctx, r.db, func(txn datastore.DB) error {
		return r.deleteSignificantNotesByIIRIDTx(ctx, txn, iirID)
	})
}

func (r *Repository) deleteSignificantNotesByIIRIDTx(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	query := `DELETE FROM significant_notes WHERE iir_id = ?`
	_, err := tx.ExecContext(ctx, query, iirID)
	if err != nil {
		return fmt.Errorf("failed to delete significant notes: %w", err)
	}
	return nil
}
