package students

type StudentProfileView struct {
	StudentRecordID int    `db:"student_record_id" json:"studentRecordId"`
	FirstName       string `db:"first_name" json:"firstName"`
	MiddleName      string `db:"middle_name" json:"middleName"`
	LastName        string `db:"last_name" json:"lastName"`
	Email           string `db:"email" json:"email"`
	Course          string `db:"course" json:"course"`
	YearLevel       int    `db:"year_level" json:"yearLevel"`
	Section         string `db:"section" json:"section"`
}

// StudentRecord model - matches the student_records table
type StudentRecord struct {
	ID     int `db:"student_record_id" json:"id"`
	UserID int `db:"user_id" json:"userId"`
}

// StudentProfile model - matches the student_profiles table
type StudentProfile struct {
	ID                  int      `db:"student_profile_id" json:"id"`
	StudentRecordID     int      `db:"student_record_id" json:"studentRecordId"`
	GenderID            int      `db:"gender_id" json:"genderId"`
	CivilStatusTypeID   int      `db:"civil_status_type_id" json:"civilStatusTypeId"`
	ReligionTypeID      int      `db:"religion_type_id" json:"religionTypeId"`
	HeightCm            *float64 `db:"height_cm" json:"heightCm"`
	WeightKg            *float64 `db:"weight_kg" json:"weightKg"`
	StudentNumber       string   `db:"student_number" json:"studentNumber"`
	Course              string   `db:"course" json:"course"`
	YearLevel           int      `db:"year_level" json:"yearLevel"`
	Section             *string  `db:"section" json:"section"`
	GoodMoralStatus     bool     `db:"good_moral_status" json:"goodMoralStatus"`
	HasDerogatoryRecord bool     `db:"has_derogatory_record" json:"hasDerogatoryRecord"`
	PlaceOfBirth        *string  `db:"place_of_birth" json:"placeOfBirth"`
	BirthDate           *string  `db:"birth_date" json:"birthDate"`
	MobileNo            *string  `db:"mobile_no" json:"mobileNo"`
}

// EnrollmentReason
type EnrollmentReason struct {
	ID   int    `db:"reason_id" json:"id"`
	Text string `db:"reason_text" json:"text"`
}

// StudentSelectedReason (Junction Table)
type StudentSelectedReason struct {
	StudentRecordID int     `db:"student_record_id" json:"studentRecordId"`
	ReasonID        int     `db:"reason_id" json:"reasonId"`
	OtherReasonText *string `db:"other_reason_text" json:"otherReasonText"`
}

// Guardian model
type Guardian struct {
	ID                 int     `db:"guardian_id" json:"id"`
	EducationalLevelID int     `db:"educational_level_id" json:"educationalLevelId"`
	BirthDate          *string `db:"birth_date" json:"birthDate"`
	LastName           string  `db:"last_name" json:"lastName"`
	FirstName          string  `db:"first_name" json:"firstName"`
	MiddleName         *string `db:"middle_name" json:"middleName"`
	Occupation         *string `db:"occupation" json:"occupation"`
	MaidenName         *string `db:"maiden_name" json:"maidenName"`
	CompanyName        *string `db:"company_name" json:"companyName"`
	ContactNumber      *string `db:"contact_number" json:"contactNumber"`
}

type GuardianInfoView struct {
	Guardian
	RelationshipTypeID int  `db:"relationship_type_id" json:"relationshipTypeId"`
	IsPrimaryContact   bool `db:"is_primary_contact" json:"isPrimaryContact"`
}

// StudentGuardian model for the junction table
type StudentGuardian struct {
	StudentRecordID    int  `db:"student_record_id" json:"studentRecordId"`
	GuardianID         int  `db:"guardian_id" json:"guardianId"`
	RelationshipTypeID int  `db:"relationship_type_id" json:"relationshipTypeId"`
	IsPrimaryContact   bool `db:"is_primary_contact" json:"isPrimaryContact"`
}

