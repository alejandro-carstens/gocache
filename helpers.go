package gocache

import (
	"bytes"
	"encoding/json"
	"math"
	"strconv"
)

func encode(item interface{}) (string, error) {
	value, err := json.Marshal(item)

	return string(value), err
}

func simpleDecode(value string) (string, error) {
	err := json.Unmarshal([]byte(value), &value)

	return string(value), err
}

func decode(value string, entity interface{}) (interface{}, error) {
	err := json.Unmarshal([]byte(value), &entity)

	return entity, err
}

func isNumeric(s interface{}) bool {
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

func getTaggedManyKey(prefix string, key string) string {
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

func isStringNumeric(value string) bool {
	_, err := strconv.ParseFloat(value, 64)

	return err == nil
}

func stringToFloat64(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

func isFloat(value float64) bool {
	return value != math.Trunc(value)
}

func isCacheMissedError(err error) bool {
	haystack := []string{MAP_NIL_ERROR_RESPONSE, MEMCACHE_NIL_ERROR_RESPONSE, REDIS_NIL_ERROR_RESPONSE}

	return inStringSlice(err.Error(), haystack)
}

func inStringSlice(needle string, haystack []string) bool {
	for _, value := range haystack {
		if needle == value {
			return true
		}
	}

	return false
}
