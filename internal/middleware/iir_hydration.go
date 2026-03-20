package middleware

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
)

// HydrateStudentContext extracts student IIR ID from database and
// sets it in the Gin context. Only applies to Student role users.
// Day One students (no IIR record) can still proceed.
func HydrateStudentContext(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract userID and roleID from context (set by
		// AuthMiddleware)
		userIDVal, exists := c.Get("userID")
		if !exists {
			c.Next()
			return
		}

		roleIDVal, exists := c.Get("roleID")
		if !exists {
			c.Next()
			return
		}

		userID, ok := userIDVal.(int)
		if !ok {
			c.Next()
			return
		}

		roleID, ok := roleIDVal.(int)
		if !ok {
			c.Next()
			return
		}

		// Only hydrate for Student role
		if constants.RoleID(roleID) !=
			constants.StudentRoleID {
			c.Next()
			return
		}

		// Query iir_records table to find IIR ID by user_id
		var iirID int
		err := db.QueryRow(`
			SELECT id FROM iir_records WHERE user_id = ?
		`, userID).Scan(&iirID)

		if err == sql.ErrNoRows {
			// Day One student - no IIR record yet
			c.Next()
			return
		}

		if err != nil {
			log.Printf(
				"[HydrateStudentContext] {Database Query IIR "+
					"Lookup}: %v",
				err)
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": "Internal server error"},
			)
			return
		}

		// Set IIR ID in context for downstream handlers
		c.Set("iirID", iirID)
		c.Next()
	}
}
