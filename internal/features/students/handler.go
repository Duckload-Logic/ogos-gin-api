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

// ========================================
// |                                      |
// |      RETRIEVE HANDLER FUNCTIONS      |
// |                                      |
// ========================================

// HandleListStudents godoc
// @Summary      List Students
// @Description  Retrieves a paginated list of students with optional filters.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        course     query     string  false  "Filter by course"
// @Param        year_level query     int     false  "Filter by year level"
// @Param        gender_id  query     int     false  "Filter by gender ID"
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

// HandleGetStudentRecordID godoc
// @Summary      Get Student Record ID
// @Description  Retrieves the student record ID associated with a given user ID.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        userID   path      int  true  "User ID"
// @Success      200      {object}  map[string]int "Returns {studentRecordID: <id>}"
// @Failure      400      {object}  map[string]string "Invalid user ID"
// @Failure      500      {object}  map[string]string "Failed to get student record ID"
// @Router       /students/record-id/{userID} [get]
func (h *Handler) HandleGetStudentRecordID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	studentRecordID, err := h.service.GetStudentRecordByStudentID(
		c.Request.Context(), userID,
	)
	if err != nil {
		fmt.Println("Error getting student record ID:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get student record ID"},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"studentRecordID": studentRecordID})
}

// HandleGetStudentEnrollmentReasons godoc
// @Summary      Get Student Enrollment Reasons
// @Description  Retrieves the enrollment reasons selected by a specific student.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        studentRecordID path      int  true  "Student Record ID"
// @Success      200             {array}   StudentSelectedReason
// @Failure      400             {object}  map[string]string "Invalid student record ID"
// @Failure      500             {object}  map[string]string "Failed to get enrollment reasons"
// @Router       /students/profile/enrollment-reasons/{studentRecordID} [get]
func (h *Handler) HandleGetStudentEnrollmentReasons(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid student record ID"},
		)
		return
	}

	resp, err := h.service.GetStudentEnrollmentReasons(c.Request.Context(), studentRecordID)
	if err != nil {
		fmt.Println("Error getting student enrollment reasons:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get student enrollment reasons"},
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleGetStudent godoc
// @Summary      Get Student Profile
// @Description  Retrieves the complete profile information for a student using their User ID.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        userID   path      int  true  "User ID"
// @Param        include_family     query     bool  false  "Include family information"
// @Param        include_health     query     bool  false  "Include health information"
// @Param        include_education  query     bool  false  "Include educational background"
// @Param        include_address    query     bool  false  "Include address information"
// @Param        include_finance    query     bool  false  "Include financial information"
// @Success      200      {object}  StudentProfileResponse
// @Failure      400      {object}  map[string]string "Invalid user ID"
// @Failure      500      {object}  map[string]string "Failed to get student profile"
// @Router       /students/record/{userID} [get]
func (h *Handler) HandleGetStudent(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student record ID"})
		return
	}

	// Bind query parameters
	var includeReq GetStudentRequest
	if err := c.ShouldBindQuery(&includeReq); err != nil {
		// If no query params, that's OK - we'll just return base profile
		fmt.Println("No include params or error:", err)
	}

	studentRecordID, err := h.service.GetStudentRecordByStudentID(
		c.Request.Context(), userID,
	)
	if err != nil {
		fmt.Println("Error getting student record ID:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get student record ID"},
		)
		return
	}

	// Get base profile
	profile, err := h.service.GetBaseProfile(c.Request.Context(), studentRecordID)
	if err != nil {
		fmt.Println("Error getting base profile:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get base profile"},
		)
		return
	}

	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	// Create response
	response := StudentProfileResponse{
		StudentProfile: profile,
	}

	// If include params are requested, fetch additional data
	if includeReq.IncludeFamily || includeReq.IncludeHealth ||
		includeReq.IncludeEducation || includeReq.IncludeAddress ||
		includeReq.IncludeFinance {

		compResponse := ComprehensiveProfileResponse{
			StudentProfile: profile,
		}

		if includeReq.IncludeFamily {
			family, err := h.service.GetFamilyInfo(
				c.Request.Context(), studentRecordID,
			)
			if err == nil && family != nil {
				compResponse.Family = family
			}
		}

		if includeReq.IncludeGuardians {
			guardians, err := h.service.GetGuardiansInfo(
				c.Request.Context(), studentRecordID,
			)

			if err == nil && guardians != nil {
				compResponse.Guardians = guardians
			}
		}

		if includeReq.IncludeFinance {
			finance, err := h.service.GetFinanceInfo(
				c.Request.Context(), studentRecordID,
			)
			if err == nil && finance != nil {
				compResponse.Finance = finance
			}
		}

		if includeReq.IncludeEducation {
			education, err := h.service.GetEducationInfo(
				c.Request.Context(), studentRecordID,
			)
			if err == nil && education != nil {
				compResponse.Education = education
			}
		}

		if includeReq.IncludeAddress {
			address, err := h.service.GetAddressInfo(
				c.Request.Context(), studentRecordID,
			)
			if err == nil && address != nil {
				compResponse.Address = address
			}
		}

		if includeReq.IncludeHealth {
			health, err := h.service.GetHealthInfo(
				c.Request.Context(), studentRecordID,
			)
			if err == nil && health != nil {
				compResponse.Health = health
			}
		}

		c.JSON(http.StatusOK, compResponse)
		return
	}

	// Return only base profile if no include params
	c.JSON(http.StatusOK, response)
}

