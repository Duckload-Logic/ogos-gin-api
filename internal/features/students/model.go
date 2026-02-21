package students

// Core Student Records
type IIRRecord struct {
	ID          int    `db:"id" json:"id"`
	UserID      int    `db:"user_id" json:"userId"`
	IsSubmitted bool   `db:"is_submitted" json:"isSubmitted"`
	CreatedAt   string `db:"created_at" json:"createdAt"`
	UpdatedAt   string `db:"updated_at" json:"updatedAt"`
}

type StudentProfile struct {
	ID                int      `db:"id" json:"id"`
	IIRID             int      `db:"iir_id" json:"studentRecordId"`
	StudentNumber     string   `db:"student_number" json:"studentNumber"`
	GenderID          int      `db:"gender_id" json:"genderId"`
	CivilStatusTypeID int      `db:"civil_status_type_id" json:"civilStatusTypeId"`
	Religion          string   `db:"religion" json:"religion"`
	HeightFt          *float64 `db:"height_ft" json:"heightFt"`
	WeightKg          *float64 `db:"weight_kg" json:"weightKg"`
	Complexion        string   `db:"complexion" json:"complexion"`
	HighSchoolGWA     *float64 `db:"high_school_gwa" json:"highSchoolGWA"`
	Course            string   `db:"course" json:"course"`
	YearLevel         string   `db:"year_level" json:"yearLevel"`
	Section           string   `db:"section" json:"section"`
	PlaceOfBirth      *string  `db:"place_of_birth" json:"placeOfBirth"`
	DateOfBirth       *string  `db:"date_of_birth" json:"dateOfBirth"`
	ContactNumber     *string  `db:"contact_no" json:"contactNumber"`
}

// Enrollment and Reasons
type EnrollmentReason struct {
	ID   int    `db:"id" json:"id"`
	Text string `db:"reason_text" json:"text"`
}

type StudentSelectedReason struct {
	IIRID           int     `db:"iir_id" json:"studentRecordId"`
	ReasonID        int     `db:"reason_id" json:"reasonId"`
	OtherReasonText *string `db:"other_reason_text" json:"otherReasonText"`
}

// Family and Related Persons
type RelatedPerson struct {
	ID               int     `db:"id" json:"id"`
	AddressID        int     `db:"address_id" json:"addressId"`
	EducationalLevel string  `db:"educational_level" json:"educationalLevel"`
	DateOfBirth      *string `db:"date_of_birth" json:"dateOfBirth"`
	LastName         string  `db:"last_name" json:"lastName"`
	FirstName        string  `db:"first_name" json:"firstName"`
	MiddleName       *string `db:"middle_name" json:"middleName"`
	Occupation       *string `db:"occupation" json:"occupation"`
	EmployerName     *string `db:"employer_name" json:"employerName"`
	EmployerAddress  *string `db:"employer_address" json:"employerAddress"`
	IsLiving         bool    `db:"is_living" json:"isLiving"`
}

type StudentRelatedPerson struct {
	IIRID        int `db:"iir_id" json:"studentRecordId"`
	PersonID     int `db:"related_person_id" json:"personId"`
	Relationship int `db:"relationship" json:"relationship"` // "Father", "Mother", "Guardian", "Uncle", "Aunt", "Sibling", "Other"

	// Roles
	IsParent           bool `db:"is_parent" json:"isParent"`
	IsGuardian         bool `db:"is_guardian" json:"isGuardian"`
	IsEmergencyContact bool `db:"is_emergency_contact" json:"isEmergencyContact"`
}

type FamilyBackground struct {
	ID                    int     `db:"id" json:"id"`
	IIRID                 int     `db:"iir_id" json:"studentRecordId"`
	ParentalStatus        string  `db:"parental_status" json:"parentalStatus"`
	ParentalStatusDetails *string `db:"parental_status_details" json:"parentalStatusDetails"`
	Brothers              int     `db:"brothers" json:"brothers"`
	Sisters               int     `db:"sisters" json:"sisters"`
	EmployedSiblings      int     `db:"employed_siblings" json:"employedSiblings"`
	OrdinalPosition       int     `db:"ordinal_position" json:"ordinalPosition"`
	HaveQuietPlaceToStudy bool    `db:"have_quiet_place_to_study" json:"haveQuietPlaceToStudy"`
	IsSharingRoom         bool    `db:"is_sharing_room" json:"isSharingRoom"`
	RoomSharingDetails    *string `db:"room_sharing_details" json:"roomSharingDetails"`
	NatureOfResidence     string  `db:"nature_of_residence" json:"natureOfResidence"`
}

