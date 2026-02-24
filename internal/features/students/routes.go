package students

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(db *sqlx.DB, r *gin.RouterGroup, h *Handler) {
	// Root group: /api/v1/students
	studentRoutes := r.Group("/students")
	studentRoutes.Use(middleware.AuthMiddleware())

	// Define lookups
	userResourceLookup := middleware.OwnershipMiddleware(db, "userID")
	iirResourceLookup := middleware.OwnershipMiddleware(db, "iirID")
	// inventoryRecordLookup := middleware.OwnershipMiddleware(db, "inventoryRecordID")

	lookupRoutes := studentRoutes.Group("/lookups")
	{
		lookupRoutes.GET("/genders", h.HandleGetGenders)
		lookupRoutes.GET("/religions", h.HandleGetReligions)
		lookupRoutes.GET("/parental-status-types", h.HandleGetParentalStatusTypes)
		lookupRoutes.GET("/enrollment-reasons", h.HandleGetEnrollmentReasons)
		lookupRoutes.GET("/income-ranges", h.HandleGetIncomeRanges)
		lookupRoutes.GET("/support-types", h.HandleGetStudentSupportTypes)
		lookupRoutes.GET("/support-types/siblings", h.HandleGetSiblingSupportTypes)
		lookupRoutes.GET("/courses", h.HandleGetCourses)
		lookupRoutes.GET("/civil-statuses", h.HandleGetCivilStatusTypes)
		lookupRoutes.GET("/nature-of-residence-types", h.HandleGetNatureOfResidenceTypes)
		lookupRoutes.GET("/student-relationship-types", h.HandleGetStudentRelationshipTypes)
	}

	inventoryRoutes := studentRoutes.Group("/inventory")

	counselorRoutes := inventoryRoutes.Group("/")
	counselorRoutes.Use(middleware.RoleMiddleware(
		int(constants.CounselorRoleID),
	))
	{
		counselorRoutes.GET("/records", h.HandleListStudents)
	}

	userRoutes := inventoryRoutes.Group("/")
	userRoutes.Use(middleware.RoleMiddleware(
		int(constants.StudentRoleID),
		int(constants.CounselorRoleID),
	))
	{
		userRoutes.GET("/records/user/:userID", userResourceLookup, h.HandleGetStudentIIRByUserID)
		userRoutes.GET("/records/iir/:iirID", iirResourceLookup, h.HandleGetStudentIIRByIIRID)
		userRoutes.GET("/records/iir/:iirID/profile", iirResourceLookup, h.HandleGetStudentProfile)
		userRoutes.GET("/records/iir/:iirID/basic-info", iirResourceLookup, h.HandleGetStudentBasicInfo)
		userRoutes.GET("/records/iir/:iirID/enrollment-reasons", iirResourceLookup, h.HandleGetStudentEnrollmentReasons)
		userRoutes.GET("/records/iir/:iirID/personal-info", iirResourceLookup, h.HandleGetStudentPersonalInfo)
		userRoutes.GET("/records/iir/:iirID/addresses", iirResourceLookup, h.HandleGetStudentAddresses)
		userRoutes.GET("/records/iir/:iirID/family-background", iirResourceLookup, h.HandleGetStudentFamilyBackground)
		userRoutes.GET("/records/iir/:iirID/related-persons", iirResourceLookup, h.HandleGetStudentRelatedPersons)
		userRoutes.GET("/records/iir/:iirID/education", iirResourceLookup, h.HandleGetEducationalBackground)
		userRoutes.GET("/records/iir/:iirID/finance", iirResourceLookup, h.HandleGetStudentFinancialInfo)
		userRoutes.GET("/records/iir/:iirID/health", iirResourceLookup, h.HandleGetStudentHealthRecord)
		userRoutes.GET("/records/iir/:iirID/consultations", iirResourceLookup, h.HandleGetStudentConsultations)
		userRoutes.GET("/records/iir/:iirID/activities", iirResourceLookup, h.HandleGetStudentActivities)
		userRoutes.GET("/records/iir/:iirID/subject-preferences", iirResourceLookup, h.HandleGetStudentSubjectPreferences)
		userRoutes.GET("/records/iir/:iirID/hobbies", iirResourceLookup, h.HandleGetStudentHobbies)
		userRoutes.GET("/records/iir/:iirID/test-results", iirResourceLookup, h.HandleGetStudentTestResults)
		userRoutes.GET("/records/iir/:iirID/significant-notes", iirResourceLookup, h.HandleGetStudentSignificantNotes)
	}

	postRoutes := inventoryRoutes.Group("/")
	postRoutes.Use(middleware.RoleMiddleware(
		int(constants.StudentRoleID),
	))
	{
		postRoutes.POST("/records/iir", h.HandleSubmitIIR)
	}
}
