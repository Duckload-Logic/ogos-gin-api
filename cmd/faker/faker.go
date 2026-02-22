package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

// lookup slices (IDs)
var (
	genderIDs             []int
	civilStatusIDs        []int
	religionIDs           []int
	courseIDs             []int
	enrollmentReasonIDs   []int
	supportTypeIDs        []int
	incomeRangeIDs        []int
	parentalStatusIDs     []int
	educationalLevelIDs   []int
	relationshipTypeIDs   map[string]int
	natureOfResidenceIDs  []int
	siblingSupportTypeIDs []int
	activityOptionIDs     []int
)

func main() {
	// ---------- CONFIGURATION ----------
	numStudents := 50  // number of students to generate
	numCounselors := 1 // number of counselors (admins)
	_ = godotenv.Load()
	dsn := buildDSNFromEnv()
	// -----------------------------------

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Fatal("db not reachable:", err)
	}

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

	// create students
	var wg sync.WaitGroup
	errChan := make(chan error, numStudents)
	for i := 0; i < numStudents; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			createStudent(i)
		}(i)
	}

	wg.Wait()

	if len(errChan) > 0 {
		log.Fatal("errors occurred during data generation:", <-errChan)
		return
	}

	fmt.Println("Dummy data generation completed successfully.")
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
	rows, err = db.Query("SELECT id FROM civil_status_types")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		civilStatusIDs = append(civilStatusIDs, id)
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
	rows, err = db.Query("SELECT id FROM educational_levels")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		educationalLevelIDs = append(educationalLevelIDs, id)
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
	tx, err := db.Begin()
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
	`, firstName, randomMiddleName(), lastName, email, fakePasswordHash(), true)
	if err != nil {
		log.Fatal(err)
	}
	userID, _ := res.LastInsertId()

	_, err = tx.Exec(`
		INSERT INTO counselor_profiles (user_id, license_number, specialization, is_available)
		VALUES (?, ?, ?, ?)
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
	tx, err := db.Begin()
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

	// 2. iir_records
	res, err = tx.Exec(`
		INSERT INTO iir_records (user_id, is_submitted) VALUES (?, ?)
	`, userID, gofakeit.Bool())
	if err != nil {
		log.Fatal(err)
	}
	iirID, _ := res.LastInsertId()

	// 3. student_personal_info
	insertPersonalInfo(tx, int(iirID), birthYear, dob)

	// 4. selected reasons
	insertSelectedReasons(tx, int(iirID))

	// 5. addresses (residential & provincial)
	resAddr1 := insertAddress(tx)
	resAddr2 := insertAddress(tx)
	insertStudentAddress(tx, int(iirID), resAddr1, "Residential")
	insertStudentAddress(tx, int(iirID), resAddr2, "Provincial")

	// 6. related persons (father, mother, guardian)
	fatherID := insertRelatedPerson(tx, nil)
	motherID := insertRelatedPerson(tx, nil)

	// Randomly decide who is the guardian: 0=father, 1=mother, 2=separate person
	guardianChoice := rand.Intn(3)
	var guardianID int

	switch guardianChoice {
	case 0: // Father is guardian
		guardianID = fatherID
		// Update father's address to resAddr2
		_, err = tx.Exec("UPDATE related_persons SET address_id = ? WHERE id = ?", resAddr2, fatherID)
		if err != nil {
			log.Fatal(err)
		}
	case 1: // Mother is guardian
		guardianID = motherID
		// Update mother's address to resAddr2
		_, err = tx.Exec("UPDATE related_persons SET address_id = ? WHERE id = ?", resAddr2, motherID)
		if err != nil {
			log.Fatal(err)
		}
	case 2: // Separate guardian
		guardianID = insertRelatedPerson(tx, &resAddr2)
	}

	// 7. link related persons
	if guardianChoice == 0 {
		// Father is also guardian
		linkRelatedPerson(tx, int(iirID), fatherID, "Father", true, true, true, true)
		linkRelatedPerson(tx, int(iirID), motherID, "Mother", true, false, false, true)
	} else if guardianChoice == 1 {
		// Mother is also guardian
		linkRelatedPerson(tx, int(iirID), fatherID, "Father", true, false, false, true)
		linkRelatedPerson(tx, int(iirID), motherID, "Mother", true, true, true, true)
	} else {
		// Separate guardian
		linkRelatedPerson(tx, int(iirID), fatherID, "Father", true, false, false, true)
		linkRelatedPerson(tx, int(iirID), motherID, "Mother", true, false, false, true)
		linkRelatedPerson(tx, int(iirID), guardianID, "Guardian", false, true, true, true)
	}

	// 8. family background
	familyBgID := insertFamilyBackground(tx, int(iirID))

	// 9. sibling supports (if employed siblings > 0)
	insertSiblingSupports(tx, familyBgID)

	// 10. educational background
	ebID := insertEducationalBackground(tx, int(iirID))

	// 11. school details for each educational level
	insertSchoolDetails(tx, ebID, birthYear)

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

	// 22. appointment (20% chance)
	if rand.Float32() < 0.2 {
		insertAppointment(tx, int(userID))
	}

	tx.Commit()
	fmt.Printf("Created student %d (iirID=%d)\n", index+1, iirID)
}

