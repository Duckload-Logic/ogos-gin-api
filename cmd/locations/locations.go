package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// JSONStructure represents the ph_locations.json format
// { "REGION": { "CITY": ["barangay1", "barangay2", ...], ... }, ... }
type JSONStructure map[string]map[string][]string

// AddressSeeder handles seeding of address data from JSON
type AddressSeeder struct {
	db         *sqlx.DB
	regionsMap sync.Map // thread-safe map for region name -> ID
	citiesMap  sync.Map // thread-safe map for city name -> ID
	mu         *sync.Mutex
}

var db *sqlx.DB

func NewAddressSeeder(db *sqlx.DB) *AddressSeeder {
	return &AddressSeeder{
		db: db,
		mu: &sync.Mutex{},
	}
}

func main() {
	_ = godotenv.Load()
	dsn := buildDSNFromEnv()

	var err error
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}
	defer db.Close()

	seeder := NewAddressSeeder(db)

	// Try to find ph_locations.json in common locations
	jsonFile := "ph_locations.json"
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		jsonFile = "cmd/locations/ph_locations.json"
	}

	if err := seeder.SeedAddresses(jsonFile); err != nil {
		log.Fatal("failed to seed addresses:", err)
	}
	fmt.Println("✓ Address seeding completed successfully.")
}

// SeedAddresses loads JSON file and seeds all address data to database
func (s *AddressSeeder) SeedAddresses(jsonFile string) error {
	fmt.Println("Starting address seeding from JSON...")
	start := time.Now()

	// Clear existing data
	fmt.Println("Clearing existing address data...")
	if err := s.clearData(); err != nil {
		return fmt.Errorf("failed to clear existing data: %w", err)
	}

	// Load JSON file
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	// Parse JSON
	var locations JSONStructure
	if err := json.Unmarshal(data, &locations); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Step 1: Seed regions
	fmt.Println("Seeding regions...")
	regionCount, err := s.seedRegions(locations)
	if err != nil {
		return fmt.Errorf("failed to seed regions: %w", err)
	}

	// Step 2: Seed cities
	fmt.Println("Seeding cities...")
	cityCount, err := s.seedCities(locations)
	if err != nil {
		return fmt.Errorf("failed to seed cities: %w", err)
	}

	// Step 3: Seed barangays
	fmt.Println("Seeding barangays...")
	barangayCount, err := s.seedBarangays(locations)
	if err != nil {
		return fmt.Errorf("failed to seed barangays: %w", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Address seeding completed in %v\n", elapsed)
	fmt.Printf("Total: %d regions, %d cities, %d barangays\n", regionCount, cityCount, barangayCount)

	return nil
}

// clearData truncates all address tables
func (s *AddressSeeder) clearData() error {
	// Disable foreign key checks temporarily
	if _, err := s.db.Exec("SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		return err
	}
	defer s.db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	// Truncate tables
	for _, table := range []string{"barangays", "cities", "regions"} {
		if _, err := s.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table)); err != nil {
			return fmt.Errorf("failed to truncate %s: %w", table, err)
		}
	}

	return nil
}

// seedRegions extracts and inserts regions from JSON into database
func (s *AddressSeeder) seedRegions(locations JSONStructure) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(locations) == 0 {
		return 0, nil
	}

	// Batch insert regions (region names are keys in JSON)
	placeholders := make([]string, 0, len(locations))
	args := make([]interface{}, 0, len(locations)*2)

	for regionName := range locations {
		// Use region name as both code and name (JSON doesn't have separate codes)
		placeholders = append(placeholders, "(?)")
		args = append(args, regionName)
	}

	query := "INSERT INTO regions (name) VALUES " + strings.Join(placeholders, ", ") +
		" ON DUPLICATE KEY UPDATE name=VALUES(name)"

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to batch insert regions: %w", err)
	}

	// Load all regions into map for city seeding
	var regions []struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}
	err = s.db.Select(&regions, "SELECT id, name FROM regions")
	if err != nil {
		return 0, fmt.Errorf("failed to load regions: %w", err)
	}

	for _, region := range regions {
		s.regionsMap.Store(region.Name, int64(region.ID))
	}

	return len(locations), nil
}

