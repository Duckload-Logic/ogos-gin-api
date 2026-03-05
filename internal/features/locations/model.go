package locations

import "time"

type Region struct {
	ID   int    `db:"id" json:"id"`
	Code string `db:"code" json:"code"`
	Name string `db:"name" json:"name"`
}

type Province struct {
	ID         int    `db:"id" json:"id"`
	Code       string `db:"code" json:"code"`
	Name       string `db:"name" json:"name"`
	RegionCode string `db:"region_code" json:"regionCode"`
}

type City struct {
	ID           int     `db:"id" json:"id"`
	Code         string  `db:"code" json:"code"`
	Name         string  `db:"name" json:"name"`
	ProvinceCode *string `db:"province_code" json:"provinceCode,omitempty"`
	Type         *string `db:"type" json:"type,omitempty"`
	ZipCode      *string `db:"zip_code" json:"zipCode,omitempty"`
	District     *string `db:"district" json:"district,omitempty"`
}

type Barangay struct {
	ID       int    `db:"id" json:"id"`
	Code     string `db:"code" json:"code"`
	Name     string `db:"name" json:"name"`
	CityCode int    `db:"city_code" json:"cityId"`
}

type Address struct {
	ID           int       `db:"id" json:"id"`
	RegionCode   string    `db:"region_code" json:"regionCode"`
	ProvinceCode string    `db:"province_code" json:"provinceCode"`
	CityCode     string    `db:"city_code" json:"cityCode"`
	BarangayCode string    `db:"barangay_code" json:"barangayCode"`
	StreetDetail *string   `db:"street_detail" json:"streetDetail"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}
