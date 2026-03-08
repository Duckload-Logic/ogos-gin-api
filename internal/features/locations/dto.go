package locations

import "github.com/olazo-johnalbert/duckload-api/internal/core/structs"

type AddressDTO struct {
	ID           int       `json:"id"`
	StreetDetail string    `json:"streetDetail"`
	Region       Region    `json:"region"`
	Province     *Province `json:"province,omitempty"`
	City         City      `json:"city"`
	Barangay     Barangay  `json:"barangay"`
}

type ProvinceDTO struct {
	ID         int                    `json:"id,omitempty"`
	Code       structs.NullableString `json:"code,omitempty"`
	Name       structs.NullableString `json:"name,omitempty"`
	RegionCode structs.NullableString `json:"regionCode,omitempty"`
}
