package main

import (
	"database/sql"
	"encoding/json"
	"reflect"
)

// User : maps to row from users table
type User struct {
	ID      int64
	Name    string
	Role    string
	Email   string
	Phone   string
	Created string
	Updated string
}

// Shift : maps to row from shifts table
type Shift struct {
	ID       int64     `json:"id"`
	Manager  int64     `json:"manager"`
	Employee NullInt64 `json:"employee"`
	Break    float64   `json:"break"`
	Start    string    `json:"startTime"`
	End      string    `json:"endTime"`
	Created  string    `json:"createdAt"`
	Updated  string    `json:"updatedAt"`
}

// Roster : maps to row from join on shifts and users table
type Roster struct {
	ID       int64
	Manager  int64
	Employee NullInt64
	Break    NullFloat64
	Start    string
	End      string
	Created  string
	Updated  string
	Name     string
	Email    NullString
	Phone    NullString
}

// Hours : holds summary data from getHours route
type Hours struct {
	Shifts     []Shift
	TotalHours int
}

// NullInt64 : allow for null value in employee field for shifts
type NullInt64 sql.NullInt64

// NullFloat64 : allow for null value in break field for shifts
type NullFloat64 sql.NullFloat64

// NullString : allow for null value in phone and email field for users
type NullString sql.NullString

// Scan : check for null values and set Valid bool - Int
func (ni *NullInt64) Scan(value interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ni = NullInt64{i.Int64, false}
	} else {
		*ni = NullInt64{i.Int64, true}
	}
	return nil
}

// Scan : check for null value and set Valid bool - Float
func (nf *NullFloat64) Scan(value interface{}) error {
	var f sql.NullFloat64
	if err := f.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*nf = NullFloat64{f.Float64, false}
	} else {
		*nf = NullFloat64{f.Float64, true}
	}

	return nil
}

// Scan : check for null value and set Valid bool - String
func (ns *NullString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ns = NullString{s.String, false}
	} else {
		*ns = NullString{s.String, true}
	}

	return nil
}

// MarshalJSON : allow for EITHER null or populated value in employee field for shifts
func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// UnmarshalJSON for NullInt64
func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.Int64)
	ni.Valid = (err == nil)
	return err
}

// MarshalJSON for NullFloat64
func (nf *NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

// UnmarshalJSON for NullFloat64
func (nf *NullFloat64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nf.Float64)
	nf.Valid = (err == nil)
	return err
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON for NullString
func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = (err == nil)
	return err
}
