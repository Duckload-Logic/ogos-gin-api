package students

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/middleware"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

func RegisterRoutes(
	db *sqlx.DB,
	rg *gin.RouterGroup,
	h *Handler,
	redis *datastore.RedisClient,
) {
	// Root group: /api/v1/students
	routes := rg.Group("/students")
	routes.Use(middleware.AuthMiddleware(redis))
	routes.Use(middleware.HydrateStudentContext(db))

	// Define lookups
	userResourceLookup := middleware.OwnershipMiddleware(db, "userID")
	iirResourceLookup := middleware.OwnershipMiddleware(db, "iirID")
	// inventoryRecordLookup := middleware.OwnershipMiddleware(db, "inventoryRecordID")

	lookupRoutes := routes.Group("/lookups")
	{
		lookupRoutes.GET("/genders", h.GetGenders)
		lookupRoutes.GET("/religions", h.GetReligions)
		lookupRoutes.GET("/parental-status-types", h.GetParentalStatusTypes)
		lookupRoutes.GET("/enrollment-reasons", h.GetEnrollmentReasons)
		lookupRoutes.GET("/income-ranges", h.GetIncomeRanges)
		lookupRoutes.GET("/activity-options", h.GetActivityOptions)
		lookupRoutes.GET("/support-types", h.GetStudentSupportTypes)
		lookupRoutes.GET("/support-types/siblings", h.GetSiblingSupportTypes)
		lookupRoutes.GET("/courses", h.GetCourses)
		lookupRoutes.GET("/civil-statuses", h.GetCivilStatusTypes)
		lookupRoutes.GET(
			"/nature-of-residence-types",
			h.GetNatureOfResidenceTypes,
		)
		lookupRoutes.GET(
			"/student-relationship-types",
			h.GetStudentRelationshipTypes,
		)
	}

	inventoryRoutes := routes.Group("/inventory")

	counselorRoutes := inventoryRoutes.Group("/")
	counselorRoutes.Use(middleware.RoleMiddleware(
		int(constants.CounselorRoleID),
	))
	{
		counselorRoutes.GET("/records", h.GetStudentList)
	}

	userRoutes := inventoryRoutes.Group("/")
	userRoutes.Use(middleware.RoleMiddleware(
		int(constants.StudentRoleID),
		int(constants.CounselorRoleID),
	))
	{
		userRoutes.GET(
			"/records/user/:userID",
			userResourceLookup,
			h.GetStudentIIRByUserID,
		)
		userRoutes.GET(
			"/records/iir/:iirID",
			iirResourceLookup,
			h.GetStudentIIRByIIRID,
		)
		userRoutes.GET(
			"/records/iir/:iirID/profile",
			iirResourceLookup,
			h.GetStudentProfile,
		)
		userRoutes.GET(
			"/records/iir/:iirID/basic-info",
			iirResourceLookup,
			h.GetStudentBasicInfo,
		)
		userRoutes.GET(
			"/records/iir/:iirID/enrollment-reasons",
			iirResourceLookup,
			h.GetStudentEnrollmentReasons,
		)
		userRoutes.GET(
			"/records/iir/:iirID/personal-info",
			iirResourceLookup,
			h.GetStudentPersonalInfo,
		)
		userRoutes.GET(
			"/records/iir/:iirID/addresses",
			iirResourceLookup,
			h.GetStudentAddresses,
		)
		userRoutes.GET(
			"/records/iir/:iirID/family-background",
			iirResourceLookup,
			h.GetStudentFamilyBackground,
		)
		userRoutes.GET(
			"/records/iir/:iirID/related-persons",
			iirResourceLookup,
			h.GetStudentRelatedPersons,
		)
		userRoutes.GET(
			"/records/iir/:iirID/education",
			iirResourceLookup,
			h.GetEducationalBackground,
		)
		userRoutes.GET(
			"/records/iir/:iirID/finance",
			iirResourceLookup,
			h.GetStudentFinancialInfo,
		)
		userRoutes.GET(
			"/records/iir/:iirID/health",
			iirResourceLookup,
			h.GetStudentHealthRecord,
		)
		userRoutes.GET(
			"/records/iir/:iirID/consultations",
			iirResourceLookup,
			h.GetStudentConsultations,
		)
		userRoutes.GET(
			"/records/iir/:iirID/activities",
			iirResourceLookup,
			h.GetStudentActivities,
		)
		userRoutes.GET(
			"/records/iir/:iirID/subject-preferences",
			iirResourceLookup,
			h.GetStudentSubjectPreferences,
		)
		userRoutes.GET(
			"/records/iir/:iirID/hobbies",
			iirResourceLookup,
			h.GetStudentHobbies,
		)
		userRoutes.GET(
			"/records/iir/:iirID/test-results",
			iirResourceLookup,
			h.GetStudentTestResults,
		)
	}

	studentRoutes := inventoryRoutes.Group("/")
	studentRoutes.Use(middleware.RoleMiddleware(
		int(constants.StudentRoleID),
	))
	{
		studentRoutes.GET("/records/iir/draft", h.GetIIRDraft)
		studentRoutes.POST("/records/iir/draft", h.PostIIRDraft)

		studentRoutes.POST("/records/iir", h.PostIIR)
	}
}
