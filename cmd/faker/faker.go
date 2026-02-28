package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var db *sqlx.DB

// lookup slices (IDs)
var (
	genderIDs                []int
	civilStatusIDs           []int
	civilStatusByName        map[string]int
	religionIDs              []int
	courseIDs                []int
	enrollmentReasonIDs      []int
	supportTypeIDs           []int
	incomeRangeIDs           []int
	parentalStatusIDs        []int
	educationalLevelIDs      []int
	educationalLevelByName   map[string]int
	relationshipTypeIDs      map[string]int
	natureOfResidenceIDs     []int
	siblingSupportTypeIDs    []int
	activityOptionIDs        []int
	timeSlotIDs              []int
	appointmentStatusIDs     []int
	appointmentStatusByName  map[string]int
	appointmentCategoryIDs   []int
	appointmentSlotMu        sync.Mutex
	reservedAppointmentSlots = make(map[string]struct{})
)

func main() {
	// ---------- CONFIGURATION ----------
	numStudents := 500 // number of students to generate
	numCounselors := 1 // number of counselors (admins)
	numWorkers := 100  // number of concurrent student workers
	_ = godotenv.Load()
	dsn := buildDSNFromEnv()
	startTime := time.Now()
	// -----------------------------------

	var err error
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}
	defer db.Close()

	// load all lookup data from reference tables
	loadLookups()

	// seed random generator
	gofakeit.Seed(time.Now().UnixNano())

	// clear existing student data (optional but keeps the run idempotent)
	clearStudentData()

	// create counselors
	for i := 0; i < numCounselors; i++ {
		createCounselor()
	}

	// create students (concurrently with a worker pool)
	jobs := make(chan int, numStudents)
	var wg sync.WaitGroup

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				createStudent(index)
			}
		}()
	}

	for i := 0; i < numStudents; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	fmt.Println("Dummy data generation completed successfully.")
	log.Println("Time taken:", time.Since(startTime))
}

// ----------------------------------------------------------------------
// LOOKUP LOADING
// ----------------------------------------------------------------------
func loadLookups() {
	// genders
	rows, err := db.Query("SELECT id FROM genders")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		genderIDs = append(genderIDs, id)
	}

	// civil status
	civilStatusByName = make(map[string]int)
	rows, err = db.Query("SELECT id, status_name FROM civil_status_types")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var statusName string
		if err := rows.Scan(&id, &statusName); err != nil {
			log.Fatal(err)
		}
		civilStatusIDs = append(civilStatusIDs, id)
		civilStatusByName[strings.ToLower(statusName)] = id
	}

	// religions
	rows, err = db.Query("SELECT id FROM religions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		religionIDs = append(religionIDs, id)
	}

	// courses
	rows, err = db.Query("SELECT id FROM courses")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		courseIDs = append(courseIDs, id)
	}

	// enrollment reasons
	rows, err = db.Query("SELECT id FROM enrollment_reasons")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		enrollmentReasonIDs = append(enrollmentReasonIDs, id)
	}

	// student support types
	rows, err = db.Query("SELECT id FROM student_support_types")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		supportTypeIDs = append(supportTypeIDs, id)
	}

	// income ranges
	rows, err = db.Query("SELECT id FROM income_ranges")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		incomeRangeIDs = append(incomeRangeIDs, id)
	}

	// parental status
	rows, err = db.Query("SELECT id FROM parental_status_types")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		parentalStatusIDs = append(parentalStatusIDs, id)
	}

	// educational levels
	educationalLevelByName = make(map[string]int)
	rows, err = db.Query("SELECT id, level_name FROM educational_levels")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var levelName string
		if err := rows.Scan(&id, &levelName); err != nil {
			log.Fatal(err)
		}
		educationalLevelIDs = append(educationalLevelIDs, id)
		educationalLevelByName[strings.ToLower(levelName)] = id
	}

	// relationship types (map by name for easy lookup)
	relationshipTypeIDs = make(map[string]int)
	rows, err = db.Query("SELECT id, relationship_name FROM student_relationship_types")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		relationshipTypeIDs[name] = id
	}

	// nature of residence
	rows, err = db.Query("SELECT id FROM nature_of_residence_types")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		natureOfResidenceIDs = append(natureOfResidenceIDs, id)
	}

	// sibling support types
	rows, err = db.Query("SELECT id FROM sibling_support_types")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		siblingSupportTypeIDs = append(siblingSupportTypeIDs, id)
	}

	// activity options
	rows, err = db.Query("SELECT id FROM activity_options WHERE is_active = TRUE")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		activityOptionIDs = append(activityOptionIDs, id)
	}

	// appointment time slots
	rows, err = db.Query("SELECT id FROM time_slots")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		timeSlotIDs = append(timeSlotIDs, id)
	}

	// appointment statuses
	appointmentStatusByName = make(map[string]int)
	rows, err = db.Query("SELECT id, name FROM statuses")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		appointmentStatusIDs = append(appointmentStatusIDs, id)
		appointmentStatusByName[strings.ToLower(name)] = id
	}

	// appointment categories (fallback to concern_categories if present in local schema)
	rows, err = db.Query("SELECT id FROM appointment_categories")
	if err != nil {
		rows, err = db.Query("SELECT id FROM concern_categories")
		if err != nil {
			log.Fatal(err)
		}
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		appointmentCategoryIDs = append(appointmentCategoryIDs, id)
	}
}

