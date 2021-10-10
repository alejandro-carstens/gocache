package gocache

import "errors"

// New new-ups an instance of Store
func New(config *Config) (Cache, error) {
	var connector cacheConnector
	if config.Local != nil {
		connector = new(localConnector)
	}
	if config.Memcache != nil {
		connector = new(memcacheConnector)
	}
	if config.Redis != nil {
		connector = new(redisConnector)
	}
	if connector == nil {
		return nil, errors.New("invalid or empty config specified")
	}
	if err := connector.validate(config); err != nil {
		return nil, err
	}

	return connector.connect(config)
}

type (
	// cacheConnector represents the connector methods to be implemented
	cacheConnector interface {
		// Connect is responsible for connecting with the caching store
		connect(config *Config) (Cache, error)
		// validate verifies that the given params
		// are valid for establishing a connection
		validate(config *Config) error
	}
	// store represents the caching methods to be implemented
	store interface {
		// GetString gets a string value from the store
		GetString(key string) (string, error)
		// Put puts a value in the given store for a predetermined amount of time in seconds
		Put(key string, value interface{}, seconds int) error
		// Increment increments an integer counter by a given value
		Increment(key string, value int64) (int64, error)
		// Decrement decrements an integer counter by a given value
		Decrement(key string, value int64) (int64, error)
		// Forget forgets/evicts a given key-value pair from the store
		Forget(key string) (bool, error)
		// Forever puts a value in the given store until it is forgotten/evicted
		Forever(key string, value interface{}) error
		// Flush flushes the store
		Flush() (bool, error)
		// GetInt64 gets an int value from the store
		GetInt64(key string) (int64, error)
		// GetFloat64 gets a float value from the store
		GetFloat64(key string) (float64, error)
		// GetPrefix gets the cache key prefix
		GetPrefix() string
		// Many gets many values from the store
		Many(keys []string) (map[string]string, error)
		// PutMany puts many values in the given store until they are forgotten/evicted
		PutMany(values map[string]string, seconds int) error
		// Get gets the struct representation of a value from the store
		Get(key string, entity interface{}) error
		// Close closes the c releasing all open resources
		Close() error
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
		Lock(name, owner string, seconds int64) Lock
	}
	// TaggedCache represents the methods a tagged-caching store needs to implement
	TaggedCache interface {
		store
		// TagFlush flushes the tags of the TaggedCache
		TagFlush() error
		// GetTags returns the TaggedCache Tags
		GetTags() tagSet
	}
	// Lock represents the methods to be implemented by a cache lock
	Lock interface {
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
	}
)
