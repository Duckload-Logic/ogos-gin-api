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
				"/record/enrollment-reasons/:studentRecordID",
				h.HandleGetStudentEnrollmentReasons,
			)
			studentRecordGroup.GET(
				"/record/base/:studentRecordID", h.HandleGetBaseProfile,
			)
			studentRecordGroup.GET(
				"/record/family/:studentRecordID", h.HandleGetFamilyInfo,
			)
			studentRecordGroup.GET(
				"/record/parents/:studentRecordID", h.HandleGetParentsInfo,
			)
			studentRecordGroup.GET(
				"/record/education/:studentRecordID", h.HandleGetEducationInfo,
			)
			studentRecordGroup.GET(
				"/record/address/:studentRecordID", h.HandleGetAddressInfo,
			)
			studentRecordGroup.GET(
				"/record/health/:studentRecordID", h.HandleGetHealthInfo,
			)
			studentRecordGroup.GET(
				"/record/finance/:studentRecordID", h.HandleGetFinanceInfo,
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
			"/onboarding/:userID",
			userLookup,
			h.HandleCreateStudentRecord,
		)

		studentRecordGroup := studentOnly.Group("/")
		studentRecordGroup.Use(studentRecordLookup)
		{
			studentRecordGroup.PUT(
				"/onboarding/enrollment-reasons/:studentRecordID",
				h.HandleSaveEnrollmentReasons,
			)

			studentRecordGroup.PUT(
				"/onboarding/base/:studentRecordID",
				h.HandleSaveBaseProfile,
			)

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
