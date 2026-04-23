package m2mclients

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

func (r *Repository) getDB(tx datastore.DB) datastore.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *Repository) WithTransaction(
	ctx context.Context,
	fn func(datastore.DB) error,
) error {
	return datastore.RunInTransaction(ctx, r.db, fn)
}

func (r *Repository) Create(
	ctx context.Context,
	tx datastore.DB,
	client M2MClient,
) (int, error) {
	dbModel := MapM2MClientToDB(client)
	cols, vals := datastore.GetInsertStatement(&dbModel, []string{})
	query := fmt.Sprintf(
		`INSERT INTO m2m_clients (%s) VALUES (%s)`,
		cols, vals,
	)
	result, err := r.getDB(tx).NamedExecContext(ctx, query, dbModel)
	if err != nil {
		return 0, fmt.Errorf("failed to create m2m client: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return int(id), nil
}

func (r *Repository) GetByClientID(
	ctx context.Context,
	clientID string,
) (*M2MClient, error) {
	var client M2MClientDB
	query := `SELECT * FROM m2m_clients WHERE client_id = ?`
	err := r.db.GetContext(ctx, &client, query, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get m2m client by client_id: %w", err)
	}
	domainClient := MapM2MClientToDomain(client)
	return &domainClient, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	tx datastore.DB,
	id int,
) (*M2MClient, error) {
	var client M2MClientDB
	query := `SELECT * FROM m2m_clients WHERE id = ?`
	err := r.getDB(tx).GetContext(ctx, &client, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get m2m client by id: %w", err)
	}
	domainClient := MapM2MClientToDomain(client)
	return &domainClient, nil
}

func (r *Repository) GetActiveByUserID(
	ctx context.Context,
	userID string,
) (*M2MClient, error) {
	var client M2MClientDB
	query := `SELECT * FROM m2m_clients WHERE user_id = ? AND is_active = 1 LIMIT 1`
	err := r.db.GetContext(ctx, &client, query, userID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get active m2m client by user_id: %w",
			err,
		)
	}
	domainClient := MapM2MClientToDomain(client)
	return &domainClient, nil
}

func (r *Repository) List(
	ctx context.Context,
	userID string,
	includeRevoked bool,
) ([]M2MClient, error) {
	var clients []M2MClientDB
	query := `SELECT * FROM m2m_clients WHERE 1=1`
	args := []interface{}{}
	if userID != "" {
		query += ` AND user_id = ?`
		args = append(args, userID)
	}
	if !includeRevoked {
		query += ` AND is_active = ?`
		args = append(args, 1)
	}
	query += ` ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &clients, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list m2m clients: %w", err)
	}
	return MapM2MClientsToDomain(clients), nil
}

func (r *Repository) Revoke(
	ctx context.Context,
	tx datastore.DB,
	id int,
) error {
	query := `UPDATE m2m_clients SET is_active = 0 WHERE id = ?`
	_, err := r.getDB(tx).ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to revoke m2m client: %w", err)
	}
	return nil
}

func (r *Repository) UpdateSecret(
	ctx context.Context,
	tx datastore.DB,
	id int,
	secretHash string,
) error {
	query := `UPDATE m2m_clients SET client_secret_hash = ? WHERE id = ?`
	_, err := r.getDB(tx).ExecContext(ctx, query, secretHash, id)
	if err != nil {
		return fmt.Errorf("failed to update client secret: %w", err)
	}
	return nil
}

func (r *Repository) UpdateVerificationStatus(
	ctx context.Context,
	tx datastore.DB,
	id int,
	isVerified bool,
) error {
	query := `UPDATE m2m_clients SET is_verified = ? WHERE id = ?`
	_, err := r.getDB(tx).ExecContext(ctx, query, isVerified, id)
	if err != nil {
		return fmt.Errorf("failed to update verification status: %w", err)
	}
	return nil
}

func (r *Repository) TouchLastUsed(
	ctx context.Context,
	tx datastore.DB,
	id int,
) error {
	query := `UPDATE m2m_clients SET last_used_at = NOW() WHERE id = ?`
	_, err := r.getDB(tx).ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to update last_used_at: %w", err)
	}
	return nil
}
