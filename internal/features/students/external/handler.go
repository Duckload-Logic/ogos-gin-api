package external

import (
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

// HandleGetStudentByEmail godoc
// @Summary Get student by email
// @Description Get student information by email address
// @Tags External Students
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param email path string true "Email address of the student"
// @Success 200 {object} OGOSStudentDTO
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /students/external/by-email/{email} [get]
func (h *Handler) HandleGetStudentByEmail(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email parameter is required"})
		return
	}

	student, err := h.service.GetStudentByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrInternalServerError})
		return
	}

	if student == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrNotFound})
		return
	}

	c.JSON(http.StatusOK, student)
}

// HandleGetPersonalInfoByStudentNumber godoc
// @Summary Get personal information by student number
// @Description Get personal information of a student by their student number
// @Tags External Students
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param studentNumber path string true "Student number of the student"
// @Success 200 {object} OGOSStudentPersonalInfoDTO
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /students/external/personal-info/{studentNumber} [get]
func (h *Handler) HandleGetPersonalInfoByStudentNumber(c *gin.Context) {
	studentNumber := c.Param("studentNumber")
	if studentNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_number parameter is required"})
		return
	}

	student, err := h.service.GetPersonalInfoByStudentNumber(c.Request.Context(), studentNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrInternalServerError})
		return
	}

	if student == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrNotFound})
		return
	}

	c.JSON(http.StatusOK, student)
}

// HandleGetAddressByStudentNumber godoc
// @Summary Get student addresses by student number
// @Description Get all addresses of a student by their student number
// @Tags External Students
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param studentNumber path string true "Student number of the student"
// @Success 200 {array} OGOSStudentAddressDTO
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /students/external/addresses/{studentNumber} [get]
func (h *Handler) HandleGetAddressByStudentNumber(c *gin.Context) {
	studentNumber := c.Param("studentNumber")
	if studentNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_number parameter is required"})
		return
	}

	studentAddresses, err := h.service.GetAddressByStudentNumber(c.Request.Context(), studentNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": constants.ErrInternalServerError})
		return
	}

	if len(studentAddresses) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrNotFound})
		return
	}

	c.JSON(http.StatusOK, studentAddresses)
}
