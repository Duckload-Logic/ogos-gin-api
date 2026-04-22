package files

import "time"

// FileDB represents the database model for a file.
type FileDB struct {
	ID        string    `db:"id"`
	FileName  string    `db:"file_name"`
	FileURL   string    `db:"file_url"`
	FileType  string    `db:"file_type"`
	FileSize  int64     `db:"file_size"`
	MimeType  string    `db:"mime_type"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}

// OCRResultDB represents the database model for OCR results.
type OCRResultDB struct {
	FileID         string    `db:"file_id"`
	RawText        string    `db:"raw_text"`
	StructuredData string    `db:"structured_data"`
	EngineV        string    `db:"engine_v"`
	Confidence     float64   `db:"confidence"`
	CreatedAt      time.Time `db:"created_at"`
}
