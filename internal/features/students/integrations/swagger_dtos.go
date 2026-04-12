package integrations

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
)

// ListStudentsSuccessResponse is a flat JSend response for HandleListStudents.
type ListStudentsSuccessResponse struct {
	Status response.JSendStatus     `json:"status" example:"success"`
	Data   OGOSListStudentsResponse `json:"data"`
}

// StudentSuccessResponse is a flat JSend response for GetStudentByStudentNumber.
type StudentSuccessResponse struct {
	Status string         `json:"status" example:"success"`
	Data   OGOSStudentDTO `json:"data"`
}

// StudentPersonalInfoSuccessResponse is a flat JSend response for GetPersonalInfoByStudentNumber.
type StudentPersonalInfoSuccessResponse struct {
	Status response.JSendStatus       `json:"status" example:"success"`
	Data   OGOSStudentPersonalInfoDTO `json:"data"`
}

// StudentAddressSuccessResponse is a flat JSend response for GetAddressByStudentNumber.
type StudentAddressSuccessResponse struct {
	Status response.JSendStatus    `json:"status" example:"success"`
	Data   []OGOSStudentAddressDTO `json:"data"`
}
