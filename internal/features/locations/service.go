package locations

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetRegions(ctx context.Context) ([]Region, error) {
	return s.repo.GetRegions(ctx)
}

func (s *Service) GetCitiesByRegion(ctx context.Context, regionID int) ([]City, error) {
	return s.repo.GetCitiesByRegion(ctx, regionID)
}

func (s *Service) GetBarangaysByCity(ctx context.Context, cityID int) ([]Barangay, error) {
	return s.repo.GetBarangaysByCity(ctx, cityID)
}

func (s *Service) GetAddressByID(ctx context.Context, addressID int) (AddressDTO, error) {
	addr, err := s.repo.GetAddressByID(ctx, addressID)
	if err != nil {
		return AddressDTO{}, err
	}

	barangay, err := s.repo.GetBarangayByID(ctx, addr.BarangayID)
	if err != nil {
		return AddressDTO{}, err
	}
	city, err := s.repo.GetCityByID(ctx, addr.CityID)
	if err != nil {
		return AddressDTO{}, err
	}
	region, err := s.repo.GetRegionByID(ctx, addr.RegionID)
	if err != nil {
		return AddressDTO{}, err
	}

	addrDTO := AddressDTO{
		ID:           addr.ID,
		StreetDetail: *addr.StreetDetail,
		Barangay:     *barangay,
		City:         *city,
		Region:       *region,
	}

	return addrDTO, nil
}

func (s *Service) GetCityByID(ctx context.Context, cityID int) (*City, error) {
	return s.repo.GetCityByID(ctx, cityID)
}

func (s *Service) GetRegionByID(ctx context.Context, regionID int) (*Region, error) {
	return s.repo.GetRegionByID(ctx, regionID)
}

func (s *Service) GetBarangayByID(ctx context.Context, barangayID int) (*Barangay, error) {
	return s.repo.GetBarangayByID(ctx, barangayID)
}

func (s *Service) SaveAddress(ctx context.Context, tx *sqlx.Tx, addr *Address) (int, error) {
	return s.repo.UpsertAddress(ctx, tx, addr)
}
