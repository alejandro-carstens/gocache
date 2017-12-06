package cache

import (
	"encoding/json"
)

func Encode(item interface{}) (string, error) {
	value, err := json.Marshal(item)

	if err != nil {
		panic(err)
	}

	return string(value), err
}

func SimpleDecode(value string) (string, error) {
	err := json.Unmarshal([]byte(value), &value)

	if err != nil {
		panic(err)
	}

	return value, err
}

func Decode(value string, entity interface{}) (interface{}, error) {
	err := json.Unmarshal([]byte(value), &entity)

	if err != nil {
		panic(err)
	}

	return entity, err
}

func IsNumeric(s interface{}) bool {
	switch s.(type) {
	case int:
		return true
	case int32:
		return true
	case float32:
		return true
	case float64:
		return true
	default:
		return false
	}
}
