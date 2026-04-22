package analytics

import "time"

// DemographicStat represents a pure business entity for analytics.
type DemographicStat struct {
	Category    string
	MaleCount   int
	FemaleCount int
	Total       int
	RankPos     int
}

// StudentPersonalInfo represents pure student analytics data.
type StudentPersonalInfo struct {
	ID              int
	IIRID           string
	StudentNumber   string
	GenderID        int
	CivilStatusID   int
	ReligionID      int
	HeightFt        float64
	WeightKg        float64
	Complexion      string
	HighSchoolGWA   float64
	CourseID        int
	YearLevel       int
	Section         int
	PlaceOfBirth    string
	DateOfBirth     time.Time
	IsEmployed      bool
	EmployerName    *string
	EmployerAddress *string
	MobileNumber    string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type StudentFinances struct {
	ID                         int
	IIRID                      string
	MonthlyFamilyIncomeRangeID *int
	WeeklyAllowance            float64
}

type FamilyBackground struct {
	IIRID                 string
	ParentalStatusID      int
	OrdinalPosition       int
	HaveQuietPlaceToStudy bool
}

type EducationalBackground struct {
	ID               int
	IIRID            string
	EducationLevelID int
	SchoolDetailID   int
	GWA              float64
}
