package cache

// Interface to which all the cache conenctors should comply with
type CacheConnectorInterface interface {
	Connect(params map[string]interface{}) (StoreInterface, error)

	validate(params map[string]interface{}) (map[string]interface{}, error)
}

// Interface to which all the caching stores should comply with
type CacheInterface interface {
	Get(key string) (interface{}, error)

	Put(key string, value interface{}, minutes int) error

	Increment(key string, value int64) (int64, error)

	Decrement(key string, value int64) (int64, error)

	Forget(key string) (bool, error)

	Forever(key string, value interface{}) error

	Flush() (bool, error)

	GetInt(key string) (int64, error)

	GetFloat(key string) (float64, error)

	GetPrefix() string

	Many(keys []string) (map[string]interface{}, error)

	PutMany(values map[string]interface{}, minutes int) error

	GetStruct(key string, entity interface{}) (interface{}, error)
}

// Interface to which all the tagged stores should comply with
type TagsInterface interface {
	Tags(names []string) TaggedStoreInterface
}

// Interface to which all the stores implementint caching and tagging should comply with
type StoreInterface interface {
	CacheInterface

	TagsInterface
}

// Interface to which all the stores that implement tagged caching should comply with
type TaggedStoreInterface interface {
	CacheInterface

	TagFlush() error
}
