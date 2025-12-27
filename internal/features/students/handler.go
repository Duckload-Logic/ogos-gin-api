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

// HandleGetBaseProfile godoc
// @Summary      Get Student Base Profile
// @Description  Retrieves the basic profile information for a student using their User ID.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        userID   path      int  true  "User ID"
// @Success      200      {object}  StudentRecord
// @Failure      400      {object}  map[string]string "Invalid user ID"
// @Failure      500      {object}  map[string]string "Failed to get base profile"
// @Router       /students/profile/base/{userID} [get]

// HandleGetBaseProfile
func (h *Handler) HandleGetBaseProfile(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	resp, err := h.service.GetBaseProfile(c.Request.Context(), userID)
	if err != nil {
		fmt.Println("Error getting base profile:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get base profile"},
		)
		return
	}

	c.JSON(http.StatusOK, resp)
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

// HandleGetFamilyInfo
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

// HandleGetGuardiansInfo
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

// HandleGetPrimaryGuardianInfo
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

// HandleGetEducationInfo
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

// HandleGetAddressInfo
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

// HandleGetHealthInfo
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
// |       UPSERT HANDLER FUNCTIONS       |
// |                                      |
// ========================================

// HandleSaveBaseProfile godoc
// @Summary      Create or Update Base Profile
// @Description  Creates a new student record or updates the basic profile information.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        request body      CreateStudentRecordRequest true "Base Profile Data"
// @Success      200     {object}  map[string]interface{}     "Returns {student_record_id: <id>}"
// @Failure      400     {object}  map[string]string          "Invalid input"
// @Failure      500     {object}  map[string]string          "Internal Server Error"
// @Router       /students/onboarding/base [post]

// HandleSaveBaseProfile
func (h *Handler) HandleSaveBaseProfile(c *gin.Context) {
	var req CreateStudentRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.SaveBaseProfile(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"student_record_id": id})
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

// HandleSaveFamilyInfo
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

// HandleSaveEducationInfo
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

// HandleSaveAddressInfo
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

// HandleSaveHealthInfo
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
