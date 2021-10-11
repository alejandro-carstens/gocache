package gocache

import (
	"bytes"
	"encoding/json"
	"strconv"
)

func encode(item interface{}) (string, error) {
	value, err := json.Marshal(item)

	return string(value), err
}

func simpleDecode(value string) (string, error) {
	err := json.Unmarshal([]byte(value), &value)

	return value, err
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

func getTaggedManyKey(prefix, key string) string {
	var (
		sub   string
		subs  []string
		runs  = bytes.Runes([]byte(key))
		count = len(prefix) + 41
	)
	for i, run := range runs {
		sub = sub + string(run)
		if (i+1)%count == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == len(runs) {
			subs = append(subs, sub)
		}
	}

	if len(subs) == 0 {
		return ""
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
