package students

import "time"

// Lookup models
type Gender struct {
	ID         int    `db:"id" json:"id"`
	GenderName string `db:"gender_name" json:"name"`
}

type ParentalStatusType struct {
	ID         int    `db:"id" json:"id"`
	StatusName string `db:"status_name" json:"name"`
}

type StudentSupportType struct {
	ID              int    `db:"id" json:"id"`
	SupportTypeName string `db:"support_type_name" json:"name"`
}

type EnrollmentReason struct {
	ID   int    `db:"id" json:"id"`
	Text string `db:"reason_text" json:"text"`
}

type IncomeRange struct {
	ID        int    `db:"id" json:"id"`
	RangeText string `db:"range_text" json:"text"`
}

type EducationalLevel struct {
	ID        int    `db:"id" json:"id"`
	LevelName string `db:"level_name" json:"name"`
}

type Course struct {
	ID         int    `db:"id" json:"id"`
	Code       string `db:"code" json:"code"`
	CourseName string `db:"course_name" json:"name"`
}

type CivilStatusType struct {
	ID         int    `db:"id" json:"id"`
	StatusName string `db:"status_name" json:"name"`
}

type Religion struct {
	ID           int    `db:"id" json:"id"`
	ReligionName string `db:"religion_name" json:"name"`
}

type StudentRelationshipType struct {
	ID               int    `db:"id" json:"id"`
	RelationshipName string `db:"relationship_name" json:"name"`
}

type NatureOfResidenceType struct {
	ID                int    `db:"id" json:"id"`
	ResidenceTypeName string `db:"residence_type_name" json:"name"`
}

type SibilingSupportType struct {
	ID          int    `db:"id" json:"id"`
	SupportName string `db:"name" json:"name"`
}

type ActivityOption struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`         // "Academic" or "Extra-Curricular"
	Category string `db:"category" json:"category"` // e.g., "Basketball", "Student Government", "Volunteering"
	IsActive bool   `db:"is_active" json:"isActive"`
}

type StudentBasicInfoView struct {
	ID         int     `json:"id"`
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName,omitempty"`
	LastName   string  `json:"lastName"`
	Email      string  `json:"email"`
}

type StudentProfileView struct {
	IIRID         int     `db:"iir_id" json:"iirId"`
	UserID        int     `db:"user_id" json:"userId"`
	FirstName     string  `db:"first_name" json:"firstName"`
	MiddleName    *string `db:"middle_name" json:"middleName,omitempty"`
	LastName      string  `db:"last_name" json:"lastName"`
	Email         string  `db:"email" json:"email"`
	StudentNumber string  `db:"student_number" json:"studentNumber"`
	GenderID      int     `db:"gender_id" json:"genderId"`
	CourseID      int     `db:"course" json:"course"`
	Section       int     `db:"section" json:"section"`
	YearLevel     int     `db:"year_level" json:"yearLevel"`
}

