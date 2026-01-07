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

	// Shared access routes for Counselors, Front Desk, and Students
	sharedAccessGroup := studentRoutes.Group("/")
	sharedAccessGroup.Use(middleware.RoleMiddleware(
		int(constants.CounselorRoleID),
		int(constants.FrontDeskRoleID),
		int(constants.StudentRoleID),
	))
	{
		// Retrieve Routes
		sharedAccessGroup.GET(
			"/records",
			h.HandleListStudents,
		)

		sharedAccessGroup.GET(
			"/record/:userID",
			userLookup,
			h.HandleGetStudentRecord,
		)

		sharedAccessGroup.GET(
			"/record/progress/:userID",
			userLookup,
			h.HandleGetStudentRecordProgress,
		)

		sharedAccessGroup.GET(
			"/:userID",
			userLookup,
			h.HandleGetStudent,
		)

		studentRecordGroup := sharedAccessGroup.Group("/")
		studentRecordGroup.Use(studentRecordLookup)
		{
			studentRecordGroup.GET(
				"/profile/family/:studentRecordID", h.HandleGetFamilyInfo,
			)
			studentRecordGroup.GET(
				"/profile/guardians/:studentRecordID", h.HandleGetGuardiansInfo,
			)
			studentRecordGroup.GET(
				"/record/family/:studentRecordID", h.HandleGetFamilyInfo,
			)
			studentRecordGroup.GET(
				"/record/parents/:studentRecordID", h.HandleGetParentsInfo,
			)
			studentRecordGroup.GET(
				"/profile/education/:studentRecordID", h.HandleGetEducationInfo,
			)
			studentRecordGroup.GET(
				"/profile/address/:studentRecordID", h.HandleGetAddressInfo,
			)
			studentRecordGroup.GET(
				"/profile/health/:studentRecordID", h.HandleGetHealthInfo,
			)
		}
	}

	// Student-only access routes
	studentOnly := studentRoutes.Group("/")
	studentOnly.Use(middleware.RoleMiddleware(
		int(constants.StudentRoleID),
	))
	{
		studentOnly.POST(
			"/onboarding/base/:userID",
			userLookup,
			h.HandleSaveBaseProfile,
		)

		studentRecordGroup := studentOnly.Group("/")
		studentRecordGroup.Use(studentRecordLookup)
		{
			studentRecordGroup.PUT(
				"/onboarding/family/:studentRecordID",
				h.HandleSaveFamilyInfo,
			)

			studentRecordGroup.PUT(
				"/onboarding/education/:studentRecordID",
				h.HandleSaveEducationInfo,
			)

			studentRecordGroup.PUT(
				"/onboarding/address/:studentRecordID",
				h.HandleSaveAddressInfo,
			)

			studentRecordGroup.PUT(
				"/onboarding/health/:studentRecordID",
				h.HandleSaveHealthInfo,
			)

			studentRecordGroup.PUT(
				"/onboarding/finance/:studentRecordID",
				h.HandleSaveFinanceInfo,
			)

			studentRecordGroup.POST(
				"/onboarding/complete/:studentRecordID",
				h.HandleCompleteOnboarding,
			)
		}
	}
}
