package utils

import (
	"database/sql"
	"encoding/json"
	"time"
)

// NullTime is an alias for sql.NullTime with custom JSON unmarshaling
type NullTime struct {
	sql.NullTime
}

// UnmarshalJSON unmarshals a JSON string into NullTime
func (nt *NullTime) UnmarshalJSON(b []byte) error {
	// Unquote the time string
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	// If the string is empty, set NullTime as null
	if s == "" {
		nt.Valid = false
		return nil
	}

	// Parse the time string
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}

	// Set the time and mark it as valid
	nt.Time = t
	nt.Valid = true
	return nil
}
