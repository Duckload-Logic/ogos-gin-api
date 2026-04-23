package students

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// Gender represents the domain entity for student gender.
type Gender struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ParentalStatusType represents the domain entity for parental status types.
type ParentalStatusType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// StudentSupportType represents the domain entity for student support types.
type StudentSupportType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// EnrollmentReason represents the domain entity for enrollment reasons.
type EnrollmentReason struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

// IncomeRange represents the domain entity for family income ranges.
type IncomeRange struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

// EducationalLevel represents the domain entity for educational levels.
type EducationalLevel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Course represents the domain entity for academic courses.
type Course struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// CivilStatusType represents the domain entity for civil status types.
type CivilStatusType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Religion represents the domain entity for religious affiliations.
type Religion struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// StudentRelationshipType represents the domain entity for
// relationships (e.g., Father, Mother).
type StudentRelationshipType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// NatureOfResidenceType represents the domain entity for residence types.
type NatureOfResidenceType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// SibilingSupportType represents the domain entity for sibling support.
type SibilingSupportType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// StudentStatus represents the domain entity for academic status.
type StudentStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ActivityOption represents the domain entity for student activities.
type ActivityOption struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	IsActive bool   `json:"isActive"`
}

// StudentBasicInfoView represents a simplified view of student
// identification info.
type StudentBasicInfoView struct {
	UserID     string
	Email      string
	FirstName  string
	MiddleName structs.NullableString
	LastName   string
	SuffixName structs.NullableString
}

// StudentProfileView represents a summary view of a student's profile.
type StudentProfileView struct {
	IIRID         string
	UserID        string
	FirstName     string
	MiddleName    structs.NullableString
	LastName      string
	SuffixName    structs.NullableString
	Email         string
	StudentNumber string
	GenderID      int
	CourseID      int
	Section       int
	YearLevel     int
	StatusID      int
	StatusName    string
}

// IIRDraft represents a student's in-progress Initial Interview Record.
type IIRDraft struct {
	ID        int
	UserID    string
	Data      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IIRRecord represents a submitted Initial Interview Record entry.
type IIRRecord struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	IsSubmitted bool      `json:"isSubmitted"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// StudentSelectedReason represents the domain mapping for enrollment reasons.
type StudentSelectedReason struct {
	IIRID           string
	ReasonID        int
	OtherReasonText *string
}

// StudentPersonalInfo represents detailed student personal parameters.
type StudentPersonalInfo struct {
	ID              int
	IIRID           string
	StudentNumber   string
	GenderID        int
	CivilStatusID   int
	ReligionID      int
	HeightFt        float64
	WeightKg        float64
	Complexion      string
	HighSchoolGWA   float64
	CourseID        int
	YearLevel       int
	Section         int
	PlaceOfBirth    string
	DateOfBirth     string
	IsEmployed      bool
	EmployerName    structs.NullableString
	EmployerAddress structs.NullableString
	MobileNumber    string
	TelephoneNumber structs.NullableString
	StatusID        int
	GraduationYear  structs.NullableInt64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// EmergencyContact represents the domain entity for emergency notifications.
type EmergencyContact struct {
	ID             int
	IIRID          string
	FirstName      string
	MiddleName     structs.NullableString
	LastName       string
	SuffixName     structs.NullableString
	ContactNumber  string
	RelationshipID int
	AddressID      int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// StudentAddress represents the domain mapping for student locations.
type StudentAddress struct {
	ID          int
	IIRID       string
	AddressID   int
	AddressType string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// RelatedPerson represents a person related to the student
// (parent, guardian, etc.).
type RelatedPerson struct {
	ID               int
	EducationalLevel string
	DateOfBirth      string
	LastName         string
	FirstName        string
	MiddleName       structs.NullableString
	SuffixName       structs.NullableString
	Occupation       structs.NullableString
	EmployerName     structs.NullableString
	EmployerAddress  structs.NullableString
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// StudentRelatedPerson represents the mapping between a student
// and a related person.
type StudentRelatedPerson struct {
	IIRID           string
	RelatedPersonID int
	RelationshipID  int
	IsParent        bool
	IsGuardian      bool
	IsLiving        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// FamilyBackground represents details about the student's family structure.
type FamilyBackground struct {
	ID                    int
	IIRID                 string
	ParentalStatusID      int
	ParentalStatusDetails structs.NullableString
	Brothers              int
	Sisters               int
	EmployedSiblings      int
	OrdinalPosition       int
	HaveQuietPlaceToStudy bool
	IsSharingRoom         bool
	RoomSharingDetails    structs.NullableString
	NatureOfResidenceId   int
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// StudentSiblingSupport maps family backgrounds to support types.
type StudentSiblingSupport struct {
	FamilyBackgroundID int
	SupportTypeID      int
}

// EducationalBackground represents the student's prior schooling status.
type EducationalBackground struct {
	ID                 int
	IIRID              string
	NatureOfSchooling  string
	InterruptedDetails structs.NullableString
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// SchoolDetails represents specific schooling history entries.
type SchoolDetails struct {
	ID                 int
	EBID               int
	EducationalLevelID int
	SchoolName         string
	SchoolAddress      string
	SchoolType         string
	YearStarted        int
	YearCompleted      int
	Awards             structs.NullableString
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// StudentHealthRecord represents physical wellness parameters.
type StudentHealthRecord struct {
	ID                      int
	IIRID                   string
	VisionHasProblem        bool
	VisionDetails           structs.NullableString
	HearingHasProblem       bool
	HearingDetails          structs.NullableString
	SpeechHasProblem        bool
	SpeechDetails           structs.NullableString
	GeneralHealthHasProblem bool
	GeneralHealthDetails    structs.NullableString
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

// StudentConsultation represents record of professional consultations.
type StudentConsultation struct {
	ID               int
	IIRID            string
	ProfessionalType string
	HasConsulted     bool
	WhenDate         structs.NullableString
	ForWhat          structs.NullableString
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// StudentFinance represents student financial background data.
type StudentFinance struct {
	ID                         int
	IIRID                      string
	MonthlyFamilyIncomeRangeID int
	OtherIncomeDetails         structs.NullableString
	WeeklyAllowance            float64
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}

// StudentFinancialSupport maps financial support types to student finance.
type StudentFinancialSupport struct {
	StudentFinanceID int
	SupportTypeID    int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// StudentActivity represents student interests and club involvements.
type StudentActivity struct {
	ID                 int
	IIRID              string
	OptionID           int
	OtherSpecification structs.NullableString
	Role               string
	RoleSpecification  structs.NullableString
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// StudentSubjectPreference represents academic subject interests.
type StudentSubjectPreference struct {
	ID          int
	IIRID       string
	SubjectName string
	IsFavorite  bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// StudentHobby represents extracurricular interests.
type StudentHobby struct {
	ID           int
	IIRID        string
	HobbyName    string
	PriorityRank int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TestResult represents academic assessment outcome records.
type TestResult struct {
	ID          int
	IIRID       string
	TestDate    string
	TestName    string
	RawScore    string
	Percentile  string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// StudentCOR represents the Certificate of Registration domain entity.
type StudentCOR struct {
	FileID     string
	StudentID  string
	ValidFrom  structs.NullableTime
	ValidUntil structs.NullableTime
}
