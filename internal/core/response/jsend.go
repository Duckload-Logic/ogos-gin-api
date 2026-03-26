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

// JSendResponse is the generic wrapper for all API responses following the JSend specification.
type JSendResponse struct {
	Status  JSendStatus `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
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
