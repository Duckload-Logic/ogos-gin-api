package structs

import (
	"database/sql"
	"encoding/json"
)

type NullableString sql.NullString
type NullableInt64 sql.NullInt64

func (ns *NullableString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		ns.String = ""
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	ns.String = s
	ns.Valid = true
	return nil
}

func (ns NullableString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func FromSqlNull(ns sql.NullString) NullableString {
	return NullableString{String: ns.String, Valid: ns.Valid}
}

func ToSqlNull(ns NullableString) sql.NullString {
	return sql.NullString{String: ns.String, Valid: ns.Valid}
}

func (ni *NullableInt64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ni.Valid = false
		ni.Int64 = 0
		return nil
	}

	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	ni.Int64 = i
	ni.Valid = true
	return nil
}

func (ni NullableInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(ni.Int64)
}

func FromSqlNullInt64(ni sql.NullInt64) NullableInt64 {
	return NullableInt64{Int64: ni.Int64, Valid: ni.Valid}
}

func ToSqlNullInt64(ni NullableInt64) sql.NullInt64 {
	return sql.NullInt64{Int64: ni.Int64, Valid: ni.Valid}
}
