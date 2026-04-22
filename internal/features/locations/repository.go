package locations

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetDB() *sqlx.DB {
	return r.db
}

func (r *Repository) WithTransaction(
	ctx context.Context,
	fn func(datastore.DB) error,
) error {
	return datastore.RunInTransaction(ctx, r.db, fn)
}

func (r *Repository) GetRegions(ctx context.Context) ([]Region, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM regions ORDER BY name",
		datastore.GetColumns(RegionDB{}),
	)
	var dbModels []RegionDB
	err := r.db.SelectContext(ctx, &dbModels, query)
	if err != nil {
		return nil, err
	}

	regions := make([]Region, len(dbModels))
	for i, m := range dbModels {
		regions[i] = m.ToDomain()
	}
	return regions, nil
}

func (r *Repository) GetProvincesByRegion(
	ctx context.Context,
	regionCode string,
) ([]Province, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM provinces WHERE region_code = ? ORDER BY name",
		datastore.GetColumns(ProvinceDB{}),
	)
	var dbModels []ProvinceDB
	err := r.db.SelectContext(ctx, &dbModels, query, regionCode)
	if err != nil {
		return nil, err
	}

	provinces := make([]Province, len(dbModels))
	for i, m := range dbModels {
		provinces[i] = m.ToDomain()
	}
	return provinces, nil
}

func (r *Repository) GetCitiesByProvince(
	ctx context.Context,
	provinceCode string,
) ([]City, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM cities WHERE province_code = ? ORDER BY name",
		datastore.GetColumns(CityDB{}),
	)
	var dbModels []CityDB
	err := r.db.SelectContext(ctx, &dbModels, query, provinceCode)
	if err != nil {
		return nil, err
	}

	cities := make([]City, len(dbModels))
	for i, m := range dbModels {
		cities[i] = m.ToDomain()
	}
	return cities, nil
}

func (r *Repository) GetCitiesByRegion(
	ctx context.Context,
	regionCode string,
) ([]City, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM cities WHERE region_code = ? ORDER BY name",
		datastore.GetColumns(CityDB{}),
	)
	var dbModels []CityDB
	err := r.db.SelectContext(ctx, &dbModels, query, regionCode)
	if err != nil {
		return nil, err
	}

	cities := make([]City, len(dbModels))
	for i, m := range dbModels {
		cities[i] = m.ToDomain()
	}
	return cities, nil
}

func (r *Repository) GetBarangaysByCity(
	ctx context.Context,
	cityCode string,
) ([]Barangay, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM barangays WHERE city_code = ? ORDER BY name",
		datastore.GetColumns(BarangayDB{}),
	)
	var dbModels []BarangayDB
	err := r.db.SelectContext(ctx, &dbModels, query, cityCode)
	if err != nil {
		return nil, err
	}

	barangays := make([]Barangay, len(dbModels))
	for i, m := range dbModels {
		barangays[i] = m.ToDomain()
	}
	return barangays, nil
}

func (r *Repository) GetAddressByID(
	ctx context.Context,
	addressID int,
) (*Address, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM addresses WHERE id = ?",
		datastore.GetColumns(AddressDB{}),
	)
	var dbModel AddressDB
	err := r.db.GetContext(ctx, &dbModel, query, addressID)
	if err != nil {
		return nil, err
	}
	domainModel := dbModel.ToDomain()
	return &domainModel, nil
}

func (r *Repository) GetCityByCode(
	ctx context.Context,
	cityCode string,
) (*City, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM cities WHERE code = ?",
		datastore.GetColumns(CityDB{}),
	)
	var dbModel CityDB
	err := r.db.GetContext(ctx, &dbModel, query, cityCode)
	if err != nil {
		return nil, err
	}
	domainModel := dbModel.ToDomain()
	return &domainModel, nil
}

func (r *Repository) GetRegionByCode(
	ctx context.Context,
	regionCode string,
) (*Region, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM regions WHERE code = ?",
		datastore.GetColumns(RegionDB{}),
	)
	var dbModel RegionDB
	err := r.db.GetContext(ctx, &dbModel, query, regionCode)
	if err != nil {
		return nil, err
	}
	domainModel := dbModel.ToDomain()
	return &domainModel, nil
}

func (r *Repository) GetBarangayByCode(
	ctx context.Context,
	barangayCode string,
) (*Barangay, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM barangays WHERE code = ?",
		datastore.GetColumns(BarangayDB{}),
	)
	var dbModel BarangayDB
	err := r.db.GetContext(ctx, &dbModel, query, barangayCode)
	if err != nil {
		return nil, err
	}
	domainModel := dbModel.ToDomain()
	return &domainModel, nil
}

func (r *Repository) GetProvinceByCode(
	ctx context.Context,
	provinceCode string,
) (*Province, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM provinces WHERE code = ?",
		datastore.GetColumns(ProvinceDB{}),
	)
	var dbModel ProvinceDB
	err := r.db.GetContext(ctx, &dbModel, query, provinceCode)
	if err != nil {
		return nil, err
	}
	domainModel := dbModel.ToDomain()
	return &domainModel, nil
}

func (r *Repository) UpsertAddress(
	ctx context.Context,
	tx datastore.DB,
	addr *Address,
) (int, error) {
	dbModel := addr.ToPersistence()
	exclude := []string{"created_at", "updated_at"}
	cols, vals := datastore.GetInsertStatement(AddressDB{}, exclude)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		AddressDB{},
		exclude,
	)

	query := fmt.Sprintf(`
		INSERT INTO addresses (%s)
		VALUES (%s)
		ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)
	result, err := tx.NamedExecContext(ctx, query, &dbModel)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert address: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf(
			"failed to get last insert ID for address: %w",
			err,
		)
	}

	return int(lastID), nil
}
