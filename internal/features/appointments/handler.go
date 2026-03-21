package appointments

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// getIIRIDFromContext extracts iirID from context or aborts
// with Forbidden status if not found.
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

// GetAppointmentCategoryList godoc
// @Summary      Get appointment categories
// @Description  Retrieves available appointment categories.
// @Tags         Appointments
// @Produce      json
// @Success      200  {object} []AppointmentCategory
// @Failure      500  {object} map[string]string
// @Router       /appointments/lookups/categories [get]
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
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to retrieve categories"},
		)
		return
	}

	c.JSON(http.StatusOK, categories)
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
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid query parameters"},
		)
		return
	}

	if req.StartDate == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "start_date parameter required",
			},
		)
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
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to retrieve statistics"},
		)
		return
	}

	c.JSON(http.StatusOK, dsc)
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
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid request format"},
		)
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
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to create appointment"},
		)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Appointment created successfully",
		"id":      appt.ID,
	})
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
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid ID format"},
		)
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
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to retrieve appointment"},
		)
		return
	}

	if appt == nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{"error": "Appointment not found"},
		)
		return
	}

	c.JSON(http.StatusOK, appt)
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
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid query parameters"},
		)
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
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to retrieve appointments"},
		)
		return
	}

	c.JSON(http.StatusOK, appts)
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
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "date parameter required"},
		)
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
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to retrieve time slots"},
		)
		return
	}

	c.JSON(http.StatusOK, slots)
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
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to retrieve statuses"},
		)
		return
	}

	c.JSON(http.StatusOK, statuses)
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
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid query parameters"},
		)
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
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to retrieve appointments"},
		)
		return
	}

	c.JSON(http.StatusOK, appointments)
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
			c.JSON(
				http.StatusForbidden,
				gin.H{
					"error": "Please complete your IIR profile",
				},
			)
			return
		}
		iirID, ok := iirIDVal.(string)
		if !ok {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "Internal server error"},
			)
			return
		}
		iirIDPtr = &iirID
	}

	var req ListAppointmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid query parameters"},
		)
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
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to retrieve statistics"},
		)
		return
	}

	c.JSON(http.StatusOK, stats)
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
func (h *Handler) PatchAppointment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid ID format"},
		)
		return
	}

	var req AppointmentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid request format"},
		)
		return
	}

	if err := h.service.UpdateAppointment(
		c.Request.Context(),
		id,
		req,
	); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(
				http.StatusNotFound,
				gin.H{"error": "Appointment not found"},
			)
			return
		}
		log.Printf(
			"[PatchAppointment] {Update Appointment}: %v",
			err,
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to update appointment"},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Appointment updated successfully",
	})
}
