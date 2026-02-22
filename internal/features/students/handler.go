package students

import (
	"fmt"
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

func (h *Handler) HandleGetGenders(c *gin.Context) {
	genders, err := h.service.GetGenders(c.Request.Context())
	if err != nil {
		fmt.Println("Error getting genders:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get genders"})
		return
	}

	c.JSON(http.StatusOK, genders)
}

func (h *Handler) HandleGetParentalStatusTypes(c *gin.Context) {
	statuses, err := h.service.GetParentalStatusTypes(c.Request.Context())
	if err != nil {
		fmt.Println("Error getting parental status types:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get parental status types"})
		return
	}

	c.JSON(http.StatusOK, statuses)
}

func (h *Handler) HandleGetEnrollmentReasons(c *gin.Context) {
	reasons, err := h.service.GetEnrollmentReasons(c.Request.Context())
	if err != nil {
		fmt.Println("Error getting enrollment reasons:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get enrollment reasons"})
		return
	}

	c.JSON(http.StatusOK, reasons)
}

func (h *Handler) HandleGetIncomeRanges(c *gin.Context) {
	ranges, err := h.service.GetIncomeRanges(c.Request.Context())
	if err != nil {
		fmt.Println("Error getting income ranges:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get income ranges"})
		return
	}

	c.JSON(http.StatusOK, ranges)
}

func (h *Handler) HandleGetStudentSupportTypes(c *gin.Context) {
	supportTypes, err := h.service.GetStudentSupportTypes(c.Request.Context())
	if err != nil {
		fmt.Println("Error getting student support types:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student support types"})
		return
	}

	c.JSON(http.StatusOK, supportTypes)
}

func (h *Handler) HandleGetSiblingSupportTypes(c *gin.Context) {
	supportTypes, err := h.service.GetSiblingSupportTypes(c.Request.Context())
	if err != nil {
		fmt.Println("Error getting sibling support types:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sibling support types"})
		return
	}

	c.JSON(http.StatusOK, supportTypes)
}

func (h *Handler) HandleGetEducationalLevels(c *gin.Context) {
	levels, err := h.service.GetEducationalLevels(c.Request.Context())
	if err != nil {
		fmt.Println("Error getting educational levels:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get educational levels"})
		return
	}

	c.JSON(http.StatusOK, levels)
}

func (h *Handler) HandleGetCourses(c *gin.Context) {
	courses, err := h.service.GetCourses(c.Request.Context())
	if err != nil {
		fmt.Println("Error getting courses:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get courses"})
		return
	}

	c.JSON(http.StatusOK, courses)
}

func (h *Handler) HandleGetCivilStatusTypes(c *gin.Context) {
	civilStatusTypes, err := h.service.GetCivilStatusTypes(c.Request.Context())
	if err != nil {
		fmt.Println("Error getting civil status types:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get civil status types"})
		return
	}

	c.JSON(http.StatusOK, civilStatusTypes)
}

func (h *Handler) HandleGetReligions(c *gin.Context) {
	religions, err := h.service.GetReligions(c.Request.Context())
	if err != nil {
		fmt.Println("Error getting religions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get religions"})
		return
	}

	c.JSON(http.StatusOK, religions)
}

func (h *Handler) HandleGetStudentRelationshipTypes(c *gin.Context) {
	relationshipTypes, err := h.service.GetStudentRelationshipTypes(c.Request.Context())
	if err != nil {
		fmt.Println("Error getting student relationship types:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student relationship types"})
		return
	}

	c.JSON(http.StatusOK, relationshipTypes)
}

func (h *Handler) HandleGetNatureOfResidenceTypes(c *gin.Context) {
	types, err := h.service.GetNatureOfResidenceTypes(c.Request.Context())
	if err != nil {
		fmt.Println("Error getting nature of residence types:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get nature of residence types"})
		return
	}

	c.JSON(http.StatusOK, types)
}

