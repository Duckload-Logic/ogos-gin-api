package notifications

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// MapNotificationToDomain converts DB model to domain model.
func MapNotificationToDomain(db NotificationDB) Notification {
	return Notification{
		ID:         db.ID,
		ReceiverID: structs.FromSqlNull(db.ReceiverID),
		ActorID:    structs.FromSqlNull(db.ActorID),
		TargetID:   structs.FromSqlNull(db.TargetID),
		TargetType: structs.FromSqlNull(db.TargetType),
		Title:      db.Title,
		Message:    db.Message,
		Type:       db.Type,
		IsRead:     db.IsRead,
		CreatedAt:  db.CreatedAt,
		UpdatedAt:  db.UpdatedAt,
	}
}

// MapNotificationToDB converts domain model to DB model.
func MapNotificationToDB(d Notification) NotificationDB {
	return NotificationDB{
		ID:         d.ID,
		ReceiverID: structs.ToSqlNull(d.ReceiverID),
		ActorID:    structs.ToSqlNull(d.ActorID),
		TargetID:   structs.ToSqlNull(d.TargetID),
		TargetType: structs.ToSqlNull(d.TargetType),
		Title:      d.Title,
		Message:    d.Message,
		Type:       d.Type,
		IsRead:     d.IsRead,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
	}
}

// MapNotificationsToDomain maps a slice of DB models to domain models.
func MapNotificationsToDomain(db []NotificationDB) []Notification {
	domain := make([]Notification, len(db))
	for i := range db {
		domain[i] = MapNotificationToDomain(db[i])
	}
	return domain
}
