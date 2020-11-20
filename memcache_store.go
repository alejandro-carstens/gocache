package gocache

import (
	"errors"

	"github.com/bradfitz/gomemcache/memcache"
)

// MemcacheNilErrorResponse is the gomemcache nil response error
const MemcacheNilErrorResponse = "memcache: cache miss"

// MemcacheStore is the representation of the memcache caching store
type MemcacheStore struct {
	client *memcache.Client
	prefix string
}

// Put puts a value in the given store for a predetermined amount of time in mins.
func (ms *MemcacheStore) Put(key string, value interface{}, minutes int) error {
	item, err := ms.item(key, value, minutes)
	if err != nil {
		return err
	}

	return ms.client.Set(item)
}

// Forever puts a value in the given store until it is forgotten/evicted
func (ms *MemcacheStore) Forever(key string, value interface{}) error {
	return ms.Put(key, value, 0)
}

// GetFloat64 gets a float value from the store
func (ms *MemcacheStore) GetFloat64(key string) (float64, error) {
	value, err := ms.get(key)
	if err != nil {
		return 0.0, err
	}
	if !isStringNumeric(value) {
		return 0.0, errors.New("invalid numeric value")
	}

	return stringToFloat64(value)
}

// GetInt64 gets an int value from the store
func (ms *MemcacheStore) GetInt64(key string) (int64, error) {
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
	newValue, err := ms.client.Increment(ms.GetPrefix()+key, uint64(value))
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
	newValue, err := ms.client.Decrement(ms.GetPrefix()+key, uint64(value))
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
	return ms.prefix
}

// PutMany puts many values in the given store until they are forgotten/evicted
func (ms *MemcacheStore) PutMany(values map[string]string, minutes int) error {
	for key, value := range values {
		err := ms.Put(key, value, minutes)
		if err != nil {
			return err
		}
	}

	return nil
}

// Many gets many values from the store
func (ms *MemcacheStore) Many(keys []string) (map[string]string, error) {
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

// Forget forgets/evicts a given key-value pair from the store
func (ms *MemcacheStore) Forget(key string) (bool, error) {
	if err := ms.client.Delete(ms.GetPrefix() + key); err != nil {
		return false, err
	}

	return true, nil
}

// Flush flushes the store
func (ms *MemcacheStore) Flush() (bool, error) {
	if err := ms.client.DeleteAll(); err != nil {
		return false, err
	}

	return true, nil
}

// Get gets the struct representation of a value from the store
func (ms *MemcacheStore) Get(key string, entity interface{}) error {
	value, err := ms.get(key)
	if err != nil {
		return err
	}
	_, err = decode(value, entity)

	return err
}

// Close closes the client releasing all open resources
func (ms *MemcacheStore) Close() error {
	return nil
}

// Tags returns the taggedCache for the given store
func (ms *MemcacheStore) Tags(names ...string) TaggedCache {
	return &taggedCache{
		store: ms,
		tags: tagSet{
			store: ms,
			names: names,
		},
	}
}

func (ms *MemcacheStore) Lock(name, owner string, seconds int64) Lock {
	return &memcacheLock{
		client:  ms.client,
		name:    name,
		owner:   owner,
		seconds: seconds,
	}
}

func (ms *MemcacheStore) get(key string) (string, error) {
	item, err := ms.client.Get(ms.GetPrefix() + key)
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
