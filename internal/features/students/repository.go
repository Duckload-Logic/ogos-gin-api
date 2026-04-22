package students

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) RepositoryInterface {
	return &Repository{db: db}
}

func (r *Repository) GetDB() *sqlx.DB {
	return r.db
}

func (r *Repository) BeginTx(ctx context.Context) (datastore.DB, error) {
	return r.db.BeginTxx(ctx, nil)
}

func (r *Repository) WithTransaction(
	ctx context.Context,
	fn func(datastore.DB) error,
) error {
	return datastore.RunInTransaction(ctx, r.db, fn)
}

// Lookup
func (r *Repository) GetGenders(ctx context.Context) ([]Gender, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM genders ORDER BY id
	`, datastore.GetColumns(GenderDB{}))

	var dbModels []GenderDB
	err := r.db.SelectContext(ctx, &dbModels, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get genders: %w", err)
	}

	genders := make([]Gender, len(dbModels))
	for i, v := range dbModels {
		genders[i] = v.ToDomain()
	}

	return genders, nil
}

func (r *Repository) GetParentalStatusTypes(
	ctx context.Context,
) ([]ParentalStatusType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM parental_status_types ORDER BY id
	`, datastore.GetColumns(ParentalStatusTypeDB{}))

	var dbModels []ParentalStatusTypeDB
	err := r.db.SelectContext(ctx, &dbModels, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get parental status types: %w", err)
	}

	statuses := make([]ParentalStatusType, len(dbModels))
	for i, v := range dbModels {
		statuses[i] = v.ToDomain()
	}

	return statuses, nil
}

func (r *Repository) GetEnrollmentReasons(
	ctx context.Context,
) ([]EnrollmentReason, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM enrollment_reasons ORDER BY id
	`, datastore.GetColumns(EnrollmentReasonDB{}))

	var dbModels []EnrollmentReasonDB
	err := r.db.SelectContext(ctx, &dbModels, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment reasons: %w", err)
	}

	reasons := make([]EnrollmentReason, len(dbModels))
	for i, v := range dbModels {
		reasons[i] = v.ToDomain()
	}

	return reasons, nil
}

func (r *Repository) GetIncomeRanges(
	ctx context.Context,
) ([]IncomeRange, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM income_ranges ORDER BY id
	`, datastore.GetColumns(IncomeRangeDB{}))

	var dbModels []IncomeRangeDB
	err := r.db.SelectContext(ctx, &dbModels, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get income ranges: %w", err)
	}

	ranges := make([]IncomeRange, len(dbModels))
	for i, v := range dbModels {
		ranges[i] = v.ToDomain()
	}

	return ranges, nil
}

func (r *Repository) GetStudentSupportTypes(
	ctx context.Context,
) ([]StudentSupportType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_support_types ORDER BY id
	`, datastore.GetColumns(StudentSupportTypeDB{}))

	var dbModels []StudentSupportTypeDB
	err := r.db.SelectContext(ctx, &dbModels, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get student support types: %w", err)
	}

	supportTypes := make([]StudentSupportType, len(dbModels))
	for i, v := range dbModels {
		supportTypes[i] = v.ToDomain()
	}

	return supportTypes, nil
}

func (r *Repository) GetSiblingSupportTypes(
	ctx context.Context,
) ([]SibilingSupportType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM sibling_support_types ORDER BY id
	`, datastore.GetColumns(SibilingSupportTypeDB{}))

	var dbModels []SibilingSupportTypeDB
	err := r.db.SelectContext(ctx, &dbModels, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get sibling support types: %w", err)
	}

	supportTypes := make([]SibilingSupportType, len(dbModels))
	for i, v := range dbModels {
		supportTypes[i] = v.ToDomain()
	}

	return supportTypes, nil
}

func (r *Repository) GetEducationalLevels(
	ctx context.Context,
) ([]EducationalLevel, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM educational_levels ORDER BY id
	`, datastore.GetColumns(EducationalLevelDB{}))

	var dbModels []EducationalLevelDB
	err := r.db.SelectContext(ctx, &dbModels, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get educational levels: %w", err)
	}

	levels := make([]EducationalLevel, len(dbModels))
	for i, v := range dbModels {
		levels[i] = v.ToDomain()
	}

	return levels, nil
}

// GetStudentStatuses retrieves all available student statuses.
func (r *Repository) GetStudentStatuses(
	ctx context.Context,
) ([]StudentStatus, error) {
	var dbModels []StudentStatusDB
	err := r.db.SelectContext(
		ctx,
		&dbModels,
		"SELECT id, status_name FROM student_statuses",
	)

	statuses := make([]StudentStatus, len(dbModels))
	for i, v := range dbModels {
		statuses[i] = v.ToDomain()
	}

	return statuses, err
}

func (r *Repository) GetCourses(ctx context.Context) ([]Course, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM courses ORDER BY id
	`, datastore.GetColumns(CourseDB{}))

	var dbModels []CourseDB
	err := r.db.SelectContext(ctx, &dbModels, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}

	courses := make([]Course, len(dbModels))
	for i, v := range dbModels {
		courses[i] = v.ToDomain()
	}

	return courses, nil
}

