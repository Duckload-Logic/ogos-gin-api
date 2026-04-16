package analytics

import (
	"context"
	"math"

	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Service struct {
	repo  RepositoryInterface
	redis *datastore.RedisClient
}

func NewService(
	repo RepositoryInterface,
	redis *datastore.RedisClient,
) *Service {
	return &Service{repo: repo, redis: redis}
}

func (s *Service) GetDashboard(
	ctx context.Context,
) (*DashboardResponseDTO, error) {
	total, err := s.repo.GetTotalStudents(ctx)
	if err != nil {
		return nil, err
	}

	dashboard := &DashboardResponseDTO{
		TotalStudents:        total,
		AgeDistribution:      []DemographicStatDTO{},
		CivilStatus:          []DemographicStatDTO{},
		Religions:            []DemographicStatDTO{},
		CityAddress:          []DemographicStatDTO{},
		MonthlyIncome:        []DemographicStatDTO{},
		OrdinalPosition:      []DemographicStatDTO{},
		FatherEducation:      []DemographicStatDTO{},
		MotherEducation:      []DemographicStatDTO{},
		ParentsMaritalStatus: []DemographicStatDTO{},
		HighSchoolGWA:        []DemographicStatDTO{},
		Elementary:           []DemographicStatDTO{},
		JuniorHigh:           []DemographicStatDTO{},
		SeniorHigh:           []DemographicStatDTO{},
		NatureOfSchooling:    []DemographicStatDTO{},
		QuietStudyPlace:      []DemographicStatDTO{},
	}

	if total > 0 {
		// Demographic data
		rawAges, _ := s.repo.GetAgeStats(ctx)
		dashboard.AgeDistribution = s.mapToDTO(rawAges, total)

		rawCivilStatus, _ := s.repo.GetCivilStatusStats(ctx)
		dashboard.CivilStatus = s.mapToDTO(rawCivilStatus, total)

		rawReligions, _ := s.repo.GetReligionStats(ctx)
		dashboard.Religions = s.mapToDTO(rawReligions, total)

		rawCityAddress, _ := s.repo.GetCityAddressStats(ctx)
		dashboard.CityAddress = s.mapToDTO(rawCityAddress, total)

		// Economic/Social data
		rawMonthlyIncome, _ := s.repo.GetMonthlyIncomeStats(ctx)
		dashboard.MonthlyIncome = s.mapToDTO(rawMonthlyIncome, total)

		rawOrdinalPosition, _ := s.repo.GetOrdinalPositionStats(ctx)
		dashboard.OrdinalPosition = s.mapToDTO(rawOrdinalPosition, total)

		rawQuietPlace, _ := s.repo.GetQuietStudyPlaceStats(ctx)
		dashboard.QuietStudyPlace = s.mapToDTO(rawQuietPlace, total)

		// Family data
		rawFatherEd, _ := s.repo.GetFatherEducationStats(ctx)
		dashboard.FatherEducation = s.mapToDTO(rawFatherEd, total)

		rawMotherEd, _ := s.repo.GetMotherEducationStats(ctx)
		dashboard.MotherEducation = s.mapToDTO(rawMotherEd, total)

		rawParentsMarital, _ := s.repo.GetParentsMaritalStatusStats(ctx)
		dashboard.ParentsMaritalStatus = s.mapToDTO(rawParentsMarital, total)

		// Academic data
		rawHSGWA, _ := s.repo.GetHSGWAStats(ctx)
		dashboard.HighSchoolGWA = s.mapToDTO(rawHSGWA, total)

		rawElem, _ := s.repo.GetElementaryStats(ctx)
		dashboard.Elementary = s.mapToDTO(rawElem, total)

		rawJHS, _ := s.repo.GetJuniorHighStats(ctx)
		dashboard.JuniorHigh = s.mapToDTO(rawJHS, total)

		rawSHS, _ := s.repo.GetSeniorHighStats(ctx)
		dashboard.SeniorHigh = s.mapToDTO(rawSHS, total)

		rawNature, _ := s.repo.GetNatureOfSchoolingStats(ctx)
		dashboard.NatureOfSchooling = s.mapToDTO(rawNature, total)
	}

	return dashboard, nil
}

func (s *Service) GetAdminDashboard(
	ctx context.Context,
	timeRange string,
) (*AdminDashboardResponseDTO, error) {
	totalStudents, err := s.repo.GetTotalStudents(ctx)
	if err != nil {
		return nil, err
	}

	totalReports, err := s.repo.GetTotalReports(ctx)
	if err != nil {
		return nil, err
	}

	totalAppointments, err := s.repo.GetTotalAppointments(ctx)
	if err != nil {
		return nil, err
	}

	totalSlips, err := s.repo.GetTotalSlips(ctx)
	if err != nil {
		return nil, err
	}

	monthlyVisitors, err := s.repo.GetMonthlyVisitorStats(ctx, timeRange)
	if err != nil {
		return nil, err
	}

	// Count live sessions (session: prefix)
	liveSessions := 0
	if s.redis != nil && s.redis.Client != nil {
		keys, err := s.redis.Client.Keys(ctx, "session:*").Result()
		if err == nil {
			liveSessions = len(keys)
		}
	}

	return &AdminDashboardResponseDTO{
		TotalStudents:     totalStudents,
		TotalReports:      totalReports,
		TotalAppointments: totalAppointments,
		TotalSlips:        totalSlips,
		LiveSessions:      liveSessions,
		MonthlyVisitors:   monthlyVisitors,
	}, nil
}

func (s *Service) mapToDTO(
	rawStats []AggregatedStatModel,
	totalStudents int,
) []DemographicStatDTO {
	dtos := make([]DemographicStatDTO, 0)

	if totalStudents == 0 {
		return dtos
	}

	for _, stat := range rawStats {
		dto := DemographicStatDTO{
			Category:    stat.Category,
			MaleCount:   stat.MaleCount,
			FemaleCount: stat.FemaleCount,
			Total:       stat.Total,
			Rank:        stat.RankPos,
		}

		dto.TotalPct = s.calculatePercentage(stat.Total, totalStudents)
		dto.MalePct = s.calculatePercentage(stat.MaleCount, totalStudents)
		dto.FemalePct = s.calculatePercentage(stat.FemaleCount, totalStudents)

		dtos = append(dtos, dto)
	}

	return dtos
}

func (s *Service) calculatePercentage(count, total int) float64 {
	if total == 0 {
		return 0
	}
	percent := (float64(count) / float64(total)) * 100
	return math.Round(percent*100) / 100
}
