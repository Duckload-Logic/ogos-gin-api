package appointments

import (
	"fmt"
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

// CreateAppointment godoc
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
		// terminal error message
		fmt.Println("Error creating appointment:", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create appointment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Appointment created successfully",
		"data":    appt,
	})
}
