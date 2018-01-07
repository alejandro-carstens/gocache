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
	tagKey, err := this.taggedItemKey(key)

	if err != nil {
		return tagKey, err
	}

	return this.Store.Get(tagKey)
}

func (this *TaggedCache) Put(key string, value interface{}, minutes int) error {
	tagKey, err := this.taggedItemKey(key)

	if err != nil {
		return err
	}

	return this.Store.Put(tagKey, value, minutes)
}

func (this *TaggedCache) Increment(key string, value int64) (int64, error) {
	tagKey, err := this.taggedItemKey(key)

	if err != nil {
		return 0, err
	}

	return this.Store.Increment(tagKey, value)
}

func (this *TaggedCache) Decrement(key string, value int64) (int64, error) {
	tagKey, err := this.taggedItemKey(key)

	if err != nil {
		return 0, err
	}

	return this.Store.Decrement(tagKey, value)
}

func (this *TaggedCache) Forget(key string) (bool, error) {
	tagKey, err := this.taggedItemKey(key)

	if err != nil {
		return false, err
	}

	return this.Store.Forget(tagKey)
}

func (this *TaggedCache) Forever(key string, value interface{}) error {
	tagKey, err := this.taggedItemKey(key)

	if err != nil {
		return err
	}

	return this.Store.Forever(tagKey, value)
}

func (this *TaggedCache) Flush() (bool, error) {
	return this.Store.Flush()
}

func (this *TaggedCache) Many(keys []string) (map[string]interface{}, error) {
	taggedKeys := make([]string, len(keys))
	values := make(map[string]interface{})

	for i, key := range keys {
		tagKey, err := this.taggedItemKey(key)

		if err != nil {
			return values, err
		}

		taggedKeys[i] = tagKey
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

func (this *TaggedCache) PutMany(values map[string]interface{}, minutes int) error {
	taggedMap := make(map[string]interface{})

	for key, value := range values {
		tagKey, err := this.taggedItemKey(key)

		if err != nil {
			return err
		}

		taggedMap[tagKey] = value
	}

	return this.Store.PutMany(taggedMap, minutes)
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

func (this *TaggedCache) taggedItemKey(key string) (string, error) {
	h := sha1.New()

	namespace, err := this.Tags.GetNamespace()

	if err != nil {
		return namespace, err
	}

	h.Write(([]byte(namespace)))

	return this.GetPrefix() + hex.EncodeToString(h.Sum(nil)) + ":" + key, nil
}

func (this *TaggedCache) GetStruct(key string, entity interface{}) (interface{}, error) {
	tagKey, err := this.taggedItemKey(key)

	if err != nil {
		return tagKey, err
	}

	return this.Store.GetStruct(tagKey, entity)
}

func (this *TaggedCache) TagFlush() error {
	return this.Tags.Reset()
}
