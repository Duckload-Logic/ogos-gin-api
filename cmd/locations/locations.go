package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// JSON file structures (matches output of cmd/psgc)
type PSGCRegion struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type PSGCProvince struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	RegionCode string `json:"regionCode"`
}

type PSGCCity struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	ZipCode      string `json:"zipCode"`
	District     string `json:"district"`
	RegionCode   string `json:"regionCode"`
	ProvinceCode string `json:"provinceCode,omitempty"`
}

type PSGCBarangay struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	CityCode string `json:"cityCode"`
}

type PSGCData struct {
	Regions   []PSGCRegion   `json:"regions"`
	Provinces []PSGCProvince `json:"provinces"`
	Cities    []PSGCCity     `json:"cities"`
	Barangays []PSGCBarangay `json:"barangays"`
}

type AddressSeeder struct {
	db *sqlx.DB
}

var db *sqlx.DB

func NewAddressSeeder(db *sqlx.DB) *AddressSeeder {
	return &AddressSeeder{db: db}
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

	// Find psgc_data.json
	jsonFile := "psgc_data.json"
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		jsonFile = "cmd/locations/psgc_data.json"
	}

	seeder := NewAddressSeeder(db)
	if err := seeder.SeedAddresses(jsonFile); err != nil {
		log.Fatal("failed to seed addresses:", err)
	}
	fmt.Println("Address seeding completed successfully.")
}

// SeedAddresses reads the PSGC JSON file and seeds all location tables.
func (s *AddressSeeder) SeedAddresses(jsonFile string) error {
	fmt.Printf("Loading PSGC data from %s...\n", jsonFile)
	start := time.Now()

	raw, err := os.ReadFile(jsonFile)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", jsonFile, err)
	}

	var data PSGCData
	if err := json.Unmarshal(raw, &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	fmt.Printf("Loaded %d regions, %d provinces, %d cities, %d barangays\n",
		len(data.Regions), len(data.Provinces), len(data.Cities), len(data.Barangays))

	// Clear existing data
	fmt.Println("Clearing existing address data...")
	if err := s.clearData(); err != nil {
		return fmt.Errorf("failed to clear existing data: %w", err)
	}

	// Seed regions
	fmt.Println("Seeding regions...")
	if err := s.seedRegions(data.Regions); err != nil {
		return fmt.Errorf("failed to seed regions: %w", err)
	}

	// Seed provinces
	fmt.Println("Seeding provinces...")
	if err := s.seedProvinces(data.Provinces); err != nil {
		return fmt.Errorf("failed to seed provinces: %w", err)
	}

	var regionRows []struct {
		ID   int    `db:"id"`
		Code string `db:"code"`
	}
	if err := s.db.Select(&regionRows, "SELECT id, code FROM regions"); err != nil {
		return fmt.Errorf("failed to load region IDs: %w", err)
	}

	// Seed cities
	fmt.Println("Seeding cities...")
	if err := s.seedCities(data.Cities); err != nil {
		return fmt.Errorf("failed to seed cities: %w", err)
	}

	// Build city code -> DB id map
	cityIDMap := make(map[string]int)
	var cityRows []struct {
		ID   int    `db:"id"`
		Code string `db:"code"`
	}
	if err := s.db.Select(&cityRows, "SELECT id, code FROM cities"); err != nil {
		return fmt.Errorf("failed to load city IDs: %w", err)
	}
	for _, c := range cityRows {
		cityIDMap[c.Code] = c.ID
	}

	// Seed barangays
	fmt.Println("Seeding barangays...")
	if err := s.seedBarangays(data.Barangays); err != nil {
		return fmt.Errorf("failed to seed barangays: %w", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Seeding completed in %v\n", elapsed)
	return nil
}

func (s *AddressSeeder) clearData() error {
	if _, err := s.db.Exec("SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		return err
	}
	defer s.db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	for _, table := range []string{"barangays", "cities", "provinces", "regions"} {
		if _, err := s.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table)); err != nil {
			return fmt.Errorf("failed to truncate %s: %w", table, err)
		}
	}
	return nil
}