func (r *Repository) GetCivilStatusTypes(
	ctx context.Context,
) ([]CivilStatusType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM civil_status_types ORDER BY id
	`, datastore.GetColumns(CivilStatusTypeDB{}))

	var dbModels []CivilStatusTypeDB
	err := r.db.SelectContext(ctx, &dbModels, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get civil status types: %w", err)
	}

	statuses := make([]CivilStatusType, len(dbModels))
	for i, v := range dbModels {
		statuses[i] = v.ToDomain()
	}

	return statuses, nil
}

func (r *Repository) GetReligions(ctx context.Context) ([]Religion, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM religions ORDER BY id
	`, datastore.GetColumns(ReligionDB{}))

	var dbModels []ReligionDB
	err := r.db.SelectContext(ctx, &dbModels, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get religions: %w", err)
	}

	religions := make([]Religion, len(dbModels))
	for i, v := range dbModels {
		religions[i] = v.ToDomain()
	}

	return religions, nil
}

func (r *Repository) GetStudentRelationshipTypes(
	ctx context.Context,
) ([]StudentRelationshipType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_relationship_types ORDER BY id
	`, datastore.GetColumns(StudentRelationshipTypeDB{}))

	var dbModels []StudentRelationshipTypeDB
	err := r.db.SelectContext(ctx, &dbModels, query)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student relationship types: %w",
			err,
		)
	}

	relationships := make([]StudentRelationshipType, len(dbModels))
	for i, v := range dbModels {
		relationships[i] = v.ToDomain()
	}

	return relationships, nil
}

func (r *Repository) GetNatureOfResidenceTypes(
	ctx context.Context,
) ([]NatureOfResidenceType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM nature_of_residence_types ORDER BY id
	`, datastore.GetColumns(NatureOfResidenceTypeDB{}))

	var dbModels []NatureOfResidenceTypeDB
	err := r.db.SelectContext(ctx, &dbModels, query)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get nature of residence types: %w",
			err,
		)
	}

	residences := make([]NatureOfResidenceType, len(dbModels))
	for i, v := range dbModels {
		residences[i] = v.ToDomain()
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
	statusID int,
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
		statusID,
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
	courseID, genderID, yearLevel, statusID int,
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

	if statusID > 0 {
		query += " AND spi.status_id = ?"
		args = append(args, statusID)
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
	ctx context.Context,
	search string,
	offset, limit int,
	orderBy string,
	courseID, genderID, yearLevel, statusID int,
) ([]StudentProfileView, error) {
	query := `
        SELECT
			iir.id as iir_id,
			u.id as user_id,
			u.first_name,
			u.middle_name,
			u.last_name,
			u.suffix_name,
			spi.gender_id,
			u.email,
			spi.student_number,
			spi.course_id,
			spi.section,
			spi.year_level,
			spi.status_id,
			ss.status_name
		FROM iir_records iir
		JOIN users u ON iir.user_id = u.id
		JOIN student_personal_info spi ON iir.id = spi.iir_id
		LEFT JOIN student_statuses ss ON spi.status_id = ss.id
		WHERE iir.is_submitted = TRUE
    `

	query, args := r.applyStudentFilters(
		query,
		nil,
		search,
		courseID,
		genderID,
		yearLevel,
		statusID,
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
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	var dbModels []StudentProfileViewDB
	err := r.db.SelectContext(ctx, &dbModels, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list students: %w", err)
	}

	views := make([]StudentProfileView, len(dbModels))
	for i, v := range dbModels {
		views[i] = v.ToDomain()
	}

	return views, nil
}

func (r *Repository) GetStudentBasicInfo(
	ctx context.Context,
	iirID string,
) (*StudentBasicInfoView, error) {
	query := `
		SELECT
			u.id as user_id,
			u.email,
			u.first_name as first_name,
			u.middle_name as middle_name,
			u.last_name as last_name,
			u.suffix_name as suffix_name
		FROM users u
		JOIN iir_records iir ON u.id = iir.user_id
		WHERE iir.id = ?
	`

	var dbModel StudentBasicInfoViewDB
	err := r.db.GetContext(ctx, &dbModel, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student basic info: %w", err)
	}

	domainModel := dbModel.ToDomain()
	return &domainModel, nil
}

func (r *Repository) GetIIRDraftByUserID(
	ctx context.Context,
	userID string,
) (*IIRDraft, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM iir_drafts WHERE user_id = ? LIMIT 1
	`, datastore.GetColumns(IIRDraftDB{}))

	var dbModel IIRDraftDB
	err := r.db.GetContext(ctx, &dbModel, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	domainModel := dbModel.ToDomain()
	return &domainModel, nil
}

func (r *Repository) GetStudentIIRByUserID(
	ctx context.Context,
	userID string,
) (*IIRRecord, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM iir_records WHERE user_id = ? LIMIT 1
	`, datastore.GetColumns(IIRRecordDB{}))

	var dbModel IIRRecordDB
	err := r.db.GetContext(ctx, &dbModel, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	domainModel := dbModel.ToDomain()
	return &domainModel, nil
}

func (r *Repository) GetStudentIIR(
	ctx context.Context,
	iirID string,
) (*IIRRecord, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM iir_records WHERE id = ? LIMIT 1
	`, datastore.GetColumns(IIRRecordDB{}))

	var dbModel IIRRecordDB
	err := r.db.GetContext(ctx, &dbModel, query, iirID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	domainModel := dbModel.ToDomain()
	return &domainModel, nil
}

func (r *Repository) GetStudentPersonalInfo(
	ctx context.Context,
	iirID string,
) (*StudentPersonalInfo, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM student_personal_info
		WHERE iir_id = ?
	`, datastore.GetColumns(StudentPersonalInfoDB{}))

	var dbModel StudentPersonalInfoDB
	err := r.db.GetContext(ctx, &dbModel, query, iirID)
	if err != nil {
		return nil, err
	}

	domainModel := dbModel.ToDomain()

	return &domainModel, nil
}

func (r *Repository) GetEmergencyContactByIIRID(
	ctx context.Context,
	iirID string,
) (*EmergencyContact, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM emergency_contacts WHERE iir_id = ? LIMIT 1
	`, datastore.GetColumns(EmergencyContactDB{}))

	var dbModel EmergencyContactDB
	err := r.db.GetContext(ctx, &dbModel, query, iirID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get emergency contact: %w", err)
	}

	domainModel := dbModel.ToDomain()
	return &domainModel, nil
}

func (r *Repository) GetGenderByID(
	ctx context.Context,
	genderID int,
) (*Gender, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM genders WHERE id = ?
	`, datastore.GetColumns(GenderDB{}))

	var dbModel GenderDB
	err := r.db.GetContext(ctx, &dbModel, query, genderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gender by ID: %w", err)
	}

	gender := dbModel.ToDomain()

	return &gender, nil
}

