package locations

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// ServiceInterface defines the business logic for location-based operations.
type ServiceInterface interface {
	GetRegions(ctx context.Context) ([]Region, error)
	GetProvincesByRegion(ctx context.Context, regionCode string) ([]Province, error)
	GetCitiesByProvince(ctx context.Context, provinceCode string) ([]City, error)
	GetCitiesByRegion(ctx context.Context, regionCode string) ([]City, error)
	GetBarangaysByCity(ctx context.Context, cityCode string) ([]Barangay, error)
	GetAddressByID(ctx context.Context, addressID int) (AddressDTO, error)
	GetCityByCode(ctx context.Context, cityCode string) (*City, error)
	GetRegionByCode(ctx context.Context, regionCode string) (*Region, error)
	GetBarangayByCode(ctx context.Context, barangayCode string) (*Barangay, error)
	GetProvinceByCode(ctx context.Context, provinceCode string) (*Province, error)
	SaveAddress(ctx context.Context, tx *sqlx.Tx, addr *Address) (int, error)
}

// RepositoryInterface defines the data access layer for location-based operations.
type RepositoryInterface interface {
	GetRegions(ctx context.Context) ([]Region, error)
	GetProvincesByRegion(ctx context.Context, regionCode string) ([]Province, error)
	GetCitiesByProvince(ctx context.Context, provinceCode string) ([]City, error)
	GetCitiesByRegion(ctx context.Context, regionCode string) ([]City, error)
	GetBarangaysByCity(ctx context.Context, cityCode string) ([]Barangay, error)
	GetAddressByID(ctx context.Context, addressID int) (*Address, error)
	GetCityByCode(ctx context.Context, cityCode string) (*City, error)
	GetRegionByCode(ctx context.Context, regionCode string) (*Region, error)
	GetBarangayByCode(ctx context.Context, barangayCode string) (*Barangay, error)
	GetProvinceByCode(ctx context.Context, provinceCode string) (*Province, error)
	UpsertAddress(ctx context.Context, tx *sqlx.Tx, addr *Address) (int, error)
}