func (s *AddressSeeder) seedRegions(regions []PSGCRegion) error {
	if len(regions) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(regions))
	args := make([]interface{}, 0, len(regions)*2)
	for _, r := range regions {
		placeholders = append(placeholders, "(?, ?)")
		args = append(args, r.Code, r.Name)
	}

	query := "INSERT INTO regions (code, name) VALUES " + strings.Join(placeholders, ", ") +
		" ON DUPLICATE KEY UPDATE name=VALUES(name)"
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *AddressSeeder) seedProvinces(provinces []PSGCProvince) error {
	if len(provinces) == 0 {
		return nil
	}

	batchSize := 500
	for i := 0; i < len(provinces); i += batchSize {
		end := i + batchSize
		if end > len(provinces) {
			end = len(provinces)
		}
		batch := provinces[i:end]

		placeholders := make([]string, 0, len(batch))
		args := make([]interface{}, 0, len(batch)*3)
		for _, p := range batch {
			placeholders = append(placeholders, "(?, ?, ?)")
			args = append(args, p.Code, p.Name, p.RegionCode)
		}

		query := "INSERT INTO provinces (code, name, region_code) VALUES " + strings.Join(placeholders, ", ") +
			" ON DUPLICATE KEY UPDATE name=VALUES(name)"
		if _, err := s.db.Exec(query, args...); err != nil {
			return fmt.Errorf("failed to batch insert provinces: %w", err)
		}
	}
	return nil
}

func (s *AddressSeeder) seedCities(cities []PSGCCity) error {
	if len(cities) == 0 {
		return nil
	}

	batchSize := 500
	for i := 0; i < len(cities); i += batchSize {
		end := i + batchSize
		if end > len(cities) {
			end = len(cities)
		}
		batch := cities[i:end]

		placeholders := make([]string, 0, len(batch))
		args := make([]interface{}, 0, len(batch)*8)
		for _, c := range batch {
			var provCode interface{}
			if c.ProvinceCode != "" {
				provCode = c.ProvinceCode
			}
			var regCode interface{}
			if c.RegionCode != "" {
				regCode = c.RegionCode
			}
			placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?)")
			args = append(args, c.Code, c.Name, provCode, c.Type, c.ZipCode, c.District, regCode)
		}

		query := "INSERT INTO cities (code, name, province_code, `type`, zip_code, district, region_code) VALUES " +
			strings.Join(placeholders, ", ") +
			" ON DUPLICATE KEY UPDATE name=VALUES(name), `type`=VALUES(`type`), zip_code=VALUES(zip_code), district=VALUES(district), region_code=VALUES(region_code)"
		if _, err := s.db.Exec(query, args...); err != nil {
			return fmt.Errorf("failed to batch insert cities: %w", err)
		}
	}
	return nil
}

func (s *AddressSeeder) seedBarangays(barangays []PSGCBarangay) error {
	if len(barangays) == 0 {
		return nil
	}

	batchSize := 2000
	for i := 0; i < len(barangays); i += batchSize {
		end := i + batchSize
		if end > len(barangays) {
			end = len(barangays)
		}
		batch := barangays[i:end]

		placeholders := make([]string, 0, len(batch))
		args := make([]interface{}, 0, len(batch)*3)
		for _, b := range batch {
			placeholders = append(placeholders, "(?, ?, ?)")
			args = append(args, b.Code, b.Name, b.CityCode)
		}

		if len(placeholders) == 0 {
			continue
		}

		query := "INSERT INTO barangays (code, name, city_code) VALUES " +
			strings.Join(placeholders, ", ") +
			" ON DUPLICATE KEY UPDATE name=VALUES(name)"
		if _, err := s.db.Exec(query, args...); err != nil {
			return fmt.Errorf("failed to batch insert barangays: %w", err)
		}

		if (i+batchSize)%10000 < batchSize {
			fmt.Printf("  Inserted %d/%d barangays...\n", end, len(barangays))
		}
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
