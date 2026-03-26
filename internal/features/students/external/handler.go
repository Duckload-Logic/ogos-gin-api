package external

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
)

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

// HandleListStudents godoc
// @Summary List students
// @Description Get a paginated list of students with optional filters
// @Tags External Students
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param search query string false "Search term for student name or number"
// @Param course_id query int false "Filter by course ID"
// @Param gender_id query int false "Filter by gender ID"
// @Param year_level query int false "Filter by year level"
// @Param order_by query string false "Order by field (first_name, last_name,student_number, created_at, updated_at, year_level, course_id)"
// @Param page query int false "Page number for pagination"
// @Param page_size query int false "Number of items per page for pagination"
// @Success 200 {object} OGOSListStudentsResponse
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /students/external [get]
func (h *Handler) GetStudents(c *gin.Context) {
	var req OGOSListStudentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("[GetStudents] {Bind Query}: %v", err)
		response.SendFail(c, gin.H{"error": "invalid query parameters"})
		return
	}

	resp, err := h.service.ListStudents(c.Request.Context(), req)
	if err != nil {
		log.Printf("[GetStudents] {Service List}: %v", err)
		response.SendError(
			c,
			string(constants.ErrInternalServerError),
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, resp)
}

// GetStudentByUserID godoc
// @Summary Get student by user ID
// @Description Get student information by user ID
// @Tags External Students
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param userID path string true "User ID of the student"
// @Success 200 {object} OGOSStudentDTO
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /students/external/by-id/{userID} [get]
func (h *Handler) GetStudentByUserID(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		response.SendFail(c, gin.H{"error": "userID parameter is required"})
		return
	}

	student, err := h.service.GetStudentByUserID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("[GetStudentByUserID] {Service Get}: %v", err)
		response.SendError(
			c,
			string(constants.ErrInternalServerError),
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	if student == nil {
		response.SendFail(
			c,
			gin.H{"error": constants.ErrNotFound},
			http.StatusNotFound,
		)
		return
	}

	response.SendSuccess(c, student)
}

// GetPersonalInfoByStudentNumber godoc
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
func (h *Handler) GetPersonalInfoByStudentNumber(c *gin.Context) {
	studentNumber := c.Param("studentNumber")
	if studentNumber == "" {
		response.SendFail(
			c,
			gin.H{"error": "student_number parameter is required"},
		)
		return
	}

	student, err := h.service.GetPersonalInfoByStudentNumber(
		c.Request.Context(),
		studentNumber,
	)
	if err != nil {
		log.Printf("[GetPersonalInfoByStudentNumber] {Service Get}: %v", err)
		response.SendError(
			c,
			string(constants.ErrInternalServerError),
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	if student == nil {
		response.SendFail(
			c,
			gin.H{"error": constants.ErrNotFound},
			http.StatusNotFound,
		)
		return
	}

	response.SendSuccess(c, student)
}

// GetAddressByStudentNumber godoc
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
func (h *Handler) GetAddressByStudentNumber(c *gin.Context) {
	studentNumber := c.Param("studentNumber")
	if studentNumber == "" {
		response.SendFail(
			c,
			gin.H{"error": "student_number parameter is required"},
		)
		return
	}

	studentAddresses, err := h.service.GetAddressByStudentNumber(
		c.Request.Context(),
		studentNumber,
	)
	if err != nil {
		log.Printf("[GetAddressByStudentNumber] {Service Get}: %v", err)
		response.SendError(
			c,
			string(constants.ErrInternalServerError),
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	if len(studentAddresses) == 0 {
		response.SendFail(
			c,
			gin.H{"error": constants.ErrNotFound},
			http.StatusNotFound,
		)
		return
	}

	response.SendSuccess(c, studentAddresses)
}
