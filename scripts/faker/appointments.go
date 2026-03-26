package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/hash"
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/slips"
)

func insertAdmissionSlip(
	ctx context.Context,
	tx *sqlx.Tx,
	iirID string,
) string {
	// More realistic status distribution (pending less likely for
	// historical data)
	statusName := chooseAdmissionSlipStatus()
	statusID, ok := admissionSlipStatusesByName[statusName]
	if !ok {
		log.Printf(
			"[Seeder] {Insert AdmissionSlip}: status '%s' not found",
			statusName)
		statusID = randomChoice(admissionSlipStatusIDs).(int)
	}
	categoryID := randomChoice(admissionSlipCategoryIDs).(int)

	// Realistic reasons for admission slips
	reason := generateRealisticAdmissionReason()

	daysAgo := rand.Intn(7) + 2 // 2-8 days ago
	dateOfAbsence := time.Now().AddDate(0, 0, -daysAgo).Format(
		"2006-01-02")

	// Date needed: submission should be a few days to weeks after
	// absence
	daysAfterAbsence := rand.Intn(7) + 1
	targetDate := time.Now().AddDate(0, 0, -daysAgo+daysAfterAbsence)

	// If targetDate is before "Now", set it to "Now"
	if targetDate.Before(time.Now()) {
		targetDate = time.Now()
	}

	dateNeeded := targetDate.Format("2006-01-02")

	// Use PDF as primary extension (more realistic for official
	// documents)
	extensions := []string{
		".pdf", ".pdf", ".pdf", ".jpg", ".jpeg",
		".png",
	}
	ext := extensions[rand.Intn(len(extensions))]
	basePath := "./uploads"
	subFolder := gofakeit.UUID()

	// More realistic file names for admission documents
	baseFileName := generateAdmissionFileName()
	readableFileName := baseFileName + ext
	fileName := hash.GetSHA256Hash(readableFileName, 16) + ext

	// DISK PATH (Where Go writes the bytes)
	fullStoragePath := filepath.Join(basePath, "slips", subFolder,
		fileName)

	dbURL := fmt.Sprintf("/slips/%s/%s", subFolder, fileName)
	// Create directory and file
	dir := filepath.Dir(fullStoragePath)
	os.MkdirAll(dir, os.ModePerm)
	f, _ := os.Create(fullStoragePath)
	// Write realistic dummy content
	content := fmt.Sprintf(
		"ADMISSION SLIP / EXCUSE SLIP\nStudent ID: %s\nDate: %s\n"+
			"Reason: %s\n\n[Document content created for admission "+
			"purposes]",
		iirID, dateOfAbsence, reason)
	f.WriteString(content)
	f.Close()

	// Admin notes more realistic based on status
	adminNotes := generateAdmissionNotes(statusName)
	admissionSlipID := uuid.New().String()

	slip := &slips.Slip{
		ID:            admissionSlipID,
		IIRID:         iirID,
		Reason:        reason,
		CategoryID:    categoryID,
		DateOfAbsence: dateOfAbsence,
		DateNeeded:    dateNeeded,
		StatusID:      statusID,
		AdminNotes:    adminNotes,
	}

	_, err := slipsRepo.CreateSlip(ctx, tx, slip)
	if err != nil {
		log.Printf("[Seeder] {Insert AdmissionSlip}: %v", err)
		return ""
	}

	attachment := &slips.SlipAttachment{
		ID:       uuid.New().String(),
		SlipID:   &admissionSlipID,
		FileName: readableFileName,
		FileURL:  dbURL,
	}

	err = slipsRepo.SaveSlipAttachment(ctx, tx, attachment)
	if err != nil {
		log.Printf("[Seeder] {Insert SlipAttachment}: %v", err)
		return ""
	}

	return admissionSlipID
}

