package analytics

// MapDemographicStatToDomain converts DB model to domain model.
func MapDemographicStatToDomain(db DemographicStatDB) DemographicStat {
	return DemographicStat(db)
}

// MapDemographicStatsToDomain maps a slice of DB models to domain models.
func MapDemographicStatsToDomain(db []DemographicStatDB) []DemographicStat {
	domain := make([]DemographicStat, len(db))
	for i := range db {
		domain[i] = MapDemographicStatToDomain(db[i])
	}
	return domain
}
