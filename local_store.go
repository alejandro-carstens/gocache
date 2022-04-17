package gocache

import (
	"errors"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

var _ Cache = &LocalStore{}

// NewLocalStore validates the passed in config and creates a Cache implementation of type *LocalStore
func NewLocalStore(cnf *LocalConfig) (*LocalStore, error) {
	if err := cnf.validate(); err != nil {
		return nil, nil
	}

	return &LocalStore{
		c:                 cache.New(cnf.DefaultExpiration, cnf.DefaultInterval),
		defaultExpiration: cnf.DefaultExpiration,
		defaultInterval:   cnf.DefaultInterval,
		prefix: prefix{
			val: cnf.Prefix,
		},
	}, nil
}

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
	if !isInterfaceNumericString(value) && !isNumeric(value) {
		return 0, errors.New("invalid numeric value")
	}

	return interfaceToFloat64(value)
}

// GetFloat32 gets a float32 value from the store
func (s *LocalStore) GetFloat32(key string) (float32, error) {
	value, valid := s.c.Get(s.k(key))
	if !valid {
		return 0, ErrNotFound
	}
	if !isInterfaceNumericString(value) && !isNumeric(value) {
		return 0, errors.New("invalid numeric value")
	}

	return interfaceToFloat32(value)
}

// GetInt64 gets an int value from the store
func (s *LocalStore) GetInt64(key string) (int64, error) {
	value, valid := s.c.Get(s.k(key))
	if !valid {
		return 0, ErrNotFound
	}
	if !isInterfaceNumericString(value) && !isNumeric(value) {
		return 0, errors.New("invalid numeric value")
	}

	return interfaceToInt64(value)
}

// GetInt gets an int value from the store
func (s *LocalStore) GetInt(key string) (int, error) {
	value, valid := s.c.Get(s.k(key))
	if !valid {
		return 0, ErrNotFound
	}
	if !isInterfaceNumericString(value) && !isNumeric(value) {
		return 0, errors.New("invalid numeric value")
	}

	return interfaceToInt(value)
}

// GetUint64 gets an uint64 value from the store
func (s *LocalStore) GetUint64(key string) (uint64, error) {
	value, valid := s.c.Get(s.k(key))
	if !valid {
		return 0, ErrNotFound
	}
	if !isInterfaceNumericString(value) && !isNumeric(value) {
		return 0, errors.New("invalid numeric value")
	}

	return interfaceToUint64(value)
}

// GetBool gets a bool value from the store
func (s *LocalStore) GetBool(key string) (bool, error) {
	value, valid := s.c.Get(s.k(key))
	if !valid {
		return false, ErrNotFound
	}

	return stringToBool(fmt.Sprint(value)), nil
}

// Increment increments an integer counter by a given value
func (s *LocalStore) Increment(key string, value int64) (int64, error) {
	if _, valid := s.c.Get(s.k(key)); !valid {
		if err := s.Forever(key, value); err != nil {
			return 0, err
		}

		return value, nil
	}
	if err := s.c.Increment(s.k(key), value); err != nil {
		return 0, err
	}

	return s.GetInt64(key)
}

// Decrement decrements an integer counter by a given value
func (s *LocalStore) Decrement(key string, value int64) (int64, error) {
	if _, valid := s.c.Get(s.k(key)); !valid {
		if err := s.Forever(key, -1*value); err != nil {
			return 0, err
		}

		return value, nil
	}
	if err := s.c.Decrement(s.k(key), value); err != nil {
		return 0, err
	}

	return s.GetInt64(key)
}

// Put puts a value in the given store for a predetermined amount of time in seconds.
func (s *LocalStore) Put(key string, value interface{}, duration time.Duration) error {
	if isNumeric(value) {
		s.c.Set(s.k(key), value, duration)

		return nil
	}

	val, err := encode(value)
	if err != nil {
		return err
	}

	s.c.Set(s.k(key), val, duration)

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
func (s *LocalStore) Forget(keys ...string) (bool, error) {
	for _, key := range keys {
		if _, valid := s.c.Get(s.k(key)); !valid {
			return false, nil
		}

		s.c.Delete(s.k(key))
	}

	return true, nil
}

// PutMany puts many values in the given store until they are forgotten/evicted
func (s *LocalStore) PutMany(entries ...Entry) error {
	for _, entry := range entries {
		if err := s.Put(entry.Key, entry.Value, entry.Duration); err != nil {
			return err
		}
	}

	return nil
}

// Many gets many values from the store
func (s *LocalStore) Many(keys ...string) (Items, error) {
	items := Items{}
	for _, key := range keys {
		val, valid := s.c.Get(s.k(key))
		if !valid {
			items[key] = Item{
				key: key,
				err: ErrNotFound,
			}

			continue
		}

		items[key] = Item{
			key:   key,
			value: fmt.Sprint(val),
		}
	}

	return items, nil
}

// Get gets the struct representation of a value from the store
func (s *LocalStore) Get(key string, entity interface{}) error {
	value, valid := s.c.Get(s.Prefix() + key)
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
func (s *LocalStore) Lock(name, owner string, duration time.Duration) Lock {
	return &localLock{
		c:        s.c,
		name:     name,
		owner:    owner,
		duration: duration,
	}
}

// Exists checks if an entry exists in the cache for the given key
func (s *LocalStore) Exists(key string) (bool, error) {
	_, valid := s.c.Get(s.k(key))

	return valid, nil
}