// ----------------------------------------------------------------------
// CLEAR EXISTING STUDENT DATA (keep reference tables and seeded admin)
// ----------------------------------------------------------------------
func clearStudentData() {
	tables := []string{
		"admission_slips",
		"student_hobbies",
		"student_subject_preferences",
		"student_activities",
		"student_financial_supports",
		"student_finances",
		"significant_notes",
		"test_results",
		"student_consultations",
		"student_health_records",
		"school_details",
		"educational_backgrounds",
		"student_sibling_supports",
		"family_backgrounds",
		"student_related_persons",
		"related_persons",
		"student_addresses",
		"addresses",
		"student_selected_reasons",
		"emergency_contacts",
		"student_personal_info",
		"iir_records",
		"appointments",
		"counselor_profiles",
		"users",
	}

	// disable FK checks temporarily
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	for _, tbl := range tables {
		// delete only rows that belong to students (role_id=1) or their dependents
		// we'll just truncate for simplicity, but you can refine if you need to keep seeded counselor
		if tbl == "users" {
			db.Exec("DELETE FROM users WHERE role_id = 1")
		} else {
			db.Exec(fmt.Sprintf("DELETE FROM %s", tbl))
		}
	}
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}

// ----------------------------------------------------------------------
// CREATE A COUNSELOR (ADMIN)
// ----------------------------------------------------------------------
func createCounselor() {
	tx, err := db.Beginx()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	email := "counselor@university.edu"

	res, err := tx.Exec(`
		INSERT INTO users (role_id, first_name, middle_name, last_name, email, password_hash, is_active)
		VALUES (2, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			id = LAST_INSERT_ID(id),
			first_name = VALUES(first_name),
			middle_name = VALUES(middle_name),
			last_name = VALUES(last_name),
			password_hash = VALUES(password_hash),
			is_active = VALUES(is_active)
	`, firstName, randomMiddleName(), lastName, email, fakePasswordHash(), true)
	if err != nil {
		log.Fatal(err)
	}
	userID, _ := res.LastInsertId()

	_, err = tx.Exec(`
		INSERT INTO counselor_profiles (user_id, license_number, specialization, is_available)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			license_number = VALUES(license_number),
			specialization = VALUES(specialization),
			is_available = VALUES(is_available)
	`, userID, gofakeit.Regex("[A-Z]{3}-[0-9]{6}"), gofakeit.JobTitle(), true)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
	fmt.Printf("Created counselor: %s %s\n", firstName, lastName)
}

// ----------------------------------------------------------------------
// CREATE A SINGLE STUDENT WITH ALL RELATED RECORDS
// ----------------------------------------------------------------------
func createStudent(index int) {
	tx, err := db.Beginx()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	// generate core data needed later
	dob := gofakeit.DateRange(time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC), time.Now().AddDate(-18, 0, 0))
	birthYear := dob.Year()
	studentEmail := fmt.Sprintf("student%d@university.edu", index+1) // guarantee unique

	// 1. users
	res, err := tx.Exec(`
		INSERT INTO users (role_id, first_name, middle_name, last_name, email, password_hash, is_active)
		VALUES (1, ?, ?, ?, ?, ?, ?)
	`, gofakeit.FirstName(), randomMiddleName(), gofakeit.LastName(), studentEmail, fakePasswordHash(), true)
	if err != nil {
		log.Fatal(err)
	}
	userID, _ := res.LastInsertId()

	if rand.Float32() < 0.7 {
		res, err = tx.Exec(`
			INSERT INTO iir_records (user_id, is_submitted) VALUES (?, ?)
		`, userID, true)
		if err != nil {
			log.Fatal(err)
		}
		iirID, _ := res.LastInsertId()

		// 3. selected reasons
		insertSelectedReasons(tx, int(iirID))

		// 4. addresses (residential & provincial)
		resAddr1 := insertAddress(tx)
		resAddr2 := insertAddress(tx)
		insertStudentAddress(tx, int(iirID), resAddr1, "Residential")
		insertStudentAddress(tx, int(iirID), resAddr2, "Provincial")

		// 5. related persons (father, mother, optional guardian)
		father := insertRelatedPerson(tx)
		mother := insertRelatedPerson(tx)
		father.AddressID = &resAddr1
		mother.AddressID = &resAddr1

		guardianScenario := pickGuardianScenario()
		var guardian *relatedPersonSeed

		switch guardianScenario {
		case "father_guardian":
			father.AddressID = &resAddr2
		case "mother_guardian":
			mother.AddressID = &resAddr2
		case "separate_guardian":
			g := insertRelatedPerson(tx)
			g.AddressID = &resAddr2
			guardian = &g
		}

		// 6. link related persons
		switch guardianScenario {
		case "no_guardian":
			linkRelatedPerson(tx, int(iirID), father.ID, "Father", true, false, true)
			linkRelatedPerson(tx, int(iirID), mother.ID, "Mother", true, false, true)
		case "father_guardian":
			linkRelatedPerson(tx, int(iirID), father.ID, "Father", true, true, true)
			linkRelatedPerson(tx, int(iirID), mother.ID, "Mother", true, false, true)
		case "mother_guardian":
			linkRelatedPerson(tx, int(iirID), father.ID, "Father", true, false, true)
			linkRelatedPerson(tx, int(iirID), mother.ID, "Mother", true, true, true)
		case "separate_guardian":
			linkRelatedPerson(tx, int(iirID), father.ID, "Father", true, false, true)
			linkRelatedPerson(tx, int(iirID), mother.ID, "Mother", true, false, true)
			linkRelatedPerson(tx, int(iirID), guardian.ID, "Guardian", false, true, true)
		}

		// 7. student_personal_info
		emergency := deriveEmergencyContact(father, mother, guardian, guardianScenario, resAddr1, resAddr2)
		emergencyContactID := insertEmergencyContact(tx, int(iirID), emergency)
		insertPersonalInfo(tx, int(iirID), dob, index, emergencyContactID)

		// 8. family background
		familyBgID := insertFamilyBackground(tx, int(iirID))

		// 9. sibling supports (if employed siblings > 0)
		insertSiblingSupports(tx, familyBgID)

		// 10. educational background
		ebID := insertEducationalBackground(tx, int(iirID))

		// 11. school details for each educational level
		insertSchoolDetails(tx, ebID, birthYear, index)

		// 12. health records
		insertHealthRecords(tx, int(iirID))

		// 13. consultations
		insertConsultations(tx, int(iirID))

		// 14. test results
		insertTestResults(tx, int(iirID))

		// 15. significant notes
		insertSignificantNotes(tx, int(iirID))

		// 16. finances
		sfID := insertStudentFinances(tx, int(iirID))

		// 17. financial supports
		insertFinancialSupports(tx, sfID)

		// 18. activities
		insertActivities(tx, int(iirID))

		// 19. subject preferences
		insertSubjectPreferences(tx, int(iirID))

		// 20. hobbies
		insertHobbies(tx, int(iirID))

		// 21. admission slip (30% chance)
		if rand.Float32() < 0.3 {
			insertAdmissionSlip(tx, int(iirID))
		}

		// 22. appointment (30% chance)
		if rand.Float32() < 0.3 {
			for i := 0; i < rand.Intn(5)+1; i++ { // up to 15 appointments per student
				insertAppointment(tx, int(userID))
			}
		}

		tx.Commit()
		fmt.Printf("Created student %d (iirID=%d)\n", index+1, iirID)
		return
	}

	tx.Commit()
	fmt.Printf("Created student %d (no IIR)\n", index+1)
}

