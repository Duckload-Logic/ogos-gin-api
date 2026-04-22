package users

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// User mappers
func (m UserDB) ToDomain() User {
	return User{
		ID:           m.ID,
		RoleID:       m.RoleID,
		FirstName:    m.FirstName,
		MiddleName:   structs.FromSqlNull(m.MiddleName),
		LastName:     m.LastName,
		SuffixName:   structs.FromSqlNull(m.SuffixName),
		Email:        m.Email,
		PasswordHash: structs.FromSqlNull(m.PasswordHash),
		AuthType:     m.AuthType,
		IsActive:     m.IsActive == 1,
		CreatedAt:    structs.FromSqlNullTime(m.CreatedAt),
		UpdatedAt:    structs.FromSqlNullTime(m.UpdatedAt),
	}
}

func (d User) ToPersistence() UserDB {
	isActive := 0
	if d.IsActive {
		isActive = 1
	}
	return UserDB{
		ID:           d.ID,
		RoleID:       d.RoleID,
		FirstName:    d.FirstName,
		MiddleName:   structs.ToSqlNull(d.MiddleName),
		LastName:     d.LastName,
		SuffixName:   structs.ToSqlNull(d.SuffixName),
		Email:        d.Email,
		PasswordHash: structs.ToSqlNull(d.PasswordHash),
		AuthType:     d.AuthType,
		IsActive:     isActive,
		CreatedAt:    structs.ToSqlNullTime(d.CreatedAt),
		UpdatedAt:    structs.ToSqlNullTime(d.UpdatedAt),
	}
}

// Role mappers
func (m RoleDB) ToDomain() Role {
	return Role(m)
}

func (d Role) ToPersistence() RoleDB {
	return RoleDB(d)
}

// ProfilePicture mappers
func (m ProfilePictureDB) ToDomain() ProfilePicture {
	return ProfilePicture(m)
}

func (d ProfilePicture) ToPersistence() ProfilePictureDB {
	return ProfilePictureDB(d)
}
