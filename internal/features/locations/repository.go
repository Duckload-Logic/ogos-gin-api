package locations

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetRegions(ctx context.Context) ([]Region, error) {
	query := fmt.Sprintf("SELECT %s FROM regions", database.GetColumns(Region{}))
	var regions []Region
	err := r.db.SelectContext(ctx, &regions, query)
	return regions, err
}

func (r *Repository) GetCitiesByRegion(ctx context.Context, regionID int) ([]City, error) {
	query := fmt.Sprintf("SELECT %s FROM cities WHERE region_id = ?", database.GetColumns(City{}))
	var cities []City
	err := r.db.SelectContext(ctx, &cities, query, regionID)
	return cities, err
}

func (r *Repository) GetBarangaysByCity(ctx context.Context, cityID int) ([]Barangay, error) {
	query := fmt.Sprintf("SELECT %s FROM barangays WHERE city_id = ?", database.GetColumns(Barangay{}))
	var barangays []Barangay
	err := r.db.SelectContext(ctx, &barangays, query, cityID)
	return barangays, err
}

func (r *Repository) GetAddressByID(ctx context.Context, addressID int) (*Address, error) {
	query := fmt.Sprintf("SELECT %s FROM addresses WHERE id = ?", database.GetColumns(Address{}))
	var address Address
	err := r.db.GetContext(ctx, &address, query, addressID)
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *Repository) GetCityByID(ctx context.Context, cityID int) (*City, error) {
	query := fmt.Sprintf("SELECT %s FROM cities WHERE id = ?", database.GetColumns(City{}))
	var city City
	err := r.db.GetContext(ctx, &city, query, cityID)
	if err != nil {
		return nil, err
	}
	return &city, nil
}

func (r *Repository) GetRegionByID(ctx context.Context, regionID int) (*Region, error) {
	query := fmt.Sprintf("SELECT %s FROM regions WHERE id = ?", database.GetColumns(Region{}))
	var region Region
	err := r.db.GetContext(ctx, &region, query, regionID)
	if err != nil {
		return nil, err
	}
	return &region, nil
}

func (r *Repository) GetBarangayByID(ctx context.Context, barangayID int) (*Barangay, error) {
	query := fmt.Sprintf("SELECT %s FROM barangays WHERE id = ?", database.GetColumns(Barangay{}))
	var barangay Barangay
	err := r.db.GetContext(ctx, &barangay, query, barangayID)
	if err != nil {
		return nil, err
	}
	return &barangay, nil
}

func (r *Repository) UpsertAddress(ctx context.Context, tx *sqlx.Tx, addr *Address) (int, error) {
	if tx != nil {
		return r.upsertAddressTx(ctx, tx, addr)
	}

	var id int
	err := database.RunInTransaction(ctx, r.db, func(txn *sqlx.Tx) error {
		var err error
		id, err = r.upsertAddressTx(ctx, txn, addr)
		return err
	})
	return id, err
}

func (r *Repository) upsertAddressTx(ctx context.Context, tx *sqlx.Tx, addr *Address) (int, error) {
	cols, vals := database.GetInsertStatement(Address{}, []string{"created_at", "updated_at"})
	updateCols := database.GetOnDuplicateKeyUpdateStatement(Address{}, []string{"created_at", "updated_at"})

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
		return 0, fmt.Errorf("failed to get last insert ID for address: %w", err)
	}

	return int(lastID), nil
}
