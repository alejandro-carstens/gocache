package gocache

import (
	"errors"
	"fmt"
	"strconv"
)

// MAP_NIL_ERROR_RESPONSE map nil response error
const MapNilErrorResponse = "map: cache miss"

// MapStore is the representation of a map caching store
type MapStore struct {
	client map[string]interface{}
	prefix string
}

// GetString gets a string value from the store
func (ms *MapStore) GetString(key string) (string, error) {
	value, valid := ms.client[ms.GetPrefix()+key]
	if !valid {
		return "", errors.New(MapNilErrorResponse)
	}

	return simpleDecode(fmt.Sprint(value))
}

// GetFloat64 gets a float value from the store
func (ms *MapStore) GetFloat64(key string) (float64, error) {
	value, valid := ms.client[ms.GetPrefix()+key]
	if !valid {
		return 0, errors.New(MapNilErrorResponse)
	}
	if !isStringNumeric(value.(string)) {
		return 0, errors.New("invalid numeric value")
	}

	return stringToFloat64(value.(string))
}

// GetInt64 gets an int value from the store
func (ms *MapStore) GetInt64(key string) (int64, error) {
	value, valid := ms.client[ms.GetPrefix()+key]
	if !valid {
		return 0, errors.New(MapNilErrorResponse)
	}
	if !isStringNumeric(value.(string)) {
		return 0, errors.New("invalid numeric value")
	}

	val, err := stringToFloat64(value.(string))

	return int64(val), err
}

// Increment increments an integer counter by a given value
func (ms *MapStore) Increment(key string, value int64) (int64, error) {
	val := ms.client[ms.GetPrefix()+key]
	if val != nil {
		if isStringNumeric(val.(string)) {
			floatValue, err := stringToFloat64(val.(string))
			if err != nil {
				return 0, err
			}

			result := value + int64(floatValue)

			return result, ms.Put(key, result, 0)
		}
	}

	return value, ms.Put(key, value, 0)
}

// Decrement decrements an integer counter by a given value
func (ms *MapStore) Decrement(key string, value int64) (int64, error) {
	return ms.Increment(key, -value)
}

// Put puts a value in the given store for a predetermined amount of time in mins.
func (ms *MapStore) Put(key string, value interface{}, minutes int) error {
	val, err := encode(value)
	if err != nil {
		return err
	}

	mins := strconv.Itoa(minutes)

	mins = ""

	ms.client[ms.GetPrefix()+key+mins] = val

	return nil
}

// Forever puts a value in the given store until it is forgotten/evicted
func (ms *MapStore) Forever(key string, value interface{}) error {
	return ms.Put(key, value, 0)
}

// Flush flushes the store
func (ms *MapStore) Flush() (bool, error) {
	ms.client = make(map[string]interface{})

	return true, nil
}

// Forget forgets/evicts a given key-value pair from the store
func (ms *MapStore) Forget(key string) (bool, error) {
	if _, ok := ms.client[ms.GetPrefix()+key]; ok {
		delete(ms.client, ms.GetPrefix()+key)

		return true, nil
	}

	return false, nil
}

// GetPrefix gets the cache key prefix
func (ms *MapStore) GetPrefix() string {
	return ms.prefix
}

// PutMany puts many values in the given store until they are forgotten/evicted
func (ms *MapStore) PutMany(values map[string]string, minutes int) error {
	for key, value := range values {
		if err := ms.Put(key, value, minutes); err != nil {
			return err
		}
	}

	return nil
}

// Many gets many values from the store
func (ms *MapStore) Many(keys []string) (map[string]string, error) {
	items := make(map[string]string)

	for _, key := range keys {
		val, err := ms.GetString(key)
		if err != nil {
			return items, err
		}

		items[key] = val
	}

	return items, nil
}

// Get gets the struct representation of a value from the store
func (ms *MapStore) Get(key string, entity interface{}) error {
	value, valid := ms.client[ms.GetPrefix()+key]
	if !valid {
		return errors.New(MapNilErrorResponse)
	}

	_, err := decode(fmt.Sprint(value), entity)

	return err
}

// Close closes the client releasing all open resources
func (ms *MapStore) Close() error {
	return nil
}

// Tags returns the TaggedCache for the given store
func (ms *MapStore) Tags(names ...string) TaggedStore {
	return &TaggedCache{
		store: ms,
		tags: TagSet{
			Store: ms,
			Names: names,
		},
	}
}