// HandleListStudents godoc
// @Summary      List Students
// @Description  Retrieves a paginated list of students with optional filters.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        course     query     string  false  "Filter by course"
// @Param        year_level query     int     false  "Filter by year level"
// @Param		gender_id query     int     false  "Filter by gender ID"
// @Param        page       query     int     false  "Page number"        default(1)
// @Param        page_size  query     int     false  "Number of records per page" default(20)
// @Success      200        {object}  ListStudentsResponse
// @Failure      500        {object}  map[string]string "Failed to list students"
// @Router       /students [get]
func (h *Handler) HandleListStudents(c *gin.Context) {
	var req ListStudentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.ListStudents(c.Request.Context(), req)
	if err != nil {
		fmt.Println("Error listing students:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to list students"},
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) HandleGetStudentProfile(c *gin.Context) {
	iirID, err := strconv.Atoi(c.Param("iirID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IIR ID"})
		return
	}

	profile, err := h.service.GetStudentProfile(c.Request.Context(), iirID)
	if err != nil {
		fmt.Println("Error getting student profile:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *Handler) HandleGetStudentBasicInfo(c *gin.Context) {
	iirID, err := strconv.Atoi(c.Param("iirID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IIR ID"})
		return
	}

	basicInfo, err := h.service.GetStudentBasicInfo(c.Request.Context(), iirID)
	if err != nil {
		fmt.Println("Error getting student basic info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student basic info"})
		return
	}

	c.JSON(http.StatusOK, basicInfo)
}

func (h *Handler) HandleGetStudentIIRByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	iir, err := h.service.GetStudentIIRByUserID(c.Request.Context(), userID)
	if err != nil {
		fmt.Println("Error getting student IIR by user ID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student IIR by user ID"})
		return
	}

	c.JSON(http.StatusOK, iir)
}

func (h *Handler) HandleGetStudentIIRByIIRID(c *gin.Context) {
	iirID, err := strconv.Atoi(c.Param("iirID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IIR ID"})
		return
	}

	iir, err := h.service.GetStudentIIR(c.Request.Context(), iirID)
	if err != nil {
		fmt.Println("Error getting student IIR by IIR ID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student IIR by IIR ID"})
		return
	}

	c.JSON(http.StatusOK, iir)
}

func (h *Handler) HandleGetStudentEnrollmentReasons(c *gin.Context) {
	iirID, err := strconv.Atoi(c.Param("iirID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IIR ID"})
		return
	}

	reasons, err := h.service.GetStudentEnrollmentReasons(c.Request.Context(), iirID)
	if err != nil {
		fmt.Println("Error getting student enrollment reasons:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student enrollment reasons"})
		return
	}

	c.JSON(http.StatusOK, reasons)
}

func (h *Handler) HandleGetStudentPersonalInfo(c *gin.Context) {
	iirID, err := strconv.Atoi(c.Param("iirID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IIR ID"})
		return
	}

	personalInfo, err := h.service.GetStudentPersonalInfo(c.Request.Context(), iirID)
	if err != nil {
		fmt.Println("Error getting student personal info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student personal info"})
		return
	}

	c.JSON(http.StatusOK, personalInfo)
}

func (h *Handler) HandleGetStudentAddresses(c *gin.Context) {
	iirID, err := strconv.Atoi(c.Param("iirID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IIR ID"})
		return
	}

	addresses, err := h.service.GetStudentAddresses(c.Request.Context(), iirID)
	if err != nil {
		fmt.Println("Error getting student addresses:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student addresses"})
		return
	}

	c.JSON(http.StatusOK, addresses)
}

func (h *Handler) HandleGetStudentFamilyBackground(c *gin.Context) {
	iirID, err := strconv.Atoi(c.Param("iirID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IIR ID"})
		return
	}

	familyBackground, err := h.service.GetStudentFamilyBackground(c.Request.Context(), iirID)
	if err != nil {
		fmt.Println("Error getting student family background:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student family background"})
		return
	}

	c.JSON(http.StatusOK, familyBackground)
}

func (h *Handler) HandleGetStudentRelatedPersons(c *gin.Context) {
	iirID, err := strconv.Atoi(c.Param("iirID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IIR ID"})
		return
	}

	relatedPersons, err := h.service.GetStudentRelatedPersons(c.Request.Context(), iirID)
	if err != nil {
		fmt.Println("Error getting student related persons:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student related persons"})
		return
	}

	c.JSON(http.StatusOK, relatedPersons)
}

