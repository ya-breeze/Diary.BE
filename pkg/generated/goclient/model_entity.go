/*
Diary - OpenAPI 3.0

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.0.1
Contact: ilya.korolev@outlook.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package goclient

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// checks if the Entity type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Entity{}

// Entity struct for Entity
type Entity struct {
	Id string `json:"id"`
}

type _Entity Entity

// NewEntity instantiates a new Entity object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewEntity(id string) *Entity {
	this := Entity{}
	this.Id = id
	return &this
}

// NewEntityWithDefaults instantiates a new Entity object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewEntityWithDefaults() *Entity {
	this := Entity{}
	return &this
}

// GetId returns the Id field value
func (o *Entity) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *Entity) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *Entity) SetId(v string) {
	o.Id = v
}

func (o Entity) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Entity) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	return toSerialize, nil
}

func (o *Entity) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"id",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err
	}

	for _, requiredProperty := range requiredProperties {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varEntity := _Entity{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varEntity)

	if err != nil {
		return err
	}

	*o = Entity(varEntity)

	return err
}

type NullableEntity struct {
	value *Entity
	isSet bool
}

func (v NullableEntity) Get() *Entity {
	return v.value
}

func (v *NullableEntity) Set(val *Entity) {
	v.value = val
	v.isSet = true
}

func (v NullableEntity) IsSet() bool {
	return v.isSet
}

func (v *NullableEntity) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableEntity(val *Entity) *NullableEntity {
	return &NullableEntity{value: val, isSet: true}
}

func (v NullableEntity) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableEntity) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
