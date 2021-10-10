package gocache

import (
	"errors"

	"github.com/bradfitz/gomemcache/memcache"
)

var _ Cache = &MemcacheStore{}

// MemcacheStore is the representation of the memcache caching store
type MemcacheStore struct {
	client *memcache.Client
	prefix string
}

// Put puts a value in the given store for a predetermined amount of time in seconds
func (s *MemcacheStore) Put(key string, value interface{}, seconds int) error {
	item, err := s.item(key, value, seconds)
	if err != nil {
		return err
	}

	return s.client.Set(item)
}

// Forever puts a value in the given store until it is forgotten/evicted
func (s *MemcacheStore) Forever(key string, value interface{}) error {
	return s.Put(key, value, 0)
}

// GetFloat64 gets a float value from the store
func (s *MemcacheStore) GetFloat64(key string) (float64, error) {
	value, err := s.get(key)
	if err != nil {
		return 0.0, err
	}
	if !isStringNumeric(value) {
		return 0.0, errors.New("invalid numeric value")
	}

	return stringToFloat64(value)
}

// GetInt64 gets an int value from the store
func (s *MemcacheStore) GetInt64(key string) (int64, error) {
	value, err := s.get(key)
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
func (s *MemcacheStore) GetString(key string) (string, error) {
	value, err := s.get(key)
	if err != nil {
		return "", err
	}

	return value, nil
}

// Increment increments an integer counter by a given value
func (s *MemcacheStore) Increment(key string, value int64) (int64, error) {
	newValue, err := s.client.Increment(s.GetPrefix()+key, uint64(value))
	if err != nil {
		if !errors.Is(err, memcache.ErrCacheMiss) {
			return value, err
		}
		if err = s.Put(key, value, 0); err != nil {
			return 0, err
		}

		return value, nil
	}

	return int64(newValue), nil
}

// Decrement decrements an integer counter by a given value
func (s *MemcacheStore) Decrement(key string, value int64) (int64, error) {
	newValue, err := s.client.Decrement(s.GetPrefix()+key, uint64(value))
	if err != nil {
		if !errors.Is(err, memcache.ErrCacheMiss) {
			return value, err
		}
		if err = s.Put(key, 0, 0); err != nil {
			return 0, err
		}

		return 0, nil
	}

	return int64(newValue), nil
}

// GetPrefix gets the cache key prefix
func (s *MemcacheStore) GetPrefix() string {
	return s.prefix
}

// PutMany puts many values in the given store until they are forgotten/evicted
func (s *MemcacheStore) PutMany(values map[string]string, seconds int) error {
	for key, value := range values {
		if err := s.Put(key, value, seconds); err != nil {
			return err
		}
	}

	return nil
}

// Many gets many values from the store
func (s *MemcacheStore) Many(keys []string) (map[string]string, error) {
	items := make(map[string]string)
	for _, key := range keys {
		val, err := s.GetString(key)
		if err != nil {
			return items, err
		}

		items[key] = val
	}

	return items, nil
}

// Forget forgets/evicts a given key-value pair from the store
func (s *MemcacheStore) Forget(key string) (bool, error) {
	if err := s.client.Delete(s.GetPrefix() + key); err != nil {
		return false, err
	}

	return true, nil
}

// Flush flushes the store
func (s *MemcacheStore) Flush() (bool, error) {
	if err := s.client.DeleteAll(); err != nil {
		return false, err
	}

	return true, nil
}

// Get gets the struct representation of a value from the store
func (s *MemcacheStore) Get(key string, entity interface{}) error {
	value, err := s.get(key)
	if err != nil {
		return err
	}
	_, err = decode(value, entity)

	return err
}

// Close closes the c releasing all open resources
func (s *MemcacheStore) Close() error {
	return nil
}

// Tags returns the taggedCache for the given store
func (s *MemcacheStore) Tags(names ...string) TaggedCache {
	return &taggedCache{
		store: s,
		tags: tagSet{
			store: s,
			names: names,
		},
	}
}

// Lock returns a memcache implementation of the Lock interface
func (s *MemcacheStore) Lock(name, owner string, seconds int64) Lock {
	return &memcacheLock{
		client:  s.client,
		name:    name,
		owner:   owner,
		seconds: seconds,
	}
}

func (s *MemcacheStore) get(key string) (string, error) {
	item, err := s.client.Get(s.GetPrefix() + key)
	if err != nil {
		return "", err
	}

	return s.getItemValue(item.Value), nil
}

func (s *MemcacheStore) getItemValue(itemValue []byte) string {
	value, err := simpleDecode(string(itemValue))
	if err != nil {
		return string(itemValue)
	}

	return value
}

func (s *MemcacheStore) item(key string, value interface{}, seconds int) (*memcache.Item, error) {
	val, err := encode(value)
	if err != nil {
		return nil, err
	}

	return &memcache.Item{
		Key:        s.GetPrefix() + key,
		Value:      []byte(val),
		Expiration: int32(seconds),
	}, nil
}
