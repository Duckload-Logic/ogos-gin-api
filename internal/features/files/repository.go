package files

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) RepositoryInterface {
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

func (r *Repository) GetFileByID(
	ctx context.Context,
	id string,
) (*File, error) {
	var file FileDB
	query := `SELECT * FROM files WHERE id = ?`
	err := r.db.GetContext(ctx, &file, query, id)
	if err != nil {
		return nil, err
	}

	domainFile := MapFileToDomain(file)
	return &domainFile, nil
}

func (r *Repository) GetFilesByIDs(
	ctx context.Context,
	ids []string,
) ([]File, error) {
	var files []FileDB
	query, args, err := sqlx.In(`SELECT * FROM files WHERE id IN (?)`, ids)
	if err != nil {
		return nil, err
	}

	query = r.db.Rebind(query)
	err = r.db.SelectContext(ctx, &files, query, args...)
	if err != nil {
		return nil, err
	}

	domainFiles := make([]File, len(files))
	for i, f := range files {
		domainFiles[i] = MapFileToDomain(f)
	}
	return domainFiles, nil
}

func (r *Repository) Create(
	ctx context.Context,
	tx datastore.DB,
	file File,
) (string, error) {
	dbModel := MapFileToDB(file)
	cols, vals := datastore.GetInsertStatement(&dbModel, []string{})
	query := fmt.Sprintf(`INSERT INTO files (%s) VALUES (%s)`, cols, vals)
	_, err := tx.NamedExecContext(ctx, query, dbModel)
	if err != nil {
		return "", err
	}

	return file.ID, nil
}

func (r *Repository) CreateBulk(
	ctx context.Context,
	tx datastore.DB,
	files []File,
) ([]string, error) {
	var ids []string
	for _, file := range files {
		id, err := r.Create(ctx, tx, file)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *Repository) Delete(
	ctx context.Context,
	tx datastore.DB,
	id string,
) error {
	query := `DELETE FROM files WHERE id = ?`
	_, err := tx.ExecContext(ctx, query, id)
	return err
}
