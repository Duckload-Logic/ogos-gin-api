package excuseslips

import (
	"net/http"
	"strconv"

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

// GetAllExcuseSlips godoc
// @Summary      Get all excuse slips
// @Description  Retrieves a list of all submitted excuse slips.
// @Tags         ExcuseSlips
// @Produce      json
// @Success      200  {object}  map[string]interface{} "Returns {data: []ExcuseSlip}"
// @Failure      500  {object}  map[string]string      "Internal Server Error"
// @Router       /excuseslips [get]
func (h *Handler) GetAll(c *gin.Context) {
    slips, err := h.service.GetAllExcuseSlips(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve excuse slips"})
        return
    }

	if slips == nil {
        slips = []*ExcuseSlip{}
    }

    c.JSON(http.StatusOK, gin.H{
        "data": slips,
    })
}

// GetExcuseSlipByID godoc
// @Summary      Get an excuse slip by ID
// @Description  Retrieves details of a specific excuse slip.
// @Tags         ExcuseSlips
// @Produce      json
// @Param        id   path      int  true  "Excuse Slip ID"
// @Success      200  {object}  map[string]interface{} "Returns {data: ExcuseSlip}"
// @Failure      400  {object}  map[string]string      "Invalid ID format"
// @Failure      404  {object}  map[string]string      "Excuse slip not found"
// @Failure      500  {object}  map[string]string      "Internal Server Error"
// @Router       /excuseslips/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    slip, err := h.service.GetExcuseSlipByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve excuse slip"})
        return
    }

    if slip == nil {
         c.JSON(http.StatusNotFound, gin.H{"error": "Excuse slip not found"})
         return
    }

    c.JSON(http.StatusOK, gin.H{
        "data": slip,
    })
}

type UpdateStatusRequest struct {
    Status string `json:"status" binding:"required"`
}

// UpdateStatus godoc
// @Summary      Update excuse slip status
// @Description  Approve or Reject an excuse slip.
// @Tags         ExcuseSlips
// @Accept       json
// @Produce      json
// @Param        id      path      int                  true  "Excuse Slip ID"
// @Param        body    body      UpdateStatusRequest  true  "New Status"
// @Success      200     {object}  map[string]string    "Status updated successfully"
// @Failure      400     {object}  map[string]string    "Invalid input"
// @Failure      404     {object}  map[string]string    "Record not found"
// @Failure      500     {object}  map[string]string    "Internal Server Error"
// @Router       /excuseslips/{id}/status [patch]
func (h *Handler) UpdateStatus(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    var req UpdateStatusRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err = h.service.UpdateExcuseSlipStatus(c.Request.Context(), id, req.Status)
    if err != nil {
        if err.Error() == "sql: no rows in result set" { 
             c.JSON(http.StatusNotFound, gin.H{"error": "Excuse slip not found"})
             return
        }
        if err.Error() == "invalid status: must be 'Pending', 'Approved', or 'Rejected'" {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

// DeleteExcuseSlip godoc
// @Summary      Delete an excuse slip
// @Description  Removes the excuse slip record and the uploaded file.
// @Tags         ExcuseSlips
// @Param        id   path      int  true  "Excuse Slip ID"
// @Success      200  {object}  map[string]string "Message"
// @Failure      404  {object}  map[string]string "Not Found"
// @Failure      500  {object}  map[string]string "Internal Error"
// @Router       /excuseslips/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    err = h.service.DeleteExcuseSlip(c.Request.Context(), id)
    if err != nil {
        if err.Error() == "excuse slip not found" || err.Error() == "sql: no rows in result set" {
            c.JSON(http.StatusNotFound, gin.H{"error": "Excuse slip not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete excuse slip"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Excuse slip deleted successfully"})
}