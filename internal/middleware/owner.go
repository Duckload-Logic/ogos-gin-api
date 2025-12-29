// owner.go
package middleware

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
)

// OwnershipMiddleware - Direct database access version
func OwnershipMiddleware(db *sql.DB, paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		loggedInUserID := c.MustGet("userID").(int)
		roleID := c.MustGet("roleID").(int)

		// Allow counselors and front desk to bypass
		if roleID == int(constants.CounselorRoleID) ||
			roleID == int(constants.FrontDeskRoleID) {
			c.Next()
			return
		}

		// For students, check ownership
		if roleID == int(constants.StudentRoleID) {
			resourceID, err := strconv.Atoi(c.Param(paramName))
			if err != nil {
				c.AbortWithStatusJSON(
					http.StatusBadRequest,
					gin.H{"error": "Invalid resource ID"},
				)
				return
			}

			// Direct database query to check ownership
			owns, err := checkStudentOwnership(
				db, loggedInUserID, paramName, resourceID,
			)
			if err != nil || !owns {
				c.AbortWithStatusJSON(
					http.StatusForbidden,
					gin.H{"error": "Access denied"},
				)
				return
			}
		}

		c.Next()
	}
}

// Direct database query - ONE function to rule them all
func checkStudentOwnership(
	db *sql.DB, userID int,
	paramName string, resourceID int,
) (bool, error) {
	switch paramName {
	case "userID":
		// Simple check: user can only access their own userID
		return userID == resourceID, nil

	case "studentRecordID":
		// Check if student_record belongs to user
		query := `
			SELECT EXISTS(
				SELECT 1 FROM student_records 
				WHERE student_record_id = ? AND user_id = ?
			)`
		var exists bool
		err := db.QueryRow(query, resourceID, userID).Scan(&exists)
		return exists, err
	default:
		// Unknown resource type - deny by default
		return false, nil
	}
}
