package cache

import (
	"bytes"
	"encoding/json"
	"math"
	"strconv"
)

func Encode(item interface{}) (string, error) {
	value, err := json.Marshal(item)

	return string(value), err
}

func SimpleDecode(value string) (string, error) {
	err := json.Unmarshal([]byte(value), &value)

	return string(value), err
}

func Decode(value string, entity interface{}) (interface{}, error) {
	err := json.Unmarshal([]byte(value), &entity)

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

func GetTaggedManyKey(prefix string, key string) string {
	count := len(prefix) + 41

	sub := ""
	subs := []string{}

	runs := bytes.Runes([]byte(key))

	for i, run := range runs {
		sub = sub + string(run)
		if (i+1)%count == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == len(runs) {
			subs = append(subs, sub)
		}
	}

	return subs[1]
}

func IsStringNumeric(value string) bool {
	_, err := strconv.ParseFloat(value, 64)

	return err == nil
}

func StringToFloat64(value string) float64 {
	val, err := strconv.ParseFloat(value, 64)

	if err != nil {
		panic(err)
	}

	return val
}

func IsFloat(value float64) bool {
	return value != math.Trunc(value)
}
