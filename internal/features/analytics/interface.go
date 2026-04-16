package analytics

import "context"

type ServiceInterface interface {
	GetDashboard(ctx context.Context) (*DashboardResponseDTO, error)
	GetAdminDashboard(ctx context.Context, timeRange string) (*AdminDashboardResponseDTO, error)
}

type RepositoryInterface interface {
	GetTotalStudents(ctx context.Context) (int, error)
	GetAgeStats(ctx context.Context) ([]AggregatedStatModel, error)
	GetCivilStatusStats(ctx context.Context) ([]AggregatedStatModel, error)
	GetReligionStats(ctx context.Context) ([]AggregatedStatModel, error)
	GetCityAddressStats(ctx context.Context) ([]AggregatedStatModel, error)
	GetMonthlyIncomeStats(ctx context.Context) ([]AggregatedStatModel, error)
	GetOrdinalPositionStats(ctx context.Context) ([]AggregatedStatModel, error)
	GetFatherEducationStats(ctx context.Context) ([]AggregatedStatModel, error)
	GetMotherEducationStats(ctx context.Context) ([]AggregatedStatModel, error)
	GetParentsMaritalStatusStats(
		ctx context.Context,
	) ([]AggregatedStatModel, error)
	GetQuietStudyPlaceStats(ctx context.Context) ([]AggregatedStatModel, error)
	GetHSGWAStats(ctx context.Context) ([]AggregatedStatModel, error)
	GetElementaryStats(ctx context.Context) ([]AggregatedStatModel, error)
	GetJuniorHighStats(ctx context.Context) ([]AggregatedStatModel, error)
	GetSeniorHighStats(ctx context.Context) ([]AggregatedStatModel, error)
	GetNatureOfSchoolingStats(
		ctx context.Context,
	) ([]AggregatedStatModel, error)
	GetTotalReports(ctx context.Context) (int, error)
	GetTotalAppointments(ctx context.Context) (int, error)
	GetTotalSlips(ctx context.Context) (int, error)
	GetMonthlyVisitorStats(ctx context.Context, timeRange string) ([]MonthlyVisitorStatDTO, error)
}
