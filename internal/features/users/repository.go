package users

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetDB() *sqlx.DB {
	return r.db
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
	`, datastore.GetColumns(User{}))

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
	`, datastore.GetColumns(Role{}))
	err := r.db.GetContext(ctx, &role, query, roleID)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repository) GetUserByEmail(
	ctx context.Context, email string, authType string,
) (*User, error) {
	var user User

	query := fmt.Sprintf(`
		SELECT %s
		FROM users
		WHERE email = ? AND auth_type = ?
		LIMIT 1
	`, datastore.GetColumns(User{}))

	err := r.db.GetContext(ctx, &user, query, email, authType)
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
	ctx context.Context,
	tx datastore.DB,
	user User,
) error {
	// id is the primary key, we should NOT update it on duplicate
	// password_hash might be empty for IDP users, we don't want to overwrite it
	exclude := []string{"updated_at"}
	cols, vals := datastore.GetInsertStatement(User{}, exclude)
	onDuplicateKeyStmt := datastore.GetOnDuplicateKeyUpdateStatement(
		User{},
		exclude,
	)
	query := fmt.Sprintf(`
			INSERT INTO users (id, %s)
			VALUES (:id, %s)
			ON DUPLICATE KEY UPDATE %s
		`, cols, vals, onDuplicateKeyStmt)

	_, err := tx.NamedExecContext(ctx, query, user)
	return err
}

func (r *Repository) BlockUser(
	ctx context.Context,
	tx datastore.DB,
	userID string,
) error {
	query := `UPDATE users SET is_active = 0 WHERE id = ?`
	_, err := tx.ExecContext(ctx, query, userID)
	return err
}

func (r *Repository) UnblockUser(
	ctx context.Context,
	tx datastore.DB,
	userID string,
) error {
	query := `UPDATE users SET is_active = 1 WHERE id = ?`
	_, err := tx.ExecContext(ctx, query, userID)
	return err
}

func (r *Repository) GetUserIDsByRole(
	ctx context.Context,
	roleID int,
) ([]string, error) {
	var userIDs []string
	query := `SELECT id FROM users WHERE role_id = ?`
	err := r.db.SelectContext(ctx, &userIDs, query, roleID)
	return userIDs, err
}
