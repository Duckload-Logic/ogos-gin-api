package appointments

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleGetConcernCategories(c *gin.Context) {
	categories, err := h.service.GetConcernCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve concern categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *Handler) HandleGetDailyStatusCount(c *gin.Context) {
	var req ListAppointmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get daily status count"})
		return
	}

	if req.StartDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request should have start_date parameter"})
		return
	}

	dsc, err := h.service.GetDailyStatusCount(c, req.StartDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get daily status count"})
		return
	}

	c.JSON(http.StatusOK, dsc)
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
func (h *Handler) HandleCreateAppointment(c *gin.Context) {
	userID := c.MustGet("userID")
	var req AppointmentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err.Error())
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
		"id":      appt.ID,
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
// @Param        start_date    query     string  false  "Filter by start date (YYYY-MM-DD)"
// @Param        end_date    query     string  false  "Filter by end date (YYYY-MM-DD)"
// @Success      200     {array}   Appointment
// @Failure      500     {object}  map[string]string "Internal Server Error"
// @Router       /appointments [get]
func (h *Handler) HandleListAppointments(c *gin.Context) {
	var req ListAppointmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appts, err := h.service.ListAppointments(c.Request.Context(), req)
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

func (h *Handler) HandleGetAppointmentStatuses(c *gin.Context) {
	statuses, err := h.service.GetAppointmentStatuses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve appointment statuses"})
		return
	}

	c.JSON(http.StatusOK, statuses)
}

// GetAppointments gets all appointments for the current user
func (h *Handler) HandleGetAppointmentsByUserID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req ListAppointmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appointments, err := h.service.GetAppointmentsByUserID(c.Request.Context(), userID.(int), req)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

func (h *Handler) HandleGetAppointmentStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	roleID, exists := c.Get("roleID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found"})
		return
	}

	var userIDPtr *int
	if roleID.(int) == int(constants.StudentRoleID) {
		id := userID.(int)
		userIDPtr = &id
	}

	var req ListAppointmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats, err := h.service.GetAppointmentStats(c.Request.Context(), req, userIDPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointment stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
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

	c.JSON(http.StatusOK, appointment)
}

func (h *Handler) HandleUpdateAppointment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid appointment ID: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	var req AppointmentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Status.ID == 0 {
		log.Printf("Status update detected for appointment ID %d\n", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status not provided"})
		return
	}

	if err := h.service.UpdateAppointmentStatus(c.Request.Context(), id, req); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
			return
		}

		log.Printf("Error updating appointment status: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment updated successfully."})
}

func (h *Handler) HandleUpdateAppointmentStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	var req AppointmentDTO
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
