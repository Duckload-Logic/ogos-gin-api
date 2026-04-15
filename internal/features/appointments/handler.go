package appointments

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type Handler struct {
	service ServiceInterface
}

// NewHandler creates a new appointments handler.
func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

// getIIRIDFromContext extracts iirID from context or aborts
// with Forbidden status if not found.
func getIIRIDFromContext(c *gin.Context) (string, bool) {
	iirIDVal, exists := c.Get("iirID")
	if !exists {
		response.SendFail(c, gin.H{
			"error": "Please complete your IIR profile",
		}, http.StatusForbidden)
		return "", false
	}

	iirID, ok := iirIDVal.(string)
	if !ok {
		response.SendError(
			c,
			"Internal server error",
			http.StatusInternalServerError,
			nil,
		)
		return "", false
	}

	return iirID, true
}

// GetAppointmentCategoryList godoc
// @Summary      Get appointment categories
// @Description  Retrieves available appointment categories.
// @Tags         Appointments
// @Produce      json
// @Success      200  {object} []AppointmentCategory
// @Failure      500  {object} map[string]string
// @Router       /appointments/lookups/categories [get]
// GetAppointmentCategoryList retrieves all appointment concern categories.
func (h *Handler) GetAppointmentCategoryList(
	c *gin.Context,
) {
	categories, err := h.service.GetConcernCategories(
		c.Request.Context(),
	)
	if err != nil {
		log.Printf(
			"[GetAppointmentCategoryList] {Fetch Categories}: %v",
			err,
		)
		response.SendError(
			c,
			"Failed to retrieve categories",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, categories)
}

// GetDailyStatusCountList godoc
// @Summary      Get daily status count
// @Description  Retrieves appointment status counts by date.
// @Tags         Appointments
// @Produce      json
// @Param        start_date query string true "Start date (YYYY-MM-DD)"
// @Success      200  {object} []DailyStatusCount
// @Failure      400  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /appointments/calendar/stats [get]
func (h *Handler) GetDailyStatusCountList(c *gin.Context) {
	var req ListAppointmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid query parameters"})
		return
	}

	if req.StartDate == "" {
		response.SendFail(c, gin.H{
			"error": "start_date parameter required",
		})
		return
	}

	dsc, err := h.service.GetDailyStatusCount(
		c,
		req.StartDate,
	)
	if err != nil {
		log.Printf(
			"[GetDailyStatusCountList] {Fetch Daily Stats}: %v",
			err,
		)
		response.SendError(
			c,
			"Failed to retrieve statistics",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, dsc)
}

// PostAppointment godoc
// @Summary      Book a new appointment
// @Description  Schedules a new appointment for a student.
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Param        request body      AppointmentDTO true "Appointment Details"
// @Success      201     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]string
// @Failure      403     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /appointments [post]
func (h *Handler) PostAppointment(c *gin.Context) {
	iirID, ok := getIIRIDFromContext(c)
	if !ok {
		return
	}

	var req AppointmentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid request format"})
		return
	}

	appt, err := h.service.CreateAppointment(
		c.Request.Context(),
		iirID,
		req,
	)
	if err != nil {
		log.Printf(
			"[PostAppointment] {Create Appointment}: %v",
			err,
		)
		response.SendError(
			c,
			"Failed to create appointment",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{
		"message": "Appointment created successfully",
		"id":      appt.ID,
	}, http.StatusCreated)
}

// GetAppointmentByID godoc
// @Summary      Get appointment details
// @Description  Retrieves the details of a specific appointment.
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Appointment ID"
// @Success      200  {object}  Appointment
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /appointments/id/{id} [get]
func (h *Handler) GetAppointmentByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendFail(c, gin.H{"error": "Invalid ID format"})
		return
	}

	appt, err := h.service.GetAppointmentByID(
		c.Request.Context(),
		id,
	)
	if err != nil {
		log.Printf(
			"[GetAppointmentByID] {Fetch Appointment}: %v",
			err,
		)
		response.SendError(
			c,
			"Failed to retrieve appointment",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	if appt == nil {
		response.SendFail(
			c,
			gin.H{"error": "Appointment not found"},
			http.StatusNotFound,
		)
		return
	}

	response.SendSuccess(c, appt)
}

// GetAppointmentList godoc
// @Summary      Get all appointments
// @Description  Retrieves a list of all appointments.
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Param        status     query string false "Filter by status"
// @Param        start_date query string false "Filter by start date"
// @Param        end_date   query string false "Filter by end date"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /appointments [get]
func (h *Handler) GetAppointmentList(c *gin.Context) {
	var req ListAppointmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid query parameters"})
		return
	}

	appts, err := h.service.ListAppointments(
		c.Request.Context(),
		req,
	)
	if err != nil {
		log.Printf(
			"[GetAppointmentList] {Fetch All Appointments}: %v",
			err,
		)
		response.SendError(
			c,
			"Failed to retrieve appointments",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, appts)
}

// GetAvailableTimeSlotList godoc
// @Summary      Get available time slots
// @Description  Retrieves available time slots for a date.
// @Tags         Appointments
// @Produce      json
// @Param        date query string true "Date (YYYY-MM-DD)"
// @Success      200  {object} []AvailableTimeSlotView
// @Failure      400  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /appointments/lookups/slots [get]
func (h *Handler) GetAvailableTimeSlotList(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		response.SendFail(c, gin.H{"error": "date parameter required"})
		return
	}

	slots, err := h.service.GetAvailableTimeSlots(
		c.Request.Context(),
		date,
	)
	if err != nil {
		log.Printf(
			"[GetAvailableTimeSlotList] {Fetch Slots}: %v",
			err,
		)
		response.SendError(
			c,
			"Failed to retrieve time slots",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, slots)
}

// GetAppointmentStatusList godoc
// @Summary      Get appointment statuses
// @Description  Retrieves available appointment statuses.
// @Tags         Appointments
// @Produce      json
// @Success      200  {object} []AppointmentStatus
// @Failure      500  {object} map[string]string
// @Router       /appointments/lookups/statuses [get]
func (h *Handler) GetAppointmentStatusList(c *gin.Context) {
	statuses, err := h.service.GetAppointmentStatuses(
		c.Request.Context(),
	)
	if err != nil {
		log.Printf(
			"[GetAppointmentStatusList] {Fetch Statuses}: %v",
			err,
		)
		response.SendError(
			c,
			"Failed to retrieve statuses",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, statuses)
}

// GetAppointmentListByIIR godoc
// @Summary      Get student's appointments
// @Description  Retrieves appointments for the authenticated student.
// @Tags         Appointments
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      403  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /appointments/me [get]
func (h *Handler) GetAppointmentListByIIR(c *gin.Context) {
	iirID, ok := getIIRIDFromContext(c)
	if !ok {
		return
	}

	var req ListAppointmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid query parameters"})
		return
	}

	appointments, err := h.service.GetAppointmentsByIIRID(
		c.Request.Context(),
		iirID,
		req,
	)
	if err != nil {
		log.Printf(
			"[GetAppointmentListByIIR] {Fetch Appointments}: %v",
			err,
		)
		response.SendError(
			c,
			"Failed to retrieve appointments",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, appointments)
}

// GetAppointmentStatsList godoc
// @Summary      Get appointment statistics
// @Description  Retrieves appointment status counts.
// @Tags         Appointments
// @Produce      json
// @Success      200  {object} []StatusCount
// @Failure      500  {object} map[string]string
// @Router       /appointments/stats [get]
func (h *Handler) GetAppointmentStatsList(c *gin.Context) {
	iirIDVal, exists := c.Get("iirID")
	roleID := c.MustGet("roleID").(int)

	var iirIDPtr *string
	if roleID == int(constants.StudentRoleID) {
		if !exists {
			response.SendFail(c, gin.H{
				"error": "Please complete your IIR profile",
			}, http.StatusForbidden)
			return
		}
		iirID, ok := iirIDVal.(string)
		if !ok {
			response.SendError(
				c,
				"Internal server error",
				http.StatusInternalServerError,
				nil,
			)
			return
		}
		iirIDPtr = &iirID
	}

	var req ListAppointmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid query parameters"})
		return
	}

	stats, err := h.service.GetAppointmentStats(
		c.Request.Context(),
		req,
		iirIDPtr,
	)
	if err != nil {
		log.Printf(
			"[GetAppointmentStatsList] {Fetch Stats}: %v",
			err,
		)
		response.SendError(
			c,
			"Failed to retrieve statistics",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, stats)
}

// PatchAppointment godoc
// @Summary      Update appointment
// @Description  Updates appointment details (reschedule).
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Param        id   path      int            true  "Appointment ID"
// @Param        body body      AppointmentDTO true  "Updated appointment"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /appointments/id/{id} [patch]
type CancelAppointmentRequest struct {
	Reason string `json:"reason"`
}

func (h *Handler) PostCancelAppointment(c *gin.Context) {
	id := c.Param("id")
	userID := audit.ExtractUserID(c.Request.Context())

	// Validate ownership
	ownerID, err := h.service.GetUserIDByAppointmentID(c.Request.Context(), id)
	if err != nil {
		response.SendError(c, "Failed to verify ownership", http.StatusInternalServerError, nil)
		return
	}
	if ownerID != userID {
		response.SendFail(c, gin.H{"error": "Access denied"}, http.StatusForbidden)
		return
	}

	// Fetch appointment to check status
	appt, err := h.service.GetAppointmentByID(c.Request.Context(), id)
	if err != nil {
		response.SendError(c, "Failed to fetch appointment", http.StatusInternalServerError, nil)
		return
	}

	// Only allow cancellation of Pending or Scheduled appointments
	statusName := strings.ToLower(appt.Status.Name)
	if statusName != "pending" && statusName != "scheduled" {
		response.SendFail(c, gin.H{"error": "Only pending or scheduled appointments can be cancelled"})
		return
	}

	// Fetch all statuses to find 'Cancelled'
	statuses, err := h.service.GetAppointmentStatuses(c.Request.Context())
	if err != nil {
		response.SendError(c, "Failed to fetch statuses", http.StatusInternalServerError, nil)
		return
	}

	var cancelStatusID int
	for _, s := range statuses {
		if strings.ToLower(s.Name) == "cancelled" {
			cancelStatusID = s.ID
			break
		}
	}

	if cancelStatusID == 0 {
		response.SendError(c, "Cancelled status not found", http.StatusInternalServerError, nil)
		return
	}

	var req CancelAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// We allow empty body for backward compatibility or simple cancellations
		// but if body is present, it should be valid JSON
		if err.Error() != "EOF" {
			log.Printf("[PostCancelAppointment] {BindJSON}: %v", err)
		}
	}

	// Update appointment status
	updateReq := *appt
	updateReq.Status.ID = cancelStatusID

	if req.Reason != "" {
		if updateReq.AdminNotes.Valid && updateReq.AdminNotes.String != "" {
			updateReq.AdminNotes.String += "\nStudent Cancellation: " + req.Reason
		} else {
			updateReq.AdminNotes = structs.StringToNullableString("Student Cancellation: " + req.Reason)
		}
	}

	if err := h.service.UpdateAppointment(c.Request.Context(), id, updateReq); err != nil {
		log.Printf("[PostCancelAppointment] {Update}: %v", err)
		response.SendError(c, "Failed to cancel appointment", http.StatusInternalServerError, nil)
		return
	}

	response.SendSuccess(c, gin.H{"message": "Appointment cancelled successfully"})
}

// PatchAppointment godoc
func (h *Handler) PatchAppointment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendFail(c, gin.H{"error": "Invalid ID format"})
		return
	}

	var req AppointmentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.service.UpdateAppointment(
		c.Request.Context(),
		id,
		req,
	); err != nil {
		if err == sql.ErrNoRows {
			response.SendFail(
				c,
				gin.H{"error": "Appointment not found"},
				http.StatusNotFound,
			)
			return
		}
		log.Printf(
			"[PatchAppointment] {Update Appointment}: %v",
			err,
		)
		response.SendError(
			c,
			"Failed to update appointment",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{
		"message": "Appointment updated successfully",
	})
}
