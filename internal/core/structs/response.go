package structs

type JSendStatus string

const (
	StatusSuccess JSendStatus = "success"
	StatusFail    JSendStatus = "fail"
	StatusError   JSendStatus = "error"
)

// JSendResponse is the generic wrapper for all API responses following the JSend specification.
type JSendResponse struct {
	Status  JSendStatus `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// NewSuccessResponse creates a successful JSend response.
func NewSuccessResponse(data interface{}) JSendResponse {
	return JSendResponse{
		Status: StatusSuccess,
		Data:   data,
	}
}

// NewFailResponse creates a JSend response for client-side errors (e.g., validation).
func NewFailResponse(data interface{}) JSendResponse {
	return JSendResponse{
		Status: StatusFail,
		Data:   data,
	}
}

// NewErrorResponse creates a JSend response for server-side errors.
func NewErrorResponse(
	message string,
	code int,
	data interface{},
) JSendResponse {
	return JSendResponse{
		Status:  StatusError,
		Message: message,
		Code:    code,
		Data:    data,
	}
}
