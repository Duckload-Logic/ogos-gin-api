package locations

import "time"

type Region struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type City struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	RegionID int    `db:"region_id" json:"regionId"`
}

type Barangay struct {
	ID     int    `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	CityID int    `db:"city_id" json:"cityId"`
}

type Address struct {
	ID           int       `db:"id" json:"id"`
	RegionID     int       `db:"region_id" json:"regionId"`
	CityID       int       `db:"city_id" json:"cityId"`
	BarangayID   int       `db:"barangay_id" json:"barangayId"`
	StreetDetail *string   `db:"street_detail" json:"streetDetail"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}
