package slips

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/builders"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
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
	userEmail := c.MustGet("userEmail").(string)
	var req CreateSlipRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	// Get all files under the key "files"
	files := form.File["files"]

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one file is required"})
		return
	}

	slip, err := h.service.SubmitExcuseSlip(c.Request.Context(), userEmail, req, files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Excuse slip submitted successfully",
		"slipId":  slip.ID,
	})
}

func (h *Handler) HandleGetUrgentSlips(c *gin.Context) {
	var req ListSlipRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slips, err := h.service.GetUrgentSlips(c.Request.Context(), &req)
	if err != nil {
		log.Println("Error retrieving urgent excuse slips:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve urgent excuse slips"})
		return
	}

	c.JSON(http.StatusOK, slips)
}

func (h *Handler) GetSlipStats(c *gin.Context) {
	var req ListSlipRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userEmail := c.MustGet("userEmail")
	roleID := c.MustGet("roleID").(int)

	var userEmailPtr *string
	if roleID == int(constants.StudentRoleID) {
		email := userEmail.(string)
		userEmailPtr = &email
	}
	stats, err := h.service.GetSlipStats(c.Request.Context(), userEmailPtr, &req)
	if err != nil {
		log.Println("Error retrieving slip stats:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve slip statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *Handler) HandleGetSlipStatuses(c *gin.Context) {
	ctx := c.Request.Context()
	statuses, err := h.service.GetSlipStatuses(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve slip statuses"})
		return
	}

	c.JSON(http.StatusOK, statuses)
}

func (h *Handler) HandleGetSlipCategories(c *gin.Context) {
	ctx := c.Request.Context()
	categories, err := h.service.GetSlipCategories(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve slip categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetAllExcuseSlips godoc
// @Summary      Get all excuse slips
// @Description  Retrieves a list of all submitted excuse slips.
// @Tags         ExcuseSlips
// @Produce      json
// @Success      200  {object}  map[string]interface{} "Returns {data: []Slip}"
// @Failure      500  {object}  map[string]string      "Internal Server Error"
// @Router       /excuseslips [get]
func (h *Handler) GetAll(c *gin.Context) {
	var req ListSlipRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slips, err := h.service.GetAllExcuseSlips(c.Request.Context(), req)
	if err != nil {
		log.Println("Error retrieving excuse slips:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve excuse slips"})
		return
	}

	c.JSON(http.StatusOK, slips)
}

func (h *Handler) GetUserSlips(c *gin.Context) {
	userEmail := c.MustGet("userEmail").(string)
	var req ListSlipRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slips, err := h.service.GetExcuseSlipsByUserEmail(c.Request.Context(), userEmail, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve excuse slips"})
		return
	}

	c.JSON(http.StatusOK, slips)
}

func (h *Handler) GetSlipAttachments(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	attachments, err := h.service.GetSlipAttachments(c.Request.Context(), id)
	if err != nil {
		log.Println("Error retrieving attachments:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attachments"})
		return
	}

	c.JSON(http.StatusOK, attachments)
}

func (h *Handler) HandleDownloadAttachment(c *gin.Context) {
	attachmentIDParam := c.Param("attachmentId")
	attachmentID, err := strconv.Atoi(attachmentIDParam)
	if err != nil {
		log.Println("Error: Invalid attachment ID format:", attachmentIDParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attachment ID format"})
		return
	}

	file, err := h.service.GetAttachmentFile(c.Request.Context(), attachmentID)
	if err != nil {
		log.Println("Error retrieving attachment from DB (ID:", attachmentID, "):", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attachment metadata"})
		return
	}

	if file == nil {
		log.Println("Error: Attachment not found in DB (ID:", attachmentID, ")")
		c.JSON(http.StatusNotFound, gin.H{"error": "Attachment not found"})
		return
	}

	parts := strings.Split(strings.TrimPrefix(file.FileURL, "/slips/"), "/")
	if len(parts) != 2 {
		log.Println("Error: Invalid attachment path format:", file.FileURL)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attachment path format"})
		return
	}

	folderHash := parts[0]
	fileName := parts[1]

	realPath := builders.BuildFileURL("slips", folderHash, fileName)

	_, err = os.Stat(realPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Error: File not found at path:", realPath)
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found on server"})
			return
		}
		log.Println("Error: Cannot access file at path:", realPath, "-", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot access file"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.FileName))
	c.File(realPath)
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
// func (h *Handler) GetByID(c *gin.Context) {
// 	idParam := c.Param("id")
// 	id, err := strconv.Atoi(idParam)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
// 		return
// 	}

// 	slip, err := h.service.GetExcuseSlipByID(c.Request.Context(), id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve excuse slip"})
// 		return
// 	}

// 	if slip == nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Excuse slip not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data": slip,
// 	})
// }

// UpdateStatus godoc
// @Summary      Update excuse slip status
// @Description  Approve, Reject, or request revision of an excuse slip with optional admin notes.
// @Tags         ExcuseSlips
// @Accept       json
// @Produce      json
// @Param        id      path      int                  true  "Excuse Slip ID"
// @Param        body    body      UpdateStatusRequest  true  "New Status and optional admin notes"
// @Success      200     {object}  map[string]string    "Status updated successfully"
// @Failure      400     {object}  map[string]string    "Invalid input"
// @Failure      404     {object}  map[string]string    "Record not found"
// @Failure      500     {object}  map[string]string    "Internal Server Error"
// @Router       /slips/id/{id}/status [patch]
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

	err = h.service.UpdateExcuseSlipStatus(c.Request.Context(), id, req.Status, req.AdminNotes)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Excuse slip not found"})
			return
		}
		if strings.Contains(err.Error(), "invalid status") {
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
// func (h *Handler) Delete(c *gin.Context) {
// 	idParam := c.Param("id")
// 	id, err := strconv.Atoi(idParam)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
// 		return
// 	}

// 	err = h.service.DeleteExcuseSlip(c.Request.Context(), id)
// 	if err != nil {
// 		if err.Error() == "excuse slip not found" || err.Error() == "sql: no rows in result set" {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Excuse slip not found"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete excuse slip"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Excuse slip deleted successfully"})
// }
