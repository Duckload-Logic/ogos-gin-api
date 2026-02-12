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
	studentRecordLookup := middleware.OwnershipMiddleware(db, "studentRecordID")

	// Admin-only access
	adminOnly := studentRoutes.Group("/")
	adminOnly.Use(middleware.RoleMiddleware(
		int(constants.CounselorRoleID),
	))
	{
		adminOnly.DELETE("/record/:studentRecordID", studentRecordLookup, h.HandleDeleteStudentRecord)
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
		sharedRetrieve.GET("/records", h.HandleListStudents)
		sharedRetrieve.GET("/record/:userID", userLookup, h.HandleGetStudentRecord)
		sharedRetrieve.GET("/record/progress/:userID", userLookup, h.HandleGetStudentRecordProgress)
		sharedRetrieve.GET("/:userID", userLookup, h.HandleGetStudent)

		// Student record detail retrieval
		recordDetails := sharedRetrieve.Group("/")
		recordDetails.Use(studentRecordLookup)
		{
			recordDetails.GET("/record/enrollment-reasons/:studentRecordID", h.HandleGetStudentEnrollmentReasons)
			recordDetails.GET("/record/base/:studentRecordID", h.HandleGetBaseProfile)
			recordDetails.GET("/record/family/:studentRecordID", h.HandleGetFamilyInfo)
			recordDetails.GET("/record/parents/:studentRecordID", h.HandleGetParentsInfo)
			recordDetails.GET("/record/education/:studentRecordID", h.HandleGetEducationInfo)
			recordDetails.GET("/record/address/:studentRecordID", h.HandleGetAddressInfo)
			recordDetails.GET("/record/health/:studentRecordID", h.HandleGetHealthInfo)
			recordDetails.GET("/record/finance/:studentRecordID", h.HandleGetFinanceInfo)
		}
	}

	// Student-only access - Onboarding (Create/Update)
	studentOnboarding := studentRoutes.Group("/")
	studentOnboarding.Use(middleware.RoleMiddleware(
		int(constants.StudentRoleID),
	))
	{
		// Initialize student record
		studentOnboarding.POST("/onboarding/:userID", userLookup, h.HandleCreateStudentRecord)

		// Update student information
		onboardingSteps := studentOnboarding.Group("/")
		onboardingSteps.Use(studentRecordLookup)
		{
			onboardingSteps.PUT("/onboarding/enrollment-reasons/:studentRecordID", h.HandleSaveEnrollmentReasons)
			onboardingSteps.PUT("/onboarding/base/:studentRecordID", h.HandleSaveBaseProfile)
			onboardingSteps.PUT("/onboarding/family/:studentRecordID", h.HandleSaveFamilyInfo)
			onboardingSteps.PUT("/onboarding/education/:studentRecordID", h.HandleSaveEducationInfo)
			onboardingSteps.PUT("/onboarding/address/:studentRecordID", h.HandleSaveAddressInfo)
			onboardingSteps.PUT("/onboarding/health/:studentRecordID", h.HandleSaveHealthInfo)
			onboardingSteps.PUT("/onboarding/finance/:studentRecordID", h.HandleSaveFinanceInfo)
			onboardingSteps.POST("/onboarding/complete/:studentRecordID", h.HandleCompleteOnboarding)
		}
	}
}
