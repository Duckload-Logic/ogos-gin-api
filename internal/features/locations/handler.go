package locations

import (
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

func (h *Handler) HandleGetRegions(c *gin.Context) {
	regions, err := h.service.GetRegions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve regions"})
		return
	}

	c.JSON(http.StatusOK, regions)
}

func (h *Handler) HandleGetCitiesByRegion(c *gin.Context) {
	regionID, err := strconv.Atoi(c.Param("regionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region ID"})
		return
	}

	cities, err := h.service.GetCitiesByRegion(c.Request.Context(), regionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cities"})
		return
	}

	c.JSON(http.StatusOK, cities)
}

func (h *Handler) HandleGetBarangaysByCity(c *gin.Context) {
	cityID, err := strconv.Atoi(c.Param("cityID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid city ID"})
		return
	}

	barangays, err := h.service.GetBarangaysByCity(c.Request.Context(), cityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve barangays"})
		return
	}

	c.JSON(http.StatusOK, barangays)
}
