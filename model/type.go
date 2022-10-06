package model

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
)

type NullTime struct {
	sql.NullTime
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
		return fmt.Errorf("null: couldn't unmarshal JSON Time: %w", err)
	}

	t.Valid = true
	return nil
}