// FamilyBackground model
type FamilyBackground struct {
	ID                    int      `db:"family_background_id" json:"id"`
	StudentRecordID       int      `db:"student_record_id" json:"studentRecordId"`
	ParentalStatusID      int      `db:"parental_status_id" json:"parentalStatusId"`
	ParentalStatusDetails *string  `db:"parental_status_details" json:"parentalStatusDetails"`
	SiblingsBrothers      int      `db:"siblings_brothers" json:"siblingsBrothers"`
	SiblingSisters        int      `db:"sibling_sisters" json:"siblingSisters"`
	MonthlyFamilyIncome   *float64 `db:"monthly_family_income" json:"monthlyFamilyIncome"`
}

// EducationalBackground model
type EducationalBackground struct {
	ID                 int     `db:"educational_background_id" json:"id"`
	StudentRecordID    int     `db:"student_record_id" json:"studentRecordId"`
	EducationalLevelID int     `db:"educational_level_id" json:"educationalLevelId"`
	SchoolName         string  `db:"school_name" json:"schoolName"`
	Location           *string `db:"location" json:"location"`
	SchoolType         string  `db:"school_type" json:"schoolType"`
	YearCompleted      string  `db:"year_completed" json:"yearCompleted"`
	Awards             *string `db:"awards" json:"awards"`
}

// StudentAddress model - Fixed to match schema
type StudentAddress struct {
	ID              int     `db:"student_address_id" json:"id"`
	StudentRecordID int     `db:"student_record_id" json:"studentRecordId"`
	AddressType     string  `db:"address_type" json:"addressType"` // Changed from AddressTypeID
	RegionName      *string `db:"region_name" json:"regionName"`
	ProvinceName    *string `db:"province_name" json:"provinceName"`
	CityName        *string `db:"city_name" json:"cityName"`
	BarangayName    *string `db:"barangay_name" json:"barangayName"`
	StreetLotBlk    *string `db:"street_lot_blk" json:"streetLotBlk"`
	UnitNo          *string `db:"unit_no" json:"unitNo"`
	BuildingName    *string `db:"building_name" json:"buildingName"`
}

// StudentHealthRecord model - Fixed to use enums directly
type StudentHealthRecord struct {
	ID                    int     `db:"health_id" json:"id"`
	StudentRecordID       int     `db:"student_record_id" json:"studentRecordId"`
	VisionRemark          string  `db:"vision_remark" json:"visionRemark"`                // Changed from ID
	HearingRemark         string  `db:"hearing_remark" json:"hearingRemark"`              // Changed from ID
	MobilityRemark        string  `db:"mobility_remark" json:"mobilityRemark"`            // Changed from ID
	SpeechRemark          string  `db:"speech_remark" json:"speechRemark"`                // Changed from ID
	GeneralHealthRemark   string  `db:"general_health_remark" json:"generalHealthRemark"` // Changed from ID
	ConsultedProfessional *string `db:"consulted_professional" json:"consultedProfessional"`
	ConsultationReason    *string `db:"consultation_reason" json:"consultationReason"`
	DateStarted           *string `db:"date_started" json:"dateStarted"` // Changed from string to sql.NullTime
	NumberOfSessions      *int64  `db:"num_sessions" json:"numberOfSessions"`
	DateConcluded         *string `db:"date_concluded" json:"dateConcluded"` // Changed from string to sql.NullTime
}

// Add to model.go
type StudentFinance struct {
	ID                     int      `db:"finance_id" json:"id"`
	StudentRecordID        int      `db:"student_record_id" json:"studentRecordId"`
	IsEmployed             *bool    `db:"is_employed" json:"isEmployed"`
	SupportsStudies        *bool    `db:"supports_studies" json:"supportsStudies"`
	SupportsFamily         *bool    `db:"supports_family" json:"supportsFamily"`
	FinancialSupportTypeID int      `db:"financial_support_type_id" json:"financialSupportTypeId"`
	WeeklyAllowance        *float64 `db:"weekly_allowance" json:"weeklyAllowance"`
}
