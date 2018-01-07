package cache

type TaggedStoreInterface interface {
	CacheInterface

	TagFlush() error
}
