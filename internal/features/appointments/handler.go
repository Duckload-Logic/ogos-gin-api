package appointments

import (
	"database/sql"
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

// Create godoc
// @Summary      Book a new appointment
// @Description  Schedules a new appointment for a student.
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Param        request body      CreateAppointmentRequest true "Appointment Details"
// @Success      201     {object}  map[string]interface{}     "Returns {message: 'Success', data: Appointment}"
// @Failure      400     {object}  map[string]string          "Invalid input"
// @Failure      500     {object}  map[string]string          "Internal Server Error"
// @Router       /appointments [post]
func (h *Handler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appt, err := h.service.CreateAppointment(c.Request.Context(), userID.(int), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create an appointment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Appointment created successfully",
		"data":    appt,
	})
}

// HandleGetAppointment godoc
// @Summary      Get Appointment Details
// @Description  Retrieves the details of a specific appointment by its ID.
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Appointment ID"
// @Success      200  {object}  Appointment
// @Failure      400  {object}  map[string]string "Invalid ID format"
// @Failure      404  {object}  map[string]string "Appointment not found"
// @Failure      500  {object}  map[string]string "Internal Server Error"
// @Router       /appointments/{id} [get]
func (h *Handler) HandleGetAppointment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	appt, err := h.service.GetAppointmentByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve appointment"})
		return
	}

	c.JSON(http.StatusOK, appt)
}

// HandleListAppointments godoc
// @Summary      List All Appointments
// @Description  Retrieves a list of all appointments. Optionally filter by status or date via query params.
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Param        status  query     string  false  "Filter by status (e.g., Pending, Confirmed)"
// @Param        date    query     string  false  "Filter by date (YYYY-MM-DD)"
// @Success      200     {array}   Appointment
// @Failure      500     {object}  map[string]string "Internal Server Error"
// @Router       /appointments [get]
func (h *Handler) HandleListAppointments(c *gin.Context) {
	status := c.Query("status")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	appts, err := h.service.ListAppointments(c.Request.Context(), status, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list appointments"})
		return
	}

	c.JSON(http.StatusOK, appts)
}

func (h *Handler) HandleGetAvailableTimeSlots(c *gin.Context) {
	date := c.Query("date")
	slots, err := h.service.GetAvailableTimeSlots(c.Request.Context(), date)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve available time slots"})
		return
	}

	c.JSON(http.StatusOK, slots)
}

// GetAppointments gets all appointments for the current user
func (h *Handler) GetAppointments(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	appointments, err := h.service.GetAppointmentsByUserID(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": appointments})
}

// GetAppointmentByID gets a specific appointment
func (h *Handler) GetAppointmentByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	appointment, err := h.service.GetAppointmentByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Appointment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": appointment})
}

func (h *Handler) HandleUpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateAppointmentStatus(c.Request.Context(), id, req); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment updated successfully."})
}