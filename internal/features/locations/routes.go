package locations

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	routes := r.Group("/locations")
	{
		routes.GET("/regions", h.HandleGetRegions)
		routes.GET("/regions/:regionCode/provinces", h.HandleGetProvincesByRegion)
		routes.GET("/regions/:regionCode/cities", h.HandleGetCitiesByRegion)
		routes.GET("/provinces/:provinceCode/cities", h.HandleGetCitiesByProvince)
		routes.GET("/cities/:cityCode/barangays", h.HandleGetBarangaysByCity)
	}
}
