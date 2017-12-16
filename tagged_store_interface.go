package cache

type TaggedStoreInterface interface {
	Get(key string) interface{}

	Put(key string, value interface{}, minutes int)

	Increment(key string, value int64) int64

	Decrement(key string, value int64) int64

	Forget(key string) bool

	Flush() bool

	GetPrefix() string

	Many(keys []string) map[string]interface{}

	PutMany(values map[string]interface{}, minutes int)

	GetStruct(key string, entity interface{}) (interface{}, error)

	Forever(key string, value interface{})

	TagFlush()
}
