package gocache

import (
	"errors"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

var _ Cache = &LocalStore{}

// LocalStore is the representation of a map caching store
type LocalStore struct {
	prefix
	c                 *cache.Cache
	defaultExpiration time.Duration
	defaultInterval   time.Duration
}

// GetString gets a string value from the store
func (s *LocalStore) GetString(key string) (string, error) {
	value, valid := s.c.Get(s.k(key))
	if !valid {
		return "", ErrNotFound
	}

	return simpleDecode(fmt.Sprint(value))
}

// GetFloat64 gets a float value from the store
func (s *LocalStore) GetFloat64(key string) (float64, error) {
	value, valid := s.c.Get(s.k(key))
	if !valid {
		return 0, ErrNotFound
	}
	if !isStringNumeric(value.(string)) {
		return 0, errors.New("invalid numeric value")
	}

	return stringToFloat64(value.(string))
}

// GetInt64 gets an int value from the store
func (s *LocalStore) GetInt64(key string) (int64, error) {
	value, valid := s.c.Get(s.k(key))
	if !valid {
		return 0, ErrNotFound
	}
	if !isStringNumeric(value.(string)) {
		return 0, errors.New("invalid numeric value")
	}

	val, err := stringToFloat64(value.(string))

	return int64(val), err
}

// Increment increments an integer counter by a given value
func (s *LocalStore) Increment(key string, value int64) (int64, error) {
	val, valid := s.c.Get(s.k(key))
	if valid && isStringNumeric(val.(string)) {
		floatVal, err := stringToFloat64(val.(string))
		if err != nil {
			return 0, err
		}

		result := value + int64(floatVal)

		return result, s.Put(key, result, -1)
	}

	return value, s.Put(key, value, -1)
}

// Decrement decrements an integer counter by a given value
func (s *LocalStore) Decrement(key string, value int64) (int64, error) {
	return s.Increment(key, -value)
}

// Put puts a value in the given store for a predetermined amount of time in seconds.
func (s *LocalStore) Put(key string, value interface{}, seconds int) error {
	val, err := encode(value)
	if err != nil {
		return err
	}

	s.c.Set(s.GetPrefix()+key, val, time.Duration(seconds)*time.Second)

	return nil
}

// Forever puts a value in the given store until it is forgotten/evicted
func (s *LocalStore) Forever(key string, value interface{}) error {
	return s.Put(key, value, -1)
}

// Flush flushes the store
func (s *LocalStore) Flush() (bool, error) {
	s.c = cache.New(s.defaultExpiration, s.defaultInterval)

	return true, nil
}

// Forget forgets/evicts a given key-value pair from the store
func (s *LocalStore) Forget(key string) (bool, error) {
	if _, valid := s.c.Get(s.k(key)); !valid {
		return false, nil
	}

	s.c.Delete(s.k(key))

	return true, nil
}

// PutMany puts many values in the given store until they are forgotten/evicted
func (s *LocalStore) PutMany(values map[string]string, minutes int) error {
	for key, value := range values {
		if err := s.Put(key, value, minutes); err != nil {
			return err
		}
	}

	return nil
}

// Many gets many values from the store
func (s *LocalStore) Many(keys []string) (map[string]string, error) {
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

// Get gets the struct representation of a value from the store
func (s *LocalStore) Get(key string, entity interface{}) error {
	value, valid := s.c.Get(s.GetPrefix() + key)
	if !valid {
		return ErrNotFound
	}

	_, err := decode(fmt.Sprint(value), entity)

	return err
}

// Close closes the c releasing all open resources
func (s *LocalStore) Close() error {
	return nil
}

// Tags returns the taggedCache for the given store
func (s *LocalStore) Tags(names ...string) TaggedCache {
	return &taggedCache{
		store: s,
		tags: tagSet{
			store: s,
			names: names,
		},
	}
}

// Lock returns a map implementation of the Lock interface
func (s *LocalStore) Lock(name, owner string, seconds int64) Lock {
	return &localLock{
		c:       s.c,
		name:    name,
		owner:   owner,
		seconds: seconds,
	}
}
