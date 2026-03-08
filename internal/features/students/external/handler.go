package external

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

func (h *Handler) GetStudentByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email query parameter is required"})
		return
	}

	student, err := h.service.GetStudentByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get student by email: %v", err)})
	}

	if student == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("student not found for email: %s", email)})
	}
	c.JSON(http.StatusOK, student)
}

func (h *Handler) GetPersonalInfoByStudentNumber(c *gin.Context) {
	studentNumber := c.Query("student_number")
	if studentNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_number query parameter is required"})
		return
	}

	student, err := h.service.GetPersonalInfoByStudentNumber(c.Request.Context(), studentNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get personal info by student number: %v", err)})
	}

	if student == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("personal info not found for student number: %s", studentNumber)})
	}

	c.JSON(http.StatusOK, student)
}

func (h *Handler) GetAddressByStudentNumber(c *gin.Context) {
	studentNumber := c.Query("student_number")
	if studentNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_number query parameter is required"})
		return
	}

	studentAddresses, err := h.service.GetAddressByStudentNumber(c.Request.Context(), studentNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get addresses by student number: %v", err)})
	}

	if len(studentAddresses) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("no addresses found for student number: %s", studentNumber)})
	}

	c.JSON(http.StatusOK, studentAddresses)
}
