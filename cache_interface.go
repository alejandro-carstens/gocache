package cache

type CacheInterface interface {
	Get(key string) interface{}

	Put(key string, value interface{}, minutes int)

	Increment(key string, value int64) (int64, error)

	Decrement(key string, value int64) (int64, error)

	Forget(key string) (bool, error)

	Forever(key string, value interface{})

	Flush() (bool, error)

	GetInt(key string) (int64, error)

	GetFloat(key string) (float64, error)

	GetPrefix() string

	Many(keys []string) map[string]interface{}

	PutMany(values map[string]interface{}, minutes int)

	GetStruct(key string, entity interface{}) (interface{}, error)
}
