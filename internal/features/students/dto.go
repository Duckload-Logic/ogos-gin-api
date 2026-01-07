package students

import "github.com/olazo-johnalbert/duckload-api/internal/core/request"

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

type GetStudentRequest struct {
	IncludeParams
}

type ComprehensiveProfileResponse struct {
	EnrollmentReasons []StudentSelectedReason  `json:"enrollmentReasons,omitempty"`
	StudentProfile    *StudentProfile          `json:"studentProfile"`
	EmergencyContact  *StudentEmergencyContact `json:"emergencyContact,omitempty"`
	Family            *FamilyBackground        `json:"family,omitempty"`
	Parents           []ParentInfoView         `json:"parents,omitempty"`
	Health            *StudentHealthRecord     `json:"health,omitempty"`
	Education         []EducationalBackground  `json:"education,omitempty"`
	Addresses         []StudentAddress         `json:"addresses,omitempty"`
	Finance           *StudentFinance          `json:"finance,omitempty"`
}

type GetEmergencyContactResponse struct {
	EmergencyContact *StudentEmergencyContact `json:"emergencyContact"`
}

// Combined response for student profile
type StudentProfileResponse struct {
	StudentProfile *StudentProfile `json:"studentProfile"`
}

type CreateStudentRecordRequest struct {
	GenderID          int                            `json:"genderId" binding:"required"`
	CivilStatusTypeID int                            `json:"civilStatusTypeId" binding:"required"`
	Religion          string                         `json:"religion" binding:"required"`
	HeightFt          float64                        `json:"heightFt" binding:"required"`
	WeightKg          float64                        `json:"weightKg" binding:"required"`
	StudentNumber     string                         `json:"studentNumber,omitempty" `
	Course            string                         `json:"course" binding:"required"`
	HighSchoolGWA     float64                        `json:"highSchoolGWA" binding:"required"`
	PlaceOfBirth      string                         `json:"placeOfBirth" binding:"required"`
	BirthDate         string                         `json:"birthDate" binding:"required"`
	ContactNo         string                         `json:"contactNo" binding:"required"`
	EmergencyContact  *UpdateEmergencyContactRequest `json:"emergencyContact" binding:"required"`
	Addresses         []StudentAddressDTO            `json:"addresses" binding:"required,min=1"`
}

type UpdateEmergencyContactRequest struct {
	ParentID                     *int   `json:"parentId,omitempty"`
	EmergencyContactName         string `json:"emergencyContactName" binding:"required"`
	EmergencyContactPhone        string `json:"emergencyContactPhone" binding:"required"`
	EmergencyContactRelationship string `json:"emergencyContactRelationship" binding:"required"`
}

type UpdateEnrollmentReasonsRequest struct {
	EnrollmentReasonIDs []int  `json:"enrollmentReasonIds"`
	OtherReasonText     string `json:"otherReasonText,omitempty"`
}

type UpdateFamilyRequest struct {
	// Family Background fields
	ParentalStatusID      int    `json:"parentalStatusId" binding:"required"`
	ParentalStatusDetails string `json:"parentalStatusDetails,omitempty"`
	SiblingsBrothers      *int   `json:"siblingsBrothers" binding:"required"`
	SiblingSisters        *int   `json:"siblingSisters" binding:"required"`
	MonthlyFamilyIncome   string `json:"monthlyFamilyIncome" binding:"required"`
	GuardianName          string `json:"guardianName" binding:"required"`
	GuardianAddress       string `json:"guardianAddress" binding:"required"`

	UpdateFinanceRequest
	// Parent Data
	Parents []ParentDTO `json:"parents" binding:"required,dive"`
}

type ParentDTO struct {
	FirstName        string `json:"firstName" binding:"required"`
	LastName         string `json:"lastName" binding:"required"`
	MiddleName       string `json:"middleName,omitempty"`
	EducationalLevel string `json:"educationalLevel" binding:"required"`
	BirthDate        string `json:"birthDate" binding:"required"` // Format: YYYY-MM-DD
	Occupation       string `json:"occupation" binding:"required"`
	CompanyName      string `json:"companyName,omitempty"`
	Relationship     int    `json:"relationship" binding:"required"`
}

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

type UpdateFinanceRequest struct {
	EmployedFamilyMembersCount int     `json:"employedFamilyMembersCount" binding:"required"`
	SupportsStudiesCount       int     `json:"supportsStudiesCount"`
	SupportsFamilyCount        int     `json:"supportsFamilyCount"`
	FinancialSupport           string  `json:"financialSupport" binding:"required"`
	WeeklyAllowance            float64 `json:"weeklyAllowance" binding:"required"`
}
