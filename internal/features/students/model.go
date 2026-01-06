package students

import (
	"database/sql"
)

// StudentRecord model
type StudentRecord struct {
	ID                  int             `db:"student_record_id" json:"id"`
	UserID              int             `db:"user_id" json:"userId"`
	GenderID            int             `db:"gender_id" json:"gender_id"`
	CivilStatusTypeID   int             `db:"civil_status_type_id" json:"civilStatusTypeId"`
	ReligionTypeID      int             `db:"religion_type_id" json:"religionTypeId"`
	HeightCm            sql.NullFloat64 `db:"height_cm" json:"heightCm"`
	WeightKg            sql.NullFloat64 `db:"weight_kg" json:"weightKg"`
	StudentNumber       string          `db:"student_number" json:"studentNumber"`
	Course              string          `db:"course" json:"course"`
	YearLevel           int             `db:"year_level" json:"yearLevel"`
	Section             sql.NullString  `db:"section" json:"section"`
	GoodMoralStatus     bool            `db:"good_moral_status" json:"goodMoralStatus"`
	HasDerogatoryRecord bool            `db:"has_derogatory_record" json:"hasDerogatoryRecord"`
	PlaceOfBirth        sql.NullString  `db:"place_of_birth" json:"place_of_birth"`
	BirthDate           sql.NullTime    `db:"birth_date" json:"birth_date"`
	MobileNo            sql.NullString  `db:"mobile_no" json:"mobile_no"`
}

// Guardian model
type Guardian struct {
	ID                 int            `db:"guardian_id" json:"id"`
	EducationalLevelID int            `db:"educational_level_id" json:"educationalLevelId"`
	BirthDate          sql.NullTime   `db:"birth_date" json:"birthDate"`
	LastName           string         `db:"last_name" json:"lastName"`
	FirstName          string         `db:"first_name" json:"firstName"`
	MiddleName         sql.NullString `db:"middle_name" json:"middleName"`
	Occupation         sql.NullString `db:"occupation" json:"occupation"`
	MaidenName         sql.NullString `db:"maiden_name" json:"maidenName"`
	CompanyName        sql.NullString `db:"company_name" json:"companyName"`
	RelationshipTypeID int            `db:"relationship_type_id" json:"relationshipTypeId"`
	IsPrimaryContact   bool           `db:"is_primary_contact" json:"isPrimaryContact"`
	ContactNumber      sql.NullString `db:"contact_number" json:"contactNumber"`
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
	ID                    int             `db:"family_background_id" json:"id"`
	StudentRecordID       int             `db:"student_record_id" json:"studentRecordId"`
	ParentalStatusID      int             `db:"parental_status_id" json:"parentalStatusId"`
	ParentalStatusDetails sql.NullString  `db:"parental_status_details" json:"parentalStatusDetails"`
	SiblingsBrothers      int             `db:"siblings_brothers" json:"siblingsBrothers"`
	SiblingSisters        int             `db:"sibling_sisters" json:"siblingSisters"`
	MonthlyFamilyIncome   sql.NullFloat64 `db:"monthly_family_income" json:"monthlyFamilyIncome"`
}

// EducationalBackground model
type EducationalBackground struct {
	ID                 int            `db:"educational_background_id" json:"id"`
	StudentRecordID    int            `db:"student_record_id" json:"studentRecordId"`
	EducationalLevelID int            `db:"educational_level_id" json:"educationalLevelId"`
	SchoolName         string         `db:"school_name" json:"schoolName"`
	Location           sql.NullString `db:"location" json:"location"`
	SchoolType         string         `db:"school_type" json:"schoolType"`
	YearCompleted      string         `db:"year_completed" json:"yearCompleted"`
	Awards             sql.NullString `db:"awards" json:"awards"`
}

// StudentAddress model
type StudentAddress struct {
	ID              int            `db:"student_address_id" json:"id"`
	StudentRecordID int            `db:"student_record_id" json:"studentRecordId"`
	AddressTypeID   int            `db:"address_type_id" json:"addressTypeId"`
	RegionName      sql.NullString `db:"region_name" json:"regionName"`
	ProvinceName    sql.NullString `db:"province_name" json:"provinceName"`
	CityName        sql.NullString `db:"city_name" json:"cityName"`
	BarangayName    sql.NullString `db:"barangay_name" json:"barangayName"`
	StreetLotBlk    sql.NullString `db:"street_lot_blk" json:"streetLotBlk"`
	UnitNo          sql.NullString `db:"unit_no" json:"unitNo"`
	BuildingName    sql.NullString `db:"building_name" json:"buildingName"`
}

type StudentHealthRecord struct {
	ID                    int            `db:"health_id" json:"id"`
	StudentRecordID       int            `db:"student_record_id" json:"studentRecordId"`
	VisionRemarkID        int            `db:"vision_remark_id" json:"visionRemarkId"`
	HearingRemarkID       int            `db:"hearing_remark_id" json:"hearingRemarkId"`
	MobilityRemarkID      int            `db:"mobility_remark_id" json:"mobilityRemarkId"`
	SpeechRemarkID        int            `db:"speech_remark_id" json:"speechRemarkId"`
	GeneralHealthRemarkID int            `db:"general_health_remark_id" json:"generalHealthId"`
	ConsultedProfessional sql.NullString `db:"consulted_professional" json:"consultedProfessional"`
	ConsultationReason    sql.NullString `db:"consultation_reason" json:"consultationReason"`
	DateStarted           sql.NullString `db:"date_started" json:"dateStarted"`
	NumberOfSessions      sql.NullInt64  `db:"num_sessions" json:"numberOfSessions"`
	DateConcluded         sql.NullString `db:"date_concluded" json:"dateConcluded"`
}
