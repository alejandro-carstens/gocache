package gocache

import (
	"errors"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

var _ Cache = &MemcacheStore{}

// NewMemcacheStore validates the passed in config and creates a Cache implementation of type *MemcacheStore
func NewMemcacheStore(cnf *MemcacheConfig) (*MemcacheStore, error) {
	if err := cnf.validate(); err != nil {
		return nil, err
	}

	client := memcache.New(cnf.Servers...)
	if cnf.MaxIdleConns > 0 {
		client.MaxIdleConns = cnf.MaxIdleConns
	}

	client.Timeout = cnf.Timeout

	return &MemcacheStore{
		client: client,
		prefix: prefix{
			val: cnf.Prefix,
		},
	}, nil
}

// MemcacheStore is the representation of the memcache caching store
type MemcacheStore struct {
	prefix
	client *memcache.Client
}

// Put puts a value in the given store for a predetermined amount of time in seconds
func (s *MemcacheStore) Put(key string, value interface{}, duration time.Duration) error {
	item, err := s.item(key, value, duration)
	if err != nil {
		return err
	}

	return s.client.Set(item)
}

// Forever puts a value in the given store until it is forgotten/evicted
func (s *MemcacheStore) Forever(key string, value interface{}) error {
	return s.Put(key, value, 0)
}

// GetFloat64 gets a float64 value from the store
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

// GetFloat32 gets a float32 value from the store
func (s *MemcacheStore) GetFloat32(key string) (float32, error) {
	value, err := s.get(key)
	if err != nil {
		return 0.0, err
	}
	if !isStringNumeric(value) {
		return 0.0, errors.New("invalid numeric value")
	}

	return stringToFloat32(value)
}

// GetInt64 gets an int64 value from the store
func (s *MemcacheStore) GetInt64(key string) (int64, error) {
	value, err := s.get(key)
	if err != nil {
		return 0, err
	}
	if !isStringNumeric(value) {
		return 0, errors.New("invalid numeric value")
	}

	return stringToInt64(value)
}

// GetInt gets an int value from the store
func (s *MemcacheStore) GetInt(key string) (int, error) {
	value, err := s.get(key)
	if err != nil {
		return 0, err
	}
	if !isStringNumeric(value) {
		return 0, errors.New("invalid numeric value")
	}

	return stringToInt(value)
}

// GetUint64 gets an uint64 value from the store
func (s *MemcacheStore) GetUint64(key string) (uint64, error) {
	value, err := s.get(key)
	if err != nil {
		return 0, err
	}
	if !isStringNumeric(value) {
		return 0, errors.New("invalid numeric value")
	}

	return stringToUint64(value)
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
	newValue, err := s.client.Increment(s.k(key), uint64(value))
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
	newValue, err := s.client.Decrement(s.k(key), uint64(value))
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

// PutMany puts many values in the given store until they are forgotten/evicted
func (s *MemcacheStore) PutMany(entries ...Entry) error {
	for _, entry := range entries {
		if err := s.Put(entry.Key, entry.Value, entry.Duration); err != nil {
			return err
		}
	}

	return nil
}

// Many gets many values from the store
func (s *MemcacheStore) Many(keys ...string) (Items, error) {
	items := Items{}
	for _, key := range keys {
		val, err := s.get(key)
		if err != nil {
			return nil, err
		}

		items[key] = Item{
			key:   key,
			value: val,
		}
	}

	return items, nil
}

// Forget forgets/evicts a given key-value pair from the store
func (s *MemcacheStore) Forget(keys ...string) (bool, error) {
	for _, key := range keys {
		if err := s.client.Delete(s.k(key)); err != nil {
			return false, err
		}
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
func (s *MemcacheStore) Lock(name, owner string, duration time.Duration) Lock {
	return &memcacheLock{
		client:   s.client,
		name:     name,
		owner:    owner,
		duration: duration,
	}
}

func (s *MemcacheStore) get(key string) (string, error) {
	item, err := s.client.Get(s.k(key))
	if err != nil {
		return "", checkErrNotFound(err)
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

func (s *MemcacheStore) item(key string, value interface{}, duration time.Duration) (*memcache.Item, error) {
	val, err := encode(value)
	if err != nil {
		return nil, err
	}

	return &memcache.Item{
		Key:        s.k(key),
		Value:      []byte(val),
		Expiration: int32(duration.Seconds()),
	}, nil
}
