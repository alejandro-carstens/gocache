package cache

import (
	"crypto/sha1"
	"encoding/hex"
)

// TaggedCache is the representation of a tagged caching store
type TaggedCache struct {
	Store StoreInterface
	Tags  TagSet
}

func (tc *TaggedCache) Get(key string) (interface{}, error) {
	tagKey, err := tc.taggedItemKey(key)

	if err != nil {
		return tagKey, err
	}

	return tc.Store.Get(tagKey)
}

func (tc *TaggedCache) Put(key string, value interface{}, minutes int) error {
	tagKey, err := tc.taggedItemKey(key)

	if err != nil {
		return err
	}

	return tc.Store.Put(tagKey, value, minutes)
}

func (tc *TaggedCache) Increment(key string, value int64) (int64, error) {
	tagKey, err := tc.taggedItemKey(key)

	if err != nil {
		return 0, err
	}

	return tc.Store.Increment(tagKey, value)
}

func (tc *TaggedCache) Decrement(key string, value int64) (int64, error) {
	tagKey, err := tc.taggedItemKey(key)

	if err != nil {
		return 0, err
	}

	return tc.Store.Decrement(tagKey, value)
}

func (tc *TaggedCache) Forget(key string) (bool, error) {
	tagKey, err := tc.taggedItemKey(key)

	if err != nil {
		return false, err
	}

	return tc.Store.Forget(tagKey)
}

func (tc *TaggedCache) Forever(key string, value interface{}) error {
	tagKey, err := tc.taggedItemKey(key)

	if err != nil {
		return err
	}

	return tc.Store.Forever(tagKey, value)
}

func (tc *TaggedCache) Flush() (bool, error) {
	return tc.Store.Flush()
}

func (tc *TaggedCache) Many(keys []string) (map[string]interface{}, error) {
	taggedKeys := make([]string, len(keys))
	values := make(map[string]interface{})

	for i, key := range keys {
		tagKey, err := tc.taggedItemKey(key)

		if err != nil {
			return values, err
		}

		taggedKeys[i] = tagKey
	}

	results, err := tc.Store.Many(taggedKeys)

	if err != nil {
		return results, err
	}

	for i, result := range results {
		values[GetTaggedManyKey(tc.GetPrefix(), i)] = result
	}

	return values, nil
}

func (tc *TaggedCache) PutMany(values map[string]interface{}, minutes int) error {
	taggedMap := make(map[string]interface{})

	for key, value := range values {
		tagKey, err := tc.taggedItemKey(key)

		if err != nil {
			return err
		}

		taggedMap[tagKey] = value
	}

	return tc.Store.PutMany(taggedMap, minutes)
}

func (tc *TaggedCache) GetPrefix() string {
	return tc.Store.GetPrefix()
}

func (tc *TaggedCache) GetInt(key string) (int64, error) {
	return tc.Store.GetInt(key)
}

func (tc *TaggedCache) GetFloat(key string) (float64, error) {
	return tc.Store.GetFloat(key)
}

func (tc *TaggedCache) taggedItemKey(key string) (string, error) {
	h := sha1.New()

	namespace, err := tc.Tags.GetNamespace()

	if err != nil {
		return namespace, err
	}

	h.Write(([]byte(namespace)))

	return tc.GetPrefix() + hex.EncodeToString(h.Sum(nil)) + ":" + key, nil
}

func (tc *TaggedCache) GetStruct(key string, entity interface{}) (interface{}, error) {
	tagKey, err := tc.taggedItemKey(key)

	if err != nil {
		return tagKey, err
	}

	return tc.Store.GetStruct(tagKey, entity)
}

func (tc *TaggedCache) TagFlush() error {
	return tc.Tags.Reset()
}
