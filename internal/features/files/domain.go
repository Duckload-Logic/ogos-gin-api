package files

import "time"

// File represents a pure business entity for a file asset.
type File struct {
	ID        string
	FileName  string
	FileURL   string
	FileType  string
	FileSize  int64
	MimeType  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// OCRResult represents the pure result of an OCR process.
type OCRResult struct {
	FileID         string
	RawText        string
	StructuredData string
	EngineV        string
	Confidence     float64
	CreatedAt      time.Time
}
