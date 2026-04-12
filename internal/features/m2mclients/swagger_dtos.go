package m2mclients

// M2MTokenSuccessResponse is a flat JSend response for PostM2MToken and PostM2MRefresh.
type M2MTokenSuccessResponse struct {
	Status string           `json:"status" example:"success"`
	Data   M2MTokenResponse `json:"data"`
}

// M2MClientListSuccessResponse is a flat JSend response for GetM2MClients.
type M2MClientListSuccessResponse struct {
	Status string         `json:"status" example:"success"`
	Data   []M2MClientDTO `json:"data"`
}

// M2MCreateClientSuccessResponse is a flat JSend response for PostM2MClient.
type M2MCreateClientSuccessResponse struct {
	Status string                  `json:"status" example:"success"`
	Data   CreateM2MClientResponse `json:"data"`
}

// M2MSecretSuccessResponse is a flat JSend response for PostM2MSecret.
type M2MSecretSuccessResponse struct {
	Status string `json:"status" example:"success"`
	Data   struct {
		ClientSecret string `json:"clientSecret" example:"your-new-client-secret"`
	} `json:"data"`
}

// M2MMessageSuccessResponse is a flat JSend response for generic messages.
type M2MMessageSuccessResponse struct {
	Status string `json:"status" example:"success"`
	Data   struct {
		Message string `json:"message" example:"Operation successful"`
	} `json:"data"`
}