func (h *Handler) HandleGetEducationalBackground(c *gin.Context) {
	iirID, err := strconv.Atoi(c.Param("iirID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IIR ID"})
		return
	}

	educationalBackground, err := h.service.GetEducationalBackground(c.Request.Context(), iirID)
	if err != nil {
		fmt.Println("Error getting educational background:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get educational background"})
		return
	}

	c.JSON(http.StatusOK, educationalBackground)
}

func (h *Handler) HandleGetStudentFinancialInfo(c *gin.Context) {
	iirID, err := strconv.Atoi(c.Param("iirID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IIR ID"})
		return
	}

	financialInfo, err := h.service.GetStudentFinancialInfo(c.Request.Context(), iirID)
	if err != nil {
		fmt.Println("Error getting student financial info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student financial info"})
		return
	}

	c.JSON(http.StatusOK, financialInfo)
}

func (h *Handler) HandleGetStudentHealthRecord(c *gin.Context) {
	iirID, err := strconv.Atoi(c.Param("iirID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IIR ID"})
		return
	}

	healthRecord, err := h.service.GetStudentHealthRecord(c.Request.Context(), iirID)
	if err != nil {
		fmt.Println("Error getting student health record:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student health record"})
		return
	}

	c.JSON(http.StatusOK, healthRecord)
}

func (h *Handler) HandleGetStudentConsultations(c *gin.Context) {
	iirID, err := strconv.Atoi(c.Param("iirID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IIR ID"})
		return
	}

	consultations, err := h.service.GetStudentConsultations(c.Request.Context(), iirID)
	if err != nil {
		fmt.Println("Error getting student consultations:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student consultations"})
		return
	}

	c.JSON(http.StatusOK, consultations)
}

func (h *Handler) HandleGetStudentActivities(c *gin.Context) {
	iirID, err := strconv.Atoi(c.Param("iirID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IIR ID"})
		return
	}

	activities, err := h.service.GetStudentActivities(c.Request.Context(), iirID)
	if err != nil {
		fmt.Println("Error getting student activities:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student activities"})
		return
	}

	c.JSON(http.StatusOK, activities)
}

func (h *Handler) HandleGetStudentSubjectPreferences(c *gin.Context) {
	iirID, err := strconv.Atoi(c.Param("iirID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IIR ID"})
		return
	}

	subjectPreferences, err := h.service.GetStudentSubjectPreferences(c.Request.Context(), iirID)
	if err != nil {
		fmt.Println("Error getting student subject preferences:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student subject preferences"})
		return
	}

	c.JSON(http.StatusOK, subjectPreferences)
}

func (h *Handler) HandleGetStudentHobbies(c *gin.Context) {
	iirID, err := strconv.Atoi(c.Param("iirID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IIR ID"})
		return
	}

	hobbies, err := h.service.GetStudentHobbies(c.Request.Context(), iirID)
	if err != nil {
		fmt.Println("Error getting student hobbies:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student hobbies"})
		return
	}

	c.JSON(http.StatusOK, hobbies)
}

func (h *Handler) HandleGetStudentTestResults(c *gin.Context) {
	iirID, err := strconv.Atoi(c.Param("iirID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IIR ID"})
		return
	}

	testResults, err := h.service.GetStudentTestResults(c.Request.Context(), iirID)
	if err != nil {
		fmt.Println("Error getting student test results:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student test results"})
		return
	}

	c.JSON(http.StatusOK, testResults)
}

func (h *Handler) HandleGetStudentSignificantNotes(c *gin.Context) {
	iirID, err := strconv.Atoi(c.Param("iirID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IIR ID"})
		return
	}

	significantNotes, err := h.service.GetStudentSignificantNotes(c.Request.Context(), iirID)
	if err != nil {
		fmt.Println("Error getting student significant notes:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student significant notes"})
		return
	}

	c.JSON(http.StatusOK, significantNotes)
}
