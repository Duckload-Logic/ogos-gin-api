package analytics

import "time"

type AggregatedStatModel struct {
	Category    string `db:"category"`
	MaleCount   int    `db:"male_count"`
	FemaleCount int    `db:"female_count"`
	Total       int    `db:"total"`
	RankPos     int    `db:"rank_pos"`
}

// StudentPersonalInfo 
type StudentPersonalInfo struct {
	ID            int       `db:"id"`
	IIRID         int       `db:"iir_id"`
	StudentNumber string    `db:"student_number"`
	GenderID      int       `db:"gender_id"`
	CivilStatusID int       `db:"civil_status_id"`
	ReligionID    int       `db:"religion_id"`
	HeightFt      float64   `db:"height_ft"`
	WeightKg      float64   `db:"weight_kg"`
	Complexion    string    `db:"complexion"`
	HighSchoolGWA float64   `db:"high_school_gwa"`
	CourseID      int       `db:"course_id"`
	YearLevel     int       `db:"year_level"`
	Section       int       `db:"section"`
	PlaceOfBirth  string    `db:"place_of_birth"`
	DateOfBirth   time.Time `db:"date_of_birth"`
	IsEmployed    bool      `db:"is_employed"`
	EmployerName  *string   `db:"employer_name"`
	EmployerAddress *string   `db:"employer_address"`
	MobileNumber  string    `db:"mobile_number"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// StudentFinances
type StudentFinances struct {
	ID                         int      `db:"id"`
	IIRID                      int      `db:"iir_id"`
	MonthlyFamilyIncomeRangeID *int     `db:"monthly_family_income_range_id"`
	WeeklyAllowance            float64  `db:"weekly_allowance"`
}

// FamilyBackground 
type FamilyBackground struct {
	IIRID                 int  `db:"iir_id"`
	ParentalStatusID      int  `db:"parental_status_id"`
	OrdinalPosition       int  `db:"ordinal_position"`
	HaveQuietPlaceToStudy bool `db:"have_quiet_place_to_study"`
}