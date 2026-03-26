package locations

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

func RegisterRoutes(
	r *gin.RouterGroup,
	h *Handler,
	redis *datastore.RedisClient,
) {
	routes := r.Group("/locations")
	{
		routes.GET("/regions", h.GetRegions)
		routes.GET("/regions/:regionCode/provinces", h.GetProvincesByRegion)
		routes.GET("/regions/:regionCode/cities", h.GetCitiesByRegion)
		routes.GET("/provinces/:provinceCode/cities", h.GetCitiesByProvince)
		routes.GET("/cities/:cityCode/barangays", h.GetBarangaysByCity)
	}
}
