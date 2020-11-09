package gocache

// CacheConnector represents the connector methods to be implemented
type cacheConnector interface {
	// Connect is responsible for connecting with the caching store
	connect(params map[string]interface{}) (Cache, error)
	// validate verifies that the given params
	// are valid for establishing a connection
	validate(params map[string]interface{}) (map[string]interface{}, error)
}

// store represents the caching methods to be implemented
type store interface {
	// GetString gets a string value from the store
	GetString(key string) (string, error)
	// Put puts a value in the given store for a predetermined amount of time in mins.
	Put(key string, value interface{}, minutes int) error
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
	PutMany(values map[string]string, minutes int) error
	// Get gets the struct representation of a value from the store
	Get(key string, entity interface{}) error
	// Close closes the client releasing all open resources
	Close() error
}

// tags represents the tagging methods to be implemented
type tags interface {
	// Tags returns the TaggedCache for the given store
	Tags(names ...string) TaggedCache
}

// Store represents the methods a caching store needs to implement
type Cache interface {
	store
	tags
}

// TaggedCache represents the methods a tagged-caching store needs to implement
type TaggedCache interface {
	store
	// TagFlush flushes the tags of the TaggedCache
	TagFlush() error
	// GetTags returns the TaggedCache Tags
	GetTags() tagSet
}

type Lock interface {
	Acquire() (bool, error)
	ForceRelease() error
	GetCurrentOwner() (string, error)
	Release() (bool, error)
}