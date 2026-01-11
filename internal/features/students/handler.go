package students

import (
	"context"
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

// HandleGetStudentRecord godoc
// @Summary      Get Student Record
// @Description  Retrieves the student record associated with a given user ID.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        userID   path      int  true  "User ID"
// @Success      200      {object}  map[string]int "Returns {studentRecordID: <id>}"
// @Failure      400      {object}  map[string]string "Invalid user ID"
// @Failure      500      {object}  map[string]string "Failed to get student record ID"
// @Router       /students/record/{userID} [get]
func (h *Handler) HandleGetStudentRecord(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	studentRecord, err := h.service.GetStudentRecordByStudentID(
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

	c.JSON(http.StatusOK, gin.H{"studentRecord": studentRecord})
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
// @Param        userID   					path      int  true  "User ID"
// @Param 	     include_emergency_contact  query     bool  false  "Include emergency contact information"
// @Param        include_family     		query     bool  false  "Include family information"
// @Param        include_health     		query     bool  false  "Include health information"
// @Param        include_education  		query     bool  false  "Include educational background"
// @Param        include_address    		query     bool  false  "Include address information"
// @Param        include_finance    		query     bool  false  "Include financial information"
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

	studentRecord, err := h.service.GetStudentRecordByStudentID(
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

	studentRecordID := studentRecord.ID

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
		includeReq.IncludeFinance || includeReq.IncludeEmergencyContact ||
		includeReq.IncludeReasons || includeReq.IncludeParents {

		compResponse := ComprehensiveProfileResponse{
			StudentProfile: profile,
		}

		if includeReq.IncludeReasons {
			reasons, err := h.service.GetStudentEnrollmentReasons(
				c.Request.Context(), studentRecordID,
			)
			if err == nil && reasons != nil {
				compResponse.EnrollmentReasons = reasons
			}
		}

		if includeReq.IncludeEmergencyContact {
			emergencyContact, err := h.service.GetEmergencyContactInfo(
				c.Request.Context(), studentRecordID,
			)
			if err == nil && emergencyContact != nil {
				compResponse.EmergencyContact = emergencyContact
			}
		}

		if includeReq.IncludeFamily {
			family, err := h.service.GetFamilyInfo(
				c.Request.Context(), studentRecordID,
			)
			if err == nil && family != nil {
				compResponse.Family = family
			}
		}

		if includeReq.IncludeParents {
			parents, err := h.service.GetParentsInfo(
				c.Request.Context(), studentRecordID,
			)

			fmt.Println("Error getting parents info:", err)

			if err == nil && parents != nil {
				compResponse.Parents = parents
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
			addresses, err := h.service.GetAddressInfo(
				c.Request.Context(), studentRecordID,
			)
			if err == nil && addresses != nil {
				compResponse.Addresses = addresses
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

// HandleGetEmergencyContactInfo godoc
// @Summary      Get Emergency Contact Information
// @Description  Retrieves the emergency contact details for a student.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        studentRecordID path      int  true  "Student Record ID"
// @Success      200             {object}  GetEmergencyContactResponse
// @Failure      400             {object}  map[string]string "Invalid student record ID"
// @Failure      500             {object}  map[string]string "Failed to get emergency contact info"
// @Router       /students/profile/emergency-contact/{studentRecordID} [get]
func (h *Handler) HandleGetEmergencyContactInfo(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid student record ID"},
		)
		return
	}

	resp, err := h.service.GetEmergencyContactInfo(c.Request.Context(), studentRecordID)
	if err != nil {
		fmt.Println("Error getting emergency contact info:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get emergency contact info"},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"emergencyContact": resp})
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

// HandleGetParentsInfo godoc
// @Summary      Get Parents Information
// @Description  Retrieves a list of parents linked to the student record.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        studentRecordID path      int  true  "Student Record ID"
// @Success      200             {array}   Parent
// @Failure      400             {object}  map[string]string "Invalid student record ID"
// @Failure      500             {object}  map[string]string "Failed to get parents info"
// @Router       /students/profile/parents/{studentRecordID} [get]
func (h *Handler) HandleGetParentsInfo(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid student record ID"},
		)
		return
	}

	resp, err := h.service.GetParentsInfo(
		c.Request.Context(),
		studentRecordID,
	)
	if err != nil {
		fmt.Println("Error getting parents info:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get parents info"},
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

	err = h.service.SaveEmergencyContactInfo(c.Request.Context(), studentRecordID, *req.EmergencyContact)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.service.SaveAddressInfo(c.Request.Context(), studentRecordID, req.Addresses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Base profile saved successfully"})
}

// HandleSaveFamilyInfo godoc
// @Summary      Save Family Information
// @Description  Updates the family information (including parents and parents) for a student.
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

	err = h.service.SaveFinanceInfo(c.Request.Context(), studentRecordID, req.UpdateFinanceRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK, gin.H{"message": "Family information saved successfully"},
	)

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

func (h *Handler) HandleCompleteOnboarding(c *gin.Context) {
	studentRecordID, err := strconv.Atoi(c.Param("studentRecordID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student record ID"})
		return
	}

	err = h.service.CompleteOnboarding(c.Request.Context(), studentRecordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student onboarding completed successfully"})
}

// ========================================
// |                                      |
// |       ADDITIONAL HANDLER FUNCS       |
// |                                      |
// ========================================

// HandleGetStudentRecordProgress godoc
// @Summary      Get Student Record Progress
// @Description  Retrieves the completion progress of various sections in a student's record.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        userID   path      int  true  "User ID"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]string "Invalid user ID"
// @Failure      500      {object}  map[string]string "Failed to get student record progress"
// @Router       /students/record/progress/{userID} [get]
func (h *Handler) HandleGetStudentRecordProgress(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	studentRecord, err := h.service.GetStudentRecordByStudentID(c.Request.Context(), userID)
	if err != nil {
		fmt.Println("Error getting student record ID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get student record ID"})
		return
	}

	studentRecordID := studentRecord.ID
	ctx := c.Request.Context()

	// Step 1: Enrollment Reasons
	enrollmentReasons, err := h.getEnrollmentReasonsSafe(ctx, studentRecordID)
	if err != nil {
		return // Error already handled and response sent
	}

	// Step 2: Personal Information (Profile, Address, Emergency Contact)
	personalData, err := h.getPersonalInformationSafe(ctx, studentRecordID)
	if err != nil {
		return // Error already handled and response sent
	}

	// Step 3: Educational Background
	education, err := h.getEducationInfoSafe(ctx, studentRecordID)
	if err != nil {
		return // Error already handled and response sent
	}

	// Step 4: Family Background (Family, Parents, Finance)
	familyData, err := h.getFamilyInformationSafe(ctx, studentRecordID)
	if err != nil {
		return // Error already handled and response sent
	}

	// Step 5: Health Information
	health, err := h.getHealthInfoSafe(ctx, studentRecordID)
	if err != nil {
		return // Error already handled and response sent
	}

	c.JSON(http.StatusOK, gin.H{
		"isSubmitted":       studentRecord.IsSubmitted,
		"enrollmentReasons": enrollmentReasons,
		"personalInfo":      personalData,
		"education":         education,
		"family":            familyData,
		"health":            health,
	})
}

// Helper functions to fetch and handle errors consistently

// getEnrollmentReasonsSafe fetches enrollment reasons and handles errors
func (h *Handler) getEnrollmentReasonsSafe(ctx context.Context, studentRecordID int) (interface{}, error) {
	data, err := h.service.GetStudentEnrollmentReasons(ctx, studentRecordID)
	if err != nil {
		fmt.Println("Error getting enrollment reasons:", err)
		// Note: Error response is NOT sent here, caller should send it
		return nil, err
	}
	return data, nil
}

// getPersonalInformationSafe fetches personal info (profile, address, emergency contact)
func (h *Handler) getPersonalInformationSafe(ctx context.Context, studentRecordID int) (map[string]interface{}, error) {
	personalData := make(map[string]interface{})

	// Get Student Profile
	profile, err := h.service.GetBaseProfile(ctx, studentRecordID)
	if err != nil {
		fmt.Println("Error getting student profile:", err)
		return nil, err
	}
	personalData["profile"] = profile

	// Get Address Info
	addresses, err := h.service.GetAddressInfo(ctx, studentRecordID)
	if err != nil {
		fmt.Println("Error getting addresses info:", err)
		return nil, err
	}
	personalData["addresses"] = addresses
	fmt.Println(addresses)

	// Get Emergency Contact Info
	emergencyContact, err := h.service.GetEmergencyContactInfo(ctx, studentRecordID)
	if err != nil {
		fmt.Println("Error getting emergency contact info:", err)
		return nil, err
	}
	personalData["emergencyContact"] = emergencyContact

	return personalData, nil
}

// getEducationInfoSafe fetches education information
func (h *Handler) getEducationInfoSafe(ctx context.Context, studentRecordID int) (interface{}, error) {
	data, err := h.service.GetEducationInfo(ctx, studentRecordID)
	if err != nil {
		fmt.Println("Error getting educational background:", err)
		return nil, err
	}
	return data, nil
}

// getFamilyInformationSafe fetches family-related info (family, parents, finance)
func (h *Handler) getFamilyInformationSafe(ctx context.Context, studentRecordID int) (map[string]interface{}, error) {
	familyData := make(map[string]interface{})

	// Get Family Background
	family, err := h.service.GetFamilyInfo(ctx, studentRecordID)
	if err != nil {
		fmt.Println("Error getting family background:", err)
		return nil, err
	}
	familyData["background"] = family

	// Get Parents Info
	parents, err := h.service.GetParentsInfo(ctx, studentRecordID)
	if err != nil {
		fmt.Println("Error getting parents info:", err)
		return nil, err
	}
	familyData["parents"] = parents

	// Get Finance Info
	finance, err := h.service.GetFinanceInfo(ctx, studentRecordID)
	if err != nil {
		fmt.Println("Error getting finance info:", err)
		return nil, err
	}
	familyData["finance"] = finance

	return familyData, nil
}

// getHealthInfoSafe fetches health information
func (h *Handler) getHealthInfoSafe(ctx context.Context, studentRecordID int) (interface{}, error) {
	data, err := h.service.GetHealthInfo(ctx, studentRecordID)
	if err != nil {
		fmt.Println("Error getting health background:", err)
		return nil, err
	}
	return data, nil
}