// ----------------------------------------------------------------------
// INDIVIDUAL INSERT HELPERS
// ----------------------------------------------------------------------

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

func insertPersonalInfo(tx *sql.Tx, iirID, birthYear int, dob time.Time) {
	studentNumber := fmt.Sprintf("%d-%04d-TG-%d", time.Now().Year(), iirID, rand.Intn(10))
	isEmployed := gofakeit.Bool()
	var employerName, employerAddress sql.NullString
	if isEmployed {
		empName := gofakeit.Company()
		empAddr := gofakeit.Address().Address
		employerName = sql.NullString{String: empName, Valid: true}
		employerAddress = sql.NullString{String: empAddr, Valid: true}
	}

	_, err := tx.Exec(`
		INSERT INTO student_personal_info (
			iir_id, student_number, gender_id, civil_status_id, religion_id,
			height_ft, weight_kg, complexion, high_school_gwa, course_id,
			year_level, section, place_of_birth, date_of_birth,
			is_employed, employer_name, employer_address, contact_number
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, iirID, studentNumber,
		randomChoice(genderIDs), randomChoice(civilStatusIDs), randomChoice(religionIDs),
		gofakeit.Float64Range(4.5, 6.5), gofakeit.Float64Range(40, 100), gofakeit.Color(),
		gofakeit.Float64Range(75, 98), randomChoice(courseIDs),
		rand.Intn(4)+1, rand.Intn(5)+1,
		gofakeit.City(), dob.Format("2006-01-02"),
		isEmployed, employerName, employerAddress, gofakeit.Phone(),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func insertSelectedReasons(tx *sql.Tx, iirID int) {
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

func insertAddress(tx *sql.Tx) int {
	// generate a random Philippine-like address
	region := gofakeit.RandomString([]string{"NCR", "Region IV-A", "Region III", "Region VII"})
	city := gofakeit.City()
	barangay := "Barangay " + gofakeit.StreetName()
	street := gofakeit.Street()

	res, err := tx.Exec(`
		INSERT INTO addresses (region, city, barangay, street_detail)
		VALUES (?, ?, ?, ?)
	`, region, city, barangay, street)
	if err != nil {
		log.Fatal(err)
	}
	id, _ := res.LastInsertId()
	return int(id)
}

func insertStudentAddress(tx *sql.Tx, iirID, addressID int, addrType string) {
	_, err := tx.Exec(`
		INSERT INTO student_addresses (iir_id, address_id, address_type)
		VALUES (?, ?, ?)
	`, iirID, addressID, addrType)
	if err != nil {
		log.Fatal(err)
	}
}

func insertRelatedPerson(tx *sql.Tx, addressID *int) int {
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
			address_id, educational_level, date_of_birth,
			last_name, first_name, middle_name,
			occupation, employer_name, employer_address, contact_number
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, addressID, educ, dob, lname, fname, mname,
		occupation, employer, employerAddr, contact)
	if err != nil {
		log.Fatal(err)
	}
	id, _ := res.LastInsertId()
	return int(id)
}

func randomEducationalAttainment() string {
	levels := []string{"Elementary", "High School", "Vocational", "College", "Post Graduate"}
	return levels[rand.Intn(len(levels))]
}

func linkRelatedPerson(tx *sql.Tx, iirID, personID int, relType string, isParent, isGuardian, isEmergency, isLiving bool) {
	relID := relationshipTypeIDs[relType]
	_, err := tx.Exec(`
		INSERT INTO student_related_persons (
			iir_id, related_person_id, relationship_id,
			is_parent, is_guardian, is_emergency_contact, is_living
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`, iirID, personID, relID, isParent, isGuardian, isEmergency, isLiving)
	if err != nil {
		log.Fatal(err)
	}
}

func insertFamilyBackground(tx *sql.Tx, iirID int) int {
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

func insertSiblingSupports(tx *sql.Tx, familyBgID int) {
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

func insertEducationalBackground(tx *sql.Tx, iirID int) int {
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
func insertSchoolDetails(tx *sql.Tx, ebID, birthYear int) {
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

	for _, levelID := range educationalLevelIDs {
		// get level name
		var levelName string
		err := db.QueryRow("SELECT level_name FROM educational_levels WHERE id = ?", levelID).Scan(&levelName)
		if err != nil {
			log.Fatal(err)
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
		awards := sql.NullString{String: gofakeit.Sentence(3), Valid: gofakeit.Bool()}

		_, err = tx.Exec(`
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

func insertHealthRecords(tx *sql.Tx, iirID int) {
	visionProb := gofakeit.Bool()
	visionDet := nullStringIf(visionProb, gofakeit.Sentence(3))
	hearingProb := gofakeit.Bool()
	hearingDet := nullStringIf(hearingProb, gofakeit.Sentence(3))
	speechProb := gofakeit.Bool()
	speechDet := nullStringIf(speechProb, gofakeit.Sentence(3))
	genProb := gofakeit.Bool()
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

func insertConsultations(tx *sql.Tx, iirID int) {
	if rand.Float32() < 0.4 { // 40% have at least one consultation
		num := rand.Intn(3) + 1
		for i := 0; i < num; i++ {
			profType := randomChoice([]string{"Psychiatrist", "Psychologist", "Counselor"})
			has := gofakeit.Bool()
			var when, what sql.NullString
			if has {
				when = sql.NullString{String: gofakeit.Date().Format("2006-01-02"), Valid: true}
				what = sql.NullString{String: gofakeit.Sentence(8), Valid: true}
			}
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

func insertTestResults(tx *sql.Tx, iirID int) {
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

func insertSignificantNotes(tx *sql.Tx, iirID int) {
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

func insertStudentFinances(tx *sql.Tx, iirID int) int {
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

func insertFinancialSupports(tx *sql.Tx, sfID int) {
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

func insertActivities(tx *sql.Tx, iirID int) {
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

func insertSubjectPreferences(tx *sql.Tx, iirID int) {
	num := rand.Intn(6) // 0-5 subjects
	subjects := []string{"Math", "Science", "English", "History", "PE", "Art", "Music", "Computer"}
	for i := 0; i < num; i++ {
		subj := subjects[rand.Intn(len(subjects))]
		fav := gofakeit.Bool()
		_, err := tx.Exec(`
			INSERT INTO student_subject_preferences (iir_id, subject_name, is_favorite)
			VALUES (?, ?, ?)
		`, iirID, subj, fav)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func insertHobbies(tx *sql.Tx, iirID int) {
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

func insertAdmissionSlip(tx *sql.Tx, iirID int) {
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

func insertAppointment(tx *sql.Tx, userID int) {
	status := randomChoice([]string{"Pending", "Approved", "Completed", "Cancelled", "Rescheduled"})
	_, err := tx.Exec(`
		INSERT INTO appointments (
			user_id, reason, scheduled_date, scheduled_time, concern_category, status
		) VALUES (?, ?, ?, ?, ?, ?)
	`, userID,
		gofakeit.Sentence(4),
		gofakeit.FutureDate().Format("2006-01-02"),
		gofakeit.Date().Format("15:04:05"),
		gofakeit.Word(),
		status)
	if err != nil {
		log.Fatal(err)
	}
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
