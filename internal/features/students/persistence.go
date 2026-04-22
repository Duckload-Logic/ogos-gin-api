package students

import (
	"database/sql"
	"time"
)

// Lookup models
type GenderDB struct {
	ID         int    `db:"id"`
	GenderName string `db:"gender_name"`
}

type ParentalStatusTypeDB struct {
	ID         int    `db:"id"`
	StatusName string `db:"status_name"`
}

type StudentSupportTypeDB struct {
	ID              int    `db:"id"`
	SupportTypeName string `db:"support_type_name"`
}

type EnrollmentReasonDB struct {
	ID   int    `db:"id"`
	Text string `db:"reason_text"`
}

type IncomeRangeDB struct {
	ID        int    `db:"id"`
	RangeText string `db:"range_text"`
}

type EducationalLevelDB struct {
	ID        int    `db:"id"`
	LevelName string `db:"level_name"`
}

type CourseDB struct {
	ID         int    `db:"id"`
	Code       string `db:"code"`
	CourseName string `db:"course_name"`
}

type CivilStatusTypeDB struct {
	ID         int    `db:"id"`
	StatusName string `db:"status_name"`
}

type ReligionDB struct {
	ID           int    `db:"id"`
	ReligionName string `db:"religion_name"`
}

type StudentRelationshipTypeDB struct {
	ID               int    `db:"id"`
	RelationshipName string `db:"relationship_name"`
}

type NatureOfResidenceTypeDB struct {
	ID                int    `db:"id"`
	ResidenceTypeName string `db:"residence_type_name"`
}

type SibilingSupportTypeDB struct {
	ID          int    `db:"id"`
	SupportName string `db:"name"`
}

type StudentStatusDB struct {
	ID         int    `db:"id"`
	StatusName string `db:"status_name"`
}

type ActivityOptionDB struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`     // "Academic"or "Club"
	Category string `db:"category"` // e.g. "Sports"
	IsActive bool   `db:"is_active"`
}

type StudentBasicInfoViewDB struct {
	UserID     string         `db:"user_id"`
	Email      string         `db:"email"`
	FirstName  string         `db:"first_name"`
	MiddleName sql.NullString `db:"middle_name"`
	LastName   string         `db:"last_name"`
	SuffixName sql.NullString `db:"suffix_name"`
}

type StudentProfileViewDB struct {
	IIRID         string         `db:"iir_id"`
	UserID        string         `db:"user_id"`
	FirstName     string         `db:"first_name"`
	MiddleName    sql.NullString `db:"middle_name"`
	LastName      string         `db:"last_name"`
	SuffixName    sql.NullString `db:"suffix_name"`
	Email         string         `db:"email"`
	StudentNumber string         `db:"student_number"`
	GenderID      int            `db:"gender_id"`
	CourseID      int            `db:"course_id"`
	Section       int            `db:"section"`
	YearLevel     int            `db:"year_level"`
	StatusID      int            `db:"status_id"`
	StatusName    string         `db:"status_name"`
}

// Core Student Records
type IIRDraftDB struct {
	ID        int       `db:"id"`
	UserID    string    `db:"user_id"`
	Data      string    `db:"data"` // JSON string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type IIRRecordDB struct {
	ID          string    `db:"id"`
	UserID      string    `db:"user_id"`
	IsSubmitted bool      `db:"is_submitted"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Enrollment and Reasons
type StudentSelectedReasonDB struct {
	IIRID           string  `db:"iir_id"`
	ReasonID        int     `db:"reason_id"`
	OtherReasonText *string `db:"other_reason_text"`
}