func chooseAdmissionSlipStatus() string {
	statuses := []struct {
		name   string
		weight int
	}{
		{"approved", 50},
		{"pending", 25},
		{"for revision", 15},
		{"rejected", 10},
	}

	totalWeight := 0
	for _, s := range statuses {
		totalWeight += s.weight
	}

	roll := rand.Intn(totalWeight)
	for _, s := range statuses {
		roll -= s.weight
		if roll < 0 {
			return s.name
		}
	}
	return "pending"
}

func generateRealisticAdmissionReason() string {
	reasons := []string{
		"Medical consultation - Doctor's appointment",
		"Family emergency requiring immediate attention",
		"Legal matter - Court appearance required",
		"University business - Office of the Registrar",
		"Health center visit - Medical examination",
		"Dental appointment - Preventive care",
		"Job interview - Employment opportunity",
		"Scholarship office - Financial aid matters",
		"Hospital admission - Medical treatment",
		"Burial ceremony - Family death",
		"Important family matter - Guardian meeting",
		"University activity - Student organization duty",
		"Mental health counseling session",
		"Physical therapy appointment",
		"Government office work - Official business",
	}
	return reasons[rand.Intn(len(reasons))]
}

func generateAdmissionFileName() string {
	fileTypes := []string{
		"medical_certificate",
		"excuse_letter",
		"hospital_clearance",
		"doctor_note",
		"appointment_slip",
		"official_letter",
		"emergency_report",
		"proof_of_attendance",
	}
	return fileTypes[rand.Intn(len(fileTypes))]
}

func generateAdmissionNotes(status string) sql.NullString {
	switch status {
	case "approved":
		notes := []string{
			"Verified with hospital records",
			"Legitimate documentation provided",
			"Approved - excused absence",
			"Admission approved as submitted",
			"",
		}
		selected := notes[rand.Intn(len(notes))]
		if selected == "" {
			return sql.NullString{Valid: false}
		}
		return sql.NullString{String: selected, Valid: true}
	case "rejected":
		notes := []string{
			"Insufficient documentation",
			"Outside acceptable timeframe",
			"Documentation not properly verified",
			"Required supporting documents missing",
		}
		return sql.NullString{String: notes[rand.Intn(len(notes))], Valid: true}
	case "for revision":
		notes := []string{
			"Requires additional supporting documents",
			"Need clarification on dates",
			"Missing required signatures",
			"Please provide recent medical certificate",
		}
		return sql.NullString{String: notes[rand.Intn(len(notes))], Valid: true}
	default:
		return sql.NullString{Valid: false}
	}
}

func insertAppointment(ctx context.Context, tx *sqlx.Tx, iirID string) string {
	if len(timeSlotIDs) == 0 || len(appointmentCategoryIDs) == 0 ||
		len(appointmentStatusIDs) == 0 {
		log.Printf(
			"[Seeder] {Insert Appointment}: missing lookup data for "+
				"iir %s",
			iirID)
		return ""
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
	if statusName == "cancelled" || statusName == "rejected" ||
		strings.Contains(statusName, "show") {
		adminNotes = sql.NullString{
			String: gofakeit.Sentence(rand.Intn(5) + 5),
			Valid:  true,
		}
	} else {
		adminNotes = sql.NullString{Valid: false}
	}

	appointmentID := uuid.New().String()

	appt := &appointments.Appointment{
		ID:    appointmentID,
		IIRID: iirID,
		Reason: sql.NullString{
			String: gofakeit.Sentence(rand.Intn(11) + 20),
			Valid:  true,
		},
		AdminNotes:            adminNotes,
		WhenDate:              whenDate,
		TimeSlotID:            timeSlotID,
		AppointmentCategoryID: concernCategoryID,
		StatusID:              statusID,
	}

	err := appointmentsRepo.CreateAppointment(ctx, tx, appt)
	if err != nil {
		log.Printf("[Seeder] {Insert Appointment}: %v", err)
		return ""
	}

	return appointmentID
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

	log.Fatal(
		"unable to reserve unique appointment slot after multiple attempts",
	)
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
