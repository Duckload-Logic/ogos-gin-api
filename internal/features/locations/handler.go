package locations

import (
	"net/http"

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

func (h *Handler) HandleGetProvincesByRegion(c *gin.Context) {
	regionCode := c.Param("regionCode")
	if regionCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region code"})
		return
	}

	provinces, err := h.service.GetProvincesByRegion(c.Request.Context(), regionCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve provinces"})
		return
	}

	c.JSON(http.StatusOK, provinces)
}

func (h *Handler) HandleGetCitiesByProvince(c *gin.Context) {
	provinceCode := c.Param("provinceCode")
	if provinceCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid province code"})
		return
	}

	cities, err := h.service.GetCitiesByProvince(c.Request.Context(), provinceCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cities"})
		return
	}

	c.JSON(http.StatusOK, cities)
}

func (h *Handler) HandleGetCitiesByRegion(c *gin.Context) {
	regionCode := c.Param("regionCode")
	if regionCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region code"})
		return
	}

	cities, err := h.service.GetCitiesByRegion(c.Request.Context(), regionCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cities"})
		return
	}

	c.JSON(http.StatusOK, cities)
}

func (h *Handler) HandleGetBarangaysByCity(c *gin.Context) {
	cityCode := c.Param("cityCode")
	if cityCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid city code"})
		return
	}

	barangays, err := h.service.GetBarangaysByCity(c.Request.Context(), cityCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve barangays"})
		return
	}

	c.JSON(http.StatusOK, barangays)
}
