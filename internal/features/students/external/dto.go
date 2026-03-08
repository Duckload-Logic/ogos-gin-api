package external

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
)

type OGOSStudentDTO struct {
	StudentNumber string `json:"studentNumber"`

	FirstName  string                 `json:"firstName"`
	MiddleName structs.NullableString `json:"middleName,omitempty"`
	LastName   string                 `json:"lastName"`
	SuffixName string                 `json:"suffixName,omitempty"`

	Email         string `json:"email"`
	ContactNumber string `json:"contactNumber"`

	Course  students.Course `json:"course"`
	Year    int             `json:"year"`
	Section string          `json:"section"`
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
	StreetDetails string                 `json:"streetDetails"`
	Barangay      locations.Barangay     `json:"barangay"`
	City          locations.City         `json:"city"`
	Province      *locations.ProvinceDTO `json:"province,omitempty"`
	Region        locations.Region       `json:"region,omitempty"`
}
