package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type IDList []int

// Implement the sql.Scanner interface
func (s *IDList) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

// Implement the driver.Valuer interface
func (s IDList) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}
