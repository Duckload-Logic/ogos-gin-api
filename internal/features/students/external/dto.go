package external

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/request"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
)

type OGOSListStudentsRequest struct {
	request.PaginationParams
	Search    string `form:"search,omitempty"`
	CourseID  int    `form:"course_id,omitempty"`
	GenderID  int    `form:"gender_id,omitempty"`
	YearLevel int    `form:"year_level,omitempty"`
	OrderBy   string `form:"order_by,omitempty" binding:"omitempty,oneof=first_name last_name student_number created_at updated_at year_level course_id"`
}

type OGOSListStudentsResponse struct {
	Students   []OGOSStudentDTO `json:"students"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"pageSize"`
	TotalPages int              `json:"totalPages"`
}

type OGOSStudentDTO struct {
	StudentNumber string `json:"studentNumber"`

	FirstName  string                 `json:"firstName"`
	MiddleName structs.NullableString `json:"middleName,omitempty"`
	LastName   string                 `json:"lastName"`
	SuffixName string                 `json:"suffixName,omitempty"`

	Email        string `json:"email"`
	MobileNumber string `json:"mobileNumber"`

	Course    students.Course `json:"course"`
	YearLevel int             `json:"yearLevel"`
	Section   string          `json:"section"`
}

type OGOSStudentPersonalInfoDTO struct {
	StudentNumber string          `json:"studentNumber"`
	Gender        students.Gender `json:"gender"`
	DateOfBirth   string          `json:"dateOfBirth"`
	PlaceOfBirth  string          `json:"placeOfBirth"`
	HeightFt      float32         `json:"heightFt"`
	WeightKg      float32         `json:"weightKg"`
}

type OGOSStudentAddressDTO struct {
	StudentNumber string                 `json:"studentNumber"`
	AddressType   string                 `json:"addressType,omitempty"`
	StreetDetail  string                 `json:"streetDetail"`
	Barangay      locations.Barangay     `json:"barangay"`
	City          locations.City         `json:"city"`
	Province      *locations.ProvinceDTO `json:"province,omitempty"`
	Region        locations.Region       `json:"region,omitempty"`
}
