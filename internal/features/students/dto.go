package students

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/request"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
)

// List Students
type ListStudentsRequest struct {
	request.PaginationParams
	Search    string `form:"search,omitempty"`
	CourseID  int    `form:"course_id,omitempty"`
	GenderID  int    `form:"gender_id,omitempty"`
	YearLevel int    `form:"year_level,omitempty"`
	OrderBy   string `form:"order_by,omitempty" binding:"omitempty,oneof=first_name last_name student_number iir_id created_at updated_at year_level course_id"`
}

type ListStudentsResponse struct {
	Students   []StudentProfileDTO `json:"students"`
	Total      int                 `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"pageSize"`
	TotalPages int                 `json:"totalPages"`
}

type StudentProfileDTO struct {
	IIRID         int                    `json:"iirId"`
	UserID        int                    `json:"userId"`
	FirstName     string                 `json:"firstName"`
	MiddleName    structs.NullableString `json:"middleName,omitempty"`
	LastName      string                 `json:"lastName"`
	Gender        Gender                 `json:"gender"`
	Email         string                 `json:"email"`
	StudentNumber string                 `json:"studentNumber"`
	Course        Course                 `json:"course"`
	Section       int                    `json:"section"`
	YearLevel     int                    `json:"yearLevel"`
}

// Get Student
type GetStudentRequest struct {
	IncludeParams
}

type ComprehensiveProfileDTO struct {
	IIRID   int `json:"iirId,omitempty"`
	Student struct {
		BasicInfo              StudentBasicInfoViewDTO `json:"basicInfo"`
		StudentPersonalInfoDTO `json:"personalInfo"`
		Addresses              []StudentAddressDTO `json:"addresses"`
	} `json:"student"`

	Education EducationalBackgroundDTO `json:"education"`

	Family struct {
		FamilyBackgroundDTO `json:"background"`
		RelatedPersons      []RelatedPersonDTO `json:"relatedPersons"`
		Finance             StudentFinanceDTO  `json:"finance"`
	} `json:"family"`

	Health struct {
		StudentHealthRecordDTO `json:"healthRecord"`
		Consultations          []StudentConsultationDTO `json:"consultations"`
	} `json:"health"`

	Interests struct {
		Activities         []StudentActivityDTO          `json:"activities"`
		SubjectPreferences []StudentSubjectPreferenceDTO `json:"subjectPreferences"`
		Hobbies            []StudentHobbyDTO             `json:"hobbies"`
	} `json:"interests"`

	TestResults      []TestResultDTO      `json:"testResults"`
	SignificantNotes []SignificantNoteDTO `json:"significantNotes,omitempty"`
}

type StudentSelectedReasonDTO struct {
	Reason          EnrollmentReason `json:"reason"`
	OtherReasonText *string          `json:"otherReasonText,omitempty"`
}

type StudentBasicInfoViewDTO struct {
	ID         int                    `json:"id"`
	FirstName  string                 `json:"firstName"`
	MiddleName structs.NullableString `json:"middleName,omitempty"`
	LastName   string                 `json:"lastName"`
	Email      string                 `json:"email"`
}

type StudentPersonalInfoDTO struct {
	ID               int                    `json:"id,omitempty"`
	IIRID            int                    `json:"iirId,omitempty"`
	StudentNumber    string                 `json:"studentNumber" binding:"required"`
	Gender           Gender                 `json:"gender" binding:"required"`
	CivilStatus      CivilStatusType        `json:"civilStatus" binding:"required"`
	Religion         Religion               `json:"religion" binding:"required"`
	HeightFt         float64                `json:"heightFt" binding:"required"`
	WeightKg         float64                `json:"weightKg" binding:"required"`
	Complexion       string                 `json:"complexion" binding:"required"`
	HighSchoolGWA    float64                `json:"highSchoolGWA" binding:"required"`
	Course           Course                 `json:"course" binding:"required"`
	YearLevel        int                    `json:"yearLevel" binding:"required"`
	Section          int                    `json:"section" binding:"required"`
	PlaceOfBirth     string                 `json:"placeOfBirth" binding:"required"`
	DateOfBirth      string                 `json:"dateOfBirth" binding:"required"`
	IsEmployed       bool                   `json:"isEmployed"`
	EmployerName     structs.NullableString `json:"employerName,omitempty"`
	EmployerAddress  structs.NullableString `json:"employerAddress,omitempty"`
	MobileNumber     string                 `json:"mobileNumber" binding:"required"`
	TelephoneNumber  structs.NullableString `json:"telephoneNumber,omitempty"`
	EmergencyContact EmergencyContactDTO    `json:"emergencyContact,omitempty"`
}

