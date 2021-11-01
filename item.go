package gocache

import (
	"errors"
	"strconv"
	"time"
)

type (
	Entry struct {
		Key      string
		Value    interface{}
		Duration time.Duration
	}
	Item struct {
		key    string
		value  string
		tagKey string
	}
	Items map[string]Item
)

func (i Item) Key() string {
	return i.key
}

func (i Item) TagKey() string {
	return i.tagKey
}

func (i Item) String() string {
	s, err := simpleDecode(i.value)
	if err != nil {
		return i.value
	}

	return s
}

func (i Item) Uint64() (uint64, error) {
	if !isInterfaceNumericString(i.value) && !isNumeric(i.value) {
		return 0, errors.New("invalid numeric value")
	}

	return strconv.ParseUint(i.value, 10, 64)
}

func (i Item) Int() (int, error) {
	if !isInterfaceNumericString(i.value) && !isNumeric(i.value) {
		return 0, errors.New("invalid numeric value")
	}

	return strconv.Atoi(i.value)
}

func (i Item) Int64() (int64, error) {
	if !isInterfaceNumericString(i.value) && !isNumeric(i.value) {
		return 0, errors.New("invalid numeric value")
	}

	return strconv.ParseInt(i.value, 10, 64)
}

func (i Item) Float32() (float32, error) {
	if !isInterfaceNumericString(i.value) && !isNumeric(i.value) {
		return 0, errors.New("invalid numeric value")
	}

	f, err := strconv.ParseFloat(i.value, 32)
	if err != nil {
		return 0, err
	}

	return float32(f), nil
}

func (i Item) Float64() (float64, error) {
	if !isInterfaceNumericString(i.value) && !isNumeric(i.value) {
		return 0, errors.New("invalid numeric value")
	}

	return strconv.ParseFloat(i.value, 64)
}

func (i Item) Unmarshal(entity interface{}) error {
	_, err := decode(i.value, entity)

	return err
}