func (r *Repository) GetCivilStatusByID(
	ctx context.Context,
	statusID int,
) (*CivilStatusType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM civil_status_types WHERE id = ?
	`, datastore.GetColumns(CivilStatusTypeDB{}))

	var dbModel CivilStatusTypeDB
	err := r.db.GetContext(ctx, &dbModel, query, statusID)
	if err != nil {
		return nil, fmt.Errorf("failed to get civil status by ID: %w", err)
	}

	status := dbModel.ToDomain()

	return &status, nil
}

func (r *Repository) GetReligionByID(
	ctx context.Context,
	religionID int,
) (*Religion, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM religions WHERE id = ?
	`, datastore.GetColumns(ReligionDB{}))

	var dbModel ReligionDB
	err := r.db.GetContext(ctx, &dbModel, query, religionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get religion by ID: %w", err)
	}

	religion := dbModel.ToDomain()

	return &religion, nil
}

func (r *Repository) GetCourseByID(
	ctx context.Context,
	courseID int,
) (*Course, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM courses WHERE id = ?
	`, datastore.GetColumns(CourseDB{}))

	var dbModel CourseDB
	err := r.db.GetContext(ctx, &dbModel, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course by ID: %w", err)
	}

	course := dbModel.ToDomain()

	return &course, nil
}

func (r *Repository) GetStudentAddresses(
	ctx context.Context,
	iirID string,
) ([]StudentAddress, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_addresses WHERE iir_id = ?
	`, datastore.GetColumns(StudentAddressDB{}))

	var dbModels []StudentAddressDB
	err := r.db.SelectContext(ctx, &dbModels, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student addresses: %w", err)
	}

	addresses := make([]StudentAddress, len(dbModels))
	for i, v := range dbModels {
		addresses[i] = v.ToDomain()
	}

	return addresses, nil
}

func (r *Repository) GetStudentEducationalBackground(
	ctx context.Context,
	iirID string,
) (*EducationalBackground, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM educational_backgrounds WHERE iir_id = ? LIMIT 1
	`, datastore.GetColumns(EducationalBackgroundDB{}))

	var dbModel EducationalBackgroundDB
	err := r.db.GetContext(ctx, &dbModel, query, iirID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf(
			"failed to get student educational background: %w",
			err,
		)
	}

	domainModel := dbModel.ToDomain()
	return &domainModel, nil
}

func (r *Repository) GetSchoolDetailsByEBID(
	ctx context.Context,
	ebID int,
) ([]SchoolDetails, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM school_details WHERE eb_id = ?
		ORDER BY educational_level_id ASC
	`, datastore.GetColumns(SchoolDetailsDB{}))

	var dbModels []SchoolDetailsDB
	err := r.db.SelectContext(ctx, &dbModels, query, ebID)
	if err != nil {
		return nil, fmt.Errorf("failed to get school details: %w", err)
	}

	details := make([]SchoolDetails, len(dbModels))
	for i, v := range dbModels {
		details[i] = v.ToDomain()
	}

	return details, nil
}

