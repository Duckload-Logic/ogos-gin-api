package locations

type AddressDTO struct {
	ID           int       `json:"id"`
	StreetDetail string    `json:"streetDetail"`
	Region       Region    `json:"region"`
	Province     *Province `json:"province,omitempty"`
	City         City      `json:"city"`
	Barangay     Barangay  `json:"barangay"`
}
