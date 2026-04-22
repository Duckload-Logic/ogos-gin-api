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

func NewRepository(db *sqlx.DB) RepositoryInterface {
	return &Repository{db: db}
}

func (r *Repository) GetDB() *sqlx.DB {
	return r.db
}

func (r *Repository) WithTransaction(
	ctx context.Context,
	fn func(datastore.DB) error,
) error {
	return datastore.RunInTransaction(ctx, r.db, fn)
}

// =============================================
// |                                           |
// |                                           |
// |                                           |
// =============================================

// GetUserByID fetches a user by their ID and maps to Domain.
func (r *Repository) GetUserByID(
	ctx context.Context,
	userID string,
) (*User, error) {
	var dbModel UserDB

	query := fmt.Sprintf(`
		SELECT %s
		FROM users
		WHERE id = ?
		LIMIT 1
	`, datastore.GetColumns(UserDB{}))

	err := r.db.GetContext(ctx, &dbModel, query, userID)
	if err != nil {
		return nil, err
	}

	domainModel := dbModel.ToDomain()
	return &domainModel, nil
}

func (r *Repository) CheckUserWhitelist(
	ctx context.Context,
	email string,
) (int, error) {
	var roleID int
	query := `SELECT role_id FROM whitelists WHERE email = ?`
	err := r.db.GetContext(ctx, &roleID, query, email)
	return roleID, err
}

func (r *Repository) GetRoleByID(
	ctx context.Context,
	roleID int,
) (*Role, error) {
	var dbModel RoleDB
	query := fmt.Sprintf(`
		SELECT %s
		FROM user_roles
		WHERE id = ?
	`, datastore.GetColumns(RoleDB{}))

	err := r.db.GetContext(ctx, &dbModel, query, roleID)
	if err != nil {
		return nil, err
	}

	domainModel := dbModel.ToDomain()
	return &domainModel, nil
}

func (r *Repository) GetUserByEmail(
	ctx context.Context,
	email, authType string,
) (*User, error) {
	var dbModel UserDB

	query := fmt.Sprintf(`
		SELECT %s
		FROM users
		WHERE email = ? AND auth_type = ?
		LIMIT 1
	`, datastore.GetColumns(UserDB{}))

	err := r.db.GetContext(ctx, &dbModel, query, email, authType)
	if err != nil {
		return nil, err
	}

	domainModel := dbModel.ToDomain()
	return &domainModel, nil
}

// =============================================
// |                                           |
// |                                           |
// |                                           |
// =============================================

// CreateUser inserts or updates a user using the persistence model.
func (r *Repository) CreateUser(
	ctx context.Context,
	tx datastore.DB,
	user User,
) error {
	dbModel := user.ToPersistence()
	exclude := []string{"updated_at"}
	cols, vals := datastore.GetInsertStatement(UserDB{}, exclude)
	onDuplicateKeyStmt := datastore.GetOnDuplicateKeyUpdateStatement(
		UserDB{},
		exclude,
	)
	query := fmt.Sprintf(`
			INSERT INTO users (id, %s)
			VALUES (:id, %s)
			ON DUPLICATE KEY UPDATE %s
		`, cols, vals, onDuplicateKeyStmt)

	_, err := tx.NamedExecContext(ctx, query, dbModel)
	return err
}

func (r *Repository) PostProfilePicture(
	ctx context.Context,
	tx datastore.DB,
	userID string,
	fileID string,
) error {
	dbModel := ProfilePicture{
		UserID: userID,
		FileID: fileID,
	}.ToPersistence()

	query := `INSERT INTO profile_pictures (user_id, file_id) 
			  VALUES (:user_id, :file_id)`
	_, err := tx.NamedExecContext(ctx, query, &dbModel)
	if err != nil {
		return fmt.Errorf("failed to post profile picture: %w", err)
	}

	return nil
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

func (r *Repository) ListUsers(
	ctx context.Context,
	params ListUsersParams,
) ([]User, int, error) {
	var dbModels []UserDB
	var total int

	baseQuery := `FROM users WHERE 1=1`
	args := []interface{}{}

	if params.RoleID > 0 {
		baseQuery += ` AND role_id = ?`
		args = append(args, params.RoleID)
	}

	if params.Search != "" {
		baseQuery += ` AND (email LIKE ? OR first_name LIKE ? OR last_name LIKE ?)`
		like := "%" + params.Search + "%"
		args = append(args, like, like, like)
	}

	if params.Active != nil {
		baseQuery += ` AND is_active = ?`
		activeVal := 0
		if *params.Active {
			activeVal = 1
		}
		args = append(args, activeVal)
	}

	// Count total
	countQuery := `SELECT COUNT(*) ` + baseQuery
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated users
	selectQuery := fmt.Sprintf(`SELECT %s `, datastore.GetColumns(UserDB{})) +
		baseQuery + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`

	limit := params.PageSize
	offset := (params.Page - 1) * params.PageSize
	args = append(args, limit, offset)

	err = r.db.SelectContext(ctx, &dbModels, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	var users []User
	for _, m := range dbModels {
		users = append(users, m.ToDomain())
	}

	return users, total, nil
}

func (r *Repository) GetRoleDistribution(
	ctx context.Context,
) ([]RoleDistributionDTO, error) {
	query := `
		SELECT r.name as role_name, COUNT(u.id) as count
		FROM user_roles r
		LEFT JOIN users u ON u.role_id = r.id
		GROUP BY r.name
	`

	var distribution []RoleDistributionDTO
	err := r.db.SelectContext(ctx, &distribution, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get role distribution: %w", err)
	}

	return distribution, nil
}