func (r *Repository) GetEducationalLevelByID(
	ctx context.Context,
	levelID int,
) (*EducationalLevel, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM educational_levels WHERE id = ?
	`, datastore.GetColumns(EducationalLevelDB{}))
	var dbModel EducationalLevelDB
	err := r.db.GetContext(ctx, &dbModel, query, levelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get educational level by ID: %w", err)
	}

	domainModel := dbModel.ToDomain()

	return &domainModel, nil
}

func (r *Repository) GetStudentRelatedPersons(
	ctx context.Context, iirID string,
) ([]StudentRelatedPerson, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_related_persons WHERE iir_id = ?
	`, datastore.GetColumns(StudentRelatedPersonDB{}))

	var dbModels []StudentRelatedPersonDB
	err := r.db.SelectContext(ctx, &dbModels, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get related persons: %w", err)
	}

	persons := make([]StudentRelatedPerson, len(dbModels))
	for i, v := range dbModels {
		persons[i] = v.ToDomain()
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
	`, datastore.GetColumns(RelatedPersonDB{}))

	var dbModel RelatedPersonDB
	err := r.db.GetContext(ctx, &dbModel, query, personID)
	if err != nil {
		return nil, fmt.Errorf("failed to get related person by ID: %w", err)
	}

	domainModel := dbModel.ToDomain()

	return &domainModel, nil
}

func (r *Repository) GetStudentRelationshipByID(
	ctx context.Context, relationshipID int,
) (*StudentRelationshipType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_relationship_types WHERE id = ?
	`, datastore.GetColumns(StudentRelationshipTypeDB{}))

	var dbModel StudentRelationshipTypeDB
	err := r.db.GetContext(ctx, &dbModel, query, relationshipID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student relationship by ID: %w",
			err,
		)
	}

	domainModel := dbModel.ToDomain()

	return &domainModel, nil
}

func (r *Repository) GetStudentFamilyBackground(
	ctx context.Context,
	iirID string,
) (*FamilyBackground, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM family_backgrounds WHERE iir_id = ? LIMIT 1
	`, datastore.GetColumns(FamilyBackgroundDB{}))

	var dbModel FamilyBackgroundDB
	err := r.db.GetContext(ctx, &dbModel, query, iirID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf(
			"failed to get student family background: %w",
			err,
		)
	}

	domainModel := dbModel.ToDomain()
	return &domainModel, nil
}

func (r *Repository) GetParentalStatusByID(
	ctx context.Context,
	statusID int,
) (*ParentalStatusType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM parental_status_types WHERE id = ?
	`, datastore.GetColumns(ParentalStatusTypeDB{}))

	var dbModel ParentalStatusTypeDB
	err := r.db.GetContext(ctx, &dbModel, query, statusID)
	if err != nil {
		return nil, fmt.Errorf("failed to get parental status by ID: %w", err)
	}

	domainModel := dbModel.ToDomain()

	return &domainModel, nil
}

func (r *Repository) GetNatureOfResidenceByID(
	ctx context.Context,
	residenceID int,
) (*NatureOfResidenceType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM nature_of_residence_types WHERE id = ?
	`, datastore.GetColumns(NatureOfResidenceTypeDB{}))

	var dbModel NatureOfResidenceTypeDB
	err := r.db.GetContext(ctx, &dbModel, query, residenceID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get nature of residence by ID: %w",
			err,
		)
	}

	domainModel := dbModel.ToDomain()

	return &domainModel, nil
}

func (r *Repository) GetStudentSiblingSupport(
	ctx context.Context,
	fbID int,
) ([]StudentSiblingSupport, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_sibling_supports WHERE family_background_id = ?
	`, datastore.GetColumns(StudentSiblingSupportDB{}))

	var dbModels []StudentSiblingSupportDB
	err := r.db.SelectContext(ctx, &dbModels, query, fbID)
	if err != nil {
		return nil, fmt.Errorf("failed to query sibling supports: %w", err)
	}

	supports := make([]StudentSiblingSupport, len(dbModels))
	for i, v := range dbModels {
		supports[i] = v.ToDomain()
	}

	return supports, nil
}

