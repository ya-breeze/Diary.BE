// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Diary - OpenAPI 3.0
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.0.1
 * Contact: ilya.korolev@outlook.com
 */

package goserver

type Entity struct {
	Id string `json:"id"`
}

type EntityInterface interface {
	GetId() string
}

func (c *Entity) GetId() string {
	return c.Id
}

// AssertEntityRequired checks if the required fields are not zero-ed
func AssertEntityRequired(obj Entity) error {
	elements := map[string]interface{}{
		"id": obj.Id,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertEntityConstraints checks if the values respects the defined constraints
func AssertEntityConstraints(obj Entity) error {
	return nil
}
