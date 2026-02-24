package locations

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	routes := r.Group("/locations")
	{
		routes.GET("/regions", h.HandleGetRegions)
		routes.GET("/regions/:regionID/cities", h.HandleGetCitiesByRegion)
		routes.GET("/cities/:cityID/barangays", h.HandleGetBarangaysByCity)
	}
}