// HandleGetBaseProfile godoc
// @Summary      Get Student Base Profile
// @Description  Retrieves the basic profile information for a student using their User ID.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        userID   path      int  true  "User ID"
// @Success      200      {object}  StudentProfileResponse
// @Failure      400      {object}  map[string]string "Invalid user ID"
// @Failure      500      {object}  map[string]string "Failed to get base profile"
// @Router       /students/profile/base/{userID} [get]
func (h *Handler) HandleGetBaseProfile(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	profile, err := h.service.GetBaseProfile(c.Request.Context(), studentRecordID)
	if err != nil {
		fmt.Println("Error getting base profile:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get base profile"},
		)
		return
	}

	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student profile not found"})
		return
	}

	response := StudentProfileResponse{
		StudentProfile: profile,
	}

	c.JSON(http.StatusOK, response)
}

// HandleGetFamilyInfo godoc
// @Summary      Get Family Information
// @Description  Retrieves family details associated with a specific student record.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        studentRecordID path      int  true  "Student Record ID"
// @Success      200             {object}  FamilyBackground
// @Failure      400             {object}  map[string]string "Invalid student record ID"
// @Failure      500             {object}  map[string]string "Failed to get family info"
// @Router       /students/profile/family/{studentRecordID} [get]
func (h *Handler) HandleGetFamilyInfo(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid student record ID"},
		)
		return
	}

	resp, err := h.service.GetFamilyInfo(c.Request.Context(), studentRecordID)
	if err != nil {
		fmt.Println("Error getting family info:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get family info"},
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleGetGuardiansInfo godoc
// @Summary      Get Guardians Information
// @Description  Retrieves a list of guardians linked to the student record.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        studentRecordID path      int  true  "Student Record ID"
// @Success      200             {array}   Guardian
// @Failure      400             {object}  map[string]string "Invalid student record ID"
// @Failure      500             {object}  map[string]string "Failed to get guardians info"
// @Router       /students/profile/guardians/{studentRecordID} [get]
func (h *Handler) HandleGetGuardiansInfo(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid student record ID"},
		)
		return
	}

	resp, err := h.service.GetGuardiansInfo(
		c.Request.Context(),
		studentRecordID,
	)
	if err != nil {
		fmt.Println("Error getting guardians info:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get guardians info"},
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleGetPrimaryGuardianInfo godoc
// @Summary      Get Primary Guardian Information
// @Description  Retrieves details of the primary guardian for the student.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        studentRecordID path      int  true  "Student Record ID"
// @Success      200             {object}  Guardian
// @Failure      400             {object}  map[string]string "Invalid student record ID"
// @Failure      500             {object}  map[string]string "Failed to get primary guardian info"
// @Router       /students/profile/guardians/primary/{studentRecordID} [get]
func (h *Handler) HandleGetPrimaryGuardianInfo(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid student record ID"},
		)
		return
	}

	resp, err := h.service.GetPrimaryGuardianInfo(
		c.Request.Context(), studentRecordID,
	)
	if err != nil {
		fmt.Println("Error getting primary guardian info:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get primary guardian info"},
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleGetEducationInfo godoc
// @Summary      Get Education Information
// @Description  Retrieves the student's educational background and history.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        studentRecordID path      int  true  "Student Record ID"
// @Success      200             {array}   EducationalBackground
// @Failure      400             {object}  map[string]string "Invalid student record ID"
// @Failure      500             {object}  map[string]string "Failed to get education info"
// @Router       /students/profile/education/{studentRecordID} [get]
func (h *Handler) HandleGetEducationInfo(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid student record ID"},
		)
		return
	}

	resp, err := h.service.GetEducationInfo(
		c.Request.Context(), studentRecordID,
	)
	if err != nil {
		fmt.Println("Error getting education info:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get education info"},
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleGetAddressInfo godoc
// @Summary      Get Address Information
// @Description  Retrieves the residential address details for the student.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        studentRecordID path      int  true  "Student Record ID"
// @Success      200             {array}   StudentAddress
// @Failure      400             {object}  map[string]string "Invalid student record ID"
// @Failure      500             {object}  map[string]string "Failed to get address info"
// @Router       /students/profile/address/{studentRecordID} [get]
func (h *Handler) HandleGetAddressInfo(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid student record ID"},
		)
		return
	}

	resp, err := h.service.GetAddressInfo(c.Request.Context(), studentRecordID)
	if err != nil {
		fmt.Println("Error getting address info:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get address info"},
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleGetHealthInfo godoc
// @Summary      Get Health Information
// @Description  Retrieves the health and medical records for the student.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        studentRecordID path      int  true  "Student Record ID"
// @Success      200             {object}  StudentHealthRecord
// @Failure      400             {object}  map[string]string "Invalid student record ID"
// @Failure      500             {object}  map[string]string "Failed to get health info"
// @Router       /students/profile/health/{studentRecordID} [get]
func (h *Handler) HandleGetHealthInfo(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid student record ID"},
		)
		return
	}

	resp, err := h.service.GetHealthInfo(c.Request.Context(), studentRecordID)
	if err != nil {
		fmt.Println("Error getting health info:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get health info"},
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ========================================
// |                                      |
// |      FINANCE RETRIEVE HANDLERS       |
// |                                      |
// ========================================

// HandleGetFinanceInfo godoc
// @Summary      Get Finance Information
// @Description  Retrieves the financial details for a student.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        studentRecordID path      int  true  "Student Record ID"
// @Success      200             {object}  StudentFinance
// @Failure      400             {object}  map[string]string "Invalid student record ID"
// @Failure      500             {object}  map[string]string "Failed to get finance info"
// @Router       /students/profile/finance/{studentRecordID} [get]
func (h *Handler) HandleGetFinanceInfo(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid student record ID"},
		)
		return
	}

	resp, err := h.service.GetFinanceInfo(c.Request.Context(), studentRecordID)
	if err != nil {
		fmt.Println("Error getting finance info:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get finance info"},
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ========================================
// |                                      |
// |       UPSERT HANDLER FUNCTIONS       |
// |                                      |
// ========================================

// HandleCreateStudentRecord godoc
// @Summary      Create Student Record (First Step)
// @Description  Creates a student record for a user (first step in onboarding)
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        userID  path      int  true  "User ID"
// @Success      200     {object}  map[string]interface{}  "Returns {student_record_id: <id>}"
// @Failure      400     {object}  map[string]string       "Invalid user ID"
// @Failure      500     {object}  map[string]string       "Internal Server Error"
// @Router       /students/onboarding/{userID} [post]
func (h *Handler) HandleCreateStudentRecord(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	studentRecordID, err := h.service.CreateStudentRecord(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"student_record_id": studentRecordID,
		"message":           "Student record created successfully",
	})
}