// Core Student Records
type IIRRecord struct {
	ID          int       `db:"id" json:"id"`
	UserID      int       `db:"user_id" json:"userId"`
	IsSubmitted bool      `db:"is_submitted" json:"isSubmitted"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

// Enrollment and Reasons
type StudentSelectedReason struct {
	IIRID           int     `db:"iir_id" json:"iirId"`
	ReasonID        int     `db:"reason_id" json:"reasonId"`
	OtherReasonText *string `db:"other_reason_text" json:"otherReasonText"`
}

type StudentPersonalInfo struct {
	ID            int    `db:"id" json:"id"`
	IIRID         int    `db:"iir_id" json:"iirId"`
	StudentNumber string `db:"student_number" json:"studentNumber"`

	// Personal Details
	GenderID      int `db:"gender_id" json:"genderId"`
	CivilStatusID int `db:"civil_status_id" json:"civilStatusId"`
	ReligionID    int `db:"religion_id" json:"religionId"`

	// Physical Attributes
	HeightFt   float64 `db:"height_ft" json:"heightFt"`
	WeightKg   float64 `db:"weight_kg" json:"weightKg"`
	Complexion string  `db:"complexion" json:"complexion"`

	// Academic Information
	HighSchoolGWA float64 `db:"high_school_gwa" json:"highSchoolGWA"`
	CourseID      int     `db:"course_id" json:"courseId"`
	YearLevel     int     `db:"year_level" json:"yearLevel"`
	Section       int     `db:"section" json:"section"`

	// Additional Details
	PlaceOfBirth                   string    `db:"place_of_birth" json:"placeOfBirth"`
	DateOfBirth                    string    `db:"date_of_birth" json:"dateOfBirth"`
	IsEmployed                     bool      `db:"is_employed" json:"isEmployed"`
	EmployerName                   *string   `db:"employer_name" json:"employerName"`
	EmployerAddress                *string   `db:"employer_address" json:"employerAddress"`
	MobileNumber                   string    `db:"mobile_number" json:"mobileNumber"`
	TelephoneNumber                *string   `db:"telephone_number" json:"telephoneNumber"`
	EmergencyContactName           string    `db:"emergency_contact_name" json:"emergencyContactName"`
	EmergencyContactNumber         string    `db:"emergency_contact_number" json:"emergencyContactNumber"`
	EmergencyContactRelationshipID int       `db:"emergency_contact_relationship_id" json:"emergencyRelationshipId"`
	EmergencyContactAddressID      int       `db:"emergency_contact_address_id" json:"emergencyAddressId"`
	CreatedAt                      time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt                      time.Time `db:"updated_at" json:"updatedAt"`
}

// Location and Address
type Address struct {
	ID           int       `db:"id" json:"id"`
	Region       string    `db:"region" json:"region"`
	City         string    `db:"city" json:"city"`
	Barangay     string    `db:"barangay" json:"barangay"`
	StreetDetail *string   `db:"street_detail" json:"streetDetail"` // Lot/Blk/Street
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}

type StudentAddress struct {
	ID          int       `db:"id" json:"id"`
	IIRID       int       `db:"iir_id" json:"iirId"`
	AddressID   int       `db:"address_id" json:"addressId"`
	AddressType string    `db:"address_type" json:"addressType"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

// Family and Related Persons
type RelatedPerson struct {
	ID               int       `db:"id" json:"id"`
	AddressID        *int      `db:"address_id" json:"addressId,omitempty"`
	EducationalLevel string    `db:"educational_level" json:"educationalLevel"`
	DateOfBirth      *string   `db:"date_of_birth" json:"dateOfBirth"`
	LastName         string    `db:"last_name" json:"lastName"`
	FirstName        string    `db:"first_name" json:"firstName"`
	MiddleName       *string   `db:"middle_name" json:"middleName"`
	Occupation       *string   `db:"occupation" json:"occupation"`
	EmployerName     *string   `db:"employer_name" json:"employerName"`
	EmployerAddress  *string   `db:"employer_address" json:"employerAddress"`
	ContactNumber    *string   `db:"contact_number" json:"contactNumber"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time `db:"updated_at" json:"updatedAt"`
}

type StudentRelatedPerson struct {
	IIRID           int `db:"iir_id" json:"iirId"`
	RelatedPersonID int `db:"related_person_id" json:"relatedPersonId"`

	// "Father", "Mother", "Guardian", "Uncle", "Aunt", "Sibling", "Other"
	RelationshipID int `db:"relationship_id" json:"relationshipId"`

	// Roles
	IsParent   bool `db:"is_parent" json:"isParent"`
	IsGuardian bool `db:"is_guardian" json:"isGuardian"`
	IsLiving   bool `db:"is_living" json:"isLiving"`

	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type FamilyBackground struct {
	ID                    int       `db:"id" json:"id"`
	IIRID                 int       `db:"iir_id" json:"iirId"`
	ParentalStatusID      int       `db:"parental_status_id" json:"parentalStatusId"`
	ParentalStatusDetails *string   `db:"parental_status_details" json:"parentalStatusDetails"`
	Brothers              int       `db:"brothers" json:"brothers"`
	Sisters               int       `db:"sisters" json:"sisters"`
	EmployedSiblings      int       `db:"employed_siblings" json:"employedSiblings"`
	OrdinalPosition       int       `db:"ordinal_position" json:"ordinalPosition"`
	HaveQuietPlaceToStudy bool      `db:"have_quiet_place_to_study" json:"haveQuietPlaceToStudy"`
	IsSharingRoom         bool      `db:"is_sharing_room" json:"isSharingRoom"`
	RoomSharingDetails    *string   `db:"room_sharing_details" json:"roomSharingDetails"`
	NatureOfResidenceId   int       `db:"nature_of_residence_id" json:"natureOfResidenceId"`
	CreatedAt             time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt             time.Time `db:"updated_at" json:"updatedAt"`
}

type StudentSiblingSupport struct {
	FamilyBackgroundID int `db:"family_background_id" json:"familyBackgroundId"`
	SupportTypeID      int `db:"support_type_id" json:"supportTypeId"`
}

// Education and Background
type EducationalBackground struct {
	ID                 int       `db:"id" json:"id"`
	IIRID              int       `db:"iir_id" json:"iirId"`
	NatureOfSchooling  string    `db:"nature_of_schooling" json:"natureOfSchooling"`
	InterruptedDetails *string   `db:"interrupted_details" json:"interruptedDetails,omitempty"`
	CreatedAt          time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt          time.Time `db:"updated_at" json:"updatedAt"`
}

type SchoolDetails struct {
	ID                 int       `db:"id" json:"id"`
	EBID               int       `db:"eb_id" json:"ebId"`
	EducationalLevelID int       `db:"educational_level_id" json:"educationalLevelId"`
	SchoolName         string    `db:"school_name" json:"schoolName"`
	SchoolAddress      string    `db:"school_address" json:"schoolAddress"`
	SchoolType         string    `db:"school_type" json:"schoolType"`
	YearStarted        int       `db:"year_started" json:"yearStarted"`
	YearCompleted      int       `db:"year_completed" json:"yearCompleted"`
	Awards             *string   `db:"awards" json:"awards,omitempty"`
	CreatedAt          time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt          time.Time `db:"updated_at" json:"updatedAt"`
}

// Health and Wellness
type StudentHealthRecord struct {
	ID                      int       `db:"id" json:"id"`
	IIRID                   int       `db:"iir_id" json:"iirID"`
	VisionHasProblem        bool      `db:"vision_has_problem" json:"visionHasProblem"`
	VisionDetails           *string   `db:"vision_details" json:"visionDetails"`
	HearingHasProblem       bool      `db:"hearing_has_problem" json:"hearingHasProblem"`
	HearingDetails          *string   `db:"hearing_details" json:"hearingDetails"`
	SpeechHasProblem        bool      `db:"speech_has_problem" json:"speechHasProblem"`
	SpeechDetails           *string   `db:"speech_details" json:"speechDetails"`
	GeneralHealthHasProblem bool      `db:"general_health_has_problem" json:"generalHealthHasProblem"`
	GeneralHealthDetails    *string   `db:"general_health_details" json:"generalHealthDetails"`
	CreatedAt               time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt               time.Time `db:"updated_at" json:"updatedAt"`
}

type StudentConsultation struct {
	ID               int       `db:"id" json:"id"`
	IIRID            int       `db:"iir_id" json:"iirId"`
	ProfessionalType string    `db:"professional_type" json:"professionalType"` // e.g., "Psychiatrist", "Psychologist", "Counselor"
	HasConsulted     bool      `db:"has_consulted" json:"hasConsulted"`
	WhenDate         *string   `db:"when_date" json:"whenDate"`
	ForWhat          *string   `db:"for_what" json:"forWhat"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time `db:"updated_at" json:"updatedAt"`
}

// Financial Support
type StudentFinance struct {
	ID                         int       `db:"id" json:"id"`
	IIRID                      int       `db:"iir_id" json:"iirId"`
	MonthlyFamilyIncomeRangeID int       `db:"monthly_family_income_range_id" json:"monthlyFamilyIncomeRangeId"`
	OtherIncomeDetails         *string   `db:"other_income_details" json:"otherIncomeDetails"`
	WeeklyAllowance            *float64  `db:"weekly_allowance" json:"weeklyAllowance"`
	CreatedAt                  time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt                  time.Time `db:"updated_at" json:"updatedAt"`
}

type StudentFinancialSupport struct {
	StudentFinanceID int       `db:"sf_id" json:"studentFinanceId"`
	SupportTypeID    int       `db:"support_type_id" json:"supportTypeId"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time `db:"updated_at" json:"updatedAt"`
}

// Interests and Activities
type StudentActivity struct {
	ID                 int       `db:"id" json:"id"`
	IIRID              int       `db:"iir_id" json:"iirId"`
	OptionID           int       `db:"option_id" json:"optionId,omitempty"`                     // FK to interest_options for standardized activities
	OtherSpecification *string   `db:"other_specification" json:"otherSpecification,omitempty"` // For "Other" category
	Role               string    `db:"role" json:"role"`                                        // "Officer", "Member", "Other"
	RoleSpecification  *string   `db:"role_specification" json:"roleSpecification,omitempty"`   // If Role is "Other", specify
	CreatedAt          time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt          time.Time `db:"updated_at" json:"updatedAt"`
}

type StudentSubjectPreference struct {
	ID          int       `db:"id" json:"id"`
	IIRID       int       `db:"iir_id" json:"iirId"`
	SubjectName string    `db:"subject_name" json:"subjectName"`
	IsFavorite  bool      `db:"is_favorite" json:"isFavorite"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

type StudentHobby struct {
	ID           int       `db:"id" json:"id"`
	IIRID        int       `db:"iir_id" json:"iirId"`
	HobbyName    string    `db:"hobby_name" json:"hobbyName"`
	PriorityRank int       `db:"priority_rank" json:"priorityRank"` // 1 = most favorite, higher numbers = less favorite
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}

// Test Results and Assessments
type TestResult struct {
	ID          int       `db:"id" json:"id"`
	IIRID       int       `db:"iir_id" json:"iirId"`
	TestDate    string    `db:"test_date" json:"testDate"`
	TestName    string    `db:"test_name" json:"testName"`
	RawScore    string    `db:"raw_score" json:"rawScore"`
	Percentile  string    `db:"percentile" json:"percentile"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

// Significant Notes and Incidents
type SignificantNote struct {
	ID                  int       `db:"id" json:"id"`
	IIRID               int       `db:"iir_id" json:"iirId"`
	NoteDate            string    `db:"note_date" json:"noteDate"`
	IncidentDescription string    `db:"incident_description" json:"incidentDescription"`
	Remarks             string    `db:"remarks" json:"remarks"`
	CreatedAt           time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt           time.Time `db:"updated_at" json:"updatedAt"`
}
