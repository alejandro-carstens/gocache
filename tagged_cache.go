package gocache

import (
	"crypto/sha1"
	"encoding/hex"
)

// TaggedCache is the representation of a tagged caching store
type TaggedCache struct {
	store Store
	tags  TagSet
}

// Put puts a value in the given store for a predetermined amount of time in mins.
func (tc *TaggedCache) Put(key string, value interface{}, minutes int) error {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return err
	}

	return tc.store.Put(tagKey, value, minutes)
}

// Increment increments an integer counter by a given value
func (tc *TaggedCache) Increment(key string, value int64) (int64, error) {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return 0, err
	}

	return tc.store.Increment(tagKey, value)
}

// Decrement decrements an integer counter by a given value
func (tc *TaggedCache) Decrement(key string, value int64) (int64, error) {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return 0, err
	}

	return tc.store.Decrement(tagKey, value)
}

// Forget forgets/evicts a given key-value pair from the store
func (tc *TaggedCache) Forget(key string) (bool, error) {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return false, err
	}

	return tc.store.Forget(tagKey)
}

// Forever puts a value in the given store until it is forgotten/evicted
func (tc *TaggedCache) Forever(key string, value interface{}) error {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return err
	}

	return tc.store.Forever(tagKey, value)
}

// Flush flushes the store
func (tc *TaggedCache) Flush() (bool, error) {
	return tc.store.Flush()
}

// Many gets many values from the store
func (tc *TaggedCache) Many(keys []string) (map[string]string, error) {
	taggedKeys := make([]string, len(keys))
	values := make(map[string]string)

	for i, key := range keys {
		tagKey, err := tc.taggedItemKey(key)
		if err != nil {
			return values, err
		}

		taggedKeys[i] = tagKey
	}

	results, err := tc.store.Many(taggedKeys)
	if err != nil {
		return results, err
	}

	for i, result := range results {
		values[getTaggedManyKey(tc.GetPrefix(), i)] = result
	}

	return values, nil
}

// PutMany puts many values in the given store until they are forgotten/evicted
func (tc *TaggedCache) PutMany(values map[string]string, minutes int) error {
	taggedMap := make(map[string]string)

	for key, value := range values {
		tagKey, err := tc.taggedItemKey(key)
		if err != nil {
			return err
		}

		taggedMap[tagKey] = value
	}

	return tc.store.PutMany(taggedMap, minutes)
}

// GetPrefix gets the cache key prefix
func (tc *TaggedCache) GetPrefix() string {
	return tc.store.GetPrefix()
}

// GetInt64 gets an int value from the store
func (tc *TaggedCache) GetInt64(key string) (int64, error) {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return 0, err
	}

	return tc.store.GetInt64(tagKey)
}

// GetFloat64 gets a float value from the store
func (tc *TaggedCache) GetFloat64(key string) (float64, error) {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return 0, err
	}

	return tc.store.GetFloat64(tagKey)
}

// Get gets the struct representation of a value from the store
func (tc *TaggedCache) Get(key string, entity interface{}) error {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return err
	}

	return tc.store.Get(tagKey, entity)
}

func (tc *TaggedCache) Close() error {
	return tc.store.Close()
}

func (tc *TaggedCache) GetString(key string) (string, error) {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return "", err
	}

	return tc.store.GetString(tagKey)
}

// TagFlush flushes the tags of the TaggedCache
func (tc *TaggedCache) TagFlush() error {
	return tc.tags.reset()
}

func (tc *TaggedCache) taggedItemKey(key string) (string, error) {
	h := sha1.New()

	namespace, err := tc.tags.getNamespace()
	if err != nil {
		return namespace, err
	}

	h.Write([]byte(namespace))

	return tc.GetPrefix() + hex.EncodeToString(h.Sum(nil)) + ":" + key, nil
}

// GetTags returns the TaggedCache Tags
func (tc *TaggedCache) GetTags() TagSet {
	return tc.tags
}
