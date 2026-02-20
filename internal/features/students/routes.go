package students

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(db *sql.DB, r *gin.RouterGroup, h *Handler) {
	studentRoutes := r.Group("/students")
	studentRoutes.Use(middleware.AuthMiddleware())

	userLookup := middleware.OwnershipMiddleware(db, "userID")
	inventoryRecordLookup := middleware.OwnershipMiddleware(db, "inventoryRecordID")

	// Admin-only access
	adminOnly := studentRoutes.Group("/")
	adminOnly.Use(middleware.RoleMiddleware(
		int(constants.CounselorRoleID),
	))
	{
		adminOnly.DELETE("/inventory/:inventoryRecordID", inventoryRecordLookup, h.HandleDeleteInventoryRecord)
	}

	// Shared access (Counselors, Front Desk, Students) - Retrieve only
	sharedRetrieve := studentRoutes.Group("/")
	sharedRetrieve.Use(middleware.RoleMiddleware(
		int(constants.CounselorRoleID),
		int(constants.FrontDeskRoleID),
		int(constants.StudentRoleID),
	))
	{
		// List and basic retrieval
		sharedRetrieve.GET("/inventory", h.HandleListStudents)
		sharedRetrieve.GET("/inventory/:inventoryRecordID", userLookup, h.HandleGetInventoryRecord)
		sharedRetrieve.GET("/inventory/:inventoryRecordID/progress", userLookup, h.HandleGetInventoryRecordProgress)
		sharedRetrieve.GET("/:userID", userLookup, h.HandleGetStudent)

		// Student record detail retrieval
		recordDetails := sharedRetrieve.Group("/")
		recordDetails.Use(inventoryRecordLookup)
		{
			recordDetails.GET("/inventory/enrollment/:inventoryRecordID", h.HandleGetStudentEnrollmentReasons)
			recordDetails.GET("/inventory/base/:inventoryRecordID", h.HandleGetBaseProfile)
			recordDetails.GET("/inventory/family/:inventoryRecordID", h.HandleGetFamilyInfo)
			recordDetails.GET("/inventory/related-persons/:inventoryRecordID", h.HandleGetRelatedPersonsInfo)
			recordDetails.GET("/inventory/education/:inventoryRecordID", h.HandleGetEducationInfo)
			recordDetails.GET("/inventory/address/:inventoryRecordID", h.HandleGetAddressInfo)
			recordDetails.GET("/inventory/health/:inventoryRecordID", h.HandleGetHealthInfo)
			recordDetails.GET("/inventory/finance/:inventoryRecordID", h.HandleGetFinanceInfo)
		}
	}

	// Student-only access - Onboarding (Create/Update)
	studentOnboarding := studentRoutes.Group("/")
	studentOnboarding.Use(middleware.RoleMiddleware(
		int(constants.StudentRoleID),
	))
	{
		// Initialize student record
		studentOnboarding.POST("/onboarding/:userID", userLookup, h.HandleCreateInventoryRecord)

		// Update student information
		onboardingSteps := studentOnboarding.Group("/")
		onboardingSteps.Use(inventoryRecordLookup)
		{
			onboardingSteps.PUT("/onboarding/enrollment-reasons/:inventoryRecordID", h.HandleSaveEnrollmentReasons)
			onboardingSteps.PUT("/onboarding/base/:inventoryRecordID", h.HandleSaveBaseProfile)
			onboardingSteps.PUT("/onboarding/family/:inventoryRecordID", h.HandleSaveFamilyInfo)
			onboardingSteps.PUT("/onboarding/education/:inventoryRecordID", h.HandleSaveEducationInfo)
			onboardingSteps.PUT("/onboarding/address/:inventoryRecordID", h.HandleSaveAddressInfo)
			onboardingSteps.PUT("/onboarding/health/:inventoryRecordID", h.HandleSaveHealthInfo)
			onboardingSteps.PUT("/onboarding/finance/:inventoryRecordID", h.HandleSaveFinanceInfo)
			onboardingSteps.POST("/onboarding/complete/:inventoryRecordID", h.HandleCompleteOnboarding)
		}
	}
}
