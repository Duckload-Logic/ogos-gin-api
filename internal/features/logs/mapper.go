package logs

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// MapSystemLogToDomain converts DB model to domain model.
func MapSystemLogToDomain(db SystemLogDB) SystemLog {
	return SystemLog{
		ID:          db.ID,
		Level:       db.Level,
		Category:    db.Category,
		Action:      db.Action,
		Message:     db.Message,
		UserID:      structs.FromSqlNull(db.UserID),
		TargetID:    structs.FromSqlNull(db.TargetID),
		TargetType:  structs.FromSqlNull(db.TargetType),
		UserEmail:   structs.FromSqlNull(db.UserEmail),
		TargetEmail: structs.FromSqlNull(db.TargetEmail),
		IPAddress:   structs.FromSqlNull(db.IPAddress),
		UserAgent:   structs.FromSqlNull(db.UserAgent),
		Metadata:    structs.FromSqlNull(db.Metadata),
		TraceID:     structs.FromSqlNull(db.TraceID),
		CreatedAt:   db.CreatedAt,
	}
}

// MapSystemLogToDB converts domain model to DB model.
func MapSystemLogToDB(d SystemLog) SystemLogDB {
	return SystemLogDB{
		ID:          d.ID,
		Level:       d.Level,
		Category:    d.Category,
		Action:      d.Action,
		Message:     d.Message,
		UserID:      structs.ToSqlNull(d.UserID),
		TargetID:    structs.ToSqlNull(d.TargetID),
		TargetType:  structs.ToSqlNull(d.TargetType),
		UserEmail:   structs.ToSqlNull(d.UserEmail),
		TargetEmail: structs.ToSqlNull(d.TargetEmail),
		IPAddress:   structs.ToSqlNull(d.IPAddress),
		UserAgent:   structs.ToSqlNull(d.UserAgent),
		Metadata:    structs.ToSqlNull(d.Metadata),
		TraceID:     structs.ToSqlNull(d.TraceID),
		CreatedAt:   d.CreatedAt,
	}
}

// MapSystemLogsToDomain maps a slice of DB models to domain models.
func MapSystemLogsToDomain(db []SystemLogDB) []SystemLog {
	domain := make([]SystemLog, len(db))
	for i := range db {
		domain[i] = MapSystemLogToDomain(db[i])
	}
	return domain
}