// ----------------------------------------------------------------------
// INDIVIDUAL INSERT HELPERS
// ----------------------------------------------------------------------

type relatedPersonSeed struct {
	ID            int
	FirstName     string
	MiddleName    sql.NullString
	LastName      string
	ContactNumber sql.NullString
	AddressID     *int
}

type emergencyContactSeed struct {
	ID             int
	FirstName      string
	MiddleName     sql.NullString
	LastName       string
	Number         string
	RelationshipID int
	AddressID      int
}

func pickGuardianScenario() string {
	r := rand.Float32()
	if r < 0.25 {
		return "no_guardian"
	}
	if r < 0.50 {
		return "father_guardian"
	}
	if r < 0.75 {
		return "mother_guardian"
	}
	return "separate_guardian"
}

func deriveEmergencyContact(
	father, mother relatedPersonSeed,
	guardian *relatedPersonSeed,
	guardianScenario string,
	residentialAddressID, provincialAddressID int,
) emergencyContactSeed {
	switch guardianScenario {
	case "father_guardian":
		return emergencyContactSeed{
			FirstName:      father.FirstName,
			MiddleName:     father.MiddleName,
			LastName:       father.LastName,
			Number:         validContact(father.ContactNumber),
			RelationshipID: relationshipID("Father"),
			AddressID:      safeAddressID(father.AddressID, provincialAddressID),
		}
	case "mother_guardian":
		return emergencyContactSeed{
			FirstName:      mother.FirstName,
			MiddleName:     mother.MiddleName,
			LastName:       mother.LastName,
			Number:         validContact(mother.ContactNumber),
			RelationshipID: relationshipID("Mother"),
			AddressID:      safeAddressID(mother.AddressID, provincialAddressID),
		}
	case "separate_guardian":
		return emergencyContactSeed{
			FirstName:      guardian.FirstName,
			MiddleName:     guardian.MiddleName,
			LastName:       guardian.LastName,
			Number:         validContact(guardian.ContactNumber),
			RelationshipID: relationshipID("Guardian"),
			AddressID:      safeAddressID(guardian.AddressID, provincialAddressID),
		}
	default:
		// no guardian: choose one of the parents as emergency contact
		if rand.Intn(2) == 0 {
			return emergencyContactSeed{
				FirstName:      father.FirstName,
				MiddleName:     father.MiddleName,
				LastName:       father.LastName,
				Number:         validContact(father.ContactNumber),
				RelationshipID: relationshipID("Father"),
				AddressID:      safeAddressID(father.AddressID, residentialAddressID),
			}
		}
		return emergencyContactSeed{
			FirstName:      mother.FirstName,
			MiddleName:     mother.MiddleName,
			LastName:       mother.LastName,
			Number:         validContact(mother.ContactNumber),
			RelationshipID: relationshipID("Mother"),
			AddressID:      safeAddressID(mother.AddressID, residentialAddressID),
		}
	}
}

