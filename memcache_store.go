package gocache

import (
	"errors"

	"github.com/bradfitz/gomemcache/memcache"
)

// MemcacheNilErrorResponse is the gomemcache nil response error
const MemcacheNilErrorResponse = "memcache: cache miss"

// MemcacheStore is the representation of the memcache caching store
type MemcacheStore struct {
	Client memcache.Client
	Prefix string
}

// Put puts a value in the given store for a predetermined amount of time in mins.
func (ms *MemcacheStore) Put(key string, value interface{}, minutes int) error {
	item, err := ms.item(key, value, minutes)
	if err != nil {
		return err
	}

	return ms.Client.Set(item)
}

// Forever puts a value in the given store until it is forgotten/evicted
func (ms *MemcacheStore) Forever(key string, value interface{}) error {
	return ms.Put(key, value, 0)
}

// Get gets a value from the store
func (ms *MemcacheStore) Get(key string) (interface{}, error) {
	value, err := ms.get(key)
	if err != nil {
		return nil, err
	}

	return ms.processValue(value)
}

// GetFloat gets a float value from the store
func (ms *MemcacheStore) GetFloat(key string) (float64, error) {
	value, err := ms.get(key)
	if err != nil {
		return 0.0, err
	}
	if !isStringNumeric(value) {
		return 0.0, errors.New("invalid numeric value")
	}

	return stringToFloat64(value)
}

// GetInt gets an int value from the store
func (ms *MemcacheStore) GetInt(key string) (int64, error) {
	value, err := ms.get(key)

	if err != nil {
		return 0, err
	}
	if !isStringNumeric(value) {
		return 0, errors.New("invalid numeric value")
	}

	val, err := stringToFloat64(value)

	return int64(val), err
}

// GetString gets a string value from the store
func (ms *MemcacheStore) GetString(key string) (string, error) {
	value, err := ms.get(key)
	if err != nil {
		return "", err
	}

	return value, nil
}

// Increment increments an integer counter by a given value
func (ms *MemcacheStore) Increment(key string, value int64) (int64, error) {
	newValue, err := ms.Client.Increment(ms.GetPrefix()+key, uint64(value))
	if err != nil {
		if err.Error() != MemcacheNilErrorResponse {
			return value, err
		}
		if err := ms.Put(key, value, 0); err != nil {
			return 0, err
		}

		return value, nil
	}

	return int64(newValue), nil
}

// Decrement decrements an integer counter by a given value
func (ms *MemcacheStore) Decrement(key string, value int64) (int64, error) {
	newValue, err := ms.Client.Decrement(ms.GetPrefix()+key, uint64(value))
	if err != nil {
		if err.Error() != MemcacheNilErrorResponse {
			return value, err
		}
		if err := ms.Put(key, 0, 0); err != nil {
			return 0, err
		}

		return 0, nil
	}

	return int64(newValue), nil
}

// GetPrefix gets the cache key prefix
func (ms *MemcacheStore) GetPrefix() string {
	return ms.Prefix
}

// PutMany puts many values in the given store until they are forgotten/evicted
func (ms *MemcacheStore) PutMany(values map[string]interface{}, minutes int) error {
	for key, value := range values {
		err := ms.Put(key, value, minutes)
		if err != nil {
			return err
		}
	}

	return nil
}

// Many gets many values from the store
func (ms *MemcacheStore) Many(keys []string) (map[string]interface{}, error) {
	items := make(map[string]interface{})

	for _, key := range keys {
		val, err := ms.Get(key)
		if err != nil {
			return items, err
		}

		items[key] = val
	}

	return items, nil
}

// Forget forgets/evicts a given key-value pair from the store
func (ms *MemcacheStore) Forget(key string) (bool, error) {
	if err := ms.Client.Delete(ms.GetPrefix() + key); err != nil {
		return false, err
	}

	return true, nil
}

// Flush flushes the store
func (ms *MemcacheStore) Flush() (bool, error) {
	if err := ms.Client.DeleteAll(); err != nil {
		return false, err
	}

	return true, nil
}

// GetStruct gets the struct representation of a value from the store
func (ms *MemcacheStore) GetStruct(key string, entity interface{}) error {
	value, err := ms.get(key)
	if err != nil {
		return err
	}

	_, err = decode(value, entity)

	return err
}

// Tags returns the TaggedCache for the given store
func (ms *MemcacheStore) Tags(names ...string) TaggedStore {
	return &TaggedCache{
		Store: ms,
		Tags: TagSet{
			Store: ms,
			Names: names,
		},
	}
}

func (ms *MemcacheStore) get(key string) (string, error) {
	item, err := ms.Client.Get(ms.GetPrefix() + key)
	if err != nil {
		return "", err
	}

	return ms.getItemValue(item.Value), nil
}

func (ms *MemcacheStore) getItemValue(itemValue []byte) string {
	value, err := simpleDecode(string(itemValue))
	if err != nil {
		return string(itemValue)
	}

	return value
}

func (ms *MemcacheStore) processValue(value string) (interface{}, error) {
	if isStringNumeric(value) {
		floatValue, err := stringToFloat64(value)
		if err != nil {
			return floatValue, err
		}
		if isFloat(floatValue) {
			return floatValue, err
		}

		return int64(floatValue), err
	}

	return value, nil
}

func (ms *MemcacheStore) item(key string, value interface{}, minutes int) (*memcache.Item, error) {
	val, err := encode(value)
	if err != nil {
		return nil, err
	}

	return &memcache.Item{
		Key:        ms.GetPrefix() + key,
		Value:      []byte(val),
		Expiration: int32(minutes),
	}, nil
}
