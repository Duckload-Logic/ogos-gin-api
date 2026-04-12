package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/slips"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

func main() {
	// ---------- CONFIGURATION ----------
	numStudents := 50  // number of students to generate
	numCounselors := 5 // number of counselors (admins)
	numSuperAdmin := 4 // number of super admins
	numWorkers := 10   // number of concurrent student workers
	numDevelopers := 2 // number of developers
	passwordHash := fakePasswordHash()
	_ = godotenv.Load()
	dsn := buildDSNFromEnv()
	fmt.Println(dsn)
	startTime := time.Now()
	// -----------------------------------

	var err error
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}
	defer db.Close()

	usersRepo = users.NewRepository(db)
	studentsRepo = students.NewRepository(db)
	appointmentsRepo = appointments.NewRepository(db)
	slipsRepo = slips.NewRepository(db)
	locationsRepo = locations.NewRepository(db)

	// load all lookup data from reference tables
	loadLookups()

	// seed random generator
	gofakeit.Seed(time.Now().UnixNano())

	// clear existing student data (optional but keeps the run idempotent)
	clearStudentData()

	// create super admin
	superadminJobs := make(chan int, numSuperAdmin)
	var superadminWg sync.WaitGroup
	for sw := 0; sw < numSuperAdmin; sw++ {
		superadminWg.Add(1)
		go func() {
			defer superadminWg.Done()
			for index := range superadminJobs {
				createSuperAdmin(index, passwordHash)
			}
		}()
	}

	for i := 0; i < numSuperAdmin; i++ {
		superadminJobs <- i
	}
	close(superadminJobs)
	superadminWg.Wait()

	// create counselors
	counselorJobs := make(chan int, numCounselors)
	var counselorWg sync.WaitGroup
	for cw := 0; cw < numCounselors; cw++ {
		counselorWg.Add(1)
		go func() {
			defer counselorWg.Done()
			for index := range counselorJobs {
				createCounselor(index, passwordHash)
			}
		}()
	}

	for i := 0; i < numCounselors; i++ {
		counselorJobs <- i
	}
	close(counselorJobs)
	counselorWg.Wait()

	// create students (concurrently with a worker pool)
	jobs := make(chan int, numStudents)
	var wg sync.WaitGroup

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				createStudent(index, passwordHash)
			}
		}()
	}

	for i := 0; i < numStudents; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	// create developers
	developerJobs := make(chan int, numDevelopers)
	var developerWg sync.WaitGroup
	for dw := 0; dw < numDevelopers; dw++ {
		developerWg.Add(1)
		go func() {
			defer developerWg.Done()
			createDeveloper(dw, passwordHash)
		}()
	}

	for i := 0; i < numDevelopers; i++ {
		developerJobs <- i
	}
	close(developerJobs)
	developerWg.Wait()

	fmt.Println("Dummy data generation completed successfully.")
	log.Println("Time taken:", time.Since(startTime))
}

// clearStudentData removes existing student data (keeps reference tables and
// seeded admin).
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
		"notifications", // added because notifications now reference user_id
		"users",
	}

	// disable FK checks temporarily
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	for _, tbl := range tables {
		if _, err := db.Exec(
			fmt.Sprintf("TRUNCATE TABLE %s", tbl),
		); err != nil {
			log.Fatal(err)
		}
	}
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}
