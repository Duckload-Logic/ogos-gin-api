package external

import "database/sql"

type OGOSStudentView struct {
	StudentNumber string `db:"student_number"`

	FirstName  string         `db:"first_name"`
	MiddleName sql.NullString `db:"middle_name,omitempty"`
	LastName   string         `db:"last_name"`

	Email         string `db:"email"`
	ContactNumber string `db:"contact_number"`

	CourseID   int    `db:"course_id"`
	CourseCode string `db:"course_code"`
	CourseName string `db:"course_name"`
	Year       int    `db:"year"`
	Section    string `db:"section"`
}

type OGOSStudentPersonalInfoView struct {
	StudentNumber string  `db:"student_number"`
	GenderID      int     `db:"gender_id"`
	GenderName    string  `db:"gender_name"`
	DateOfBirth   string  `db:"date_of_birth"`
	PlaceOfBirth  string  `db:"place_of_birth"`
	HeightFt      float32 `db:"height_ft"`
	WeightKg      float32 `db:"weight_kg"`
}

type OGOSStudentAddressView struct {
	StudentNumber string         `db:"student_number"`
	StreetDetails string         `db:"street_details"`
	BarangayCode  string         `db:"barangay"`
	BarangayName  string         `db:"barangay_name"`
	CityCode      string         `db:"city_code"`
	CityName      string         `db:"city_name"`
	ProvinceCode  sql.NullString `db:"province_code,omitempty"`
	ProvinceName  sql.NullString `db:"province_name,omitempty"`
	RegionCode    string         `db:"region_code,omitempty"`
	RegionName    string         `db:"region_name,omitempty"`
}
