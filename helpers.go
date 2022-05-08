package gocache

import (
	"encoding/json"
	"fmt"
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

func isNumeric(i interface{}) bool {
	switch i.(type) {
	case int:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case uint:
		return true
	case uintptr:
		return true
	case uint8:
		return true
	case uint16:
		return true
	case uint32:
		return true
	case uint64:
		return true
	case float32:
		return true
	case float64:
		return true
	default:
		return false
	}
}

func isBool(i interface{}) bool {
	switch i.(type) {
	case bool:
		return true
	default:
		return false
	}
}

func isStringNumeric(value string) bool {
	_, err := strconv.ParseFloat(value, 64)

	return err == nil
}

func isStringBool(value string) bool {
	_, err := strconv.ParseBool(value)

	return err == nil
}

func stringToFloat64(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

func interfaceToFloat64(value interface{}) (float64, error) {
	return stringToFloat64(fmt.Sprint(value))
}

func stringToFloat32(value string) (float32, error) {
	n, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0, err
	}

	return float32(n), nil
}

func interfaceToFloat32(value interface{}) (float32, error) {
	return stringToFloat32(fmt.Sprint(value))
}

func stringToInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func interfaceToInt64(value interface{}) (int64, error) {
	return stringToInt64(fmt.Sprint(value))
}

func stringToInt(value string) (int, error) {
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(n), nil
}

func interfaceToInt(value interface{}) (int, error) {
	return stringToInt(fmt.Sprint(value))
}

func stringToUint64(value string) (uint64, error) {
	return strconv.ParseUint(value, 10, 64)
}

func interfaceToUint64(value interface{}) (uint64, error) {
	return stringToUint64(fmt.Sprint(value))
}

func isInterfaceNumericString(value interface{}) bool {
	str, valid := value.(string)
	if !valid {
		return false
	}

	return isStringNumeric(str)
}

func stringToBool(value string) bool {
	// If the cache val is '0' or 'false' we return false
	if len(value) > 0 && (value == "0" || value == "false" || value == `""`) {
		return false
	}
	// If the cache val is empty we return false
	if len(value) == 0 {
		return false
	}

	return true
}
