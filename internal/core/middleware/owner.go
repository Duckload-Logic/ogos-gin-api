// owner.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
)

// OwnershipMiddleware - Direct database access version
func OwnershipMiddleware(db *sqlx.DB, paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		loggedInUserID := c.MustGet("userID").(string)
		roleID := c.MustGet("roleID").(int)

		// Allow counselors and super admins to bypass
		if roleID == int(constants.CounselorRoleID) ||
			roleID == int(constants.SuperAdminRoleID) {
			c.Next()
			return
		}

		// For students, check ownership
		if roleID == int(constants.StudentRoleID) {
			paramValue := c.Param(paramName)

			// For email-based params, compare directly
			if paramName == "userID" {
				if paramValue != loggedInUserID {
					c.AbortWithStatusJSON(
						http.StatusForbidden,
						gin.H{"error": "Access denied"},
					)
					return
				}
				c.Next()
				return
			}

			// For int-based params, parse and check ownership
			resourceID := paramValue

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
	db *sqlx.DB, userID string,
	paramName string, resourceID string,
) (bool, error) {
	switch paramName {
	case "iirID":
		// Check if student_record belongs to user
		query := `
			SELECT EXISTS(
				SELECT 1 FROM iir_records
				WHERE id = ? AND user_id = ?
			)`
		var exists bool
		err := db.QueryRow(query, resourceID, userID).Scan(&exists)
		return exists, err
	default:
		return false, nil
	}
}