// Education and Background
type EducationalBackground struct {
	ID               int     `db:"id" json:"id"`
	IIRID            int     `db:"iir_id" json:"studentRecordId"`
	EducationalLevel string  `db:"educational_level" json:"educationalLevel"`
	SchoolName       string  `db:"school_name" json:"schoolName"`
	Location         *string `db:"location" json:"location"`
	SchoolType       string  `db:"school_type" json:"schoolType"`
	YearCompleted    string  `db:"year_completed" json:"yearCompleted"`
	Awards           *string `db:"awards" json:"awards"`
}

// Location and Address
type Address struct {
	ID           int     `db:"id" json:"id"`
	Region       string  `db:"region" json:"region"`
	City         string  `db:"city" json:"city"`
	Barangay     string  `db:"barangay" json:"barangay"`
	StreetDetail *string `db:"street_detail" json:"streetDetail"` // Lot/Blk/Street
}

type StudentAddress struct {
	ID          int    `db:"id" json:"id"`
	IIRID       int    `db:"iir_id" json:"studentRecordId"`
	AddressID   int    `db:"address_id" json:"addressId"`
	AddressType string `db:"address_type" json:"addressType"`
}

// Health and Wellness
type StudentHealthRecord struct {
	ID                      int     `db:"id" json:"id"`
	IIRID                   int     `db:"student_record_id" json:"studentRecordId"`
	VisionHasProblem        bool    `db:"vision_has_problem" json:"visionHasProblem"`
	VisionDetails           *string `db:"vision_details" json:"visionDetails"`
	HearingHasProblem       bool    `db:"hearing_has_problem" json:"hearingHasProblem"`
	HearingDetails          *string `db:"hearing_details" json:"hearingDetails"`
	SpeechHasProblem        bool    `db:"speech_has_problem" json:"speechHasProblem"`
	SpeechDetails           *string `db:"speech_details" json:"speechDetails"`
	GeneralHealthHasProblem bool    `db:"general_health_has_problem" json:"generalHealthHasProblem"`
	GeneralHealthDetails    *string `db:"general_health_details" json:"generalHealthDetails"`
}

// Financial Support
type StudentFinance struct {
	ID                         int      `db:"id" json:"id"`
	IIRID                      int      `db:"iir_id" json:"studentRecordId"`
	MonthlyFamilyIncomeRangeID *int     `db:"monthly_family_income_range_id" json:"monthlyFamilyIncome RangeId"`
	OtherIncomeDetails         *string  `db:"other_income_details" json:"otherIncomeDetails"`
	FinancialSupportTypeID     *int     `db:"financial_support_type_id" json:"financialSupportTypeId"`
	WeeklyAllowance            *float64 `db:"weekly_allowance" json:"weeklyAllowance"`
}

// Interests and Activities
type StudentInterest struct {
	ID              int    `db:"id" json:"id"`
	IIRID           int    `db:"iir_id" json:"studentRecordId"`
	Type            string `db:"interest_type" json:"type"` // e.g., "Academic", "Extra-Curricular"
	Name            string `db:"interest_name" json:"name"` // e.g., "Math Club", "Chess"
	IsFavorite      bool   `db:"is_favorite" json:"isFavorite"`
	IsLeastFavorite bool   `db:"is_least_favorite" json:"isLeastFavorite"`
	Rank            int    `db:"rank" json:"rank"` // For students to rank their interests (1 = most interested)
}

// Test Results and Assessments
type TestResult struct {
	ID          int    `db:"id" json:"id"`
	IIRID       int    `db:"iir_id" json:"studentRecordId"`
	TestDate    string `db:"test_date" json:"testDate"`
	TestName    string `db:"test_name" json:"testName"`
	RawScore    string `db:"raw_score" json:"rawScore"`
	Percentile  string `db:"percentile" json:"percentile"`
	Description string `db:"description" json:"description"`
}
