package models

import (
	"database/sql/driver"
	"encoding/json"
)

type ArraySegments []Segment
type ArrayUint []uint

func (array ArraySegments) Value() (driver.Value, error) {
	return anyValue(array)
}

func (array *ArraySegments) Scan(src interface{}) error {
	source, err := sourceToBytes(src)
	if err != nil {
		return err
	}

	var tmpArray ArraySegments
	err = json.Unmarshal(source, &tmpArray)
	if err != nil {
		return err
	}

	*array = tmpArray

	return nil
}

func (array ArrayUint) Value() (driver.Value, error) {
	return anyValue(array)
}

func (array *ArrayUint) Scan(src interface{}) error {
	i, err := anyScan(src)
	if err != nil {
		return err
	}

	tmpA := i.([]interface{})
	tmpB := make([]uint, len(tmpA))
	for j := range tmpA {
		tmpB[j] = uint(tmpA[j].(float64))
	}

	*array = tmpB

	return nil
}

// ...

func anyValue(value any) (string, error) {
	jsonValue, err := json.Marshal(value)
	return string(jsonValue), err
}

func sourceToBytes(src interface{}) ([]byte, error) {
	var source []byte
	switch v := src.(type) {
	case []byte:
		source = v
	case string:
		source = []byte(v)
	default:
		return nil, ErrTypeAssertionFailed
	}
	return source, nil
}

func anyScan(src interface{}) (interface{}, error) {
	source, err := sourceToBytes(src)
	if err != nil {
		return nil, err
	}

	var i interface{}
	err = json.Unmarshal(source, &i)

	return i, err
}
