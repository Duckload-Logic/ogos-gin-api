package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

// parseUsersFromCSV reads a CSV file and returns a slice of users.User objects.
// Expected Headers: id,email,role_id,first_name,middle_name,last_name,
// suffix_name,password_hash,auth_type,is_active,created_at,updated_at
func parseUsersFromCSV(filePath string) ([]users.User, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// read header
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	var usersList []users.User
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV record: %w", err)
		}

		// map record to User struct
		// 0: id
		// 1: email
		// 2: role_id
		// 3: first_name
		// 4: middle_name
		// 5: last_name
		// 6: suffix_name
		// 7: password_hash
		// 8: auth_type
		// 9: is_active
		// 10: created_at
		// 11: updated_at

		roleID, _ := strconv.Atoi(record[2])
		isActive, _ := strconv.Atoi(record[9])

		user := users.User{
			ID:           record[0],
			Email:        record[1],
			RoleID:       roleID,
			FirstName:    record[3],
			MiddleName:   stringToNullString(record[4]),
			LastName:     record[5],
			SuffixName:   stringToNullString(record[6]),
			PasswordHash: stringToNullString(record[7]),
			AuthType:     record[8],
			IsActive:     isActive,
			CreatedAt:    stringToNullTime(record[10]),
			UpdatedAt:    stringToNullTime(record[11]),
		}

		usersList = append(usersList, user)
	}

	return usersList, nil
}
