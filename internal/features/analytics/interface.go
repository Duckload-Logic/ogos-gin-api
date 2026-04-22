package analytics

import "context"

type ServiceInterface interface {
	GetDashboard(
		ctx context.Context,
		year int,
		courseID int,
	) (*DashboardResponseDTO, error)
	GetAdminDashboard(
		ctx context.Context,
		timeRange string,
		source string,
	) (*AdminDashboardResponseDTO, error)
}

type RepositoryInterface interface {
	GetTotalStudents(ctx context.Context, year int, courseID int) (int, error)
	GetGenderStats(
		ctx context.Context,
		year int,
		courseID int,
	) ([]DemographicStat, error)
	GetAgeStats(
		ctx context.Context,
		year int,
		courseID int,
	) ([]DemographicStat, error)
	GetCivilStatusStats(
		ctx context.Context,
		year int,
		courseID int,
	) ([]DemographicStat, error)
	GetReligionStats(
		ctx context.Context,
		year int,
		courseID int,
	) ([]DemographicStat, error)
	GetCityAddressStats(
		ctx context.Context,
		year int,
		courseID int,
	) ([]DemographicStat, error)
	GetMonthlyIncomeStats(
		ctx context.Context,
		year int,
		courseID int,
	) ([]DemographicStat, error)
	GetOrdinalPositionStats(
		ctx context.Context,
		year int,
		courseID int,
	) ([]DemographicStat, error)
	GetFatherEducationStats(
		ctx context.Context,
		year int,
		courseID int,
	) ([]DemographicStat, error)
	GetMotherEducationStats(
		ctx context.Context,
		year int,
		courseID int,
	) ([]DemographicStat, error)
	GetParentsMaritalStatusStats(
		ctx context.Context, year int, courseID int,
	) ([]DemographicStat, error)
	GetQuietStudyPlaceStats(
		ctx context.Context,
		year int,
		courseID int,
	) ([]DemographicStat, error)
	GetHSGWAStats(
		ctx context.Context,
		year int,
		courseID int,
	) ([]DemographicStat, error)
	GetElementaryStats(
		ctx context.Context,
		year int,
		courseID int,
	) ([]DemographicStat, error)
	GetJuniorHighStats(
		ctx context.Context,
		year int,
		courseID int,
	) ([]DemographicStat, error)
	GetSeniorHighStats(
		ctx context.Context,
		year int,
		courseID int,
	) ([]DemographicStat, error)
	GetNatureOfSchoolingStats(
		ctx context.Context, year int, courseID int,
	) ([]DemographicStat, error)
	GetTotalReports(ctx context.Context) (int, error)
	GetTotalAppointments(ctx context.Context) (int, error)
	GetTotalSlips(ctx context.Context) (int, error)
	GetMonthlyVisitorStats(
		ctx context.Context,
		timeRange string,
	) ([]MonthlyVisitorStatDTO, error)
	GetMonthlyAppointmentStats(
		ctx context.Context,
		timeRange string,
	) ([]MonthlyVisitorStatDTO, error)
}