// seedCities extracts and inserts cities from JSON into database
func (s *AddressSeeder) seedCities(locations JSONStructure) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Prepare batch insert (city names are keys in nested JSON)
	placeholders := make([]string, 0)
	args := make([]interface{}, 0)
	cityCount := 0

	for regionName, cities := range locations {
		regionID, ok := s.regionsMap.Load(regionName)
		if !ok {
			// Skip cities without matched region
			continue
		}

		for cityName := range cities {
			placeholders = append(placeholders, "(?, ?)")
			args = append(args, cityName, regionID) // Use city name as both code and name
			cityCount++
		}
	}

	if cityCount == 0 {
		return 0, nil
	}

	query := "INSERT INTO cities (name, region_id) VALUES " + strings.Join(placeholders, ", ") +
		" ON DUPLICATE KEY UPDATE name=VALUES(name)"

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to batch insert cities: %w", err)
	}

	// Load all cities into map for barangay seeding
	var cities []struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}
	err = s.db.Select(&cities, "SELECT id, name FROM cities")
	if err != nil {
		return 0, fmt.Errorf("failed to load cities: %w", err)
	}

	for _, city := range cities {
		s.citiesMap.Store(city.Name, int64(city.ID))
	}

	return cityCount, nil
}

// seedBarangays extracts and inserts barangays from JSON into database
func (s *AddressSeeder) seedBarangays(locations JSONStructure) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	barangayCount := 0
	batchSize := 5000 // Insert 5000 at a time to avoid placeholder limit
	var currentBatch []interface{}
	var currentPlaceholders []string

	for _, cities := range locations {
		for cityName, barangays := range cities {
			cityID, ok := s.citiesMap.Load(cityName)
			if !ok {
				// Skip barangays without matched city
				continue
			}

			for _, barangayName := range barangays {
				currentPlaceholders = append(currentPlaceholders, "(?, ?)")
				currentBatch = append(currentBatch, barangayName, cityID)
				barangayCount++

				// Insert when batch size reached
				if len(currentBatch) >= batchSize {
					if err := s.executeBatchInsert("barangays", currentPlaceholders, currentBatch); err != nil {
						return 0, err
					}
					currentBatch = nil
					currentPlaceholders = nil
				}
			}
		}
	}

	// Insert remaining barangays
	if len(currentBatch) > 0 {
		if err := s.executeBatchInsert("barangays", currentPlaceholders, currentBatch); err != nil {
			return 0, err
		}
	}

	return barangayCount, nil
}

// executeBatchInsert executes a batch insert for a table
func (s *AddressSeeder) executeBatchInsert(table string, placeholders []string, args []interface{}) error {
	if len(placeholders) == 0 {
		return nil
	}

	var query string
	switch table {
	case "regions":
		query = "INSERT INTO regions (name) VALUES " + strings.Join(placeholders, ", ") +
			" ON DUPLICATE KEY UPDATE name=VALUES(name)"
	case "cities":
		query = "INSERT INTO cities (name, region_id) VALUES " + strings.Join(placeholders, ", ") +
			" ON DUPLICATE KEY UPDATE name=VALUES(name)"
	case "barangays":
		query = "INSERT INTO barangays ( name, city_id) VALUES " + strings.Join(placeholders, ", ") +
			" ON DUPLICATE KEY UPDATE name=VALUES(name)"
	default:
		return fmt.Errorf("unknown table: %s", table)
	}

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to batch insert %s: %w", table, err)
	}

	return nil
}

func buildDSNFromEnv() string {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	if pass == "" {
		pass = os.Getenv("DB_PASSWORD")
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", user, pass, host, port, name)
}
