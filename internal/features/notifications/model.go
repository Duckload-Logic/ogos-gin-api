package notifications

import "time"

type NotificationModel struct {
    ID        int       `db:"id"`
    UserID    int       `db:"user_id"`
    Title     string    `db:"title"`
    Message   string    `db:"message"`
    Type      string    `db:"type"`
    IsRead    bool      `db:"is_read"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}