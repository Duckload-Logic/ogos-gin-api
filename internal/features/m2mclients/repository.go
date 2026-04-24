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
) error {
	exclude := []string{"updated_at"}
	cols, vals := datastore.GetInsertStatement(M2MClient{}, exclude)
	query := fmt.Sprintf(`INSERT INTO m2m_clients (%s) VALUES (%s)`, cols, vals)

	_, err := tx.NamedExecContext(ctx, query, client)
	return err
}

func (r *Repository) GetByClientID(
	ctx context.Context,
	clientID string,
) (*M2MClient, error) {
	var client M2MClient
	query := fmt.Sprintf(`
		SELECT %s
		FROM m2m_clients
		WHERE client_id = ? AND is_active = 1
		LIMIT 1
	`, datastore.GetColumns(M2MClient{}))

	err := r.db.GetContext(ctx, &client, query, clientID)
	return &client, err
}

func (r *Repository) GetActiveByUserID(
	ctx context.Context,
	userID string,
) (*M2MClient, error) {
	var client M2MClient
	query := fmt.Sprintf(`
		SELECT %s
		FROM m2m_clients
		WHERE user_id = ? AND is_active = 1
		LIMIT 1
	`, datastore.GetColumns(M2MClient{}))

	err := r.db.GetContext(ctx, &client, query, userID)
	return &client, err
}

func (r *Repository) DeactivateAllForUser(
	ctx context.Context,
	tx datastore.DB,
	userID string,
) error {
	query := `UPDATE m2m_clients SET is_active = 0 WHERE user_id = ?`
	_, err := tx.ExecContext(ctx, query, userID)
	return err
}

func (r *Repository) ListActive(
	ctx context.Context,
) ([]M2MClient, error) {
	var clients []M2MClient
	query := fmt.Sprintf(`
		SELECT %s
		FROM m2m_clients
		WHERE is_active = 1
	`, datastore.GetColumns(M2MClient{}))

	err := r.db.SelectContext(ctx, &clients, query)
	return clients, err
}
