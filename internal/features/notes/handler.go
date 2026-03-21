package notes

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleGetStudentSignificantNotes(
	c *gin.Context,
) {
	iirID := c.GetString("iirID")

	significantNotes, err := h.service.GetStudentSignificantNotes(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		log.Printf(
			"[GetStudentSignificantNotes] {Database Query}: %v",
			err,
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Failed to get student significant notes",
			},
		)
		return
	}

	c.JSON(http.StatusOK, significantNotes)
}

func (h *Handler) HandlePostStudentSignificantNote(
	c *gin.Context,
) {
	iirID := c.Param("iirID")
	if iirID == "" {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "IIR ID not found"},
		)
		return
	}

	var noteReq SignificantNoteDTO
	if err := c.ShouldBindJSON(&noteReq); err != nil {
		log.Printf(
			"[PostSignificantNote] {JSON Bind}: %v",
			err,
		)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid request body"},
		)
		return
	}

	err := h.service.CreateSignificantNote(
		c.Request.Context(),
		iirID,
		noteReq,
	)
	if err != nil {
		log.Printf(
			"[PostSignificantNote] {Database Insert}: %v",
			err,
		)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to save significant note"},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"message": "Significant note saved successfully"},
	)
}
