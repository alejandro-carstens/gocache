package cache

import (
	"crypto/sha1"
	"encoding/hex"
)

type TaggedCache struct {
	Store StoreInterface
	Tags  TagSet
}

func (this *TaggedCache) Get(key string) interface{} {
	return this.Store.Get(this.taggedItemKey(key))
}

func (this *TaggedCache) Put(key string, value interface{}, minutes int) {
	this.Store.Put(this.taggedItemKey(key), value, minutes)
}

func (this *TaggedCache) Increment(key string, value int64) int64 {
	return this.Store.Increment(this.taggedItemKey(key), value)
}

func (this *TaggedCache) Decrement(key string, value int64) int64 {
	return this.Store.Decrement(this.taggedItemKey(key), value)
}

func (this *TaggedCache) Forget(key string) bool {
	return this.Store.Forget(this.taggedItemKey(key))
}

func (this *TaggedCache) Forever(key string, value interface{}) {
	this.Store.Forever(this.taggedItemKey(key), value)
}

func (this *TaggedCache) Flush() bool {
	return this.Store.Flush()
}

func (this *TaggedCache) Many(keys []string) map[string]interface{} {
	taggedKeys := make([]string, len(keys))
	values := make(map[string]interface{})

	for i, key := range keys {
		taggedKeys[i] = this.taggedItemKey(key)
	}

	results := this.Store.Many(taggedKeys)

	for i, result := range results {
		values[GetTaggedManyKey(this.GetPrefix(), i)] = result
	}

	return values
}

func (this *TaggedCache) PutMany(values map[string]interface{}, minutes int) {
	taggedMap := make(map[string]interface{})

	for key, value := range values {
		taggedMap[this.taggedItemKey(key)] = value
	}

	this.Store.PutMany(taggedMap, minutes)
}

func (this *TaggedCache) GetPrefix() string {
	return this.Store.GetPrefix()
}

func (this *TaggedCache) taggedItemKey(key string) string {
	h := sha1.New()

	h.Write(([]byte(this.Tags.GetNamespace())))

	return this.GetPrefix() + hex.EncodeToString(h.Sum(nil)) + ":" + key
}
