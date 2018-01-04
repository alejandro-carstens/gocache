package cache

type TagsInterface interface {
	Tags(names []string) TaggedStoreInterface
}
