package files

import (
	"context"
	"fmt"

	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

func (r *Repository) SaveOCRResult(
	ctx context.Context,
	tx datastore.DB,
	result OCRResult,
) error {
	dbModel := MapOCRResultToDB(result)
	query := `
		INSERT INTO file_ocr_results (
			file_id, raw_text, structured_data, engine_v, confidence
		) VALUES (
			:file_id, :raw_text, :structured_data, :engine_v, :confidence
		)
		ON DUPLICATE KEY UPDATE
			raw_text = VALUES(raw_text),
			structured_data = VALUES(structured_data),
			engine_v = VALUES(engine_v),
			confidence = VALUES(confidence)
	`
	_, err := tx.NamedExecContext(ctx, query, dbModel)
	if err != nil {
		return fmt.Errorf("[OCRRepository] {SaveOCRResult}: %w", err)
	}
	return nil
}

func (r *Repository) GetOCRResultByFileID(
	ctx context.Context,
	fileID string,
) (*OCRResult, error) {
	var result OCRResultDB
	query := "SELECT * FROM file_ocr_results WHERE file_id = ?"
	err := r.db.GetContext(ctx, &result, query, fileID)
	if err != nil {
		return nil, fmt.Errorf(
			"[OCRRepository] {GetOCRResultByFileID}: %w",
			err,
		)
	}
	domainResult := MapOCRResultToDomain(result)
	return &domainResult, nil
}
