package gocache

import (
	"errors"
	"time"

	"github.com/alejandro-carstens/gocache/encoder"
)

// New new-ups an instance of Store
func New(config config, encoder encoder.Encoder) (Cache, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}
	switch config.(type) {
	case *LocalConfig:
		return NewLocalStore(config.(*LocalConfig), encoder)
	case *RedisConfig:
		return NewRedisStore(config.(*RedisConfig), encoder)
	case *MemcacheConfig:
		return NewMemcacheStore(config.(*MemcacheConfig), encoder)
	}

	return nil, errors.New("invalid or empty config specified")
}

type (
	// store represents the caching methods to be implemented
	store interface {
		// GetString gets a string value from the store
		GetString(key string) (string, error)
		// Put puts a value in the given store for a predetermined amount of time in seconds
		Put(key string, value interface{}, duration time.Duration) error
		// Add an item to the cache only if an item doesn't already exist for the given key, or if the existing item has
		// expired. If the record was successfully added true will be returned else false will be returned
		Add(key string, value interface{}, duration time.Duration) (bool, error)
		// Increment increments an integer counter by a given value
		Increment(key string, value int64) (int64, error)
		// Decrement decrements an integer counter by a given value
		Decrement(key string, value int64) (int64, error)
		// Forget forgets/evicts a given key-value pair from the store
		Forget(key string) (bool, error)
		// ForgetMany forgets/evicts a set of given key-value pair from the store
		ForgetMany(keys ...string) error
		// Forever puts a value in the given store until it is forgotten/evicted manually
		Forever(key string, value interface{}) error
		// Flush flushes the store
		Flush() (bool, error)
		// GetInt64 gets an int64 value from the store
		GetInt64(key string) (int64, error)
		// GetInt gets an int value from the store
		GetInt(key string) (int, error)
		// GetFloat64 gets a float64 value from the store
		GetFloat64(key string) (float64, error)
		// GetFloat32 gets a float32 value from the store
		GetFloat32(key string) (float32, error)
		// GetUint64 gets a uint64 value from the store
		GetUint64(key string) (uint64, error)
		// GetBool gets a bool value from the store
		GetBool(key string) (bool, error)
		// Prefix gets the cache key prefix
		Prefix() string
		// Many gets many values from the store
		Many(keys ...string) (Items, error)
		// PutMany puts many values in the given store until they are forgotten/evicted
		PutMany(entries ...Entry) error
		// Get gets the struct representation of a value from the store
		Get(key string, entity interface{}) error
		// Close closes the c releasing all open resources
		Close() error
		// Exists checks if an entry exists in the cache for the given key
		Exists(key string) (bool, error)
		// Expire allows for overriding the expiry time for a given key
		Expire(key string, duration time.Duration) error
	}
	// tags represents the tagging methods to be implemented
	tags interface {
		// Tags returns the TaggedCache for the given store
		Tags(names ...string) TaggedCache
	}
	// Cache represents the methods a caching store needs to implement
	Cache interface {
		store
		tags
		// Lock returns an implementation of the Lock interface
		Lock(name, owner string, duration time.Duration) Lock
	}
	// TaggedCache represents the methods a tagged-caching store needs to implement
	TaggedCache interface {
		store
		// TagSet returns the underlying tagged cache tag set
		TagSet() *TagSet
	}
	lock interface {
		// Acquire is responsible for acquiring a lock
		Acquire() (bool, error)
		// ForceRelease forces a cache lock release
		ForceRelease() error
		// GetCurrentOwner retrieves the current
		// owner of a given cache lock
		GetCurrentOwner() (string, error)
		// Release frees up a lock for use by
		// a different concurrent process
		Release() (bool, error)
		// Expire allows to set a new expiration time on the lock. It will
		// return true if the operation was successful, false otherwise
		Expire(time.Duration) (bool, error)
	}
	// Lock represents the methods to be implemented by a cache lock
	Lock interface {
		lock
		// Get attempts to acquire a lock. If acquired, fn will be invoked and the lock will be safely release once
		// the invocation either succeeds or errors
		Get(fn func() error) (acquired bool, err error)
		// Block will attempt to acquire a lock for the specified "wait" time. If acquired, fn will be invoked and
		// the lock will be safely release once the invocation either succeeds or errors. The interval variable
		// will be used as the wait duration between attempts to acquire the lock
		Block(interval, wait time.Duration, fn func() error) (acquired bool, err error)
	}
)
