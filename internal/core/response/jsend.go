package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSendStatus string

const (
	StatusSuccess JSendStatus = "success"
	StatusFail    JSendStatus = "fail"
	StatusError   JSendStatus = "error"
)

// JSendResponse is the generic wrapper for all API responses following the
// JSend specification.
type JSendResponse struct {
	Status  JSendStatus `json:"status"            example:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty" example:"Detailed error message"`
	Code    int         `json:"code,omitempty"    example:"500"`
}

// SuccessResponse is a Swagger-friendly wrapper for a successful JSend response.
type SuccessResponse struct {
	Status JSendStatus `json:"status" example:"success"`
	Data   interface{} `json:"data"`
}

// FailResponse is a Swagger-friendly wrapper for a client-side JSend failure.
type FailResponse struct {
	Status JSendStatus `json:"status" example:"fail"`
	Data   interface{} `json:"data"`
}

// ErrorResponse is a Swagger-friendly wrapper for a server-side JSend error.
type ErrorResponse struct {
	Status  JSendStatus `json:"status"         example:"error"`
	Message string      `json:"message"        example:"An unexpected error occurred"`
	Code    int         `json:"code"           example:"500"`
	Data    interface{} `json:"data,omitempty"`
}

// CommonErrorResponse is a flat, unified struct for all 4xx and 5xx responses.
// This is used for Swagger documentation to provide a clean, non-composite schema.
type CommonErrorResponse struct {
	Status  string      `json:"status"            example:"fail"`
	Message string      `json:"message,omitempty" example:"Detailed error message"`
	Data    interface{} `json:"data,omitempty"`
	Code    int         `json:"code,omitempty"    example:"400"`
}

// SendSuccess sends a successful JSend response.
func SendSuccess(c *gin.Context, data interface{}, code ...int) {
	statusCode := http.StatusOK
	if len(code) > 0 {
		statusCode = code[0]
	}
	c.JSON(statusCode, JSendResponse{
		Status: StatusSuccess,
		Data:   data,
	})
}

// SendFail sends a JSend response for client-side errors (e.g., validation).
func SendFail(c *gin.Context, data interface{}, code ...int) {
	statusCode := http.StatusBadRequest
	if len(code) > 0 {
		statusCode = code[0]
	}
	c.JSON(statusCode, JSendResponse{
		Status: StatusFail,
		Data:   data,
	})
}

// SendError sends a JSend response for server-side errors.
func SendError(c *gin.Context, message string, code int, data interface{}) {
	c.JSON(code, JSendResponse{
		Status:  StatusError,
		Message: message,
		Code:    code,
		Data:    data,
	})
}
