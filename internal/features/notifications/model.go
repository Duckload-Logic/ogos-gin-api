package notifications

import "time"

type NotificationModel struct {
	ID        int       `db:"id"         json:"id"`
	UserID    string    `db:"user_id"    json:"userId"`
	Title     string    `db:"title"      json:"title"`
	Message   string    `db:"message"    json:"message"`
	Type      string    `db:"type"       json:"type"`
	IsRead    bool      `db:"is_read"    json:"isRead"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
