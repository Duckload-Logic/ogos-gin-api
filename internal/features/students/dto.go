package students

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/request"
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
	IIRID         int     `json:"iirId"`
	UserID        int     `json:"userId"`
	FirstName     string  `json:"firstName"`
	MiddleName    *string `json:"middleName,omitempty"`
	LastName      string  `json:"lastName"`
	Gender        Gender  `json:"gender"`
	Email         string  `json:"email"`
	StudentNumber string  `json:"studentNumber"`
	Course        Course  `json:"course"`
	Section       int     `json:"section"`
	YearLevel     int     `json:"yearLevel"`
}

// Get Student
type GetStudentRequest struct {
	IncludeParams
}

type ComprehensiveProfileResponse struct {
	IIRID   int `json:"iirId"`
	Student struct {
		BasicInfo              StudentBasicInfoView `json:"basicInfo"`
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
	SignificantNotes []SignificantNoteDTO `json:"significantNotes"`
}

type StudentSelectedReasonDTO struct {
	Reason          EnrollmentReason `json:"reason"`
	OtherReasonText *string          `json:"otherReasonText,omitempty"`
}

type StudentPersonalInfoDTO struct {
	ID                           int                     `json:"id"`
	StudentNumber                string                  `json:"studentNumber" binding:"required"`
	Gender                       Gender                  `json:"gender" binding:"required"`
	CivilStatus                  CivilStatusType         `json:"civilStatus" binding:"required"`
	Religion                     Religion                `json:"religion" binding:"required"`
	HeightFt                     float64                 `json:"heightFt" binding:"required"`
	WeightKg                     float64                 `json:"weightKg" binding:"required"`
	Complexion                   string                  `json:"complexion" binding:"required"`
	HighSchoolGWA                float64                 `json:"highSchoolGWA" binding:"required"`
	Course                       Course                  `json:"course" binding:"required"`
	YearLevel                    int                     `json:"yearLevel" binding:"required"`
	Section                      int                     `json:"section" binding:"required"`
	PlaceOfBirth                 string                  `json:"placeOfBirth" binding:"required"`
	DateOfBirth                  string                  `json:"dateOfBirth" binding:"required"`
	MobileNumber                 string                  `json:"mobileNumber" binding:"required"`
	TelephoneNumber              *string                 `json:"telephoneNumber,omitempty"`
	EmergencyContactName         string                  `json:"emergencyContactName" binding:"required"`
	EmergencyContactNumber       string                  `json:"emergencyContactNumber" binding:"required"`
	EmergencyContactRelationship StudentRelationshipType `json:"emergencyContactRelationship" binding:"required"`
	EmergencyContactAddress      Address                 `json:"emergencyContactAddress" binding:"required"`
}

type StudentAddressDTO struct {
	ID          int       `json:"id"`
	AddressType string    `json:"addressType"`
	Address     Address   `json:"address"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type EducationalBackgroundDTO struct {
	ID                 int                `json:"id"`
	NatureOfSchooling  string             `json:"natureOfSchooling"`
	InterruptedDetails *string            `json:"interruptedDetails,omitempty"`
	School             []SchoolDetailsDTO `json:"schools"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
}

type SchoolDetailsDTO struct {
	ID               int              `json:"id"`
	EducationalLevel EducationalLevel `json:"educationalLevel"`
	SchoolName       string           `json:"schoolName"`
	SchoolAddress    string           `json:"schoolAddress,omitempty"`
	SchoolType       string           `json:"schoolType"`
	YearStarted      int              `json:"yearStarted,omitempty"`
	YearCompleted    int              `json:"yearCompleted"`
	Awards           *string          `json:"awards,omitempty"`
}

type RelatedPersonDTO struct {
	ID               int     `json:"id"`
	LastName         string  `json:"lastName"`
	FirstName        string  `json:"firstName"`
	MiddleName       *string `json:"middleName,omitempty"`
	DateOfBirth      *string `json:"dateOfBirth,omitempty"`
	EducationalLevel string  `json:"educationalLevel"`
	Occupation       *string `json:"occupation,omitempty"`
	EmployerName     *string `json:"employerName,omitempty"`
	EmployerAddress  *string `json:"employerAddress,omitempty"`
	ContactNumber    *string `json:"contactNumber,omitempty"`

	Relationship       StudentRelationshipType `json:"relationship"`
	IsParent           bool                    `json:"isParent"`
	IsGuardian         bool                    `json:"isGuardian"`
	IsEmergencyContact bool                    `json:"isEmergencyContact"`
	IsLiving           bool                    `json:"isLiving"`

	Address *Address `json:"address,omitempty"`
}

type FamilyBackgroundDTO struct {
	ID                    int                   `json:"id"`
	ParentalStatus        ParentalStatusType    `json:"parentalStatus"`
	ParentalStatusDetails *string               `json:"parentalStatusDetails,omitempty"`
	Brothers              int                   `json:"brothers"`
	Sisters               int                   `json:"sisters"`
	EmployedSiblings      int                   `json:"employedSiblings"`
	OrdinalPosition       int                   `json:"ordinalPosition"`
	HaveQuietPlaceToStudy bool                  `json:"haveQuietPlaceToStudy"`
	SiblingSupportTypes   []SibilingSupportType `json:"siblingSupportTypes"`
	IsSharingRoom         bool                  `json:"isSharingRoom"`
	RoomSharingDetails    *string               `json:"roomSharingDetails,omitempty"`
	NatureOfResidence     NatureOfResidenceType `json:"natureOfResidence"`
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
	ID                       int                  `json:"id"`
	MonthlyFamilyIncomeRange IncomeRange          `json:"monthlyFamilyIncomeRange" binding:"required"`
	OtherIncomeDetails       *string              `json:"otherIncomeDetails,omitempty"`
	FinancialSupportTypes    []StudentSupportType `json:"financialSupportTypes" binding:"required"`
	WeeklyAllowance          *float64             `json:"weeklyAllowance" binding:"required"`
}

type StudentHealthRecordDTO struct {
	ID                      int     `json:"id"`
	VisionHasProblem        bool    `json:"visionHasProblem"`
	VisionDetails           *string `json:"visionDetails,omitempty"`
	HearingHasProblem       bool    `json:"hearingHasProblem"`
	HearingDetails          *string `json:"hearingDetails,omitempty"`
	SpeechHasProblem        bool    `json:"speechHasProblem"`
	SpeechDetails           *string `json:"speechDetails,omitempty"`
	GeneralHealthHasProblem bool    `json:"generalHealthHasProblem"`
	GeneralHealthDetails    *string `json:"generalHealthDetails,omitempty"`
}

type StudentConsultationDTO struct {
	ID               int     `json:"id"`
	ProfessionalType string  `json:"professionalType"`
	HasConsulted     bool    `json:"hasConsulted"`
	WhenDate         *string `json:"whenDate,omitempty"`
	ForWhat          *string `json:"forWhat,omitempty"`
}

type StudentActivityDTO struct {
	ID                 int            `json:"id"`
	ActivityOption     ActivityOption `json:"activityOption"`
	OtherSpecification *string        `json:"otherSpecification,omitempty"`
	Role               string         `json:"role"` // "Officer", "Member", "Other"
	RoleSpecification  *string        `json:"roleSpecification,omitempty"`
}

type StudentSubjectPreferenceDTO struct {
	ID          int    `json:"id"`
	SubjectName string `json:"subjectName"`
	IsFavorite  bool   `json:"isFavorite"`
}

type StudentHobbyDTO struct {
	ID           int    `json:"id"`
	HobbyName    string `json:"hobbyName"`
	PriorityRank int    `json:"priorityRank"`
}

type TestResultDTO struct {
	ID          int    `json:"id"`
	TestDate    string `json:"testDate"`
	TestName    string `json:"testName"`
	RawScore    string `json:"rawScore"`
	Percentile  string `json:"percentile"`
	Description string `json:"description,omitempty"`
}

type SignificantNoteDTO struct {
	ID                  int       `json:"id"`
	NoteDate            string    `json:"noteDate"`
	IncidentDescription string    `json:"incidentDescription"`
	Remarks             string    `json:"remarks"`
	CreatedAt           time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt           time.Time `db:"updated_at" json:"updatedAt"`
}
