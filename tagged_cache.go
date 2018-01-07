package cache

import (
	"crypto/sha1"
	"encoding/hex"
)

type TaggedCache struct {
	Store StoreInterface
	Tags  TagSet
}

func (this *TaggedCache) Get(key string) (interface{}, error) {
	return this.Store.Get(this.taggedItemKey(key))
}

func (this *TaggedCache) Put(key string, value interface{}, minutes int) error {
	return this.Store.Put(this.taggedItemKey(key), value, minutes)
}

func (this *TaggedCache) Increment(key string, value int64) (int64, error) {
	return this.Store.Increment(this.taggedItemKey(key), value)
}

func (this *TaggedCache) Decrement(key string, value int64) (int64, error) {
	return this.Store.Decrement(this.taggedItemKey(key), value)
}

func (this *TaggedCache) Forget(key string) (bool, error) {
	return this.Store.Forget(this.taggedItemKey(key))
}

func (this *TaggedCache) Forever(key string, value interface{}) error {
	return this.Store.Forever(this.taggedItemKey(key), value)
}

func (this *TaggedCache) Flush() (bool, error) {
	return this.Store.Flush()
}

func (this *TaggedCache) Many(keys []string) (map[string]interface{}, error) {
	taggedKeys := make([]string, len(keys))
	values := make(map[string]interface{})

	for i, key := range keys {
		taggedKeys[i] = this.taggedItemKey(key)
	}

	results, err := this.Store.Many(taggedKeys)

	if err != nil {
		return results, err
	}

	for i, result := range results {
		values[GetTaggedManyKey(this.GetPrefix(), i)] = result
	}

	return values, nil
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

func (this *TaggedCache) GetInt(key string) (int64, error) {
	return this.Store.GetInt(key)
}

func (this *TaggedCache) GetFloat(key string) (float64, error) {
	return this.Store.GetFloat(key)
}

func (this *TaggedCache) taggedItemKey(key string) string {
	h := sha1.New()

	h.Write(([]byte(this.Tags.GetNamespace())))

	return this.GetPrefix() + hex.EncodeToString(h.Sum(nil)) + ":" + key
}

func (this *TaggedCache) GetStruct(key string, entity interface{}) (interface{}, error) {
	return this.Store.GetStruct(this.taggedItemKey(key), entity)
}

func (this *TaggedCache) TagFlush() {
	this.Tags.Reset()
}
