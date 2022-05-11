package gocache

import (
	"errors"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/alejandro-carstens/gocache/encoder"
)

var _ Cache = &LocalStore{}

// NewLocalStore validates the passed in config and creates a Cache implementation of type *LocalStore
func NewLocalStore(cnf *LocalConfig, encoder encoder.Encoder) (*LocalStore, error) {
	if err := cnf.validate(); err != nil {
		return nil, nil
	}

	return &LocalStore{
		prefix: prefix{
			val: cnf.Prefix,
		},
		defaultExpiration: cnf.DefaultExpiration,
		defaultInterval:   cnf.DefaultInterval,
		c:                 cache.New(cnf.DefaultExpiration, cnf.DefaultInterval),
		encoder:           encoder,
	}, nil
}

// LocalStore is the representation of a map caching store
type LocalStore struct {
	prefix
	c                 *cache.Cache
	defaultExpiration time.Duration
	defaultInterval   time.Duration
	encoder           encoder.Encoder
}

// GetString gets a string value from the store
func (s *LocalStore) GetString(key string) (string, error) {
	value, valid := s.c.Get(s.k(key))
	if !valid {
		return "", ErrNotFound
	}
	if isNumeric(value) || isBool(value) {
		return fmt.Sprint(value), nil
	}

	data, valid := value.([]byte)
	if !valid {
		return "", errors.New("cannot decode cached value")
	}

	var v string
	if err := s.encoder.Decode(data, &v); err != nil {
		return "", err
	}

	return v, nil
}

// GetFloat64 gets a float value from the store
func (s *LocalStore) GetFloat64(key string) (float64, error) {
	value, valid := s.c.Get(s.k(key))
	if !valid {
		return 0, ErrNotFound
	}
	if !isNumeric(value) {
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
	if !isNumeric(value) {
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
	if !isNumeric(value) {
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
	if !isNumeric(value) {
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
	if !isNumeric(value) {
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
	if isNumeric(value) || isBool(value) {
		return stringToBool(fmt.Sprint(value)), nil
	}

	data, valid := value.([]byte)
	if !valid {
		return false, errors.New("cannot decode cached value")
	}

	var v string
	if err := s.encoder.Decode(data, &v); err != nil {
		return false, err
	}

	return stringToBool(v), nil
}

// Increment increments an integer counter by a given value
func (s *LocalStore) Increment(key string, value int64) (int64, error) {
	if _, valid := s.c.Get(s.k(key)); !valid {
		if err := s.Forever(key, value); err != nil {
			return 0, err
		}

		return value, nil
	}

	return s.c.IncrementInt64(s.k(key), value)
}

// Decrement decrements an integer counter by a given value
func (s *LocalStore) Decrement(key string, value int64) (int64, error) {
	if _, valid := s.c.Get(s.k(key)); !valid {
		if err := s.Forever(key, -1*value); err != nil {
			return 0, err
		}

		return value, nil
	}

	return s.c.DecrementInt64(s.k(key), value)
}

// Put puts a value in the given store for a predetermined amount of time in seconds.
func (s *LocalStore) Put(key string, value interface{}, duration time.Duration) error {
	if isNumeric(value) || isBool(value) {
		s.c.Set(s.k(key), value, duration)

		return nil
	}

	val, err := s.encoder.Encode(value)
	if err != nil {
		return err
	}

	s.c.Set(s.k(key), val, duration)

	return nil
}

// Add an item to the cache only if an item doesn't already exist for the given key, or if the existing item has
// expired. If the record was successfully added true will be returned else false will be returned
func (s *LocalStore) Add(key string, value interface{}, duration time.Duration) (bool, error) {
	if isNumeric(value) || isBool(value) {
		return s.c.Add(s.k(key), value, duration) == nil, nil
	}

	val, err := s.encoder.Encode(value)
	if err != nil {
		return false, err
	}

	return s.c.Add(s.k(key), val, duration) == nil, nil
}

// Forever puts a value in the given store until it is forgotten/evicted
func (s *LocalStore) Forever(key string, value interface{}) error {
	return s.Put(key, value, -1)
}

// Flush flushes the store
func (s *LocalStore) Flush() (bool, error) {
	s.c.Flush()

	return true, nil
}

func (s *LocalStore) Forget(key string) (bool, error) {
	var exists bool
	if _, exists = s.c.Get(s.k(key)); exists {
		s.c.Delete(s.k(key))
	}

	return exists, nil
}

// ForgetMany forgets/evicts a set of given key-value pair from the store
func (s *LocalStore) ForgetMany(keys ...string) error {
	for _, key := range keys {
		s.c.Delete(s.k(key))
	}

	return nil
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
		if isNumeric(val) || isBool(val) {
			items[key] = Item{
				key:     key,
				value:   fmt.Sprint(val),
				encoder: s.encoder,
			}

			continue
		}

		data, valid := val.([]byte)
		if !valid {
			return nil, errors.New("cannot decode cached value")
		}

		items[key] = Item{
			key:     key,
			value:   string(data),
			encoder: s.encoder,
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

	data, valid := value.([]byte)
	if !valid {
		return errors.New("cannot decode cached value")
	}

	return s.encoder.Decode(data, entity)
}

// Close closes the c releasing all open resources
func (*LocalStore) Close() error {
	return nil
}

// Tags returns the taggedCache for the given store
func (s *LocalStore) Tags(names ...string) TaggedCache {
	return &taggedCache{
		store: s,
		tags: &TagSet{
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
