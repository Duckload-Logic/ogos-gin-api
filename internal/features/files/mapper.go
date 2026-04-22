package files

// MapFileToDomain converts DB model to domain model.
func MapFileToDomain(db FileDB) File {
	return File(db)
}

// MapFileToDB converts domain model to DB model.
func MapFileToDB(d File) FileDB {
	return FileDB(d)
}

// MapOCRResultToDomain converts DB model to domain model.
func MapOCRResultToDomain(db OCRResultDB) OCRResult {
	return OCRResult(db)
}

// MapOCRResultToDB converts domain model to DB model.
func MapOCRResultToDB(d OCRResult) OCRResultDB {
	return OCRResultDB(d)
}
