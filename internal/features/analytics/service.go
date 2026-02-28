package analytics

import (
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetDashboard(ctx context.Context) (*DashboardResponseDTO, error) {
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
		rawAges, _ := s.repo.GetAgeStats(ctx)
		dashboard.AgeDistribution = s.mapToDTO(rawAges, total)

		rawCivilStatus, _ := s.repo.GetCivilStatusStats(ctx)
		dashboard.CivilStatus = s.mapToDTO(rawCivilStatus, total)

		rawReligions, _ := s.repo.GetReligionStats(ctx)
		dashboard.Religions = s.mapToDTO(rawReligions, total)

		rawCityAddress, _ := s.repo.GetCityAddressStats(ctx)
		dashboard.CityAddress = s.mapToDTO(rawCityAddress, total)

		rawMonthlyIncome, _ := s.repo.GetMonthlyIncomeStats(ctx)
		dashboard.MonthlyIncome = s.mapToDTO(rawMonthlyIncome, total)

		rawOrdinalPosition, _ := s.repo.GetOrdinalPositionStats(ctx)
		dashboard.OrdinalPosition = s.mapToDTO(rawOrdinalPosition, total)

		rawFatherEd, _ := s.repo.GetFatherEducationStats(ctx)
		dashboard.FatherEducation = s.mapToDTO(rawFatherEd, total)

		rawMotherEd, _ := s.repo.GetMotherEducationStats(ctx)
		dashboard.MotherEducation = s.mapToDTO(rawMotherEd, total)

		rawParentsMarital, _ := s.repo.GetParentsMaritalStatusStats(ctx)
		dashboard.ParentsMaritalStatus = s.mapToDTO(rawParentsMarital, total)

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

		rawQuietPlace, _ := s.repo.GetQuietStudyPlaceStats(ctx)
		dashboard.QuietStudyPlace = s.mapToDTO(rawQuietPlace, total)
	}

	return dashboard, nil
}

func (s *Service) mapToDTO(rawStats []AggregatedStatModel, totalStudents int) []DemographicStatDTO {
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

		dto.TotalPct = (float64(stat.Total) / float64(totalStudents)) * 100
		dto.MalePct = (float64(stat.MaleCount) / float64(totalStudents)) * 100
		dto.FemalePct = (float64(stat.FemaleCount) / float64(totalStudents)) * 100

		dtos = append(dtos, dto)
	}

	return dtos
}