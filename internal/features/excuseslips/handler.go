package excuseslips

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ========================================
// |                                      |
// |      RETRIEVE HANDLER FUNCTIONS      |
// |                                      |
// ========================================

// SubmitExcuseSlip godoc
// @Summary      Submit an excuse slip
// @Description  Allows a student to submit an excuse slip with a supporting document (file upload).
// @Tags         ExcuseSlips
// @Accept       multipart/form-data
// @Produce      json
// @Param        studentRecordId formData int    true "Student Record ID"
// @Param        reason          formData string true "Reason for absence"
// @Param        absenceDate     formData string true "Date of absence (YYYY-MM-DD)"
// @Param        file            formData file   true "Supporting Document (Image/PDF)"
// @Success      201             {object} map[string]interface{} "Returns {message, data}"
// @Failure      400             {object} map[string]string      "Invalid input or missing file"
// @Failure      500             {object} map[string]string      "Internal Server Error"
// @Router       /excuseslips [post]
func (h *Handler) Submit(c *gin.Context) {
	var req CreateExcuseSlipRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	slip, err := h.service.SubmitExcuseSlip(c.Request.Context(), req, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Excuse slip submitted successfully",
		"data":    slip,
	})
}