func safeAddressID(addrID *int, fallback int) int {
	if addrID != nil {
		return *addrID
	}
	return fallback
}

func validContact(contact sql.NullString) string {
	if contact.Valid && strings.TrimSpace(contact.String) != "" {
		return contact.String
	}
	return gofakeit.Phone()
}

func insertEmergencyContact(tx *sqlx.Tx, iirID int, emergency emergencyContactSeed) int {
	res, err := tx.Exec(`
		INSERT INTO emergency_contacts (
			iir_id, first_name, middle_name, last_name,
			contact_number, relationship_id, address_id
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		iirID,
		emergency.FirstName,
		emergency.MiddleName,
		emergency.LastName,
		emergency.Number,
		emergency.RelationshipID,
		emergency.AddressID,
	)
	if err != nil {
		log.Fatal(err)
	}
	id, _ := res.LastInsertId()
	return int(id)
}

func randomMiddleName() sql.NullString {
	if gofakeit.Bool() {
		return sql.NullString{String: gofakeit.FirstName(), Valid: true}
	}
	return sql.NullString{Valid: false}
}

func fakePasswordHash() string {
	// static hash for all dummy users
	return "$2y$10$gxeDD.IKlEkqJmqmyVxy6eU9tFvC4ZK8KL3VZc2ex3BvNLo8DL5Dq"
}

func insertPersonalInfo(tx *sqlx.Tx, iirID int, dob time.Time, studentIndex int, emergencyContactID int) {
	studentStatus := rand.Intn(2) % 2 // 0 or 1
	isEmployed := studentIndex%2 == 0
	studentNumber := fmt.Sprintf("%d-%05d-TG-%d", time.Now().Year(), iirID, studentStatus)
	var employerName, employerAddress sql.NullString
	if isEmployed {
		empName := gofakeit.Company()
		empAddr := gofakeit.Address().Address
		employerName = sql.NullString{String: empName, Valid: true}
		employerAddress = sql.NullString{String: empAddr, Valid: true}
	}

	mobileNumber := gofakeit.Phone()
	telephoneNumber := sql.NullString{Valid: false}
	if studentIndex%3 != 0 {
		telephoneNumber = sql.NullString{String: gofakeit.Phone(), Valid: true}
	}

	civilStatusID := chooseCivilStatusID()

	_, err := tx.Exec(`
		INSERT INTO student_personal_info (
			iir_id, student_number, gender_id, civil_status_id, religion_id,
			height_ft, weight_kg, complexion, high_school_gwa, course_id,
			year_level, section, place_of_birth, date_of_birth,
			is_employed, employer_name, employer_address, mobile_number, telephone_number
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, iirID, studentNumber,
		randomChoice(genderIDs), civilStatusID, randomChoice(religionIDs),
		gofakeit.Float64Range(4.5, 6.5), gofakeit.Float64Range(40, 100), gofakeit.Color(),
		gofakeit.Float64Range(75, 98), randomChoice(courseIDs),
		rand.Intn(4)+1, rand.Intn(5)+1,
		gofakeit.City(), dob.Format("2006-01-02"),
		isEmployed, employerName, employerAddress, mobileNumber, telephoneNumber,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func insertSelectedReasons(tx *sqlx.Tx, iirID int) {
	// pick 1-3 random reasons
	num := rand.Intn(3) + 1
	selected := make(map[int]bool)
	for len(selected) < num {
		rid := randomChoice(enrollmentReasonIDs).(int)
		selected[rid] = true
	}
	for rid := range selected {
		_, err := tx.Exec(`
			INSERT INTO student_selected_reasons (iir_id, reason_id, other_reason_text)
			VALUES (?, ?, ?)
		`, iirID, rid, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func insertAddress(tx *sqlx.Tx) int {
	// Fetch a random region
	var regionID int
	err := tx.Get(&regionID, "SELECT id FROM regions ORDER BY RAND() LIMIT 1")
	if err != nil {
		log.Fatal("No regions found. Please run the address seeder first (make locations): ", err)
	}

	// Fetch a random city in that region that has barangays
	var cityID int
	err = tx.Get(&cityID, `
		SELECT DISTINCT c.id FROM cities c
		INNER JOIN barangays b ON c.id = b.city_id
		WHERE c.region_id = ?
		ORDER BY RAND() LIMIT 1
	`, regionID)
	if err != nil {
		log.Fatal("No cities with barangays found for region ID "+fmt.Sprint(regionID)+": ", err)
	}

	// Fetch a random barangay in that city
	var barangayID int
	err = tx.Get(&barangayID, "SELECT id FROM barangays WHERE city_id = ? ORDER BY RAND() LIMIT 1", cityID)
	if err != nil {
		log.Fatal("No barangays found for city ID "+fmt.Sprint(cityID)+": ", err)
	}

	street := gofakeit.Street()

	res, err := tx.Exec(`
		INSERT INTO addresses (region_id, city_id, barangay_id, street_detail)
		VALUES (?, ?, ?, ?)
	`, regionID, cityID, barangayID, street)
	if err != nil {
		log.Fatal(err)
	}
	id, _ := res.LastInsertId()
	return int(id)
}

func insertStudentAddress(tx *sqlx.Tx, iirID, addressID int, addrType string) {
	_, err := tx.Exec(`
		INSERT INTO student_addresses (iir_id, address_id, address_type)
		VALUES (?, ?, ?)
	`, iirID, addressID, addrType)
	if err != nil {
		log.Fatal(err)
	}
}

func insertRelatedPerson(tx *sqlx.Tx) relatedPersonSeed {
	fname := gofakeit.FirstName()
	lname := gofakeit.LastName()
	mname := randomMiddleName()
	dob := gofakeit.DateRange(time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC), time.Now().AddDate(-30, 0, 0))
	educ := randomEducationalAttainment()
	occupation := sql.NullString{String: gofakeit.JobTitle(), Valid: gofakeit.Bool()}
	employer := sql.NullString{String: gofakeit.Company(), Valid: occupation.Valid}
	employerAddr := sql.NullString{String: gofakeit.Address().Address, Valid: occupation.Valid}
	contact := sql.NullString{String: gofakeit.Phone(), Valid: gofakeit.Bool()}

	res, err := tx.Exec(`
		INSERT INTO related_persons (
			educational_level, date_of_birth,
			last_name, first_name, middle_name,
			occupation, employer_name, employer_address, contact_number
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, educ, dob, lname, fname, mname,
		occupation, employer, employerAddr, contact)
	if err != nil {
		log.Fatal(err)
	}
	id, _ := res.LastInsertId()

	return relatedPersonSeed{
		ID:            int(id),
		FirstName:     fname,
		MiddleName:    mname,
		LastName:      lname,
		ContactNumber: contact,
		AddressID:     nil,
	}
}

func randomEducationalAttainment() string {
	levels := []string{"Elementary", "High School", "Vocational", "College", "Post Graduate"}
	return levels[rand.Intn(len(levels))]
}

func linkRelatedPerson(tx *sqlx.Tx, iirID, personID int, relType string, isParent, isGuardian, isLiving bool) {
	relID := relationshipID(relType)
	_, err := tx.Exec(`
		INSERT INTO student_related_persons (
			iir_id, related_person_id, relationship_id,
			is_parent, is_guardian, is_living
		) VALUES (?, ?, ?, ?, ?, ?)
	`, iirID, personID, relID, isParent, isGuardian, isLiving)

	if err != nil {
		log.Fatal(err)
	}
}

func insertFamilyBackground(tx *sqlx.Tx, iirID int) int {
	brothers := rand.Intn(4)
	sisters := rand.Intn(4)
	employedSibs := rand.Intn(brothers + sisters + 1)
	ordinal := rand.Intn(brothers+sisters+1) + 1
	parentalID := randomChoice(parentalStatusIDs)
	var details sql.NullString
	if parentalID == 5 { // "Other"
		details = sql.NullString{String: gofakeit.Sentence(3), Valid: true}
	}
	quiet := gofakeit.Bool()
	sharing := gofakeit.Bool()
	var shareDetails sql.NullString
	if sharing {
		shareDetails = sql.NullString{String: "Shares with " + gofakeit.FirstName(), Valid: true}
	}
	natureID := randomChoice(natureOfResidenceIDs)

	res, err := tx.Exec(`
		INSERT INTO family_backgrounds (
			iir_id, parental_status_id, parental_status_details,
			brothers, sisters, employed_siblings, ordinal_position,
			have_quiet_place_to_study, is_sharing_room, room_sharing_details,
			nature_of_residence_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, iirID, parentalID, details,
		brothers, sisters, employedSibs, ordinal,
		quiet, sharing, shareDetails,
		natureID)
	if err != nil {
		log.Fatal(err)
	}
	id, _ := res.LastInsertId()
	return int(id)
}

func insertSiblingSupports(tx *sqlx.Tx, familyBgID int) {
	// only if employed siblings > 0, maybe add 1-2 supports
	if rand.Intn(3) == 0 { // 1/3 chance
		num := rand.Intn(2) + 1
		selected := make(map[int]bool)
		for len(selected) < num {
			stID := randomChoice(siblingSupportTypeIDs).(int)
			selected[stID] = true
		}
		for stID := range selected {
			_, err := tx.Exec(`
				INSERT INTO student_sibling_supports (family_background_id, support_type_id)
				VALUES (?, ?)
			`, familyBgID, stID)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func insertEducationalBackground(tx *sqlx.Tx, iirID int) int {
	nature := "Continuous"
	var details sql.NullString
	if rand.Float32() < 0.1 { // 10% interrupted
		nature = "Interrupted"
		details = sql.NullString{String: gofakeit.Sentence(5), Valid: true}
	}
	res, err := tx.Exec(`
		INSERT INTO educational_backgrounds (iir_id, nature_of_schooling, interrupted_details)
		VALUES (?, ?, ?)
	`, iirID, nature, details)
	if err != nil {
		log.Fatal(err)
	}
	id, _ := res.LastInsertId()
	return int(id)
}

// school details for every educational level
func insertSchoolDetails(tx *sqlx.Tx, ebID, birthYear, studentIndex int) {
	// typical schooling ages:
	// pre-elementary: 4-5
	// elementary: 6-11
	// high school: 12-15
	// vocational: 16-17 (optional)
	// college: 18-21/22
	nowYear := time.Now().Year()

	// map level name -> typical start age
	levelStartAge := map[string]int{
		"Pre-Elementary": 4,
		"Elementary":     6,
		"High School":    12,
		"Vocational":     16,
		"College":        18,
	}

	levelPath := []string{"pre-elementary", "elementary", "high school"}
	if studentIndex%3 == 0 {
		// scenario: pre-school up to vocational/SHS route
		levelPath = append(levelPath, "vocational")
	} else {
		// scenario: pre-school up to college route
		levelPath = append(levelPath, "college")
	}

	for _, levelKey := range levelPath {
		levelID, ok := educationalLevelByName[levelKey]
		if !ok {
			continue
		}

		levelName := strings.Title(levelKey)
		if levelKey == "pre-elementary" {
			levelName = "Pre-Elementary"
		}
		if levelKey == "high school" {
			levelName = "High School"
		}

		startAge := levelStartAge[levelName]
		startYear := birthYear + startAge
		// duration depends on level
		duration := 0
		switch levelName {
		case "Pre-Elementary":
			duration = 2
		case "Elementary":
			duration = 6
		case "High School":
			duration = 4
		case "Vocational":
			duration = 2
		case "College":
			duration = 4
		}
		completedYear := startYear + duration
		// if completedYear > nowYear, maybe still studying, but we'll set completed as nowYear or future (like seed)
		if completedYear > nowYear {
			completedYear = nowYear + rand.Intn(3) // still in school
		}
		schoolType := randomChoice([]string{"Private", "Public"})
		awards := randomAwards()

		_, err := tx.Exec(`
			INSERT INTO school_details (
				eb_id, educational_level_id, school_name, school_address,
				school_type, year_started, year_completed, awards
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, ebID, levelID,
			gofakeit.Company(), gofakeit.Address().Address,
			schoolType, startYear, completedYear, awards)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func insertHealthRecords(tx *sqlx.Tx, iirID int) {
	visionProb := false
	hearingProb := false
	speechProb := false
	genProb := false

	roll := rand.Float32()
	if roll >= 0.35 {
		impairments := []string{"vision", "hearing", "speech", "general"}
		rand.Shuffle(len(impairments), func(i, j int) {
			impairments[i], impairments[j] = impairments[j], impairments[i]
		})

		count := 1
		if roll >= 0.75 {
			count = 2
		}

		for _, impairment := range impairments[:count] {
			switch impairment {
			case "vision":
				visionProb = true
			case "hearing":
				hearingProb = true
			case "speech":
				speechProb = true
			case "general":
				genProb = true
			}
		}
	}

	visionDet := nullStringIf(visionProb, gofakeit.Sentence(3))
	hearingDet := nullStringIf(hearingProb, gofakeit.Sentence(3))
	speechDet := nullStringIf(speechProb, gofakeit.Sentence(3))
	genDet := nullStringIf(genProb, gofakeit.Sentence(5))

	_, err := tx.Exec(`
		INSERT INTO student_health_records (
			iir_id,
			vision_has_problem, vision_details,
			hearing_has_problem, hearing_details,
			speech_has_problem, speech_details,
			general_health_has_problem, general_health_details
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, iirID,
		visionProb, visionDet,
		hearingProb, hearingDet,
		speechProb, speechDet,
		genProb, genDet)
	if err != nil {
		log.Fatal(err)
	}
}

func nullStringIf(cond bool, val string) sql.NullString {
	if cond {
		return sql.NullString{String: val, Valid: true}
	}
	return sql.NullString{Valid: false}
}

func insertConsultations(tx *sqlx.Tx, iirID int) {
	if rand.Float32() < 0.55 { // 55% have prior consultations
		num := rand.Intn(3) + 1
		for i := 0; i < num; i++ {
			profType := randomChoice([]string{"Psychiatrist", "Psychologist", "Counselor"})
			has := true
			when := sql.NullString{String: gofakeit.Date().Format("2006-01-02"), Valid: true}
			what := sql.NullString{String: gofakeit.Sentence(8), Valid: true}
			_, err := tx.Exec(`
				INSERT INTO student_consultations (
					iir_id, professional_type, has_consulted, when_date, for_what
				) VALUES (?, ?, ?, ?, ?)
			`, iirID, profType, has, when, what)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func insertTestResults(tx *sqlx.Tx, iirID int) {
	numOfTests := 2
	for i := 0; i < numOfTests; i++ {
		_, err := tx.Exec(`
			INSERT INTO test_results (
				iir_id, test_date, test_name, raw_score, percentile, description
			) VALUES (?, ?, ?, ?, ?, ?)
		`, iirID,
			gofakeit.Date().Format("2006-01-02"),
			gofakeit.JobTitle()+" Test",
			fmt.Sprintf("%d", rand.Intn(50)+50),
			fmt.Sprintf("%d", rand.Intn(30)+70),
			gofakeit.Sentence(5))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func insertSignificantNotes(tx *sqlx.Tx, iirID int) {
	if rand.Float32() < 0.3 {
		num := rand.Intn(3) + 1
		for i := 0; i < num; i++ {
			_, err := tx.Exec(`
				INSERT INTO significant_notes (
					iir_id, note_date, incident_description, remarks
				) VALUES (?, ?, ?, ?)
			`, iirID,
				gofakeit.Date().Format("2006-01-02"),
				gofakeit.Sentence(8),
				gofakeit.Sentence(5))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func insertStudentFinances(tx *sqlx.Tx, iirID int) int {
	incRangeID := randomChoice(incomeRangeIDs)
	var otherInc sql.NullString
	if incRangeID == 10 { // Above 50k, maybe other details
		otherInc = sql.NullString{String: "Business income", Valid: gofakeit.Bool()}
	}
	allowance := gofakeit.Float64Range(100, 1500)

	res, err := tx.Exec(`
		INSERT INTO student_finances (
			iir_id, monthly_family_income_range_id, other_income_details, weekly_allowance
		) VALUES (?, ?, ?, ?)
	`, iirID, incRangeID, otherInc, allowance)
	if err != nil {
		log.Fatal(err)
	}
	id, _ := res.LastInsertId()
	return int(id)
}

func insertFinancialSupports(tx *sqlx.Tx, sfID int) {
	num := rand.Intn(3) + 1
	selected := make(map[int]bool)
	for len(selected) < num {
		stID := randomChoice(supportTypeIDs).(int)
		selected[stID] = true
	}
	for stID := range selected {
		_, err := tx.Exec(`
			INSERT INTO student_financial_supports (sf_id, support_type_id)
			VALUES (?, ?)
		`, sfID, stID)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func insertActivities(tx *sqlx.Tx, iirID int) {
	if len(activityOptionIDs) == 0 {
		return
	}
	num := rand.Intn(4) // 0-3 activities
	for i := 0; i < num; i++ {
		optID := randomChoice(activityOptionIDs)
		role := randomChoice([]string{"Officer", "Member", "Other"})
		var roleSpec sql.NullString
		if role == "Other" {
			roleSpec = sql.NullString{String: gofakeit.Word(), Valid: true}
		}
		_, err := tx.Exec(`
			INSERT INTO student_activities (
				iir_id, option_id, role, role_specification
			) VALUES (?, ?, ?, ?)
		`, iirID, optID, role, roleSpec)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func insertSubjectPreferences(tx *sqlx.Tx, iirID int) {
	if rand.Float32() < 0.2 {
		return
	}

	subjects := []string{"Math", "Science", "English", "History", "PE", "Art", "Music", "Computer"}
	num := rand.Intn(5) + 2 // 2-6 unique subjects
	selected := pickUniqueStrings(subjects, num)
	fav := gofakeit.Bool()

	for _, sub := range selected {
		_, err := tx.Exec(`
			INSERT INTO student_subject_preferences (iir_id, subject_name, is_favorite)
			VALUES (?, ?, ?)
		`, iirID, sub, fav)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func insertHobbies(tx *sqlx.Tx, iirID int) {
	num := rand.Intn(4) + 1 // 1-4 hobbies
	for i := 1; i <= num; i++ {
		// ensure unique priority rank
		rank := i
		_, err := tx.Exec(`
			INSERT INTO student_hobbies (iir_id, hobby_name, priority_rank)
			VALUES (?, ?, ?)
		`, iirID, gofakeit.Hobby(), rank)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func insertAdmissionSlip(tx *sqlx.Tx, iirID int) {
	status := randomChoice([]string{"Pending", "Approved", "Rejected"})
	_, err := tx.Exec(`
		INSERT INTO admission_slips (
			iir_id, reason, date_of_absence, file_path, excuse_slip_status
		) VALUES (?, ?, ?, ?, ?)
	`, iirID,
		gofakeit.Sentence(5),
		gofakeit.Date().Format("2006-01-02"),
		"/uploads/"+gofakeit.UUID()+".pdf",
		status)
	if err != nil {
		log.Fatal(err)
	}
}

func insertAppointment(tx *sqlx.Tx, userID int) {
	if len(timeSlotIDs) == 0 || len(appointmentCategoryIDs) == 0 || len(appointmentStatusIDs) == 0 {
		log.Printf("skipping appointment creation for user %d: missing appointment lookup data", userID)
		return
	}

	whenDate, timeSlotID := reserveAppointmentSlot()
	statusID := chooseAppointmentStatusID()
	concernCategoryID := randomChoice(appointmentCategoryIDs).(int)

	// Determine admin_notes value
	var adminNotes sql.NullString
	statusName := ""
	for name, id := range appointmentStatusByName {
		if id == statusID {
			statusName = strings.ToLower(name)
			break
		}
	}
	if statusName == "cancelled" || statusName == "rejected" || strings.Contains(statusName, "show") {
		adminNotes = sql.NullString{String: gofakeit.Sentence(rand.Intn(5) + 5), Valid: true}
	} else {
		adminNotes = sql.NullString{Valid: false}
	}

	_, err := tx.Exec(`
		       INSERT INTO appointments (
			       user_id, reason, admin_notes, when_date, time_slot_id, appointment_category_id, status_id
		       ) VALUES (?, ?, ?, ?, ?, ?, ?)
	       `, userID,
		gofakeit.Sentence(rand.Intn(11)+20),
		adminNotes,
		whenDate,
		timeSlotID,
		concernCategoryID,
		statusID)
	if err != nil {
		log.Fatal(err)
	}
}

func reserveAppointmentSlot() (string, int) {
	if len(timeSlotIDs) == 0 {
		log.Fatal("no time slots found in time_slots")
	}

	for attempts := 0; attempts < 500; attempts++ {
		when := time.Now().AddDate(0, 0, rand.Intn(180)+1)
		weekday := when.Weekday()
		if weekday == time.Saturday || weekday == time.Sunday {
			continue // skip weekends
		}
		whenDate := when.Format("2006-01-02")
		timeSlotID := randomChoice(timeSlotIDs).(int)
		key := fmt.Sprintf("%s|%d", whenDate, timeSlotID)

		appointmentSlotMu.Lock()
		_, exists := reservedAppointmentSlots[key]
		if !exists {
			reservedAppointmentSlots[key] = struct{}{}
			appointmentSlotMu.Unlock()
			return whenDate, timeSlotID
		}
		appointmentSlotMu.Unlock()
	}

	log.Fatal("unable to reserve unique appointment slot after multiple attempts")
	return "", 0
}

func chooseAppointmentStatusID() int {
	if len(appointmentStatusByName) == 0 {
		return randomChoice(appointmentStatusIDs).(int)
	}

	type weightedStatus struct {
		id     int
		weight int
	}

	weighted := make([]weightedStatus, 0)
	used := make(map[int]bool)

	add := func(name string, weight int) {
		if id, ok := appointmentStatusByName[name]; ok {
			weighted = append(weighted, weightedStatus{id: id, weight: weight})
			used[id] = true
		}
	}

	add("pending", 50)
	add("approved", 20)
	add("completed", 15)
	add("cancelled", 10)
	add("rejected", 20)
	add("rescheduled", 10)

	for _, id := range appointmentStatusIDs {
		if !used[id] {
			weighted = append(weighted, weightedStatus{id: id, weight: 3})
		}
	}

	if len(weighted) == 0 {
		return randomChoice(appointmentStatusIDs).(int)
	}

	totalWeight := 0
	for _, item := range weighted {
		totalWeight += item.weight
	}

	roll := rand.Intn(totalWeight)
	for _, item := range weighted {
		roll -= item.weight
		if roll < 0 {
			return item.id
		}
	}

	return weighted[len(weighted)-1].id
}

func chooseCivilStatusID() int {
	if len(civilStatusByName) == 0 {
		return randomChoice(civilStatusIDs).(int)
	}

	type weightedStatus struct {
		id     int
		weight int
	}

	weighted := make([]weightedStatus, 0)
	used := make(map[int]bool)

	add := func(name string, weight int) {
		if id, ok := civilStatusByName[name]; ok {
			weighted = append(weighted, weightedStatus{id: id, weight: weight})
			used[id] = true
		}
	}

	add("single", 75)
	add("married", 8)
	add("widowed", 2)
	add("divorced", 2)
	add("separated", 6)

	for _, id := range civilStatusIDs {
		if !used[id] {
			weighted = append(weighted, weightedStatus{id: id, weight: 4})
		}
	}

	if len(weighted) == 0 {
		return randomChoice(civilStatusIDs).(int)
	}

	totalWeight := 0
	for _, item := range weighted {
		totalWeight += item.weight
	}

	roll := rand.Intn(totalWeight)
	for _, item := range weighted {
		roll -= item.weight
		if roll < 0 {
			return item.id
		}
	}

	return weighted[len(weighted)-1].id
}

func randomAwards() sql.NullString {
	awardPool := []string{
		"With Honors",
		"Best in Mathematics",
		"Best in Conduct",
		"Leadership Award",
		"Academic Excellence",
		"Science Fair Winner",
	}

	r := rand.Float32()
	if r < 0.35 {
		return sql.NullString{Valid: false}
	}

	if r < 0.7 {
		return sql.NullString{String: awardPool[rand.Intn(len(awardPool))], Valid: true}
	}

	count := rand.Intn(3) + 2 // 2-4 awards
	selected := pickUniqueStrings(awardPool, count)
	return sql.NullString{String: strings.Join(selected, ", "), Valid: true}
}

func pickUniqueStrings(pool []string, count int) []string {
	if count <= 0 || len(pool) == 0 {
		return []string{}
	}
	if count >= len(pool) {
		copyPool := append([]string{}, pool...)
		rand.Shuffle(len(copyPool), func(i, j int) {
			copyPool[i], copyPool[j] = copyPool[j], copyPool[i]
		})
		return copyPool
	}

	shuffled := append([]string{}, pool...)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled[:count]
}

func relationshipID(name string) int {
	if id, ok := relationshipTypeIDs[name]; ok {
		return id
	}
	for relName, id := range relationshipTypeIDs {
		if strings.EqualFold(relName, name) {
			return id
		}
	}
	log.Fatalf("relationship type not found: %s", name)
	return 0
}

// helper: random choice from slice
func randomChoice(slice interface{}) interface{} {
	switch v := slice.(type) {
	case []int:
		return v[rand.Intn(len(v))]
	case []string:
		return v[rand.Intn(len(v))]
	default:
		panic("unsupported type")
	}
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
