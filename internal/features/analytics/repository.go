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
	err := r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM student_personal_info")
	return total, err
}

// --- PERSONAL INFORMATION ---

func (r *Repository) GetAgeStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			CAST(TIMESTAMPDIFF(YEAR, spi.date_of_birth, CURDATE()) AS CHAR) AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE spi.date_of_birth IS NOT NULL
		GROUP BY category;`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetCivilStatusStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(status_name, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN civil_status_types cs ON spi.civil_status_id = cs.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category
		ORDER BY rank_pos ASC;`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetReligionStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(religion_name, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN religions rel ON spi.religion_id = rel.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category
		ORDER BY rank_pos ASC;`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetCityAddressStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(sa.city, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN student_addresses sa ON spi.iir_id = sa.iir_id
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category;`
	return r.executeStatQuery(ctx, query)
}
// --- FAMILY & FINANCIAL BACKGROUND ---

func (r *Repository) GetMonthlyIncomeStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(ir.range_text, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN student_finances sf ON spi.iir_id = sf.iir_id
		LEFT JOIN income_ranges ir ON sf.monthly_family_income_range_id = ir.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category;`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetOrdinalPositionStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			CAST(fb.ordinal_position AS CHAR) AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN family_backgrounds fb ON spi.iir_id = fb.iir_id
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category;`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetFatherEducationStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(rp.educational_level, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN related_persons rp ON spi.iir_id = rp.iir_id
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE rp.relationship = 'Father'
		GROUP BY category;`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetMotherEducationStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(rp.educational_level, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN related_persons rp ON spi.iir_id = rp.iir_id
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE rp.relationship = 'Mother'
		GROUP BY category;`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetParentsMaritalStatusStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(pst.status_name, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN family_backgrounds fb ON spi.iir_id = fb.iir_id
		LEFT JOIN parental_status_types pst ON fb.parental_status_id = pst.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category;`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetQuietStudyPlaceStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			CASE WHEN fb.have_quiet_place_to_study = 1 THEN 'Yes' ELSE 'No' END AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN family_backgrounds fb ON spi.iir_id = fb.iir_id
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category;`
	return r.executeStatQuery(ctx, query)
}

// --- ACADEMIC BACKGROUND ---

func (r *Repository) GetHSGWAStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			CASE 
				WHEN spi.high_school_gwa >= 95.00 THEN '95.00 - 100.00'
				WHEN spi.high_school_gwa >= 90.00 THEN '90.00 - 94.99'
				WHEN spi.high_school_gwa >= 85.00 THEN '85.00 - 89.99'
				WHEN spi.high_school_gwa >= 80.00 THEN '80.00 - 84.99'
				WHEN spi.high_school_gwa >= 75.00 THEN '75.00 - 79.99'
				ELSE 'Below 75.00'
			END AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE spi.high_school_gwa IS NOT NULL AND spi.high_school_gwa > 0
		GROUP BY category;`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetElementaryStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(sd.school_type, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN educational_backgrounds eb ON spi.iir_id = eb.iir_id
		JOIN educational_levels el ON eb.education_level_id = el.id
		LEFT JOIN school_details sd ON eb.school_detail_id = sd.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE el.level_name = 'Elementary'
		GROUP BY category;`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetJuniorHighStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(sd.school_type, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN educational_backgrounds eb ON spi.iir_id = eb.iir_id
		JOIN educational_levels el ON eb.education_level_id = el.id
		LEFT JOIN school_details sd ON eb.school_detail_id = sd.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE el.level_name = 'Junior High School'
		GROUP BY category;`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetSeniorHighStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(sd.school_type, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN educational_backgrounds eb ON spi.iir_id = eb.iir_id
		JOIN educational_levels el ON eb.education_level_id = el.id
		LEFT JOIN school_details sd ON eb.school_detail_id = sd.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE el.level_name = 'Senior High School'
		GROUP BY category;`
	return r.executeStatQuery(ctx, query)
}

func (r *Repository) GetNatureOfSchoolingStats(ctx context.Context) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(eb.nature_of_schooling, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN educational_backgrounds eb ON spi.iir_id = eb.iir_id
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category;`
	return r.executeStatQuery(ctx, query)
}

// --- HELPERS ---

func (r *Repository) executeStatQuery(ctx context.Context, query string) ([]AggregatedStatModel, error) {
	stats := make([]AggregatedStatModel, 0)
	err := r.db.SelectContext(ctx, &stats, query)
	return stats, err
}

func (r *Repository) getSchoolTypeStats(ctx context.Context, levelName string) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(sd.school_type, 'Not Indicated') AS category,
			SUM(CASE WHEN g.name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN educational_backgrounds eb ON spi.iir_id = eb.iir_id
		JOIN educational_levels el ON eb.education_level_id = el.id
		JOIN school_details sd ON eb.school_detail_id = sd.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE el.name = ?
		GROUP BY category;`
	
	stats := make([]AggregatedStatModel, 0)
	err := r.db.SelectContext(ctx, &stats, query, levelName)
	return stats, err
}

func (r *Repository) getParentEducationStats(ctx context.Context, parentType string) ([]AggregatedStatModel, error) {
	query := `
		SELECT 
			COALESCE(ed.level_name, 'Not Indicated') AS category,
			SUM(CASE WHEN g.name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN family_backgrounds fb ON spi.iir_id = fb.iir_id
		LEFT JOIN educational_levels ed ON (CASE WHEN ? = 'Father' THEN fb.father_education_id ELSE fb.mother_education_id END) = ed.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		GROUP BY category;`
	
	stats := make([]AggregatedStatModel, 0)
	err := r.db.SelectContext(ctx, &stats, query, parentType)
	return stats, err
}