// HandleSaveEnrollmentReasons godoc
// @Summary      Save Enrollment Reasons
// @Description  Saves or updates the enrollment reasons for a student.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        studentRecordID path      int                             true "Student Record ID"
// @Param        request         body      UpdateEnrollmentReasonsRequest  true "Enrollment Reasons Data"
// @Success      200             {object}  map[string]string               "Message: Enrollment reasons saved successfully"
// @Failure      400             {object}  map[string]string               "Invalid input or ID"
// @Failure      500             {object}  map[string]string               "Internal Server Error"
// @Router       /students/onboarding/enrollment-reasons/{studentRecordID} [put]
func (h *Handler) HandleSaveEnrollmentReasons(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid student record ID"},
		)
		return
	}

	var req UpdateEnrollmentReasonsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.SaveEnrollmentReasons(c.Request.Context(), studentRecordID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Enrollment reasons saved successfully"})
}

// HandleSaveBaseProfile godoc
// @Summary      Create or Update Base Profile
// @Description  Creates a new student record or updates the basic profile information.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        userID   path      int                           true  "User ID"
// @Param        request  body      CreateStudentRecordRequest    true  "Base Profile Data"
// @Success      200      {object}  map[string]interface{}        "Returns {student_record_id: <id>}"
// @Failure      400      {object}  map[string]string             "Invalid input"
// @Failure      500      {object}  map[string]string             "Internal Server Error"
// @Router       /students/onboarding/base/{userID} [post]
func (h *Handler) HandleSaveBaseProfile(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student record ID"})
		return
	}

	var req CreateStudentRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.SaveBaseProfile(c.Request.Context(), studentRecordID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Base profile saved successfully"})
}

