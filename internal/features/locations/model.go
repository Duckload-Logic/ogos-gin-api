package locations

import "time"

type Region struct {
	ID   int    `db:"id" json:"id,omitempty"`
	Code string `db:"code" json:"code,omitempty"`
	Name string `db:"name" json:"name,omitempty"`
}

type Province struct {
	ID         int    `db:"id" json:"id,omitempty"`
	Code       string `db:"code" json:"code,omitempty"`
	Name       string `db:"name" json:"name,omitempty"`
	RegionCode string `db:"region_code" json:"regionCode,omitempty"`
}

type City struct {
	ID           int     `db:"id" json:"id,omitempty"`
	Code         string  `db:"code" json:"code,omitempty"`
	Name         string  `db:"name" json:"name,omitempty"`
	ProvinceCode *string `db:"province_code" json:"provinceCode,omitempty"`
	Type         *string `db:"type" json:"type,omitempty"`
	ZipCode      *string `db:"zip_code" json:"zipCode,omitempty"`
	District     *string `db:"district" json:"district,omitempty"`
	RegionCode   *string `db:"region_code" json:"regionCode,omitempty"`
}

type Barangay struct {
	ID       int    `db:"id" json:"id,omitempty"`
	Code     string `db:"code" json:"code,omitempty"`
	Name     string `db:"name" json:"name,omitempty"`
	CityCode int    `db:"city_code" json:"cityId,omitempty"`
}

type Address struct {
	ID           int       `db:"id" json:"id"`
	RegionCode   string    `db:"region_code" json:"regionCode"`
	ProvinceCode *string   `db:"province_code" json:"provinceCode,omitempty"`
	CityCode     string    `db:"city_code" json:"cityCode"`
	BarangayCode string    `db:"barangay_code" json:"barangayCode"`
	StreetDetail *string   `db:"street_detail" json:"streetDetail"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}
