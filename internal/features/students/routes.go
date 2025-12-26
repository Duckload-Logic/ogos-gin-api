package students

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	studentRoutes := r.Group("/students")
	{
		// Retrieve Routes
		studentRoutes.GET("/profile/base/:userID", h.HandleGetBaseProfile)
		studentRoutes.GET(
			"/profile/family/:studentRecordID", h.HandleGetFamilyInfo,
		)
		studentRoutes.GET(
			"/profile/guardians/:studentRecordID", h.HandleGetGuardiansInfo,
		)
		studentRoutes.GET(
			"/profile/guardians/primary/:studentRecordID",
			h.HandleGetPrimaryGuardianInfo,
		)
		studentRoutes.GET(
			"/profile/education/:studentRecordID", h.HandleGetEducationInfo,
		)
		studentRoutes.GET(
			"/profile/address/:studentRecordID", h.HandleGetAddressInfo,
		)
		studentRoutes.GET(
			"/profile/health/:studentRecordID", h.HandleGetHealthInfo,
		)

		// Save/Update Routes
		studentRoutes.POST("/onboarding/base", h.HandleSaveBaseProfile)
		studentRoutes.PUT(
			"/onboarding/family/:studentRecordID", h.HandleSaveFamilyInfo,
		)
		studentRoutes.PUT(
			"/onboarding/education/:studentRecordID", h.HandleSaveEducationInfo,
		)
		studentRoutes.PUT(
			"/onboarding/address/:studentRecordID", h.HandleSaveAddressInfo,
		)
		studentRoutes.PUT(
			"/onboarding/health/:studentRecordID", h.HandleSaveHealthInfo,
		)
	}
}
