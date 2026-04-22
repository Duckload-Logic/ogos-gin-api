package locations

import (
	"database/sql"
	"time"
)

type RegionDB struct {
	ID   int    `db:"id"   json:"id"`
	Code string `db:"code" json:"code"`
	Name string `db:"name" json:"name"`
}

type ProvinceDB struct {
	ID         int    `db:"id"          json:"id"`
	Code       string `db:"code"        json:"code"`
	Name       string `db:"name"        json:"name"`
	RegionCode string `db:"region_code" json:"regionCode"`
}

type CityDB struct {
	ID           int            `db:"id"            json:"id"`
	Code         string         `db:"code"          json:"code"`
	Name         string         `db:"name"          json:"name"`
	ProvinceCode sql.NullString `db:"province_code" json:"provinceCode"`
	Type         sql.NullString `db:"type"          json:"type"`
	ZipCode      sql.NullString `db:"zip_code"      json:"zipCode"`
	District     sql.NullString `db:"district"      json:"district"`
	RegionCode   sql.NullString `db:"region_code"   json:"regionCode"`
}

type BarangayDB struct {
	ID       int    `db:"id"        json:"id"`
	Code     string `db:"code"      json:"code"`
	Name     string `db:"name"      json:"name"`
	CityCode int    `db:"city_code" json:"cityCode"`
}

type AddressDB struct {
	ID           int            `db:"id"            json:"id"`
	RegionCode   string         `db:"region_code"   json:"regionCode"`
	ProvinceCode sql.NullString `db:"province_code" json:"provinceCode"`
	CityCode     string         `db:"city_code"     json:"cityCode"`
	BarangayCode string         `db:"barangay_code" json:"barangayCode"`
	StreetDetail sql.NullString `db:"street_detail" json:"streetDetail"`
	CreatedAt    time.Time      `db:"created_at"    json:"createdAt"`
	UpdatedAt    time.Time      `db:"updated_at"    json:"updatedAt"`
}
