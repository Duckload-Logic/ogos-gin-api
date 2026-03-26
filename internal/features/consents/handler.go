package consents

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
)

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetLatestDocument(c *gin.Context) {
	docType := c.Param("type")
	if docType != "terms" && docType != "privacy" {
		response.SendFail(c, gin.H{"error": "Invalid document type"})
		return
	}

	doc, err := h.service.GetLatestDocument(c.Request.Context(), docType)
	if err != nil {
		response.SendError(
			c,
			"Failed to fetch document",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, doc)
}

func (h *Handler) GetDocumentContent(c *gin.Context) {
	docType := c.Param("type")
	if docType != "terms" && docType != "privacy" {
		response.SendFail(c, gin.H{"error": "Invalid document type"})
		return
	}

	content, contentType, err := h.service.GetDocumentContent(
		c.Request.Context(),
		docType,
	)
	if err != nil {
		response.SendError(
			c,
			"Failed to fetch document content",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	c.Data(http.StatusOK, contentType, content)
}

func (h *Handler) GetConsentCheck(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	if userID == "" {
		response.SendFail(
			c,
			gin.H{"error": "Unauthorized"},
			http.StatusUnauthorized,
		)
		return
	}

	docIDStr := c.Param("docID")

	docID, err := strconv.Atoi(docIDStr)
	if err != nil {
		response.SendFail(c, gin.H{"error": "Invalid document ID"})
		return
	}

	accepted, err := h.service.HasUserAccepted(
		c.Request.Context(),
		userID,
		docID,
	)
	if err != nil {
		response.SendError(
			c,
			"Failed to check user consent",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{"accepted": accepted})
}

func (h *Handler) PostConsent(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	if userID == "" {
		response.SendFail(
			c,
			gin.H{"error": "Unauthorized"},
			http.StatusUnauthorized,
		)
		return
	}

	docType := c.Param("type")
	if docType != "terms" && docType != "privacy" {
		response.SendFail(c, gin.H{"error": "Invalid document type"})
		return
	}

	docIDStr := c.Param("docID")
	docID, err := strconv.Atoi(docIDStr)
	if err != nil {
		response.SendFail(c, gin.H{"error": "Invalid document ID"})
		return
	}

	err = h.service.SaveConsent(c.Request.Context(), userID, docID)
	if err != nil {
		log.Printf("[PostConsent] {Service Error}: %v", err)
		response.SendError(
			c,
			"Failed to save consent",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(
		c,
		gin.H{"message": "Consent saved successfully"},
		http.StatusAccepted,
	)
}

// Admin endpoint to view user consent history
func (h *Handler) GetConsentHistory(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		response.SendFail(c, gin.H{"error": "Invalid user ID"})
		return
	}

	consents, err := h.service.ListUserConsentHistory(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		response.SendError(
			c,
			"Failed to fetch consent history",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, consents)
}

func (h *Handler) PostDocument(c *gin.Context) {
	docType := c.PostForm("docType")

	fileHeader, err := c.FormFile("document")
	if err != nil {
		response.SendFail(c, gin.H{"error": "Document file is required"})
		return
	}

	// VALIDATION: Strictly check for .md extension
	if filepath.Ext(fileHeader.Filename) != ".md" {
		response.SendFail(
			c,
			gin.H{"error": "Only Markdown (.md) files are allowed"},
		)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.SendError(
			c,
			"Failed to read file",
			http.StatusInternalServerError,
			nil,
		)
		return
	}
	defer file.Close()

	seeker, ok := file.(io.ReadSeeker)
	if !ok {
		response.SendError(
			c,
			"Invalid file stream",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	// pass "text/markdown" as the content type
	err = h.service.UploadNewDocument(
		c.Request.Context(),
		docType,
		seeker,
		"text/markdown",
	)
	if err != nil {
		log.Printf("[PostDocument] {Upload Error}: %v", err)
		response.SendError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	response.SendSuccess(c, gin.H{
		"message": "Markdown document uploaded and activated successfully",
	})
}