type StudentPersonalInfoDB struct {
	ID    int    `db:"id"`
	IIRID string `db:"iir_id"`

	StudentNumber string `db:"student_number"`

	// Personal Details
	GenderID      int `db:"gender_id"`
	CivilStatusID int `db:"civil_status_id"`
	ReligionID    int `db:"religion_id"`

	// Physical Attributes
	HeightFt   float64 `db:"height_ft"`
	WeightKg   float64 `db:"weight_kg"`
	Complexion string  `db:"complexion"`

	// Academic Information
	HighSchoolGWA float64 `db:"high_school_gwa"`
	CourseID      int     `db:"course_id"`
	YearLevel     int     `db:"year_level"`
	Section       int     `db:"section"`

	// Additional Details
	PlaceOfBirth    string         `db:"place_of_birth"`
	DateOfBirth     string         `db:"date_of_birth"`
	IsEmployed      bool           `db:"is_employed"`
	EmployerName    sql.NullString `db:"employer_name"`
	EmployerAddress sql.NullString `db:"employer_address"`
	MobileNumber    string         `db:"mobile_number"`
	TelephoneNumber sql.NullString `db:"telephone_number"`
	StatusID        int            `db:"status_id"`
	GraduationYear  sql.NullInt64  `db:"graduation_year"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at"`
}

