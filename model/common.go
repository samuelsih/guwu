package model

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/lib/pq"
)

type NullTime struct {
	pq.NullTime
}

func (t NullTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}

	return t.Time.MarshalJSON()
}

func (t *NullTime) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		t.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &t.Time); err != nil {
		return fmt.Errorf("couldn't unmarshal JSON Time: %w", err)
	}

	t.Valid = true
	return nil
}

type NullString struct {
	sql.NullString
}

func (s NullString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(s.String)
}

func (s *NullString) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		s.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &s.String); err != nil {
		return fmt.Errorf("couldn't unmarshal JSON String: %w", err)
	}

	s.Valid = true
	return nil
}
