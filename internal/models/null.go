package models

import (
	"database/sql"
	"encoding/json"
)

func NewNullFloat64(f float64) NullFloat64 {
	var nf NullFloat64
	nf.Valid = true
	nf.Float64 = f
	return nf
}

type NullFloat64 struct {
	sql.NullFloat64
}

func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if nf.Valid {
		return json.Marshal(nf.Float64)
	}
	return json.Marshal(nil)
}

func (nf *NullFloat64) UnmarshalJSON(data []byte) error {
	var f *float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	if f != nil {
		nf.Valid = true
		nf.Float64 = *f
	} else {
		nf.Valid = false
	}
	return nil
}

func NewNullInt64(i int64) NullInt64 {
	var ni NullInt64
	ni.Valid = true
	ni.Int64 = i
	return ni
}

type NullInt64 struct {
	sql.NullInt64
}

func (nf NullInt64) MarshalJSON() ([]byte, error) {
	if nf.Valid {
		return json.Marshal(nf.Int64)
	}
	return json.Marshal(nil)
}

func (nf *NullInt64) UnmarshalJSON(data []byte) error {
	var i *int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		nf.Valid = true
		nf.Int64 = *i
	} else {
		nf.Valid = false
	}
	return nil
}
