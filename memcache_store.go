package gocache

import (
	"errors"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/alejandro-carstens/gocache/encoder"
)

var _ Cache = &MemcacheStore{}

// NewMemcacheStore validates the passed in config and creates a Cache implementation of type *MemcacheStore
func NewMemcacheStore(cnf *MemcacheConfig, encoder encoder.Encoder) (*MemcacheStore, error) {
	if err := cnf.validate(); err != nil {
		return nil, err
	}

	client := memcache.New(cnf.Servers...)
	if cnf.MaxIdleConns > 0 {
		client.MaxIdleConns = cnf.MaxIdleConns
	}

	client.Timeout = cnf.Timeout

	return &MemcacheStore{
		prefix: prefix{
			val: cnf.Prefix,
		},
		client:  client,
		encoder: encoder,
	}, nil
}

// MemcacheStore is the representation of the memcache caching store
type MemcacheStore struct {
	prefix
	client  *memcache.Client
	encoder encoder.Encoder
}

// Put puts a value in the given store for a predetermined amount of time in seconds
func (s *MemcacheStore) Put(key string, value interface{}, duration time.Duration) error {
	item, err := s.item(key, value, duration)
	if err != nil {
		return err
	}

	return s.client.Set(item)
}

// Add an item to the cache only if an item doesn't already exist for the given key, or if the existing item has
// expired. If the record was successfully added true will be returned else false will be returned
func (s *MemcacheStore) Add(key string, value interface{}, duration time.Duration) (bool, error) {
	item, err := s.item(key, value, duration)
	if err != nil {
		return false, err
	}
	if err = s.client.Add(item); errors.Is(err, memcache.ErrNotStored) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// Forever puts a value in the given store until it is forgotten/evicted
func (s *MemcacheStore) Forever(key string, value interface{}) error {
	return s.Put(key, value, 0)
}

// GetFloat64 gets a float64 value from the store
func (s *MemcacheStore) GetFloat64(key string) (float64, error) {
	value, err := s.value(key)
	if err != nil {
		return 0, err
	}
	if !isStringNumeric(value) {
		return 0, errors.New("invalid numeric value")
	}

	return stringToFloat64(value)
}

// GetFloat32 gets a float32 value from the store
func (s *MemcacheStore) GetFloat32(key string) (float32, error) {
	value, err := s.value(key)
	if err != nil {
		return 0, err
	}
	if !isStringNumeric(value) {
		return 0, errors.New("invalid numeric value")
	}

	return stringToFloat32(value)
}

// GetInt64 gets an int64 value from the store
func (s *MemcacheStore) GetInt64(key string) (int64, error) {
	value, err := s.value(key)
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
	value, err := s.value(key)
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
	value, err := s.value(key)
	if err != nil {
		return 0, err
	}
	if !isStringNumeric(value) {
		return 0, errors.New("invalid numeric val")
	}

	return stringToUint64(value)
}

// GetBool gets a bool value from the store
func (s *MemcacheStore) GetBool(key string) (bool, error) {
	value, err := s.value(key)
	if err != nil {
		return false, checkErrNotFound(err)
	}
	if isStringNumeric(value) || isStringBool(value) {
		return stringToBool(value), nil
	}
	if err = s.encoder.Decode([]byte(value), &value); err != nil {
		return false, err
	}

	return stringToBool(value), nil
}

// GetString gets a string value from the store
func (s *MemcacheStore) GetString(key string) (string, error) {
	item, err := s.client.Get(s.k(key))
	if err != nil {
		return "", checkErrNotFound(err)
	}

	var v = string(item.Value)
	if isStringNumeric(v) || isStringBool(v) {
		return v, nil
	}
	if err = s.encoder.Decode(item.Value, &v); err != nil {
		return "", err
	}

	return v, nil
}

// Increment increments an integer counter by a given value
func (s *MemcacheStore) Increment(key string, value int64) (int64, error) {
	res, err := s.client.Increment(s.k(key), uint64(value))
	if err != nil {
		if !errors.Is(err, memcache.ErrCacheMiss) {
			return 0, err
		}
		if value < 0 {
			value = 0
		}

		added, err := s.Add(key, value, 0)
		if err != nil {
			return 0, err
		}
		if !added {
			return 0, ErrFailedToAddItemEntry
		}

		return value, nil
	}

	return int64(res), nil
}

// Decrement decrements an integer counter by a given value. Please note that for memcache a new value will be
// capped at 0
func (s *MemcacheStore) Decrement(key string, value int64) (int64, error) {
	newValue, err := s.client.Decrement(s.k(key), uint64(value))
	if err != nil {
		if !errors.Is(err, memcache.ErrCacheMiss) {
			return 0, err
		}

		added, err := s.Add(key, 0, 0)
		if err != nil {
			return 0, err
		}
		if !added {
			return 0, ErrFailedToAddItemEntry
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
		item, err := s.client.Get(s.k(key))
		if err != nil {
			items[key] = Item{
				key: key,
				err: checkErrNotFound(err),
			}

			continue
		}

		v := string(item.Value)
		if isStringNumeric(v) || isStringBool(v) {
			items[key] = Item{
				key:     key,
				value:   v,
				encoder: s.encoder,
			}

			continue
		}

		items[key] = Item{
			key:     key,
			value:   v,
			encoder: s.encoder,
		}
	}

	return items, nil
}

// Forget forgets/evicts a given key-value pair from the store
func (s *MemcacheStore) Forget(key string) (bool, error) {
	if err := s.client.Delete(s.k(key)); errors.Is(err, memcache.ErrCacheMiss) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// ForgetMany forgets/evicts a set of given key-value pair from the store
func (s *MemcacheStore) ForgetMany(keys ...string) error {
	for _, key := range keys {
		if err := s.client.Delete(s.k(key)); err != nil && !errors.Is(err, memcache.ErrCacheMiss) {
			return err
		}
	}

	return nil
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
	item, err := s.client.Get(s.k(key))
	if err != nil {
		return checkErrNotFound(err)
	}

	return s.encoder.Decode(item.Value, &entity)
}

// Close closes the c releasing all open resources
func (*MemcacheStore) Close() error {
	return nil
}

// Tags returns the taggedCache for the given store
func (s *MemcacheStore) Tags(names ...string) TaggedCache {
	return &taggedCache{
		store: s,
		tags: &TagSet{
			store: s,
			names: names,
		},
	}
}

// Lock returns a memcache implementation of the Lock interface
func (s *MemcacheStore) Lock(name, owner string, duration time.Duration) Lock {
	return newMemcacheLock(s.client, name, owner, duration)
}

// Exists checks if an entry exists in the cache for the given key
func (s *MemcacheStore) Exists(key string) (bool, error) {
	_, err := s.client.Get(s.k(key))
	if err == nil {
		return true, nil
	} else if err != nil && isErrNotFound(err) {
		return false, nil
	}

	return false, err
}

// Expire implementation of the Cache interface
func (s *MemcacheStore) Expire(key string, duration time.Duration) error {
	if err := s.client.Touch(s.k(key), int32(duration.Seconds())); err != nil {
		return checkErrNotFound(err)
	}

	return nil
}

func (s *MemcacheStore) value(key string) (string, error) {
	item, err := s.client.Get(s.k(key))
	if err != nil {
		return "", checkErrNotFound(err)
	}

	return string(item.Value), nil
}

func (s *MemcacheStore) item(key string, value interface{}, duration time.Duration) (*memcache.Item, error) {
	var (
		val []byte
		err error
	)
	if isNumeric(value) || isBool(value) {
		val = []byte(fmt.Sprint(value))
	} else {
		val, err = s.encoder.Encode(value)
	}
	if err != nil {
		return nil, err
	}

	return &memcache.Item{
		Key:        s.k(key),
		Value:      val,
		Expiration: int32(duration.Seconds()),
	}, nil
}
