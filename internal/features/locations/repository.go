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

func (r *Repository) GetRegions(ctx context.Context) ([]Region, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM regions ORDER BY name",
		datastore.GetColumns(Region{}),
	)
	var regions []Region
	err := r.db.SelectContext(ctx, &regions, query)
	return regions, err
}

func (r *Repository) GetProvincesByRegion(
	ctx context.Context,
	regionCode string,
) ([]Province, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM provinces WHERE region_code = ? ORDER BY name",
		datastore.GetColumns(Province{}),
	)
	var provinces []Province
	err := r.db.SelectContext(ctx, &provinces, query, regionCode)
	return provinces, err
}

func (r *Repository) GetCitiesByProvince(
	ctx context.Context,
	provinceCode string,
) ([]City, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM cities WHERE province_code = ? ORDER BY name",
		datastore.GetColumns(City{}),
	)
	var cities []City
	err := r.db.SelectContext(ctx, &cities, query, provinceCode)
	return cities, err
}

func (r *Repository) GetCitiesByRegion(
	ctx context.Context,
	regionCode string,
) ([]City, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM cities WHERE region_code = ? ORDER BY name",
		datastore.GetColumns(City{}),
	)
	var cities []City
	err := r.db.SelectContext(ctx, &cities, query, regionCode)
	return cities, err
}

func (r *Repository) GetBarangaysByCity(
	ctx context.Context,
	cityCode string,
) ([]Barangay, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM barangays WHERE city_code = ? ORDER BY name",
		datastore.GetColumns(Barangay{}),
	)
	var barangays []Barangay
	err := r.db.SelectContext(ctx, &barangays, query, cityCode)
	return barangays, err
}

func (r *Repository) GetAddressByID(
	ctx context.Context,
	addressID int,
) (*Address, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM addresses WHERE id = ?",
		datastore.GetColumns(Address{}),
	)
	var address Address
	err := r.db.GetContext(ctx, &address, query, addressID)
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *Repository) GetCityByCode(
	ctx context.Context,
	cityCode string,
) (*City, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM cities WHERE code = ?",
		datastore.GetColumns(City{}),
	)
	var city City
	err := r.db.GetContext(ctx, &city, query, cityCode)
	if err != nil {
		return nil, err
	}
	return &city, nil
}

func (r *Repository) GetRegionByCode(
	ctx context.Context,
	regionCode string,
) (*Region, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM regions WHERE code = ?",
		datastore.GetColumns(Region{}),
	)
	var region Region
	err := r.db.GetContext(ctx, &region, query, regionCode)
	if err != nil {
		return nil, err
	}
	return &region, nil
}

func (r *Repository) GetBarangayByCode(
	ctx context.Context,
	barangayCode string,
) (*Barangay, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM barangays WHERE code = ?",
		datastore.GetColumns(Barangay{}),
	)
	var barangay Barangay
	err := r.db.GetContext(ctx, &barangay, query, barangayCode)
	if err != nil {
		return nil, err
	}
	return &barangay, nil
}

func (r *Repository) GetProvinceByCode(
	ctx context.Context,
	provinceCode string,
) (*Province, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM provinces WHERE code = ?",
		datastore.GetColumns(Province{}),
	)
	var province Province
	err := r.db.GetContext(ctx, &province, query, provinceCode)
	if err != nil {
		return nil, err
	}
	return &province, nil
}

func (r *Repository) UpsertAddress(
	ctx context.Context,
	tx *sqlx.Tx,
	addr *Address,
) (int, error) {
	if tx != nil {
		return r.upsertAddressTx(ctx, tx, addr)
	}

	var id int
	err := datastore.RunInTransaction(ctx, r.db, func(txn *sqlx.Tx) error {
		var err error
		id, err = r.upsertAddressTx(ctx, txn, addr)
		return err
	})
	return id, err
}

func (r *Repository) upsertAddressTx(
	ctx context.Context,
	tx *sqlx.Tx,
	addr *Address,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(
		Address{},
		[]string{"created_at", "updated_at"},
	)
	updateCols := datastore.GetOnDuplicateKeyUpdateStatement(
		Address{},
		[]string{"created_at", "updated_at"},
	)

	if addr.ProvinceCode != nil && *addr.ProvinceCode == "" {
		addr.ProvinceCode = nil
	}

	query := fmt.Sprintf(`
		INSERT INTO addresses (%s)
		VALUES (%s)
		ON DUPLICATE KEY UPDATE %s
	`, cols, vals, updateCols)
	result, err := tx.NamedExecContext(ctx, query, addr)
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
