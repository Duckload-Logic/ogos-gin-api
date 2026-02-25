package users

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// =============================================
// |                                           |
// |                                           |
// |                                           |
// =============================================

// GetUser
func (r *Repository) GetUser(
	ctx context.Context, userID *int, email *string,
) (*User, error) {
	var user User

	query := fmt.Sprintf(`
		SELECT %s
		FROM users
		WHERE 1=1
	`, database.GetColumns(User{}))

	var args []interface{}

	if userID != nil {
		query += " AND id = ?"
		args = append(args, *userID)
	}

	if email != nil {
		query += " AND email = ?"
		args = append(args, *email)
	}

	query += " LIMIT 1"

	err := r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetRoleByID(
	ctx context.Context, roleID int,
) (*Role, error) {
	var role Role
	query := fmt.Sprintf(`
		SELECT %s
		FROM user_roles
		WHERE id = ?
	`, database.GetColumns(Role{}))
	err := r.db.GetContext(ctx, &role, query, roleID)
	if err != nil {
		log.Printf("Error fetching role with ID %d: %v\n", roleID, err)
		return nil, err
	}
	return &role, nil
}

// =============================================
// |                                           |
// |                                           |
// |                                           |
// =============================================

// CreateUser
func (r *Repository) CreateUser(
	ctx context.Context, user User,
) (int, error) {
	var userID int

	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(User{}, []string{"updated_at"})
		onDuplicateKeyStmt := database.GetOnDuplicateKeyUpdateStatement(User{}, []string{"updated_at"})
		query := fmt.Sprintf(`
			INSERT INTO users (%s)
			VALUES (%s)
			ON DUPLICATE KEY UPDATE %s
		`, cols, vals, onDuplicateKeyStmt)

		result, err := tx.ExecContext(ctx, query,
			user.RoleID,
			user.FirstName,
			user.MiddleName,
			user.LastName,
			user.Email,
			user.PasswordHash,
		)
		if err != nil {
			return err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return err
		}

		userID = int(id)
		return nil
	})

	return userID, err
}
