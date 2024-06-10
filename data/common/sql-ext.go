package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Json array

type JsonArray[T any] []T

func (j JsonArray[T]) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JsonArray[T]) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, &j)
	}
	return errors.New("invalid src type")
}

// Json map

type JsonMap[K comparable, V any] map[K]V

func (j JsonMap[K, V]) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JsonMap[K, V]) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, &j)
	}
	return errors.New("invalid src type")
}
