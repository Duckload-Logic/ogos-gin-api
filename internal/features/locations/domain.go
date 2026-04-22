package locations

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type Region struct {
	ID   int
	Code string
	Name string
}

type Province struct {
	ID         int
	Code       string
	Name       string
	RegionCode string
}

type City struct {
	ID           int
	Code         string
	Name         string
	ProvinceCode structs.NullableString
	Type         structs.NullableString
	ZipCode      structs.NullableString
	District     structs.NullableString
	RegionCode   structs.NullableString
}

type Barangay struct {
	ID       int
	Code     string
	Name     string
	CityCode int
}

type Address struct {
	ID           int
	RegionCode   string
	ProvinceCode structs.NullableString
	CityCode     string
	BarangayCode string
	StreetDetail structs.NullableString
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
