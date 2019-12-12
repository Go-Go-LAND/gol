package gol

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullBool struct {
	sql.NullBool
}

func (rec *NullBool) Get() *bool {
	if !rec.NullBool.Valid {
		return nil
	}

	val := rec.NullBool.Bool

	return &val
}

func (rec *NullBool) Set(value bool) {
	rec.NullBool.Bool = value
	rec.NullBool.Valid = true
}

func (rec *NullBool) Delete()  {
	var val bool
	rec.NullBool.Bool = val
	rec.NullBool.Valid = false
}

func (rec NullBool) MarshalJSON() ([]byte, error) {
	if !rec.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(rec.Bool)
}

func (rec *NullBool) UnmarshalJSON(value []byte) error {
	var val *bool

	err := json.Unmarshal(value, &val)
	if err != nil {
		return err
	}

	if val == nil {
		rec.Valid = false
		return nil
	}

	rec.Bool = *val
	rec.Valid = true

	return nil
}

type NullInt32 struct {
	sql.NullInt32
}

func (rec *NullInt32) Get() *int {
	if !rec.NullInt32.Valid {
		return nil
	}

	val := int(rec.NullInt32.Int32)

	return &val
}

func (rec *NullInt32) Set(value int) {
	rec.NullInt32.Int32 = int32(value)
	rec.NullInt32.Valid = true
}

func (rec *NullInt32) Delete()  {
	var val int32
	rec.NullInt32.Int32 = val
	rec.NullInt32.Valid = false
}

func (rec NullInt32) MarshalJSON() ([]byte, error) {
	if !rec.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(rec.Int32)
}

func (rec *NullInt32) UnmarshalJSON(value []byte) error {
	var val *int32

	err := json.Unmarshal(value, &val)
	if err != nil {
		return err
	}

	if val == nil {
		rec.Valid = false
		return nil
	}

	rec.Int32 = *val
	rec.Valid = true

	return nil
}

type NullInt64 struct {
	sql.NullInt64
}

func (rec *NullInt64) Get() *int {
	if !rec.NullInt64.Valid {
		return nil
	}

	val := int(rec.NullInt64.Int64)

	return &val
}

func (rec *NullInt64) Set(value int) {
	rec.NullInt64.Int64 = int64(value)
	rec.NullInt64.Valid = true
}

func (rec *NullInt64) Delete()  {
	var val int64
	rec.NullInt64.Int64 = val
	rec.NullInt64.Valid = false
}

func (rec NullInt64) MarshalJSON() ([]byte, error) {
	if !rec.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(rec.Int64)
}

func (rec *NullInt64) UnmarshalJSON(value []byte) error {
	var val *int64

	err := json.Unmarshal(value, &val)
	if err != nil {
		return err
	}

	if val == nil {
		rec.Valid = false
		return nil
	}

	rec.Int64 = *val
	rec.Valid = true

	return nil
}

type NullFloat64 struct {
	sql.NullFloat64
}

func (rec *NullFloat64) Get() *float64 {
	if !rec.NullFloat64.Valid {
		return nil
	}

	val := rec.NullFloat64.Float64

	return &val
}

func (rec *NullFloat64) Set(value float64) {
	rec.NullFloat64.Float64 = value
	rec.NullFloat64.Valid = true
}

func (rec *NullFloat64) Delete()  {
	var val float64
	rec.NullFloat64.Float64 = val
	rec.NullFloat64.Valid = false
}

func (rec NullFloat64) MarshalJSON() ([]byte, error) {
	if !rec.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(rec.Float64)
}

func (rec *NullFloat64) UnmarshalJSON(value []byte) error {
	var val *float64

	err := json.Unmarshal(value, &val)
	if err != nil {
		return err
	}

	if val == nil {
		rec.Valid = false
		return nil
	}

	rec.Float64 = *val
	rec.Valid = true

	return nil
}

type NullString struct {
	sql.NullString
}

func (rec *NullString) Get() *string {
	if !rec.NullString.Valid {
		return nil
	}

	val := rec.NullString.String

	return &val
}

func (rec *NullString) Set(value string) {
	rec.NullString.String = value
	rec.NullString.Valid = true
}

func (rec *NullString) Delete()  {
	var val string
	rec.NullString.String = val
	rec.NullString.Valid = false
}

func (rec NullString) MarshalJSON() ([]byte, error) {
	if !rec.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(rec.String)
}

func (rec *NullString) UnmarshalJSON(value []byte) error {
	var val *string

	err := json.Unmarshal(value, &val)
	if err != nil {
		return err
	}

	if val == nil {
		rec.Valid = false
		return nil
	}

	rec.String = *val
	rec.Valid = true

	return nil
}

type NullTime struct {
	sql.NullTime
}

func (rec *NullTime) Get() *time.Time {
	if !rec.NullTime.Valid {
		return nil
	}

	val := rec.NullTime.Time

	return &val
}

func (rec *NullTime) Set(value time.Time) {
	rec.NullTime.Time = value
	rec.NullTime.Valid = true
}

func (rec *NullTime) Delete()  {
	var val time.Time
	rec.NullTime.Time = val
	rec.NullTime.Valid = false
}

func (rec NullTime) MarshalJSON() ([]byte, error) {
	if !rec.Valid {
		return json.Marshal(nil)
	}

	return json.Marshal(rec.Time)
}

func (rec *NullTime) UnmarshalJSON(value []byte) error {
	var val *time.Time

	err := json.Unmarshal(value, &val)
	if err != nil {
		return err
	}

	if val == nil {
		rec.Valid = false
		return nil
	}

	rec.Time = *val
	rec.Valid = true

	return nil
}
