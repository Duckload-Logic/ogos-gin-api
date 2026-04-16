package students

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
)

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetGenders(c *gin.Context) {
	genders, err := h.service.GetGenders(c.Request.Context())
	if err != nil {
		response.SendError(
			c,
			"Failed to get genders",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, genders)
}

func (h *Handler) GetParentalStatusTypes(c *gin.Context) {
	statuses, err := h.service.GetParentalStatusTypes(c.Request.Context())
	if err != nil {
		response.SendError(
			c,
			"Failed to get parental status types",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, statuses)
}

func (h *Handler) GetEnrollmentReasons(c *gin.Context) {
	reasons, err := h.service.GetEnrollmentReasons(c.Request.Context())
	if err != nil {
		response.SendError(
			c,
			"Failed to get enrollment reasons",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, reasons)
}

func (h *Handler) GetIncomeRanges(c *gin.Context) {
	ranges, err := h.service.GetIncomeRanges(c.Request.Context())
	if err != nil {
		response.SendError(
			c,
			"Failed to get income ranges",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, ranges)
}

func (h *Handler) GetStudentSupportTypes(c *gin.Context) {
	supportTypes, err := h.service.GetStudentSupportTypes(c.Request.Context())
	if err != nil {
		response.SendError(
			c,
			"Failed to get student support types",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, supportTypes)
}

func (h *Handler) GetSiblingSupportTypes(c *gin.Context) {
	supportTypes, err := h.service.GetSiblingSupportTypes(c.Request.Context())
	if err != nil {
		response.SendError(
			c,
			"Failed to get sibling support types",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, supportTypes)
}

func (h *Handler) GetEducationalLevels(c *gin.Context) {
	levels, err := h.service.GetEducationalLevels(c.Request.Context())
	if err != nil {
		response.SendError(
			c,
			"Failed to get educational levels",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, levels)
}

func (h *Handler) GetCourses(c *gin.Context) {
	courses, err := h.service.GetCourses(c.Request.Context())
	if err != nil {
		response.SendError(
			c,
			"Failed to get courses",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, courses)
}

func (h *Handler) GetCivilStatusTypes(c *gin.Context) {
	civilStatusTypes, err := h.service.GetCivilStatusTypes(c.Request.Context())
	if err != nil {
		response.SendError(
			c,
			"Failed to get civil status types",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, civilStatusTypes)
}

func (h *Handler) GetReligions(c *gin.Context) {
	religions, err := h.service.GetReligions(c.Request.Context())
	if err != nil {
		response.SendError(
			c,
			"Failed to get religions",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, religions)
}

func (h *Handler) GetStudentRelationshipTypes(c *gin.Context) {
	relationshipTypes, err := h.service.GetStudentRelationshipTypes(
		c.Request.Context(),
	)
	if err != nil {
		response.SendError(
			c,
			"Failed to get student relationship types",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, relationshipTypes)
}

func (h *Handler) GetNatureOfResidenceTypes(c *gin.Context) {
	types, err := h.service.GetNatureOfResidenceTypes(c.Request.Context())
	if err != nil {
		response.SendError(
			c,
			"Failed to get nature of residence types",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, types)
}

func (h *Handler) GetActivityOptions(c *gin.Context) {
	options, err := h.service.GetActivityOptions(c.Request.Context())
	if err != nil {
		response.SendError(
			c,
			"Failed to get activity options",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, options)
}

// GetStudentList godoc
// @Summary      List Students
// @Description  Retrieves a paginated list of students with optional filters.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        course     query     string  false  "Filter by course"
// @Param        year_level query     int     false  "Filter by year level"
// @Param		gender_id query     int     false  "Filter by gender ID"
// @Param        page       query     int     false  "Page number" default(1)
// @Param        page_size  query     int     false  "Number of records per page" default(20)
// @Success      200        {object}  ListStudentsResponse
// @Failure      500        {object}  map[string]string "Failed to list students"
// @Router       /students [get]
func (h *Handler) GetStudentList(c *gin.Context) {
	var req ListStudentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.ListStudents(c.Request.Context(), req)
	if err != nil {
		log.Printf("[GetStudentList] {Service Error}: %v", err)
		response.SendError(
			c,
			"Failed to list students",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, resp)
}

func (h *Handler) GetStudentProfile(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	profile, err := h.service.GetStudentProfile(c.Request.Context(), iirID)
	if err != nil {
		log.Printf("[GetStudentProfile] {Service Error}: %v", err)
		response.SendError(
			c,
			"Failed to get student profile",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, profile)
}

func (h *Handler) GetStudentBasicInfo(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	basicInfo, err := h.service.GetStudentBasicInfo(c.Request.Context(), iirID)
	if err != nil {
		response.SendError(
			c,
			"Failed to get student basic info",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, basicInfo)
}

func (h *Handler) GetIIRDraft(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	draft, err := h.service.GetIIRDraft(c.Request.Context(), userID)
	if err != nil {
		log.Printf("[GetIIRDraft] {Fetch Draft Error}: %v", err)
		response.SendError(
			c,
			"Failed to get IIR draft",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, draft)
}

func (h *Handler) GetStudentIIRByUserID(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		response.SendFail(c, gin.H{"error": "User ID is required"})
		return
	}

	iir, err := h.service.GetStudentIIRByUserID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("[GetStudentIIRByUserID] {Fetch IIR Error}: %v", err)
		response.SendError(
			c,
			"Failed to get student IIR by user ID",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, iir)
}

func (h *Handler) GetStudentIIRByIIRID(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	iir, err := h.service.GetStudentIIR(c.Request.Context(), iirID)
	if err != nil {
		response.SendError(
			c,
			"Failed to get student IIR by IIR ID",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, iir)
}

func (h *Handler) GetStudentEnrollmentReasons(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	reasons, err := h.service.GetStudentEnrollmentReasons(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		response.SendError(
			c,
			"Failed to get student enrollment reasons",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, reasons)
}

func (h *Handler) GetStudentPersonalInfo(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	personalInfo, err := h.service.GetStudentPersonalInfo(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		response.SendError(
			c,
			"Failed to get student personal info",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, personalInfo)
}

func (h *Handler) GetStudentAddresses(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	addresses, err := h.service.GetStudentAddresses(c.Request.Context(), iirID)
	if err != nil {
		response.SendError(
			c,
			"Failed to get student addresses",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, addresses)
}

func (h *Handler) GetStudentFamilyBackground(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	familyBackground, err := h.service.GetStudentFamilyBackground(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		response.SendError(
			c,
			"Failed to get student family background",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, familyBackground)
}

func (h *Handler) GetStudentRelatedPersons(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	relatedPersons, err := h.service.GetStudentRelatedPersons(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		response.SendError(
			c,
			"Failed to get student related persons",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, relatedPersons)
}

func (h *Handler) GetEducationalBackground(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	educationalBackground, err := h.service.GetEducationalBackground(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		response.SendError(
			c,
			"Failed to get educational background",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, educationalBackground)
}

func (h *Handler) GetStudentFinancialInfo(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	financialInfo, err := h.service.GetStudentFinancialInfo(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		response.SendError(
			c,
			"Failed to get student financial info",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, financialInfo)
}

func (h *Handler) GetStudentHealthRecord(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	healthRecord, err := h.service.GetStudentHealthRecord(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		response.SendError(
			c,
			"Failed to get student health record",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, healthRecord)
}

func (h *Handler) GetStudentConsultations(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	consultations, err := h.service.GetStudentConsultations(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		response.SendError(
			c,
			"Failed to get student consultations",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, consultations)
}

func (h *Handler) GetStudentActivities(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	activities, err := h.service.GetStudentActivities(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		response.SendError(
			c,
			"Failed to get student activities",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, activities)
}

func (h *Handler) GetStudentSubjectPreferences(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	subjectPreferences, err := h.service.GetStudentSubjectPreferences(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		response.SendError(
			c,
			"Failed to get student subject preferences",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, subjectPreferences)
}

func (h *Handler) GetStudentHobbies(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	hobbies, err := h.service.GetStudentHobbies(c.Request.Context(), iirID)
	if err != nil {
		response.SendError(
			c,
			"Failed to get student hobbies",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, hobbies)
}

func (h *Handler) GetStudentTestResults(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	testResults, err := h.service.GetStudentTestResults(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		response.SendError(
			c,
			"Failed to get student test results",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, testResults)
}

func (h *Handler) PostIIRDraft(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	var req ComprehensiveProfileDTO
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		log.Printf(
			"[PostIIRDraft] {JSON Decode}: %v",
			err,
		)
		response.SendFail(c, gin.H{"error": "Invalid JSON format"})
		return
	}

	draftID, err := h.service.SaveIIRDraft(c.Request.Context(), userID, req)
	if err != nil {
		log.Printf(
			"[PostIIRDraft] {Service Error}: %v",
			err,
		)
		response.SendError(
			c,
			"Failed to save IIR draft",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(
		c,
		gin.H{"id": draftID, "message": "IIR draft saved successfully"},
	)
}

func (h *Handler) PostIIR(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	var req ComprehensiveProfileDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[PostIIR] {JSON Decode}: %s", err.Error())
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	iirID, err := h.service.SubmitStudentIIR(c.Request.Context(), userID, req)
	if err != nil {
		log.Printf("[PostIIR] {Service Error}: %s", err.Error())
		response.SendError(
			c,
			"Failed to submit student IIR",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(
		c,
		gin.H{"id": iirID, "message": "Student IIR submitted successfully"},
	)
}

// GenerateIIR godoc
// @Summary      Generate Student IIR PDF
// @Description  Generates and downloads the student's Initial Interview Record
// as a PDF.
// @Tags         Students
// @Produce      application/pdf
// @Param        iirID   path      string  true  "IIR ID"
// @Success      200     {file}    binary
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /students/{iirID}/iir/download [get]
func (h *Handler) GenerateIIR(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	pdfBytes, fileName, err := h.service.GenerateIIR(c.Request.Context(), iirID)
	if err != nil {
		log.Printf("[GenerateIIR] {Service Error}: %v", err)
		response.SendError(
			c,
			"Failed to generate IIR PDF",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	c.Header(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename=\"%s\"", fileName),
	)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
