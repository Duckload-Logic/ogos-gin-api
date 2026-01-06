package appointments

import (
	"fmt"
	"net/http"
	"strconv"
    "database/sql"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

type UpdateStatusRequest struct {
    Status string `json:"status" binding:"required"`
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
    var req CreateAppointmentRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    appt, err := h.service.CreateAppointment(c.Request.Context(), req)
    if err != nil {
        fmt.Println("Error creating an appointment:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create an appointment"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Appointment created successfully",
        "data":    appt,
    })
}

// ========================================
// |                                      |
// |      RETRIEVE HANDLER FUNCTIONS      |
// |                                      |
// ========================================

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
        // FIX: Check specifically for "no rows" and return 404
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
            return
        }
        
        // Print the real error to the console for debugging
        fmt.Println("Error retrieving appointment:", err)
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
    // Basic query params binding (optional)
    status := c.Query("status")
    date := c.Query("date")

    // Assuming your service has a List method that accepts simple filters
    appts, err := h.service.ListAppointments(c.Request.Context(), status, date)
    if err != nil {
        fmt.Println("Error listing appointments:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list appointments"})
        return
    }

    c.JSON(http.StatusOK, appts)
}

// HandleGetStudentAppointments godoc
// @Summary      Get Student's Appointments
// @Description  Retrieves all appointments for a specific student ID.
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Param        studentID  path      int  true  "Student ID"
// @Success      200        {array}   Appointment
// @Failure      400        {object}  map[string]string "Invalid Student ID"
// @Failure      500        {object}  map[string]string "Internal Server Error"
// @Router       /appointments/student/{studentID} [get]
func (h *Handler) HandleGetStudentAppointments(c *gin.Context) {
    studentID, err := strconv.Atoi(c.Param("studentID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
        return
    }

    appts, err := h.service.GetAppointmentsByStudentID(c.Request.Context(), studentID)
    if err != nil {
        fmt.Println("Error retrieving student appointments:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve student appointments"})
        return
    }

    c.JSON(http.StatusOK, appts)
}

func (h *Handler) HandleUpdateStatus(c *gin.Context) {
    // 1. Validate the ID first (don't ignore the error with _)
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
        return
    }

    // 2. Bind the JSON body
    var req UpdateStatusRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 3. Call the service and handle the specific "No Rows" error
    if err := h.service.UpdateAppointmentStatus(c.Request.Context(), id, req.Status); err != nil {
        if err == sql.ErrNoRows {
            // This is the fix: Return 404 instead of 500
            c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
            return
        }
        // Return 500 for actual server crashes
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Appointment status updated."})
}