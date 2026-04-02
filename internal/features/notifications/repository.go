package notifications

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

func (r *Repository) GetByUserID(
	ctx context.Context,
	userID string,
) ([]NotificationModel, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM notifications
		WHERE receiver_id = ?
		ORDER BY created_at DESC
	`, datastore.GetColumns(NotificationModel{}))

	var results []NotificationModel
	err := r.db.SelectContext(ctx, &results, query, userID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get notifications for user %s: %w",
			userID,
			err,
		)
	}

	return results, nil
}

func (r *Repository) GetDB() *sqlx.DB {
	return r.db
}

func (r *Repository) MarkAsRead(
	ctx context.Context,
	tx datastore.DB,
	id string,
) error {
	query := `UPDATE notifications SET is_read = TRUE WHERE id = ?`
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark notification %s as read: %w", id, err)
	}
	return nil
}

func (r *Repository) Create(
	ctx context.Context,
	tx datastore.DB,
	notif *NotificationModel,
) error {
	cols, vals := datastore.GetInsertStatement(notif, []string{"created_at"})
	query := fmt.Sprintf(`
        INSERT INTO notifications (id, %s)
        VALUES (:id, %s)`, cols, vals)

	_, err := tx.NamedExecContext(ctx, query, &notif)
	if err != nil {
		return fmt.Errorf(
			"failed to create notification for user %s: %w",
			notif.ReceiverID.String,
			err,
		)
	}
	return nil
}
