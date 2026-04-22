package locations

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// Region mappers
func (m RegionDB) ToDomain() Region {
	return Region(m)
}

func (d Region) ToPersistence() RegionDB {
	return RegionDB(d)
}

// Province mappers
func (m ProvinceDB) ToDomain() Province {
	return Province(m)
}

func (d Province) ToPersistence() ProvinceDB {
	return ProvinceDB(d)
}

// City mappers
func (m CityDB) ToDomain() City {
	return City{
		ID:           m.ID,
		Code:         m.Code,
		Name:         m.Name,
		ProvinceCode: structs.FromSqlNull(m.ProvinceCode),
		Type:         structs.FromSqlNull(m.Type),
		ZipCode:      structs.FromSqlNull(m.ZipCode),
		District:     structs.FromSqlNull(m.District),
		RegionCode:   structs.FromSqlNull(m.RegionCode),
	}
}

func (d City) ToPersistence() CityDB {
	return CityDB{
		ID:           d.ID,
		Code:         d.Code,
		Name:         d.Name,
		ProvinceCode: structs.ToSqlNull(d.ProvinceCode),
		Type:         structs.ToSqlNull(d.Type),
		ZipCode:      structs.ToSqlNull(d.ZipCode),
		District:     structs.ToSqlNull(d.District),
		RegionCode:   structs.ToSqlNull(d.RegionCode),
	}
}

// Barangay mappers
func (m BarangayDB) ToDomain() Barangay {
	return Barangay(m)
}

func (d Barangay) ToPersistence() BarangayDB {
	return BarangayDB(d)
}

// Address mappers
func (m AddressDB) ToDomain() Address {
	return Address{
		ID:           m.ID,
		RegionCode:   m.RegionCode,
		ProvinceCode: structs.FromSqlNull(m.ProvinceCode),
		CityCode:     m.CityCode,
		BarangayCode: m.BarangayCode,
		StreetDetail: structs.FromSqlNull(m.StreetDetail),
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (d Address) ToPersistence() AddressDB {
	return AddressDB{
		ID:           d.ID,
		RegionCode:   d.RegionCode,
		ProvinceCode: structs.ToSqlNull(d.ProvinceCode),
		CityCode:     d.CityCode,
		BarangayCode: d.BarangayCode,
		StreetDetail: structs.ToSqlNull(d.StreetDetail),
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}
