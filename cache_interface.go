package cache

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
