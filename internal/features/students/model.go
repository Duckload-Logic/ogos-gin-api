package students

type StudentProfileView struct {
	StudentRecordID int    `db:"student_record_id" json:"studentRecordId"`
	FirstName       string `db:"first_name" json:"firstName"`
	MiddleName      string `db:"middle_name" json:"middleName"`
	LastName        string `db:"last_name" json:"lastName"`
	Email           string `db:"email" json:"email"`
	Course          string `db:"course" json:"course"`
}

// StudentRecord model - matches the student_records table
type StudentRecord struct {
	ID          int    `db:"student_record_id" json:"id"`
	UserID      int    `db:"user_id" json:"userId"`
	IsSubmitted bool   `db:"is_submitted" json:"isSubmitted"`
	CreatedAt   string `db:"created_at" json:"createdAt"`
	UpdatedAt   string `db:"updated_at" json:"updatedAt"`
}

// StudentProfile model - matches the student_profiles table
type StudentProfile struct {
	ID                int      `db:"student_profile_id" json:"id"`
	StudentRecordID   int      `db:"student_record_id" json:"studentRecordId"`
	GenderID          int      `db:"gender_id" json:"genderId"`
	CivilStatusTypeID int      `db:"civil_status_type_id" json:"civilStatusTypeId"`
	Religion          string   `db:"religion" json:"religion"`
	HeightFt          *float64 `db:"height_ft" json:"heightFt"`
	WeightKg          *float64 `db:"weight_kg" json:"weightKg"`
	StudentNumber     string   `db:"student_number" json:"studentNumber"`
	HighSchoolGWA     *float64 `db:"high_school_gwa" json:"highSchoolGWA"`
	Course            string   `db:"course" json:"course"`
	PlaceOfBirth      *string  `db:"place_of_birth" json:"placeOfBirth"`
	BirthDate         *string  `db:"birth_date" json:"birthDate"`
	ContactNo         *string  `db:"contact_no" json:"contactNo"`
}

// StudentEmergencyContact model
type StudentEmergencyContact struct {
	ID                           int    `db:"emergency_contact_id" json:"id"`
	StudentRecordID              int    `db:"student_record_id" json:"studentRecordId"`
	ParentID                     *int   `db:"parent_id" json:"parentId"` // Can be NULL
	EmergencyContactName         string `db:"emergency_contact_name" json:"emergencyContactName"`
	EmergencyContactPhone        string `db:"emergency_contact_phone" json:"emergencyContactPhone"`
	EmergencyContactRelationship string `db:"emergency_contact_relationship" json:"emergencyContactRelationship"`
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

// Parent model
type Parent struct {
	ID               int     `db:"parent_id" json:"id"`
	EducationalLevel string  `db:"educational_level" json:"educationalLevel"`
	BirthDate        *string `db:"birth_date" json:"birthDate"`
	LastName         string  `db:"last_name" json:"lastName"`
	FirstName        string  `db:"first_name" json:"firstName"`
	MiddleName       *string `db:"middle_name" json:"middleName"`
	Occupation       *string `db:"occupation" json:"occupation"`
	CompanyName      *string `db:"company_name" json:"companyName"`
}

type ParentInfoView struct {
	Parent
	Relationship string `db:"relationship" json:"relationship"`
}

// StudentParent model for the junction table
type StudentParent struct {
	StudentRecordID  int  `db:"student_record_id" json:"studentRecordId"`
	ParentID         int  `db:"parent_id" json:"parentId"`
	Relationship     int  `db:"relationship" json:"relationship"`
	IsPrimaryContact bool `db:"is_primary_contact" json:"isPrimaryContact"`
}

// FamilyBackground model
type FamilyBackground struct {
	ID                    int     `db:"family_background_id" json:"id"`
	StudentRecordID       int     `db:"student_record_id" json:"studentRecordId"`
	ParentalStatusID      int     `db:"parental_status_id" json:"parentalStatusId"`
	ParentalStatusDetails *string `db:"parental_status_details" json:"parentalStatusDetails"`
	SiblingsBrothers      int     `db:"siblings_brothers" json:"siblingsBrothers"`
	SiblingSisters        int     `db:"sibling_sisters" json:"siblingSisters"`
	MonthlyFamilyIncome   string  `db:"monthly_family_income" json:"monthlyFamilyIncome"`
	GuardianName          string  `db:"guardian_name" json:"guardianName"`
	GuardianAddress       string  `db:"guardian_address" json:"guardianAddress"`
}

// EducationalBackground model
type EducationalBackground struct {
	ID               int     `db:"educational_background_id" json:"id"`
	StudentRecordID  int     `db:"student_record_id" json:"studentRecordId"`
	EducationalLevel string  `db:"educational_level" json:"educationalLevel"`
	SchoolName       string  `db:"school_name" json:"schoolName"`
	Location         *string `db:"location" json:"location"`
	SchoolType       string  `db:"school_type" json:"schoolType"`
	YearCompleted    string  `db:"year_completed" json:"yearCompleted"`
	Awards           *string `db:"awards" json:"awards"`
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
	ID                         int      `db:"finance_id" json:"id"`
	StudentRecordID            int      `db:"student_record_id" json:"studentRecordId"`
	EmployedFamilyMembersCount *int     `db:"employed_family_members_count" json:"employedFamilyMembersCount"`
	SupportsStudiesCount       *int     `db:"supports_studies_count" json:"supportsStudiesCount"`
	SupportsFamilyCount        *int     `db:"supports_family_count" json:"supportsFamilyCount"`
	FinancialSupport           string   `db:"financial_support" json:"financialSupport"`
	WeeklyAllowance            *float64 `db:"weekly_allowance" json:"weeklyAllowance"`
}