func (r *Repository) GetSiblingSupportTypeByID(
	ctx context.Context,
	supportID int,
) (*SibilingSupportType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM sibling_support_types WHERE id = ?
	`, datastore.GetColumns(SibilingSupportTypeDB{}))

	var dbModel SibilingSupportTypeDB
	err := r.db.GetContext(ctx, &dbModel, query, supportID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get sibling support type by ID: %w",
			err,
		)
	}

	domainModel := dbModel.ToDomain()

	return &domainModel, nil
}

func (r *Repository) GetStudentFinancialInfo(
	ctx context.Context,
	iirID string,
) (*StudentFinance, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_finances WHERE iir_id = ? LIMIT 1
	`, datastore.GetColumns(StudentFinanceDB{}))

	var dbModel StudentFinanceDB
	err := r.db.GetContext(ctx, &dbModel, query, iirID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get student financial info: %w", err)
	}

	domainModel := dbModel.ToDomain()

	return &domainModel, nil
}

func (r *Repository) GetFinancialSupportTypeByFinanceID(
	ctx context.Context,
	financeID int,
) ([]StudentFinancialSupport, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_financial_supports WHERE sf_id = ?
	`, datastore.GetColumns(StudentFinancialSupportDB{}))

	var dbModels []StudentFinancialSupportDB
	err := r.db.SelectContext(ctx, &dbModels, query, financeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query financial supports: %w", err)
	}

	supports := make([]StudentFinancialSupport, len(dbModels))
	for i, v := range dbModels {
		supports[i] = v.ToDomain()
	}

	return supports, nil
}

func (r *Repository) GetIncomeRangeByID(
	ctx context.Context,
	rangeID int,
) (*IncomeRange, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM income_ranges WHERE id = ?
	`, datastore.GetColumns(IncomeRangeDB{}))

	var dbModel IncomeRangeDB
	err := r.db.GetContext(ctx, &dbModel, query, rangeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get income range by ID: %w", err)
	}

	domainModel := dbModel.ToDomain()

	return &domainModel, nil
}

func (r *Repository) GetStudentSupportByID(
	ctx context.Context,
	supportID int,
) (*StudentSupportType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_support_types WHERE id = ?
	`, datastore.GetColumns(StudentSupportTypeDB{}))

	var dbModel StudentSupportTypeDB
	err := r.db.GetContext(ctx, &dbModel, query, supportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student support by ID: %w", err)
	}

	domainModel := dbModel.ToDomain()

	return &domainModel, nil
}

func (r *Repository) GetStudentHealthRecord(
	ctx context.Context,
	iirID string,
) (*StudentHealthRecord, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_health_records WHERE iir_id = ? LIMIT 1
	`, datastore.GetColumns(StudentHealthRecordDB{}))

	var dbModel StudentHealthRecordDB
	err := r.db.GetContext(ctx, &dbModel, query, iirID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get student health record: %w", err)
	}

	domainModel := dbModel.ToDomain()
	return &domainModel, nil
}

func (r *Repository) GetActivityOptions(
	ctx context.Context,
) ([]ActivityOption, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM activity_options ORDER BY id
	`, datastore.GetColumns(ActivityOptionDB{}))

	var dbModels []ActivityOptionDB
	err := r.db.SelectContext(ctx, &dbModels, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity options: %w", err)
	}

	options := make([]ActivityOption, len(dbModels))
	for i, v := range dbModels {
		options[i] = v.ToDomain()
	}

	return options, nil
}

func (r *Repository) GetStudentConsultations(
	ctx context.Context,
	iirID string,
) ([]StudentConsultation, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_consultations WHERE iir_id = ?
	`, datastore.GetColumns(StudentConsultationDB{}))

	var dbModels []StudentConsultationDB
	err := r.db.SelectContext(ctx, &dbModels, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student consultations: %w", err)
	}

	consultations := make([]StudentConsultation, len(dbModels))
	for i, v := range dbModels {
		consultations[i] = v.ToDomain()
	}

	return consultations, nil
}

func (r *Repository) GetStudentActivities(
	ctx context.Context,
	iirID string,
) ([]StudentActivity, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_activities WHERE iir_id = ?
	`, datastore.GetColumns(StudentActivityDB{}))

	var dbModels []StudentActivityDB
	err := r.db.SelectContext(ctx, &dbModels, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student activities: %w", err)
	}

	activities := make([]StudentActivity, len(dbModels))
	for i, v := range dbModels {
		activities[i] = v.ToDomain()
	}

	return activities, nil
}

func (r *Repository) GetActivityOptionByID(
	ctx context.Context,
	optionID int,
) (*ActivityOption, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM activity_options WHERE id = ?
	`, datastore.GetColumns(ActivityOptionDB{}))

	var dbModel ActivityOptionDB
	err := r.db.GetContext(ctx, &dbModel, query, optionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity option: %w", err)
	}

	domainModel := dbModel.ToDomain()
	return &domainModel, nil
}

