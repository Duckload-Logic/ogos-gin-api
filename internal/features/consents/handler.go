package consents

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleGetLatestDocument(c *gin.Context) {
	docType := c.Param("type")
	if docType != "terms" && docType != "privacy" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document type"})
		return
	}

	doc, err := h.service.GetLatestDocument(c.Request.Context(), docType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch document"})
		return
	}

	c.JSON(http.StatusOK, doc)
}

func (h *Handler) HandleGetDocumentContent(c *gin.Context) {
	docType := c.Param("type")
	if docType != "terms" && docType != "privacy" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document type"})
		return
	}

	content, contentType, err := h.service.GetDocumentContent(c.Request.Context(), docType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch document content"})
		return
	}

	c.Data(http.StatusOK, contentType, content)
}

func (h *Handler) HandleCheckUserConsent(c *gin.Context) {
	userID := c.MustGet("userID").(string) // Get userID from context set by AuthMiddleware
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	docIDStr := c.Param("docID")

	docID, err := strconv.Atoi(docIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	accepted, err := h.service.HasUserAccepted(c.Request.Context(), userID, docID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user consent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accepted": accepted})
}

func (h *Handler) HandleSaveConsent(c *gin.Context) {
	userID := c.MustGet("userID").(string) // Get userID from context set by AuthMiddleware
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	docType := c.Param("type")
	if docType != "terms" && docType != "privacy" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document type"})
		return
	}

	docIDStr := c.Param("docID")
	docID, err := strconv.Atoi(docIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	err = h.service.SaveConsent(c.Request.Context(), userID, docID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save consent"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Consent saved successfully"})
}

// Admin endpoint to view user consent history
func (h *Handler) HandleListUserConsentHistory(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	consents, err := h.service.ListUserConsentHistory(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch consent history"})
		return
	}

	c.JSON(http.StatusOK, consents)
}

func (h *Handler) HandleUploadDocument(c *gin.Context) {
	docType := c.PostForm("docType")

	fileHeader, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document file is required"})
		return
	}

	// VALIDATION: Strictly check for .md extension
	if filepath.Ext(fileHeader.Filename) != ".md" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only Markdown (.md) files are allowed"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer file.Close()

	seeker, ok := file.(io.ReadSeeker)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid file stream"})
		return
	}

	// pass "text/markdown" as the content type
	err = h.service.UploadNewDocument(c.Request.Context(), docType, seeker, "text/markdown")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Markdown document uploaded and activated successfully"})
}
