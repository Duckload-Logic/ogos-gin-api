package students

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(db *sql.DB, r *gin.RouterGroup, h *Handler) {
	// Root group: /api/v1/students
	studentRoutes := r.Group("/students")
	studentRoutes.Use(middleware.AuthMiddleware())

	// Define lookups
	// userLookup := middleware.OwnershipMiddleware(db, "userID")
	// inventoryRecordLookup := middleware.OwnershipMiddleware(db, "inventoryRecordID")

	lookupRoutes := studentRoutes.Group("/lookups")
	{
		lookupRoutes.GET("/genders", h.HandleGetGenders)
		lookupRoutes.GET("/parental-status-types", h.HandleGetParentalStatusTypes)
		lookupRoutes.GET("/enrollment-reasons", h.HandleGetEnrollmentReasons)
		lookupRoutes.GET("/income-ranges", h.HandleGetIncomeRanges)
		lookupRoutes.GET("/support-types", h.HandleGetStudentSupportTypes)
		lookupRoutes.GET("/support-types/siblings", h.HandleGetSiblingSupportTypes)
		lookupRoutes.GET("/courses", h.HandleGetCourses)
		lookupRoutes.GET("/civil-status-types", h.HandleGetCivilStatusTypes)
		lookupRoutes.GET("/nature-of-residence-types", h.HandleGetNatureOfResidenceTypes)
		lookupRoutes.GET("/student-relationship-types", h.HandleGetStudentRelationshipTypes)
	}

	inventoryRoutes := studentRoutes.Group("/inventory")
	inventoryRoutes.Use(middleware.RoleMiddleware(
		int(constants.StudentRoleID),
		int(constants.CounselorRoleID),
	))
	{
		inventoryRoutes.GET("/records", h.HandleListStudents)
		inventoryRoutes.GET("/records/user/:userId", h.HandleGetStudentIIRByUserID)
		inventoryRoutes.GET("/records/iir/:iirId", h.HandleGetStudentIIRByIIRID)
		inventoryRoutes.GET("/records/iir/:iirId/profile", h.HandleGetStudentProfile)
		inventoryRoutes.GET("/records/iir/:iirId/enrollment-reasons", h.HandleGetStudentEnrollmentReasons)
		inventoryRoutes.GET("/records/iir/:iirId/personal-info", h.HandleGetStudentPersonalInfo)
		inventoryRoutes.GET("/records/iir/:iirId/addresses", h.HandleGetStudentAddresses)
		inventoryRoutes.GET("/records/iir/:iirId/family-background", h.HandleGetStudentFamilyBackground)
		inventoryRoutes.GET("/records/iir/:iirId/related-persons", h.HandleGetStudentRelatedPersons)
		inventoryRoutes.GET("/records/iir/:iirId/education", h.HandleGetEducationalBackground)
		inventoryRoutes.GET("/records/iir/:iirId/finance", h.HandleGetStudentFinancialInfo)
		inventoryRoutes.GET("/records/iir/:iirId/health", h.HandleGetStudentHealthRecord)
		inventoryRoutes.GET("/records/iir/:iirId/consultations", h.HandleGetStudentConsultations)
		inventoryRoutes.GET("/records/iir/:iirId/activities", h.HandleGetStudentActivities)
		inventoryRoutes.GET("/records/iir/:iirId/subject-preferences", h.HandleGetStudentSubjectPreferences)
		inventoryRoutes.GET("/records/iir/:iirId/hobbies", h.HandleGetStudentHobbies)
		inventoryRoutes.GET("/records/iir/:iirId/test-results", h.HandleGetStudentTestResults)
		inventoryRoutes.GET("/records/iir/:iirId/significant-notes", h.HandleGetStudentSignificantNotes)
	}
}