type EmergencyContactDB struct {
	ID             int            `db:"id"`
	IIRID          string         `db:"iir_id"`
	FirstName      string         `db:"first_name"`
	MiddleName     sql.NullString `db:"middle_name"`
	LastName       string         `db:"last_name"`
	SuffixName     sql.NullString `db:"suffix_name"`
	ContactNumber  string         `db:"contact_number"`
	RelationshipID int            `db:"relationship_id"`
	AddressID      int            `db:"address_id"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
}

// Location and Address

type StudentAddressDB struct {
	ID          int       `db:"id"`
	IIRID       string    `db:"iir_id"`
	AddressID   int       `db:"address_id"`
	AddressType string    `db:"address_type"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Family and Related Persons
type RelatedPersonDB struct {
	ID               int            `db:"id"`
	EducationalLevel string         `db:"educational_level"`
	DateOfBirth      string         `db:"date_of_birth"`
	LastName         string         `db:"last_name"`
	FirstName        string         `db:"first_name"`
	MiddleName       sql.NullString `db:"middle_name"`
	SuffixName       sql.NullString `db:"suffix_name"`
	Occupation       sql.NullString `db:"occupation"`
	EmployerName     sql.NullString `db:"employer_name"`
	EmployerAddress  sql.NullString `db:"employer_address"`
	CreatedAt        time.Time      `db:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"`
}

type StudentRelatedPersonDB struct {
	IIRID           string `db:"iir_id"`
	RelatedPersonID int    `db:"related_person_id"`

	// "Father", "Mother", "Guardian", "Uncle", "Aunt", "Sibling", "Other"
	RelationshipID int `db:"relationship_id"`

	// Roles
	IsParent   bool `db:"is_parent"`
	IsGuardian bool `db:"is_guardian"`
	IsLiving   bool `db:"is_living"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type FamilyBackgroundDB struct {
	ID                    int            `db:"id"`
	IIRID                 string         `db:"iir_id"`
	ParentalStatusID      int            `db:"parental_status_id"`
	ParentalStatusDetails sql.NullString `db:"parental_status_details"`
	Brothers              int            `db:"brothers"`
	Sisters               int            `db:"sisters"`
	EmployedSiblings      int            `db:"employed_siblings"`
	OrdinalPosition       int            `db:"ordinal_position"`
	HaveQuietPlaceToStudy bool           `db:"have_quiet_place_to_study"`
	IsSharingRoom         bool           `db:"is_sharing_room"`
	RoomSharingDetails    sql.NullString `db:"room_sharing_details"`
	NatureOfResidenceId   int            `db:"nature_of_residence_id"`
	CreatedAt             time.Time      `db:"created_at"`
	UpdatedAt             time.Time      `db:"updated_at"`
}

type StudentSiblingSupportDB struct {
	FamilyBackgroundID int `db:"family_background_id"`
	SupportTypeID      int `db:"support_type_id"`
}

// Education and Background
type EducationalBackgroundDB struct {
	ID                 int            `db:"id"`
	IIRID              string         `db:"iir_id"`
	NatureOfSchooling  string         `db:"nature_of_schooling"`
	InterruptedDetails sql.NullString `db:"interrupted_details"`
	CreatedAt          time.Time      `db:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at"`
}

type SchoolDetailsDB struct {
	ID                 int            `db:"id"`
	EBID               int            `db:"eb_id"`
	EducationalLevelID int            `db:"educational_level_id"`
	SchoolName         string         `db:"school_name"`
	SchoolAddress      string         `db:"school_address"`
	SchoolType         string         `db:"school_type"`
	YearStarted        int            `db:"year_started"`
	YearCompleted      int            `db:"year_completed"`
	Awards             sql.NullString `db:"awards"`
	CreatedAt          time.Time      `db:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at"`
}

// Health and Wellness
type StudentHealthRecordDB struct {
	ID                      int            `db:"id"`
	IIRID                   string         `db:"iir_id"`
	VisionHasProblem        bool           `db:"vision_has_problem"`
	VisionDetails           sql.NullString `db:"vision_details"`
	HearingHasProblem       bool           `db:"hearing_has_problem"`
	HearingDetails          sql.NullString `db:"hearing_details"`
	SpeechHasProblem        bool           `db:"speech_has_problem"`
	SpeechDetails           sql.NullString `db:"speech_details"`
	GeneralHealthHasProblem bool           `db:"general_health_has_problem"`
	GeneralHealthDetails    sql.NullString `db:"general_health_details"`
	CreatedAt               time.Time      `db:"created_at"`
	UpdatedAt               time.Time      `db:"updated_at"`
}

type StudentConsultationDB struct {
	ID               int            `db:"id"`
	IIRID            string         `db:"iir_id"`
	ProfessionalType string         `db:"professional_type"` // Counselor...
	HasConsulted     bool           `db:"has_consulted"`
	WhenDate         sql.NullString `db:"when_date"`
	ForWhat          sql.NullString `db:"for_what"`
	CreatedAt        time.Time      `db:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"`
}

// Financial Support
type StudentFinanceDB struct {
	ID              int            `db:"id"`
	IIRID           string         `db:"iir_id"`
	IncomeRangeID   int            `db:"monthly_family_income_range_id"`
	OtherIncome     sql.NullString `db:"other_income_details"`
	WeeklyAllowance float64        `db:"weekly_allowance"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at"`
}

type StudentFinancialSupportDB struct {
	StudentFinanceID int       `db:"sf_id"`
	SupportTypeID    int       `db:"support_type_id"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

// Interests and Activities
type StudentActivityDB struct {
	ID                 int            `db:"id"`
	IIRID              string         `db:"iir_id"`
	OptionID           int            `db:"option_id"`
	OtherSpecification sql.NullString `db:"other_specification"`
	Role               string         `db:"role"`
	RoleSpecification  sql.NullString `db:"role_specification"`
	CreatedAt          time.Time      `db:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at"`
}

type StudentSubjectPreferenceDB struct {
	ID          int       `db:"id"`
	IIRID       string    `db:"iir_id"`
	SubjectName string    `db:"subject_name"`
	IsFavorite  bool      `db:"is_favorite"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type StudentHobbyDB struct {
	ID           int       `db:"id"`
	IIRID        string    `db:"iir_id"`
	HobbyName    string    `db:"hobby_name"`
	PriorityRank int       `db:"priority_rank"` // 1 = favorite
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// Test Results and Assessments
type TestResultDB struct {
	ID          int       `db:"id"`
	IIRID       string    `db:"iir_id"`
	TestDate    string    `db:"test_date"`
	TestName    string    `db:"test_name"`
	RawScore    string    `db:"raw_score"`
	Percentile  string    `db:"percentile"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type StudentCORDB struct {
	FileID     string       `db:"file_id"`
	StudentID  string       `db:"student_id"`
	ValidFrom  sql.NullTime `db:"valid_from"`
	ValidUntil sql.NullTime `db:"valid_until"`
}
