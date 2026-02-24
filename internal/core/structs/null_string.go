package structs

import (
	"database/sql"
	"encoding/json"
)

type NullableString sql.NullString

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