type EmergencyContactDTO struct {
	ID            int                     `json:"id,omitempty"`
	FirstName     string                  `json:"firstName" binding:"required"`
	MiddleName    structs.NullableString  `json:"middleName,omitempty"`
	LastName      string                  `json:"lastName" binding:"required"`
	ContactNumber string                  `json:"contactNumber" binding:"required"`
	Relationship  StudentRelationshipType `json:"relationship" binding:"required"`
	Address       locations.AddressDTO    `json:"address" binding:"required"`
}

type StudentAddressDTO struct {
	ID          int                  `json:"id,omitempty"`
	AddressType string               `json:"addressType" binding:"required"`
	Address     locations.AddressDTO `json:"address" binding:"required"`
	CreatedAt   time.Time            `json:"createdAt,omitempty"`
	UpdatedAt   time.Time            `json:"updatedAt,omitempty"`
}

type EducationalBackgroundDTO struct {
	ID                 int                    `json:"id,omitempty"`
	NatureOfSchooling  string                 `json:"natureOfSchooling" binding:"required"`
	InterruptedDetails structs.NullableString `json:"interruptedDetails,omitempty"`
	School             []SchoolDetailsDTO     `json:"schools" binding:"required"`
	CreatedAt          time.Time              `json:"createdAt,omitempty"`
	UpdatedAt          time.Time              `json:"updatedAt,omitempty"`
}

type SchoolDetailsDTO struct {
	ID               int                    `json:"id,omitempty"`
	EducationalLevel EducationalLevel       `json:"educationalLevel" binding:"required"`
	SchoolName       string                 `json:"schoolName" binding:"required"`
	SchoolAddress    string                 `json:"schoolAddress,omitempty"`
	SchoolType       string                 `json:"schoolType" binding:"required"`
	YearStarted      int                    `json:"yearStarted,omitempty"`
	YearCompleted    int                    `json:"yearCompleted" binding:"required"`
	Awards           structs.NullableString `json:"awards,omitempty"`
}

type RelatedPersonDTO struct {
	ID               int                     `json:"id,omitempty"`
	LastName         string                  `json:"lastName" binding:"required"`
	FirstName        string                  `json:"firstName" binding:"required"`
	MiddleName       structs.NullableString  `json:"middleName,omitempty"`
	DateOfBirth      string                  `json:"dateOfBirth,omitempty" binding:"omitempty"`
	EducationalLevel string                  `json:"educationalLevel" binding:"required"`
	Occupation       structs.NullableString  `json:"occupation,omitempty"`
	EmployerName     structs.NullableString  `json:"employerName,omitempty"`
	EmployerAddress  structs.NullableString  `json:"employerAddress,omitempty"`
	Relationship     StudentRelationshipType `json:"relationship" binding:"required"`
	IsParent         bool                    `json:"isParent"`
	IsGuardian       bool                    `json:"isGuardian"`
	IsLiving         bool                    `json:"isLiving"`
}

type FamilyBackgroundDTO struct {
	ID                    int                    `json:"id,omitempty"`
	ParentalStatus        ParentalStatusType     `json:"parentalStatus" binding:"required"`
	ParentalStatusDetails structs.NullableString `json:"parentalStatusDetails,omitempty"`
	Brothers              int                    `json:"brothers" binding:"required"`
	Sisters               int                    `json:"sisters" binding:"required"`
	EmployedSiblings      int                    `json:"employedSiblings" binding:"required"`
	OrdinalPosition       int                    `json:"ordinalPosition" binding:"required"`
	HaveQuietPlaceToStudy bool                   `json:"haveQuietPlaceToStudy"`
	SiblingSupportTypes   []SibilingSupportType  `json:"siblingSupportTypes" binding:"required"`
	IsSharingRoom         bool                   `json:"isSharingRoom"`
	RoomSharingDetails    structs.NullableString `json:"roomSharingDetails,omitempty"`
	NatureOfResidence     NatureOfResidenceType  `json:"natureOfResidence" binding:"required"`
}

