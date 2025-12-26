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
