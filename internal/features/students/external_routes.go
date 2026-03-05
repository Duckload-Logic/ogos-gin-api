package students

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

// RegisterExternalRoutes exposes read-only student/IIR endpoints
// for external systems authenticated via API key.
// Mounted at: /api/v1/external/students
func RegisterExternalRoutes(rg *gin.RouterGroup, h *Handler, validator middleware.APIKeyValidator) {
	ext := rg.Group("/external/students")
	ext.Use(middleware.APIKeyMiddleware(validator))

	// Lookups
	lookups := ext.Group("/lookups")
	{
		lookups.GET("/genders", h.HandleGetGenders)
		lookups.GET("/religions", h.HandleGetReligions)
		lookups.GET("/parental-status-types", h.HandleGetParentalStatusTypes)
		lookups.GET("/enrollment-reasons", h.HandleGetEnrollmentReasons)
		lookups.GET("/income-ranges", h.HandleGetIncomeRanges)
		lookups.GET("/support-types", h.HandleGetStudentSupportTypes)
		lookups.GET("/support-types/siblings", h.HandleGetSiblingSupportTypes)
		lookups.GET("/courses", h.HandleGetCourses)
		lookups.GET("/civil-statuses", h.HandleGetCivilStatusTypes)
		lookups.GET("/nature-of-residence-types", h.HandleGetNatureOfResidenceTypes)
		lookups.GET("/student-relationship-types", h.HandleGetStudentRelationshipTypes)
	}

	// Inventory (read-only)
	inventory := ext.Group("/inventory")
	{
		inventory.GET("/records", h.HandleListStudents)
		inventory.GET("/records/iir/:iirID", h.HandleGetStudentIIRByIIRID)
		inventory.GET("/records/iir/:iirID/profile", h.HandleGetStudentProfile)
		inventory.GET("/records/iir/:iirID/basic-info", h.HandleGetStudentBasicInfo)
		inventory.GET("/records/iir/:iirID/enrollment-reasons", h.HandleGetStudentEnrollmentReasons)
		inventory.GET("/records/iir/:iirID/personal-info", h.HandleGetStudentPersonalInfo)
		inventory.GET("/records/iir/:iirID/addresses", h.HandleGetStudentAddresses)
		inventory.GET("/records/iir/:iirID/family-background", h.HandleGetStudentFamilyBackground)
		inventory.GET("/records/iir/:iirID/related-persons", h.HandleGetStudentRelatedPersons)
		inventory.GET("/records/iir/:iirID/education", h.HandleGetEducationalBackground)
		inventory.GET("/records/iir/:iirID/finance", h.HandleGetStudentFinancialInfo)
		inventory.GET("/records/iir/:iirID/health", h.HandleGetStudentHealthRecord)
		inventory.GET("/records/iir/:iirID/consultations", h.HandleGetStudentConsultations)
		inventory.GET("/records/iir/:iirID/activities", h.HandleGetStudentActivities)
		inventory.GET("/records/iir/:iirID/subject-preferences", h.HandleGetStudentSubjectPreferences)
		inventory.GET("/records/iir/:iirID/hobbies", h.HandleGetStudentHobbies)
		inventory.GET("/records/iir/:iirID/test-results", h.HandleGetStudentTestResults)
		inventory.GET("/records/iir/:iirID/significant-notes", h.HandleGetStudentSignificantNotes)
	}
}
