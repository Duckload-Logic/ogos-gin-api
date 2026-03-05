package apikeys

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, key APIKey) (int, error) {
	query := `
		INSERT INTO api_keys (name, key_hash, key_prefix, scopes, is_active, expires_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query,
		key.Name, key.KeyHash, key.KeyPrefix, key.Scopes, key.IsActive, key.ExpiresAt,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create api key: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return int(id), nil
}

func (r *Repository) GetByHash(ctx context.Context, keyHash string) (*APIKey, error) {
	var key APIKey
	query := `SELECT * FROM api_keys WHERE key_hash = ?`
	err := r.db.GetContext(ctx, &key, query, keyHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get api key by hash: %w", err)
	}
	return &key, nil
}

func (r *Repository) List(ctx context.Context, includeRevoked bool) ([]APIKey, error) {
	var keys []APIKey
	query := `SELECT * FROM api_keys`
	if !includeRevoked {
		query += ` WHERE is_active = 1`
	}
	query += ` ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &keys, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list api keys: %w", err)
	}
	return keys, nil
}

func (r *Repository) Revoke(ctx context.Context, id int) error {
	query := `UPDATE api_keys SET is_active = 0 WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to revoke api key: %w", err)
	}
	return nil
}

func (r *Repository) TouchLastUsed(ctx context.Context, id int) error {
	query := `UPDATE api_keys SET last_used_at = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to update last_used_at: %w", err)
	}
	return nil
}
