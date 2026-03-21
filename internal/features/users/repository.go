package users

import (
	"context"
	"fmt"

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
func (r *Repository) GetUserByID(
	ctx context.Context, userID string,
) (*User, error) {
	var user User

	query := fmt.Sprintf(`
		SELECT %s
		FROM users
		WHERE id = ?
		LIMIT 1
	`, database.GetColumns(User{}))

	err := r.db.GetContext(ctx, &user, query, userID)
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
		return nil, err
	}
	return &role, nil
}

func (r *Repository) GetUserByEmail(
	ctx context.Context, email string,
) (*User, error) {
	var user User

	query := fmt.Sprintf(`
		SELECT %s
		FROM users
		WHERE email = ?
		LIMIT 1
	`, database.GetColumns(User{}))

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// =============================================
// |                                           |
// |                                           |
// |                                           |
// =============================================

// CreateUser
func (r *Repository) CreateUser(
	ctx context.Context, user User,
) error {
	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		cols, vals := database.GetInsertStatement(User{}, []string{"updated_at"})
		onDuplicateKeyStmt := database.GetOnDuplicateKeyUpdateStatement(User{}, []string{"updated_at"})
		query := fmt.Sprintf(`
			INSERT INTO users (%s)
			VALUES (%s)
			ON DUPLICATE KEY UPDATE %s
		`, cols, vals, onDuplicateKeyStmt)

		_, err := tx.ExecContext(ctx, query,
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

		return nil
	})

	return err
}
