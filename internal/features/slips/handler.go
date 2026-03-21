package slips

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// getIIRIDFromContext extracts iirID from context or aborts with Forbidden status if not found.
func getIIRIDFromContext(c *gin.Context) (string, bool) {
	iirIDVal, exists := c.Get("iirID")
	if !exists {
		c.JSON(
			http.StatusForbidden,
			gin.H{
				"error": "Please complete your IIR profile",
			},
		)
		return "", false
	}

	iirID, ok := iirIDVal.(string)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Internal server error"},
		)
		return "", false
	}

	return iirID, true
}

// PostSlip godoc
// @Summary      Submit an excuse slip
// @Description  Allows a student to submit an excuse slip with
// @Description  supporting document (file upload).
// @Tags         ExcuseSlips
// @Accept       multipart/form-data
// @Produce      json
// @Param        reason      formData string true "Reason for absence"
// @Param        absenceDate formData string true "Date (YYYY-MM-DD)"
// @Param        files       formData file   true "Supporting Document"
// @Success      201         {object} map[string]interface{}
// @Failure      400         {object} map[string]string
// @Failure      403         {object} map[string]string
// @Failure      500         {object} map[string]string
// @Router       /slips [post]
func (h *Handler) PostSlip(c *gin.Context) {
	iirID, ok := getIIRIDFromContext(c)
	if !ok {
		return
	}

	var req CreateSlipRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid request format"},
		)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Printf(
			"[PostSlip] {Parse Multipart Form}: %v",
			err,
		)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Failed to parse form"},
		)
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "At least one file required"},
		)
		return
	}

	slip, err := h.service.SubmitExcuseSlip(
		c.Request.Context(),
		iirID,
		req,
		files,
	)
	if err != nil {
		log.Printf(
			"[PostSlip] {Submit Excuse Slip}: %v",
			err,
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to submit slip"},
		)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Excuse slip submitted successfully",
		"slipId":  slip.ID,
	})
}

// GetUrgentSlipList godoc
// @Summary      Get urgent excuse slips
// @Description  Retrieves urgent slips for counselor review.
// @Tags         ExcuseSlips
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]string
// @Router       /slips/urgent [get]
func (h *Handler) GetUrgentSlipList(c *gin.Context) {
	var req ListSlipRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid query parameters"},
		)
		return
	}

	slips, err := h.service.GetUrgentSlips(
		c.Request.Context(),
		&req,
	)
	if err != nil {
		log.Printf(
			"[GetUrgentSlipList] {Fetch Urgent Slips}: %v",
			err,
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to retrieve slips"},
		)
		return
	}

	c.JSON(http.StatusOK, slips)
}

// GetSlipStatsList godoc
// @Summary      Get slip statistics
// @Description  Retrieves slip status counts.
// @Tags         ExcuseSlips
// @Produce      json
// @Success      200  {object} []SlipStatusCount
// @Failure      500  {object} map[string]string
// @Router       /slips/stats [get]
func (h *Handler) GetSlipStatsList(c *gin.Context) {
	var req ListSlipRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid query parameters"},
		)
		return
	}

	roleID := c.MustGet("roleID").(int)
	var iirIDPtr *string

	if roleID == int(constants.StudentRoleID) {
		iirID, ok := getIIRIDFromContext(c)
		if !ok {
			return
		}
		iirIDPtr = &iirID
	}

	stats, err := h.service.GetSlipStats(
		c.Request.Context(),
		iirIDPtr,
		&req,
	)
	if err != nil {
		log.Printf(
			"[GetSlipStatsList] {Fetch Slip Stats}: %v",
			err,
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to retrieve statistics"},
		)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetSlipStatusList godoc
// @Summary      Get slip statuses
// @Description  Retrieves available slip statuses.
// @Tags         ExcuseSlips
// @Produce      json
// @Success      200  {object} []SlipStatus
// @Failure      500  {object} map[string]string
// @Router       /slips/lookups/statuses [get]
func (h *Handler) GetSlipStatusList(c *gin.Context) {
	statuses, err := h.service.GetSlipStatuses(
		c.Request.Context(),
	)
	if err != nil {
		log.Printf(
			"[GetSlipStatusList] {Fetch Statuses}: %v",
			err,
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to retrieve statuses"},
		)
		return
	}

	c.JSON(http.StatusOK, statuses)
}

// GetSlipCategoryList godoc
// @Summary      Get slip categories
// @Description  Retrieves available slip categories.
// @Tags         ExcuseSlips
// @Produce      json
// @Success      200  {object} []SlipCategory
// @Failure      500  {object} map[string]string
// @Router       /slips/lookups/categories [get]
func (h *Handler) GetSlipCategoryList(c *gin.Context) {
	categories, err := h.service.GetSlipCategories(
		c.Request.Context(),
	)
	if err != nil {
		log.Printf(
			"[GetSlipCategoryList] {Fetch Categories}: %v",
			err,
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to retrieve categories"},
		)
		return
	}

	c.JSON(http.StatusOK, categories)
}

// GetSlipList godoc
// @Summary      Get all excuse slips
// @Description  Retrieves list of all submitted slips.
// @Tags         ExcuseSlips
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]string
// @Router       /slips [get]
func (h *Handler) GetSlipList(c *gin.Context) {
	var req ListSlipRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid query parameters"},
		)
		return
	}

	slips, err := h.service.GetAllExcuseSlips(
		c.Request.Context(),
		req,
	)
	if err != nil {
		log.Printf(
			"[GetSlipList] {Fetch All Slips}: %v",
			err,
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to retrieve slips"},
		)
		return
	}

	c.JSON(http.StatusOK, slips)
}

