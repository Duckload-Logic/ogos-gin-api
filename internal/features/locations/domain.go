package locations

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type Region struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type Province struct {
	ID         int    `json:"id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	RegionCode string `json:"regionCode"`
}

type City struct {
	ID           int                    `json:"id"`
	Code         string                 `json:"code"`
	Name         string                 `json:"name"`
	ProvinceCode structs.NullableString `json:"provinceCode"`
	Type         structs.NullableString `json:"type"`
	ZipCode      structs.NullableString `json:"zipCode"`
	District     structs.NullableString `json:"district"`
	RegionCode   structs.NullableString `json:"regionCode"`
}

type Barangay struct {
	ID       int    `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	CityCode int    `json:"cityCode"`
}

type Address struct {
	ID           int                    `json:"id"`
	RegionCode   string                 `json:"regionCode"`
	ProvinceCode structs.NullableString `json:"provinceCode"`
	CityCode     string                 `json:"cityCode"`
	BarangayCode string                 `json:"barangayCode"`
	StreetDetail structs.NullableString `json:"streetDetail"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}
