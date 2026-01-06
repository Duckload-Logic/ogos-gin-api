package students

type CreateStudentRecordRequest struct {
	GenderID            int     `json:"genderId" binding:"required"`
	CivilStatusTypeID   int     `json:"civilStatusTypeId" binding:"required"`
	ReligionTypeID      int     `json:"religionTypeId" binding:"required"`
	HeightCm            float64 `json:"heightCm" binding:"required"`
	WeightKg            float64 `json:"weightKg" binding:"required"`
	StudentNumber       string  `json:"studentNumber" binding:"required"`
	Course              string  `json:"course" binding:"required"`
	YearLevel           int     `json:"yearLevel" binding:"required"`
	Section             string  `json:"section,omitempty"` // Optional
	GoodMoralStatus     bool    `json:"goodMoralStatus"`   // Booleans default to false if missing
	HasDerogatoryRecord bool    `json:"hasDerogatoryRecord"`
	PlaceOfBirth        string  `json:"placeOfBirth" binding:"required"`
	BirthDate           string  `json:"birthDate" binding:"required"`
	MobileNo            string  `json:"mobileNo" binding:"required"`
}

type UpdateFamilyRequest struct {
	// Family Background fields
	ParentalStatusID      int     `json:"parentalStatusId" binding:"required"`
	ParentalStatusDetails string  `json:"parentalStatusDetails,omitempty"`
	SiblingsBrothers      int     `json:"siblingsBrothers" binding:"required"`
	SiblingSisters        int     `json:"siblingSisters" binding:"required"`
	MonthlyFamilyIncome   float64 `json:"monthlyFamilyIncome" binding:"required"`

	// Guardian Data
	Guardians []GuardianDTO `json:"guardians" binding:"required,dive"`
}

type GuardianDTO struct {
	FirstName          string `json:"firstName" binding:"required"`
	LastName           string `json:"lastName" binding:"required"`
	MiddleName         string `json:"middleName,omitempty"`
	EducationalLevelID int    `json:"educationalLevelId" binding:"required"`
	BirthDate          string `json:"birthDate" binding:"required"` // Format: YYYY-MM-DD
	Occupation         string `json:"occupation" binding:"required"`
	MaidenName         string `json:"maidenName,omitempty"`
	CompanyName        string `json:"companyName,omitempty"`
	IsPrimary          bool   `json:"isPrimary"` // Indicates if this guardian is the primary contact
	RelationshipTypeID int    `json:"relationshipTypeId,omitempty"`
	ContactNumber      string `json:"contactNumber,omitempty"`
}

type UpdateEducationRequest struct {
	EducationalBGs []EducationalBGDTO `json:"educationalBackgrounds" binding:"required,dive"`
}

type EducationalBGDTO struct {
	EducationalLevelID int    `json:"educationalLevelId" binding:"required"`
	SchoolName         string `json:"schoolName" binding:"required"`
	Location           string `json:"location,omitempty"`
	SchoolType         string `json:"schoolType" binding:"required,oneof=Public Private"`
	YearCompleted      string `json:"yearCompleted" binding:"required"`
	Awards             string `json:"awards,omitempty"`
}

type UpdateAddressRequest struct {
	Addresses []StudentAddressDTO `json:"addresses" binding:"required,dive"`
}

type StudentAddressDTO struct {
	AddressTypeID int    `json:"addressTypeId" binding:"required"`
	RegionName    string `json:"regionName" binding:"required"`
	ProvinceName  string `json:"provinceName" binding:"required"`
	CityName      string `json:"cityName" binding:"required"`
	BarangayName  string `json:"barangayName" binding:"required"`
	StreetLotBlk  string `json:"streetLotBlk,omitempty"`
	UnitNo        string `json:"unitNo,omitempty"`
	BuildingName  string `json:"buildingName,omitempty"`
}

type UpdateHealthRecordRequest struct {
	VisionRemarkID        int     `json:"visionRemarkId" binding:"required"`
	HearingRemarkID       int     `json:"hearingRemarkId" binding:"required"`
	MobilityRemarkID      int     `json:"mobilityRemarkId" binding:"required"`
	SpeechRemarkID        int     `json:"speechRemarkId" binding:"required"`
	GeneralHealthRemarkID int     `json:"generalHealthRemarkId" binding:"required"`
	ConsultedProfessional *string `json:"consultedProfessional,omitempty"`
	ConsultationReason    *string `json:"consultationReason,omitempty"`
	DateStarted           *string `json:"dateStarted,omitempty"`
	NumberOfSessions      *int    `json:"numberOfSessions,omitempty"`
	DateConcluded         *string `json:"dateConcluded,omitempty"`
}