type EducationalBGDTO struct {
	ID               int    `json:"id"`
	EducationalLevel string `json:"educationalLevel" binding:"required"`
	SchoolName       string `json:"schoolName" binding:"required"`
	Location         string `json:"location,omitempty"`
	SchoolType       string `json:"schoolType" binding:"required,oneof=Public Private"`
	YearCompleted    string `json:"yearCompleted" binding:"required"`
	Awards           string `json:"awards,omitempty"`
}

type StudentFinanceDTO struct {
	ID                       int                    `json:"id,omitempty"`
	MonthlyFamilyIncomeRange IncomeRange            `json:"monthlyFamilyIncomeRange" binding:"required"`
	OtherIncomeDetails       structs.NullableString `json:"otherIncomeDetails,omitempty"`
	FinancialSupportTypes    []StudentSupportType   `json:"financialSupportTypes" binding:"required"`
	WeeklyAllowance          float64                `json:"weeklyAllowance" binding:"required"`
}

type StudentHealthRecordDTO struct {
	ID                      int                    `json:"id,omitempty"`
	VisionHasProblem        bool                   `json:"visionHasProblem"`
	VisionDetails           structs.NullableString `json:"visionDetails,omitempty"`
	HearingHasProblem       bool                   `json:"hearingHasProblem"`
	HearingDetails          structs.NullableString `json:"hearingDetails,omitempty"`
	SpeechHasProblem        bool                   `json:"speechHasProblem"`
	SpeechDetails           structs.NullableString `json:"speechDetails,omitempty"`
	GeneralHealthHasProblem bool                   `json:"generalHealthHasProblem"`
	GeneralHealthDetails    structs.NullableString `json:"generalHealthDetails,omitempty"`
}

type StudentConsultationDTO struct {
	ID               int                    `json:"id,omitempty"`
	ProfessionalType string                 `json:"professionalType" binding:"required"`
	HasConsulted     bool                   `json:"hasConsulted"`
	WhenDate         structs.NullableString `json:"whenDate,omitempty"`
	ForWhat          structs.NullableString `json:"forWhat,omitempty"`
}

type StudentActivityDTO struct {
	ID                 int                    `json:"id,omitempty"`
	ActivityOption     ActivityOption         `json:"activityOption" binding:"required"`
	OtherSpecification structs.NullableString `json:"otherSpecification,omitempty"`
	Role               string                 `json:"role" binding:"required"` // "Officer", "Member", "Other"
	RoleSpecification  structs.NullableString `json:"roleSpecification,omitempty"`
}

type StudentSubjectPreferenceDTO struct {
	ID          int    `json:"id,omitempty"`
	SubjectName string `json:"subjectName" binding:"required"`
	IsFavorite  bool   `json:"isFavorite"`
}

type StudentHobbyDTO struct {
	ID           int    `json:"id,omitempty"`
	HobbyName    string `json:"hobbyName" binding:"required"`
	PriorityRank int    `json:"priorityRank" binding:"required"`
}

type TestResultDTO struct {
	ID          int    `json:"id,omitempty"`
	TestDate    string `json:"testDate" binding:"required"`
	TestName    string `json:"testName" binding:"required"`
	RawScore    string `json:"rawScore" binding:"required"`
	Percentile  string `json:"percentile" binding:"required"`
	Description string `json:"description,omitempty"`
}

type SignificantNoteDTO struct {
	ID                  int       `json:"id,omitempty"`
	NoteDate            string    `json:"noteDate" binding:"required"`
	IncidentDescription string    `json:"incidentDescription" binding:"required"`
	Remarks             string    `json:"remarks" binding:"required"`
	CreatedAt           time.Time `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt           time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}
