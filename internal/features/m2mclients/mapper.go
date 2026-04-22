package m2mclients

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// MapM2MClientToDomain converts DB model to domain model.
func MapM2MClientToDomain(db M2MClientDB) M2MClient {
	return M2MClient{
		ID:                db.ID,
		UserID:            db.UserID,
		ClientName:        db.ClientName,
		ClientID:          db.ClientID,
		ClientSecretHash:  db.ClientSecretHash,
		ClientDescription: db.ClientDescription,
		Scopes:            structs.FromSqlNull(db.Scopes),
		IsActive:          db.IsActive,
		IsVerified:        db.IsVerified,
		LastUsedAt:        structs.FromSqlNullTime(db.LastUsedAt),
		ExpiresAt:         structs.FromSqlNullTime(db.ExpiresAt),
		CreatedAt:         db.CreatedAt,
		UpdatedAt:         db.UpdatedAt,
	}
}

// MapM2MClientToDB converts domain model to DB model.
func MapM2MClientToDB(d M2MClient) M2MClientDB {
	return M2MClientDB{
		ID:                d.ID,
		UserID:            d.UserID,
		ClientName:        d.ClientName,
		ClientID:          d.ClientID,
		ClientSecretHash:  d.ClientSecretHash,
		ClientDescription: d.ClientDescription,
		Scopes:            structs.ToSqlNull(d.Scopes),
		IsActive:          d.IsActive,
		IsVerified:        d.IsVerified,
		LastUsedAt:        structs.ToSqlNullTime(d.LastUsedAt),
		ExpiresAt:         structs.ToSqlNullTime(d.ExpiresAt),
		CreatedAt:         d.CreatedAt,
		UpdatedAt:         d.UpdatedAt,
	}
}

// MapM2MClientsToDomain maps a slice of DB models to domain models.
func MapM2MClientsToDomain(db []M2MClientDB) []M2MClient {
	domain := make([]M2MClient, len(db))
	for i := range db {
		domain[i] = MapM2MClientToDomain(db[i])
	}
	return domain
}
