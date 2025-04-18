package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringList []string

// Implement the sql.Scanner interface
func (s *StringList) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal StringList value '%+v'", value)
	}

	result := []string{}
	err := json.Unmarshal(bytes, &result)
	*s = StringList(result)
	return err
}

// Implement the driver.Valuer interface
func (s StringList) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}
	return json.Marshal(s)
}
