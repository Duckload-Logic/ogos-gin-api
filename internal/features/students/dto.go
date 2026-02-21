package students

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/request"
)

// List Students
type ListStudentsRequest struct {
	request.PaginationParams
	Course   string `form:"course,omitempty"`
	GenderID int    `form:"gender_id,omitempty"`
}

type ListStudentsResponse struct {
	Students   []StudentProfileView `json:"students"`
	Total      int                  `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"pageSize"`
	TotalPages int                  `json:"totalPages"`
}

type StudentProfileView struct {
	IIRID         int     `db:"iir_id" json:"iirId"`
	UserID        int     `db:"user_id" json:"userId"`
	FirstName     string  `db:"first_name" json:"firstName"`
	MiddleName    *string `db:"middle_name" json:"middleName,omitempty"`
	LastName      string  `db:"last_name" json:"lastName"`
	Email         string  `db:"email" json:"email"`
	StudentNumber string  `db:"student_number" json:"studentNumber"`
	Course        string  `db:"course" json:"course"`
	Section       string  `db:"section" json:"section"`
	YearLevel     string  `db:"year_level" json:"yearLevel"`
}

// Get Student
type GetStudentRequest struct {
	IncludeParams
}

type StudentProfileDTO = StudentProfile
type FamilyBackgroundDTO = FamilyBackground
type StudentHealthRecordDTO = StudentHealthRecord
type EducationalBackgroundDTO = EducationalBackground
type StudentAddressResponseDTO = StudentAddress
type StudentFinanceDTO = StudentFinance

type ComprehensiveProfileResponse struct {
	QuickProfile      StudentProfileDTO           `json:"quick,omitempty"`
	EnrollmentReasons []StudentSelectedReason     `json:"enrollment,omitempty"`
	StudentProfile    *StudentProfileDTO          `json:"profile,omitempty"`
	Family            *FamilyBackgroundDTO        `json:"family,omitempty"`
	RelatedPersons    []RelatedPersonInfoView     `json:"relatedPersons,omitempty"`
	Health            *StudentHealthRecordDTO     `json:"health,omitempty"`
	Education         []EducationalBackgroundDTO  `json:"education,omitempty"`
	Addresses         []StudentAddressResponseDTO `json:"addresses,omitempty"`
	Finance           *StudentFinanceDTO          `json:"finance,omitempty"`
}

type StudentProfileResponse struct {
	StudentProfile *StudentProfileDTO `json:"studentProfile"`
}

type StudentProfilePayload struct {
	GenderID          int     `json:"genderId" binding:"required"`
	CivilStatusTypeID int     `json:"civilStatusTypeId" binding:"required"`
	Religion          string  `json:"religion" binding:"required"`
	HeightFt          float64 `json:"heightFt" binding:"required"`
	WeightKg          float64 `json:"weightKg" binding:"required"`
	StudentNumber     string  `json:"studentNumber,omitempty"`
	Course            string  `json:"course" binding:"required"`
	HighSchoolGWA     float64 `json:"highSchoolGWA" binding:"required"`
	PlaceOfBirth      string  `json:"placeOfBirth" binding:"required"`
	DateOfBirth       string  `json:"dateOfBirth" binding:"required"`
	ContactNumber     string  `json:"contactNumber" binding:"required"`
}

// Create Student
type CreateInventoryRecordRequest struct {
	StudentProfilePayload
	EmergencyContact *UpdateEmergencyContactRequest `json:"emergencyContact" binding:"required"`
	Addresses        []StudentAddressDTO            `json:"addresses" binding:"required,min=1"`
}

// Update Emergency Contact
type UpdateEmergencyContactRequest struct {
	RelatedPersonID              *int    `json:"relatedPersonId,omitempty"`
	EmergencyContactFirstName    string  `json:"emergencyContactFirstName" binding:"required"`
	EmergencyContactMiddleName   *string `json:"emergencyContactMiddleName,omitempty"`
	EmergencyContactLastName     string  `json:"emergencyContactLastName" binding:"required"`
	EmergencyContactPhone        string  `json:"emergencyContactPhone" binding:"required"`
	EmergencyContactRelationship string  `json:"emergencyContactRelationship" binding:"required"`
}

// Update Enrollment Reasons
type UpdateEnrollmentReasonsRequest struct {
	EnrollmentReasonIDs []int  `json:"enrollmentReasonIds"`
	OtherReasonText     string `json:"otherReasonText,omitempty"`
}

