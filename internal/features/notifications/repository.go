package notifications

import (
    "context"
    "github.com/jmoiron/sqlx"
)

type Repository struct {
    db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
    return &Repository{db: db}
}   

func (r *Repository) GetByUserID(ctx context.Context, userID string) ([]NotificationModel, error) {
    query := `SELECT id, user_id, title, message, type, is_read, created_at FROM notifications WHERE user_id = ? ORDER BY created_at DESC`
    rows, err := r.db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []NotificationModel
    for rows.Next() {
        var n NotificationModel
        if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.Type, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
        results = append(results, n)
    }
    return results, nil
}

func (r *Repository) MarkAsRead(ctx context.Context, id string) error {
    query := `UPDATE notifications SET is_read = TRUE WHERE id = ?`
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}

func (r *Repository) Create(ctx context.Context, userID int, title, message, notifType string) error {
	query := `
		INSERT INTO notifications (user_id, title, message, type, is_read, created_at)
		VALUES (?, ?, ?, ?, FALSE, NOW())`
	
	_, err := r.db.ExecContext(ctx, query, userID, title, message, notifType)
	return err
}