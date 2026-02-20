package students

// Core Student Records
type StudentRecord struct {
	ID          int    `db:"student_record_id" json:"id"`
	UserID      int    `db:"user_id" json:"userId"`
	IsSubmitted bool   `db:"is_submitted" json:"isSubmitted"`
	CreatedAt   string `db:"created_at" json:"createdAt"`
	UpdatedAt   string `db:"updated_at" json:"updatedAt"`
}

type StudentProfile struct {
	ID                int      `db:"student_profile_id" json:"id"`
	StudentRecordID   int      `db:"student_record_id" json:"studentRecordId"`
	StudentNumber     string   `db:"student_number" json:"studentNumber"`
	GenderID          int      `db:"gender_id" json:"genderId"`
	CivilStatusTypeID int      `db:"civil_status_type_id" json:"civilStatusTypeId"`
	Religion          string   `db:"religion" json:"religion"`
	HeightFt          *float64 `db:"height_ft" json:"heightFt"`
	WeightKg          *float64 `db:"weight_kg" json:"weightKg"`
	Complexion        string   `db:"complexion" json:"complexion"`
	HighSchoolGWA     *float64 `db:"high_school_gwa" json:"highSchoolGWA"`
	Course            string   `db:"course" json:"course"`
	PlaceOfBirth      *string  `db:"place_of_birth" json:"placeOfBirth"`
	DateOfBirth       *string  `db:"birth_date" json:"dateOfBirth"`
	ContactNumber     *string  `db:"contact_no" json:"contactNumber"`
}

// Enrollment and Reasons
type EnrollmentReason struct {
	ID   int    `db:"reason_id" json:"id"`
	Text string `db:"reason_text" json:"text"`
}

type StudentSelectedReason struct {
	StudentRecordID int     `db:"student_record_id" json:"studentRecordId"`
	ReasonID        int     `db:"reason_id" json:"reasonId"`
	OtherReasonText *string `db:"other_reason_text" json:"otherReasonText"`
}

// Family and Related Persons
type RelatedPerson struct {
	ID               int     `db:"related_person_id" json:"id"`
	AddressID        int     `db:"address_id" json:"addressId"`
	EducationalLevel string  `db:"educational_level" json:"educationalLevel"`
	DateOfBirth      *string `db:"birth_date" json:"dateOfBirth"`
	LastName         string  `db:"last_name" json:"lastName"`
	FirstName        string  `db:"first_name" json:"firstName"`
	MiddleName       *string `db:"middle_name" json:"middleName"`
	Occupation       *string `db:"occupation" json:"occupation"`
	EmployerName     *string `db:"employer_name" json:"employerName"`
	EmployerAddress  *string `db:"employer_address" json:"employerAddress"`
	IsLiving         bool    `db:"is_living" json:"isLiving"`
}

type StudentRelatedPerson struct {
	StudentRecordID int `db:"student_record_id" json:"studentRecordId"`
	PersonID        int `db:"related_person_id" json:"personId"`
	Relationship    int `db:"relationship" json:"relationship"` // "Father", "Mother", "Guardian", "Uncle", "Aunt", "Sibling", "Other"

	// Roles
	IsParent           bool `db:"is_parent" json:"isParent"`
	IsGuardian         bool `db:"is_guardian" json:"isGuardian"`
	IsEmergencyContact bool `db:"is_emergency_contact" json:"isEmergencyContact"`
}

type FamilyBackground struct {
	ID                    int     `db:"family_background_id" json:"id"`
	StudentRecordID       int     `db:"student_record_id" json:"studentRecordId"`
	ParentalStatusID      int     `db:"parental_status_id" json:"parentalStatusId"`
	ParentalStatusDetails *string `db:"parental_status_details" json:"parentalStatusDetails"`
	Brothers              int     `db:"siblings_brothers" json:"brothers"`
	Sisters               int     `db:"sibling_sisters" json:"sisters"`
	MonthlyFamilyIncome   string  `db:"monthly_family_income" json:"monthlyFamilyIncome"`
}

// Education and Background
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

// Location and Address
type Address struct {
	ID           int     `db:"address_id" json:"id"`
	Region       string  `db:"region" json:"region"`
	City         string  `db:"city" json:"city"`
	Barangay     string  `db:"barangay" json:"barangay"`
	StreetDetail *string `db:"street_detail" json:"streetDetail"` // Lot/Blk/Street
}

type StudentAddress struct {
	ID              int    `db:"student_address_id" json:"id"`
	StudentRecordID int    `db:"student_record_id" json:"studentRecordId"`
	AddressID       int    `db:"address_id" json:"addressId"`
	AddressType     string `db:"address_type" json:"addressType"`
}

// Health and Wellness
type StudentHealthRecord struct {
	ID                    int     `db:"health_id" json:"id"`
	StudentRecordID       int     `db:"student_record_id" json:"studentRecordId"`
	VisionRemark          string  `db:"vision_remark" json:"visionRemark"`
	HearingRemark         string  `db:"hearing_remark" json:"hearingRemark"`
	MobilityRemark        string  `db:"mobility_remark" json:"mobilityRemark"`
	SpeechRemark          string  `db:"speech_remark" json:"speechRemark"`
	GeneralHealthRemark   string  `db:"general_health_remark" json:"generalHealthRemark"`
	ConsultedProfessional *string `db:"consulted_professional" json:"consultedProfessional"`
	ConsultationReason    *string `db:"consultation_reason" json:"consultationReason"`
	DateStarted           *string `db:"date_started" json:"dateStarted"`
	NumberOfSessions      *int64  `db:"num_sessions" json:"numberOfSessions"`
	DateConcluded         *string `db:"date_concluded" json:"dateConcluded"`
}

// Financial Support
type StudentFinance struct {
	ID                         int      `db:"finance_id" json:"id"`
	StudentRecordID            int      `db:"student_record_id" json:"studentRecordId"`
	EmployedFamilyMembersCount *int     `db:"employed_family_members_count" json:"employedFamilyMembersCount"`
	SupportsStudiesCount       *int     `db:"supports_studies_count" json:"supportsStudiesCount"`
	SupportsFamilyCount        *int     `db:"supports_family_count" json:"supportsFamilyCount"`
	FinancialSupport           string   `db:"financial_support" json:"financialSupport"`
	WeeklyAllowance            *float64 `db:"weekly_allowance" json:"weeklyAllowance"`
}

// NEW: Section V - Interests and Hobbies
type StudentInterest struct {
	ID              int    `db:"interest_id" json:"id"`
	StudentRecordID int    `db:"student_record_id" json:"studentRecordId"`
	Type            string `db:"interest_type" json:"type"` // e.g., "Academic", "Extra-Curricular"
	Name            string `db:"interest_name" json:"name"` // e.g., "Math Club", "Chess"
	IsFavorite      bool   `db:"is_favorite" json:"isFavorite"`
	Rank            int    `db:"rank" json:"rank"` // For the 1, 2, 3, 4 ranking
}

// NEW: Section VI & VII - For Guidance Use
type TestResult struct {
	ID              int    `db:"test_result_id" json:"id"`
	StudentRecordID int    `db:"student_record_id" json:"studentRecordId"`
	TestDate        string `db:"test_date" json:"testDate"`
	TestName        string `db:"test_name" json:"testName"`
	RawScore        string `db:"raw_score" json:"rawScore"`
	Percentile      string `db:"percentile" json:"percentile"`
	Description     string `db:"description" json:"description"`
}