// HandleSaveFamilyInfo godoc
// @Summary      Save Family Information
// @Description  Updates the family information (including parents and guardians) for a student.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        studentRecordID path      int                 true "Student Record ID"
// @Param        request         body      UpdateFamilyRequest true "Family Info Data"
// @Success      200             {object}  map[string]string   "Message: Family information saved successfully"
// @Failure      400             {object}  map[string]string   "Invalid input or ID"
// @Failure      500             {object}  map[string]string   "Internal Server Error"
// @Router       /students/onboarding/family/{studentRecordID} [put]
func (h *Handler) HandleSaveFamilyInfo(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid student record ID"},
		)
		return
	}

	var req UpdateFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.SaveFamilyInfo(c.Request.Context(), studentRecordID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(
		http.StatusOK, gin.H{"message": "Family information saved successfully"})
}

// HandleSaveEducationInfo godoc
// @Summary      Save Education Information
// @Description  Updates the educational background information for a student.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        studentRecordID path      int                    true "Student Record ID"
// @Param        request         body      UpdateEducationRequest true "Education Info Data"
// @Success      200             {object}  map[string]string      "Message: Education information saved successfully"
// @Failure      400             {object}  map[string]string      "Invalid input or ID"
// @Failure      500             {object}  map[string]string      "Internal Server Error"
// @Router       /students/onboarding/education/{studentRecordID} [put]
func (h *Handler) HandleSaveEducationInfo(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student record ID"})
		return
	}

	var req UpdateEducationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.SaveEducationInfo(c.Request.Context(), studentRecordID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Education information saved successfully"})
}

// HandleSaveAddressInfo godoc
// @Summary      Save Address Information
// @Description  Updates the address details (region, province, city, barangay) for a student.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        studentRecordID path      int                  true "Student Record ID"
// @Param        request         body      UpdateAddressRequest true "Address Info Data"
// @Success      200             {object}  map[string]string    "Message: Address information saved successfully"
// @Failure      400             {object}  map[string]string    "Invalid input or ID"
// @Failure      500             {object}  map[string]string    "Internal Server Error"
// @Router       /students/onboarding/address/{studentRecordID} [put]
func (h *Handler) HandleSaveAddressInfo(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student record ID"})
		return
	}

	var req UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.SaveAddressInfo(c.Request.Context(), studentRecordID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address information saved successfully"})
}

// HandleSaveHealthInfo godoc
// @Summary      Save Health Information
// @Description  Updates the health records (vision, hearing, remarks) for a student.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        studentRecordID path      int                       true "Student Record ID"
// @Param        request         body      UpdateHealthRecordRequest true "Health Info Data"
// @Success      200             {object}  map[string]string         "Message: Health information saved successfully"
// @Failure      400             {object}  map[string]string         "Invalid input or ID"
// @Failure      500             {object}  map[string]string         "Internal Server Error"
// @Router       /students/onboarding/health/{studentRecordID} [put]
func (h *Handler) HandleSaveHealthInfo(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student record ID"})
		return
	}

	var req UpdateHealthRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.SaveHealthRecord(c.Request.Context(), studentRecordID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Health information saved successfully"})
}

// HandleSaveFinanceInfo godoc
// @Summary      Save Finance Information
// @Description  Updates the financial information for a student.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        studentRecordID path      int                   true "Student Record ID"
// @Param        request         body      UpdateFinanceRequest  true "Finance Info Data"
// @Success      200             {object}  map[string]string     "Message: Finance information saved successfully"
// @Failure      400             {object}  map[string]string     "Invalid input or ID"
// @Failure      500             {object}  map[string]string     "Internal Server Error"
// @Router       /students/onboarding/finance/{studentRecordID} [put]
func (h *Handler) HandleSaveFinanceInfo(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student record ID"})
		return
	}

	var req UpdateFinanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.SaveFinanceInfo(c.Request.Context(), studentRecordID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Finance information saved successfully"})
}
