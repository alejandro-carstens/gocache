package gocache

import (
	"errors"
	"strconv"
	"time"

	"github.com/alejandro-carstens/gocache/encoder"
)

type (
	// Entry represents a cache entry or input
	Entry struct {
		Key      string
		Value    interface{}
		Duration time.Duration
	}
	// Item represents a retrieved entry from the cache
	Item struct {
		key     string
		value   string
		tagKey  string
		err     error
		encoder encoder.Encoder
	}
	// Items is an Item map keyed by the Item key
	Items map[string]Item
)

// Key returns an Item's key
func (i Item) Key() string {
	return i.key
}

// TagKey returns the actual key of an Item if it was retrieved with a tag
func (i Item) TagKey() string {
	return i.tagKey
}

// String returns the string representation of an Item's value
func (i Item) String() (string, error) {
	if i.err != nil {
		return "", ErrFailedToRetrieveEntry
	}
	if isStringNumeric(i.value) || isStringBool(i.value) {
		return i.value, nil
	}

	var v = i.value
	if err := i.encoder.Decode([]byte(v), &v); err != nil {
		return "", err
	}

	return v, nil
}

// Uint64 returns the uint64 representation of an Item's value
func (i Item) Uint64() (uint64, error) {
	if i.err != nil {
		return 0, ErrFailedToRetrieveEntry
	}
	if !isInterfaceNumericString(i.value) && !isNumeric(i.value) {
		return 0, errors.New("invalid numeric value")
	}

	return strconv.ParseUint(i.value, 10, 64)
}

// Int returns the int representation of an Item's value
func (i Item) Int() (int, error) {
	if i.err != nil {
		return 0, ErrFailedToRetrieveEntry
	}
	if !isInterfaceNumericString(i.value) && !isNumeric(i.value) {
		return 0, errors.New("invalid numeric value")
	}

	return strconv.Atoi(i.value)
}

// Bool returns the boolean representation of an Item's value
func (i Item) Bool() (bool, error) {
	if i.err != nil {
		return false, ErrFailedToRetrieveEntry
	}

	return stringToBool(i.value), nil
}

// Int64 returns the int64 representation of an Item's value
func (i Item) Int64() (int64, error) {
	if i.err != nil {
		return 0, ErrFailedToRetrieveEntry
	}
	if !isInterfaceNumericString(i.value) && !isNumeric(i.value) {
		return 0, errors.New("invalid numeric value")
	}

	return strconv.ParseInt(i.value, 10, 64)
}

// Float32 returns the float32 representation of an Item's value
func (i Item) Float32() (float32, error) {
	if i.err != nil {
		return 0, ErrFailedToRetrieveEntry
	}
	if !isInterfaceNumericString(i.value) && !isNumeric(i.value) {
		return 0, errors.New("invalid numeric value")
	}

	f, err := strconv.ParseFloat(i.value, 32)
	if err != nil {
		return 0, err
	}

	return float32(f), nil
}

// Float64 returns the float32 representation of an Item's value
func (i Item) Float64() (float64, error) {
	if i.err != nil {
		return 0, ErrFailedToRetrieveEntry
	}
	if !isInterfaceNumericString(i.value) && !isNumeric(i.value) {
		return 0, errors.New("invalid numeric value")
	}

	return strconv.ParseFloat(i.value, 64)
}

// Unmarshal decodes an Item's value to the provided entity
func (i Item) Unmarshal(entity interface{}) error {
	if i.err != nil {
		return ErrFailedToRetrieveEntry
	}

	return i.encoder.Decode([]byte(i.value), entity)
}

// Error returns the error that occurred when trying to retrieve a given Item
func (i Item) Error() error {
	return i.err
}

// EntryNotFound checks if an entry was retrieved for the given key
func (i Item) EntryNotFound() bool {
	return errors.Is(i.err, ErrNotFound)
}
