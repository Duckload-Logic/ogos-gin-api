package analytics

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetTotalStudents(ctx context.Context) (int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM student_personal_info").Scan(&total)
	return total, err
}

//personal info

func (r *Repository) GetAgeStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			CAST(TIMESTAMPDIFF(YEAR, spi.date_of_birth, CURDATE()) AS CHAR) AS category,
			SUM(CASE WHEN g.name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE spi.date_of_birth IS NOT NULL
		GROUP BY category
		ORDER BY CAST(category AS UNSIGNED) ASC;
	`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetCivilStatusStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(cs.name, 'Not Indicated') AS category,
			SUM(CASE WHEN g.name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN civil_statuses cs ON spi.civil_status_id = cs.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category
		ORDER BY rank_pos ASC;
	`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetReligionStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(rel.name, 'Not Indicated') AS category,
			SUM(CASE WHEN g.name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN religions rel ON spi.religion_id = rel.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category
		ORDER BY rank_pos ASC;
	`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetCityAddressStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(a.city, 'Not Indicated') AS category,
			SUM(CASE WHEN g.name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN addresses a ON spi.address_id = a.id 
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category
		ORDER BY rank_pos ASC;
	`
	return r.executeStatQuery(ctx, query)
}

//family & background

func (r *Repository) GetMonthlyIncomeStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(spi.monthly_income, 'Not Indicated') AS category,
			SUM(CASE WHEN g.name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category
		ORDER BY rank_pos ASC;
	`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetOrdinalPositionStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(spi.ordinal_position, 'Not Indicated') AS category,
			SUM(CASE WHEN g.name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category
		ORDER BY rank_pos ASC;
	`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetFatherEducationStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(spi.father_education, 'Not Indicated') AS category,
			SUM(CASE WHEN g.name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category
		ORDER BY rank_pos ASC;
	`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetMotherEducationStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(spi.mother_education, 'Not Indicated') AS category,
			SUM(CASE WHEN g.name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category
		ORDER BY rank_pos ASC;
	`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetParentsMaritalStatusStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(spi.parents_marital_status, 'Not Indicated') AS category,
			SUM(CASE WHEN g.name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category
		ORDER BY rank_pos ASC;
	`
	return r.executeStatQuery(ctx, query)
}

//academic background
func (r *Repository) GetHSGWAStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			CASE
				WHEN spi.high_school_gwa >= 97 THEN '97-99'
				WHEN spi.high_school_gwa >= 94 AND spi.high_school_gwa < 97 THEN '94-96'
				WHEN spi.high_school_gwa >= 91 AND spi.high_school_gwa < 94 THEN '91-93'
				WHEN spi.high_school_gwa >= 88 AND spi.high_school_gwa < 91 THEN '88-90'
				WHEN spi.high_school_gwa >= 85 AND spi.high_school_gwa < 88 THEN '85-87'
				WHEN spi.high_school_gwa >= 82 AND spi.high_school_gwa < 85 THEN '82-84'
				ELSE 'Not Indicated'
			END AS category,
			SUM(CASE WHEN g.name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category
		ORDER BY rank_pos ASC;
	`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetElementaryStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(spi.elem_school_type, 'Not Indicated') AS category,
			SUM(CASE WHEN g.name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category
		ORDER BY rank_pos ASC;
	`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetJuniorHighStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(spi.jhs_type, 'Not Indicated') AS category,
			SUM(CASE WHEN g.name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category
		ORDER BY rank_pos ASC;
	`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetSeniorHighStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(spi.shs_type, 'Not Indicated') AS category,
			SUM(CASE WHEN g.name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category
		ORDER BY rank_pos ASC;
	`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetNatureOfSchoolingStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(spi.nature_of_schooling, 'Not Indicated') AS category,
			SUM(CASE WHEN g.name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category
		ORDER BY rank_pos ASC;
	`
	return r.executeStatQuery(ctx, query)
}

//study environment

func (r *Repository) GetQuietStudyPlaceStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(spi.quiet_study_place, 'Not Indicated') AS category,
			SUM(CASE WHEN g.name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category
		ORDER BY rank_pos ASC;
	`
	return r.executeStatQuery(ctx, query)
}

//helper function

func (r *Repository) executeStatQuery(ctx context.Context, query string) ([]AggregatedStatModel, error) {
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []AggregatedStatModel
	for rows.Next() {
		var stat AggregatedStatModel
		if err := rows.Scan(&stat.Category, &stat.MaleCount, &stat.FemaleCount, &stat.Total, &stat.RankPos); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}
	return stats, nil
}