func (r *Repository) GetStudentSubjectPreferences(
	ctx context.Context,
	iirID string,
) ([]StudentSubjectPreference, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_subject_preferences WHERE iir_id = ?
	`, datastore.GetColumns(StudentSubjectPreferenceDB{}))

	var dbModels []StudentSubjectPreferenceDB
	err := r.db.SelectContext(ctx, &dbModels, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subject preferences: %w", err)
	}

	preferences := make([]StudentSubjectPreference, len(dbModels))
	for i, v := range dbModels {
		preferences[i] = v.ToDomain()
	}

	return preferences, nil
}

func (r *Repository) GetStudentHobbies(
	ctx context.Context,
	iirID string,
) ([]StudentHobby, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_hobbies WHERE iir_id = ?
	`, datastore.GetColumns(StudentHobbyDB{}))

	var dbModels []StudentHobbyDB
	err := r.db.SelectContext(ctx, &dbModels, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hobbies: %w", err)
	}

	hobbies := make([]StudentHobby, len(dbModels))
	for i, v := range dbModels {
		hobbies[i] = v.ToDomain()
	}

	return hobbies, nil
}

func (r *Repository) GetStudentTestResults(
	ctx context.Context,
	iirID string,
) ([]TestResult, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM test_results WHERE iir_id = ?
	`, datastore.GetColumns(TestResultDB{}))

	var dbModels []TestResultDB
	err := r.db.SelectContext(ctx, &dbModels, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test results: %w", err)
	}

	results := make([]TestResult, len(dbModels))
	for i, v := range dbModels {
		results[i] = v.ToDomain()
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
		IIRDraftDB{},
		[]string{"created_at", "updated_at"},
	)
	onDuplicateKey := datastore.GetOnDuplicateKeyUpdateStatement(
		IIRDraftDB{},
		[]string{"created_at", "updated_at"},
	)
	query := fmt.Sprintf(`
		INSERT INTO iir_drafts (%s)
		VALUES (%s)
		ON DUPLICATE KEY UPDATE %s
	`, cols, vals, onDuplicateKey)

	dbModel := draft.ToPersistence()
	result, err := tx.NamedExecContext(ctx, query, dbModel)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert IIR draft: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
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
		IIRRecordDB{},
		[]string{"created_at", "updated_at"},
	)
	onDuplicateKey := datastore.GetOnDuplicateKeyUpdateStatement(
		IIRRecordDB{},
		[]string{"created_at", "updated_at"},
	)
	query := fmt.Sprintf(`
		INSERT INTO iir_records (id, %s)
		VALUES (:id, %s)
		ON DUPLICATE KEY UPDATE %s
	`, cols, vals, onDuplicateKey)

	dbModel := iir.ToPersistence()
	_, err := tx.NamedExecContext(ctx, query, &dbModel)
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
		StudentPersonalInfoDB{},
		[]string{"created_at", "updated_at", "address_id"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentPersonalInfoDB{},
		[]string{"created_at", "updated_at", "iir_id"},
	)

	query := fmt.Sprintf(`
		INSERT INTO student_personal_info (%s)
		VALUES (%s)
		ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)

	dbModel := info.ToPersistence()
	_, err := tx.NamedExecContext(ctx, query, &dbModel)
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
		EmergencyContactDB{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		EmergencyContactDB{},
		[]string{"created_at", "updated_at", "iir_id"},
	)

	query := fmt.Sprintf(`
		INSERT INTO emergency_contacts (%s)
		VALUES (%s)
		ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)

	dbModel := ec.ToPersistence()
	result, err := tx.NamedExecContext(ctx, query, &dbModel)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert emergency contact: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
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
		StudentAddressDB{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentAddressDB{},
		[]string{"created_at", "updated_at", "iir_id"},
	)

	query := fmt.Sprintf(`
		INSERT INTO student_addresses (%s)
		VALUES (%s)
		ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)

	dbModel := sa.ToPersistence()
	result, err := tx.NamedExecContext(ctx, query, &dbModel)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert student address: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
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
	dbModel := StudentSelectedReasonDB{
		IIRID:           ssr.IIRID,
		ReasonID:        ssr.ReasonID,
		OtherReasonText: ssr.OtherReasonText,
	}
	_, err := tx.NamedExecContext(ctx, query, &dbModel)
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
		RelatedPersonDB{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		RelatedPersonDB{},
		[]string{"created_at", "updated_at"},
	)

	query := fmt.Sprintf(`
		INSERT INTO related_persons (%s)
		VALUES (%s)
		ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)

	dbModel := rp.ToPersistence()
	result, err := tx.NamedExecContext(ctx, query, &dbModel)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert related person: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
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
		StudentRelatedPersonDB{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentRelatedPersonDB{},
		[]string{"created_at", "updated_at"},
	)

	query := fmt.Sprintf(`
		INSERT INTO student_related_persons (%s)
		VALUES (%s)
		ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)

	dbModel := srp.ToPersistence()
	_, err := tx.NamedExecContext(ctx, query, &dbModel)
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
		FamilyBackgroundDB{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		FamilyBackgroundDB{},
		[]string{"created_at", "updated_at", "iir_id"},
	)

	query := fmt.Sprintf(`
		INSERT INTO family_backgrounds (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)

	dbModel := fb.ToPersistence()
	result, err := tx.NamedExecContext(ctx, query, &dbModel)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert family background: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
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
		StudentSiblingSupportDB{},
		[]string{"created_at", "updated_at"},
	)
	query := fmt.Sprintf(`
		INSERT INTO student_sibling_supports (%s) VALUES (%s)
	`, cols, vals)

	dbModel := sss.ToPersistence()
	_, err := tx.NamedExecContext(ctx, query, &dbModel)
	if err != nil {
		return fmt.Errorf("failed to create sibling support: %w", err)
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
		EducationalBackgroundDB{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		EducationalBackgroundDB{},
		[]string{"created_at", "updated_at", "iir_id"},
	)

	query := fmt.Sprintf(`
		INSERT INTO educational_backgrounds (%s)
		VALUES (%s) ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)

	dbModel := eb.ToPersistence()
	result, err := tx.NamedExecContext(ctx, query, &dbModel)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert educational background: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
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
		SchoolDetailsDB{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		SchoolDetailsDB{},
		[]string{"created_at", "updated_at"},
	)

	query := fmt.Sprintf(`
		INSERT INTO school_details (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)

	dbModel := sd.ToPersistence()
	result, err := tx.NamedExecContext(ctx, query, &dbModel)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert school details: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
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
		StudentHealthRecordDB{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentHealthRecordDB{},
		[]string{"created_at", "updated_at", "iir_id"},
	)

	query := fmt.Sprintf(`
		INSERT INTO student_health_records (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)

	dbModel := hr.ToPersistence()
	result, err := tx.NamedExecContext(ctx, query, &dbModel)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert health record: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
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
		StudentConsultationDB{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentConsultationDB{},
		[]string{"created_at", "updated_at", "iir_id"},
	)

	query := fmt.Sprintf(`
		INSERT INTO student_consultations (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)

	dbModel := sc.ToPersistence()
	result, err := tx.NamedExecContext(ctx, query, &dbModel)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert consultation: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
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
		StudentFinanceDB{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentFinanceDB{},
		[]string{"created_at", "updated_at", "iir_id"},
	)

	query := fmt.Sprintf(`
		INSERT INTO student_finances (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)

	dbModel := sf.ToPersistence()
	result, err := tx.NamedExecContext(ctx, query, &dbModel)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert student finance: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
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
		StudentFinancialSupportDB{},
		[]string{"created_at", "updated_at"},
	)
	query := fmt.Sprintf(`
		INSERT INTO student_financial_supports (%s) VALUES (%s)
	`, cols, vals)

	dbModel := sfs.ToPersistence()
	_, err := tx.NamedExecContext(ctx, query, &dbModel)
	if err != nil {
		return fmt.Errorf("failed to create financial support: %w", err)
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
		StudentActivityDB{},
		[]string{"created_at", "updated_at"},
	)

	query := fmt.Sprintf(`
		INSERT INTO student_activities (%s) VALUES (%s)
	`, cols, vals)

	dbModel := sa.ToPersistence()
	result, err := tx.NamedExecContext(ctx, query, &dbModel)
	if err != nil {
		return 0, fmt.Errorf("failed to create student activity: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
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
		StudentSubjectPreferenceDB{},
		[]string{"created_at", "updated_at"},
	)

	query := fmt.Sprintf(`
		INSERT INTO student_subject_preferences (%s) VALUES (%s)
	`, cols, vals)

	dbModel := ssp.ToPersistence()
	result, err := tx.NamedExecContext(ctx, query, &dbModel)
	if err != nil {
		return 0, fmt.Errorf("failed to create subject preference: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
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
		StudentHobbyDB{},
		[]string{"created_at", "updated_at"},
	)

	query := fmt.Sprintf(`
		INSERT INTO student_hobbies (%s) VALUES (%s)
	`, cols, vals)

	dbModel := sh.ToPersistence()
	result, err := tx.NamedExecContext(ctx, query, &dbModel)
	if err != nil {
		return 0, fmt.Errorf("failed to create student hobby: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
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
		TestResultDB{},
		[]string{"created_at", "updated_at"},
	)

	query := fmt.Sprintf(`
		INSERT INTO test_results (%s) VALUES (%s)
	`, cols, vals)

	dbModel := tr.ToPersistence()
	result, err := tx.NamedExecContext(ctx, query, &dbModel)
	if err != nil {
		return 0, fmt.Errorf("failed to create test result: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
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

// IsStudentLocked checks if a student's record is locked (Graduated,
// Archived, or Withdrawn).
func (r *Repository) IsStudentLocked(
	ctx context.Context,
	iirID string,
) (bool, error) {
	var statusID int
	err := r.db.GetContext(ctx, &statusID,
		"SELECT status_id FROM student_personal_info WHERE iir_id = ?", iirID)
	if err != nil {
		return false, err
	}

	// Graduated, 4: Archived, 5: Withdrawn
	if statusID == 2 || statusID == 4 || statusID == 5 {
		return true, nil
	}

	return false, nil
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

// BulkUpdateStudentStatus updates status_id (and optionally graduation_year)
// for a set of students. When SelectAllMatching is true it builds a WHERE
// clause from the request filters instead of an IN list, so every matching
// record across all pages is affected in one query.
func (r *Repository) BulkUpdateStudentStatus(
	ctx context.Context,
	req BulkUpdateStatusRequest,
) error {
	var args []interface{}
	var setClause string

	if req.GraduationYear != nil {
		setClause = "SET status_id = ?, graduation_year = ?"
		args = append(args, req.StatusID, *req.GraduationYear)
	} else {
		setClause = "SET status_id = ?"
		args = append(args, req.StatusID)
	}

	base := fmt.Sprintf(
		"UPDATE student_personal_info %s WHERE",
		setClause,
	)

	if req.SelectAllMatching {
		// Build dynamic WHERE from the supplied filter set.
		var conditions []string
		conditions = append(conditions, "status_id NOT IN (2, 4, 5)")

		if req.Filters.Search != "" {
			// Filter by name/email via subquery join.
			conditions = append(conditions,
				"iir_id IN ("+
					"SELECT i.id FROM iir_records i "+
					"JOIN users u ON u.id = i.user_id "+
					"WHERE CONCAT(u.first_name,' ',u.last_name) "+
					"LIKE ? OR u.email LIKE ?)",
			)
			like := "%" + req.Filters.Search + "%"
			args = append(args, like, like)
		}
		if req.Filters.CourseID > 0 {
			conditions = append(conditions, "course_id = ?")
			args = append(args, req.Filters.CourseID)
		}
		if req.Filters.YearLevel > 0 {
			conditions = append(conditions, "year_level = ?")
			args = append(args, req.Filters.YearLevel)
		}

		// Exclude records the counselor explicitly un-checked.
		if len(req.ExcludedIIRIDs) > 0 {
			placeholders := strings.Repeat("?,", len(req.ExcludedIIRIDs))
			placeholders = placeholders[:len(placeholders)-1]
			conditions = append(conditions,
				fmt.Sprintf("iir_id NOT IN (%s)", placeholders),
			)
			for _, id := range req.ExcludedIIRIDs {
				args = append(args, id)
			}
		}

		query := fmt.Sprintf(
			"%s %s",
			base,
			strings.Join(conditions, " AND "),
		)
		_, err := r.db.ExecContext(ctx, query, args...)
		return err
	}

	// Explicit ID list path.
	if len(req.IIRIDs) == 0 {
		return nil
	}
	placeholders := strings.Repeat("?,", len(req.IIRIDs))
	placeholders = placeholders[:len(placeholders)-1]
	for _, id := range req.IIRIDs {
		args = append(args, id)
	}

	query := fmt.Sprintf(
		"%s iir_id IN (%s) AND status_id NOT IN (2, 4, 5)",
		base,
		placeholders,
	)
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *Repository) SaveStudentCOR(
	ctx context.Context,
	tx datastore.DB,
	cor StudentCOR,
) error {
	cols, vals := datastore.GetInsertStatement(StudentCORDB{}, []string{})
	onDuplicate := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentCORDB{},
		[]string{},
	)

	query := fmt.Sprintf(`
		INSERT INTO student_cors (%s) VALUES (%s) %s`,
		cols,
		vals,
		onDuplicate,
	)

	dbModel := cor.ToPersistence()
	_, err := tx.NamedExecContext(ctx, query, &dbModel)
	if err != nil {
		return fmt.Errorf("failed to save student COR: %w", err)
	}
	return nil
}

func (r *Repository) GetStudentCORByUserID(
	ctx context.Context,
	userID string,
) (StudentCOR, error) {
	query := `
		SELECT file_id, student_id, valid_from, valid_until
		FROM student_cors
		WHERE student_id = ? AND valid_from <= NOW() AND valid_until > NOW()
	`
	var dbModel StudentCORDB
	err := r.db.GetContext(ctx, &dbModel, query, userID)
	if err != nil {
		return StudentCOR{}, fmt.Errorf("failed to get student COR: %w", err)
	}

	return dbModel.ToDomain(), nil
}

func (r *Repository) GetStudentCORsByUserID(
	ctx context.Context,
	userID string,
) ([]StudentCOR, error) {
	query := `
		SELECT file_id, student_id, valid_from, valid_until
		FROM student_cors
		WHERE student_id = ?
	`
	var dbModels []StudentCORDB
	err := r.db.SelectContext(ctx, &dbModels, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student CORs: %w", err)
	}

	cors := make([]StudentCOR, len(dbModels))
	for i, v := range dbModels {
		cors[i] = v.ToDomain()
	}

	return cors, nil
}