// GetSlipListByIIR godoc
// @Summary      Get student's excuse slips
// @Description  Retrieves slips for the authenticated student.
// @Tags         ExcuseSlips
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      403  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /slips/me [get]
func (h *Handler) GetSlipListByIIR(c *gin.Context) {
	iirID, ok := getIIRIDFromContext(c)
	if !ok {
		return
	}

	var req ListSlipRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid query parameters"},
		)
		return
	}

	slips, err := h.service.GetExcuseSlipsByIIRID(
		c.Request.Context(),
		iirID,
		req,
	)
	if err != nil {
		log.Printf(
			"[GetSlipListByIIR] {Fetch Slips by IIR}: %v",
			err,
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to retrieve slips"},
		)
		return
	}

	c.JSON(http.StatusOK, slips)
}

// GetSlipAttachmentList godoc
// @Summary      Get slip attachments
// @Description  Retrieves attachments for a specific slip.
// @Tags         ExcuseSlips
// @Produce      json
// @Param        id   path      int  true  "Slip ID"
// @Success      200  {object} []AttachmentDTO
// @Failure      400  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /slips/id/{id}/attachments [get]
func (h *Handler) GetSlipAttachmentList(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid ID format"},
		)
		return
	}

	attachments, err := h.service.GetSlipAttachments(
		c.Request.Context(),
		id,
	)
	if err != nil {
		log.Printf(
			"[GetSlipAttachmentList] {Fetch Attachments}: %v",
			err,
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to retrieve attachments"},
		)
		return
	}

	c.JSON(http.StatusOK, attachments)
}

// GetAttachmentFile godoc
// @Summary      Download attachment
// @Description  Downloads a specific attachment file.
// @Tags         ExcuseSlips
// @Param        attachmentId path int true "Attachment ID"
// @Success      200  {file} binary
// @Failure      400  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /slips/id/{id}/attachments/{attachmentId} [get]
func (h *Handler) GetAttachmentFile(c *gin.Context) {
	attachmentIDParam := c.Param("attachmentId")
	attachmentID, err := strconv.Atoi(attachmentIDParam)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid attachment ID format"},
		)
		return
	}

	attachment, err := h.service.DownloadAttachment(
		c.Request.Context(),
		attachmentID,
		c.Writer,
	)
	if err != nil {
		if strings.Contains(
			err.Error(),
			"attachment not found",
		) {
			c.JSON(
				http.StatusNotFound,
				gin.H{"error": "Attachment not found"},
			)
			return
		}
		log.Printf(
			"[GetAttachmentFile] {Download Attachment}: %v",
			err,
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to download file"},
		)
		return
	}

	c.Header(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename=%q", attachment.FileName),
	)
}

// PatchSlipStatus godoc
// @Summary      Update slip status
// @Description  Approve, reject, or request revision of slip.
// @Tags         ExcuseSlips
// @Accept       json
// @Produce      json
// @Param        id   path      int                  true  "Slip ID"
// @Param        body body      UpdateStatusRequest  true  "Status and notes"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /slips/id/{id}/status [patch]
func (h *Handler) PatchSlipStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid ID format"},
		)
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid request format"},
		)
		return
	}

	err = h.service.UpdateExcuseSlipStatus(
		c.Request.Context(),
		id,
		req.Status,
		req.AdminNotes,
	)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(
				http.StatusNotFound,
				gin.H{"error": "Slip not found"},
			)
			return
		}
		if strings.Contains(
			err.Error(),
			"invalid status",
		) {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": err.Error()},
			)
			return
		}
		log.Printf(
			"[PatchSlipStatus] {Update Status}: %v",
			err,
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to update status"},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status updated successfully",
	})
}
