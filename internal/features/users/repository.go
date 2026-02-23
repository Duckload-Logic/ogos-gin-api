package users

import (
	"context"

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

	query := `
		SELECT
			id, role_id,
			first_name, middle_name,
			last_name, email,
			password_hash, created_at,
			updated_at
		FROM users
		WHERE 1=1
	`

	var args []interface{}

	if userID != nil {
		query += " AND id = ?"
		args = append(args, *userID)
	}

	if email != nil {
		query += " AND email = ?"
		args = append(args, *email)
	}

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.RoleID,
		&user.FirstName,
		&user.MiddleName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
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
) (int, error) {
	var userID int

	err := database.RunInTransaction(ctx, r.db, func(tx *sqlx.Tx) error {
		query := `
			INSERT INTO users (
				role_id, first_name,
				middle_name, last_name,
				email, password_hash
			)
			VALUES (?,?,?,?,?,?)
		`

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
