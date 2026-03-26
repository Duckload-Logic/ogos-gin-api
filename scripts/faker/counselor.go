package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

// createSuperAdmin creates a super admin user.
func createSuperAdmin(index int, password string) {
	email := fmt.Sprintf("superadmin%d@university.edu", index+1)

	user := users.User{
		ID:           uuid.New().String(),
		RoleID:       3, // Super Admin
		FirstName:    "Super",
		MiddleName:   randomMiddleName(),
		LastName:     "Admin",
		Email:        email,
		PasswordHash: sql.NullString{Valid: true, String: password},
		AuthType:     "native",
		IsActive:     1,
	}

	err := usersRepo.CreateUser(context.Background(), db, user)
	if err != nil {
		log.Fatalf("[Seeder] {Create SuperAdmin}: %v", err)
	}

	fmt.Printf(
		"Created super admin: %s %s (%s / SuperAdmin@123)\n",
		user.FirstName,
		user.LastName,
		user.Email,
	)
}

// createCounselor creates a counselor (admin) user and profile.
func createCounselor(index int, password string) {
	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	email := fmt.Sprintf("counselor%d@university.edu", index+1)

	user := users.User{
		ID:           uuid.New().String(),
		RoleID:       2, // Counselor
		FirstName:    firstName,
		MiddleName:   randomMiddleName(),
		LastName:     lastName,
		Email:        email,
		PasswordHash: sql.NullString{Valid: true, String: password},
		AuthType:     "native",
		IsActive:     1,
	}

	err := usersRepo.CreateUser(context.Background(), db, user)
	if err != nil {
		log.Fatalf("[Seeder] {Create Counselor User}: %v", err)
	}

	// For counselor_profiles, we still use raw SQL as there's no repository
	// yet,
	// but we could also use a named exec for better parity.
	_, err = db.Exec(`
		INSERT INTO counselor_profiles (user_id, license_number, specialization, is_available)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			license_number = VALUES(license_number),
			specialization = VALUES(specialization),
			is_available = VALUES(is_available)
	`, user.ID, gofakeit.Regex("[A-Z]{3}-[0-9]{6}"), gofakeit.JobTitle(), true)
	if err != nil {
		log.Fatalf("[Seeder] {Create Counselor Profile}: %v", err)
	}

	fmt.Printf("Created counselor: %s %s\n", firstName, lastName)
}