type FamilyBackgroundPayload struct {
	ParentalStatusID      int     `json:"parentalStatusId" binding:"required"`
	ParentalStatusDetails string  `json:"parentalStatusDetails,omitempty"`
	Brothers              *int    `json:"brothers" binding:"required"`
	Sisters               *int    `json:"sisters" binding:"required"`
	MonthlyFamilyIncome   string  `json:"monthlyFamilyIncome" binding:"required"`
	GuardianFirstName     string  `json:"guardianFirstName" binding:"required"`
	GuardianLastName      string  `json:"guardianLastName" binding:"required"`
	GuardianMiddleName    *string `json:"guardianMiddleName,omitempty"`
	GuardianAddress       string  `json:"guardianAddress" binding:"required"`
}

// Update Family and Related Persons
type UpdateFamilyRequest struct {
	FamilyBackgroundPayload
	UpdateFinanceRequest
	RelatedPersons []RelatedPersonDTO `json:"relatedPersons" binding:"required,dive"`
}

type RelatedPersonPayload struct {
	FirstName        string `json:"firstName" binding:"required"`
	LastName         string `json:"lastName" binding:"required"`
	MiddleName       string `json:"middleName,omitempty"`
	EducationalLevel string `json:"educationalLevel" binding:"required"`
	DateOfBirth      string `json:"dateOfBirth" binding:"required"`
	Occupation       string `json:"occupation" binding:"required"`
	CompanyName      string `json:"companyName,omitempty"`
}

type RelatedPersonDTO struct {
	RelatedPersonPayload
	Relationship int `json:"relationship" binding:"required"`
}

type RelatedPersonInfoView struct {
	RelatedPerson
	Relationship string `db:"relationship" json:"relationship"`
}

// Update Education
type UpdateEducationRequest struct {
	EducationalBGs []EducationalBGDTO `json:"educationalBackgrounds" binding:"required,dive"`
}

type EducationalBGDTO struct {
	EducationalLevel string `json:"educationalLevel" binding:"required"`
	SchoolName       string `json:"schoolName" binding:"required"`
	Location         string `json:"location,omitempty"`
	SchoolType       string `json:"schoolType" binding:"required,oneof=Public Private"`
	YearCompleted    string `json:"yearCompleted" binding:"required"`
	Awards           string `json:"awards,omitempty"`
}

// Update Address
type UpdateAddressRequest struct {
	Addresses []StudentAddressDTO `json:"addresses" binding:"required,dive"`
}

type StudentAddressDTO struct {
	AddressType  string `json:"addressType" binding:"required"`
	RegionName   string `json:"regionName" binding:"required"`
	ProvinceName string `json:"provinceName,omitempty"`
	CityName     string `json:"cityName" binding:"required"`
	BarangayName string `json:"barangayName" binding:"required"`
	StreetLotBlk string `json:"streetLotBlk,omitempty"`
	UnitNo       string `json:"unitNo,omitempty"`
	BuildingName string `json:"buildingName,omitempty"`
}

// Update Health
type UpdateHealthRecordRequest struct {
	VisionRemark          string  `json:"visionRemark" binding:"required"`
	HearingRemark         string  `json:"hearingRemark" binding:"required"`
	MobilityRemark        string  `json:"mobilityRemark" binding:"required"`
	SpeechRemark          string  `json:"speechRemark" binding:"required"`
	GeneralHealthRemark   string  `json:"generalHealthRemark" binding:"required"`
	ConsultedProfessional *string `json:"consultedProfessional,omitempty"`
	ConsultationReason    *string `json:"consultationReason,omitempty"`
	DateStarted           *string `json:"dateStarted,omitempty"`
	NumberOfSessions      *int64  `json:"numberOfSessions,omitempty"`
	DateConcluded         *string `json:"dateConcluded,omitempty"`
}

// Update Finance
type UpdateFinanceRequest struct {
	EmployedFamilyMembersCount int     `json:"employedFamilyMembersCount" binding:"required"`
	SupportsStudiesCount       int     `json:"supportsStudiesCount"`
	SupportsFamilyCount        int     `json:"supportsFamilyCount"`
	FinancialSupport           string  `json:"financialSupport" binding:"required"`
	WeeklyAllowance            float64 `json:"weeklyAllowance" binding:"required"`
}
