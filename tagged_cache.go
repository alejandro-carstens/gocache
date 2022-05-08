package gocache

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"time"
)

var _ TaggedCache = &taggedCache{}

// taggedCache is the representation of a tagged caching store
type taggedCache struct {
	store store
	tags  tagSet
}

// Put puts a val in the given store for a predetermined amount of time in seconds
func (tc *taggedCache) Put(key string, value interface{}, duration time.Duration) error {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return err
	}

	return tc.store.Put(tagKey, value, duration)
}

// Add an item to the cache only if an item doesn't already exist for the given key, or if the existing item has
// expired. If the record was successfully added true will be returned else false will be returned
func (tc *taggedCache) Add(key string, value interface{}, duration time.Duration) (bool, error) {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return false, err
	}

	return tc.store.Add(tagKey, value, duration)
}

// Increment increments an integer counter by a given val
func (tc *taggedCache) Increment(key string, value int64) (int64, error) {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return 0, err
	}

	return tc.store.Increment(tagKey, value)
}

// Decrement decrements an integer counter by a given val
func (tc *taggedCache) Decrement(key string, value int64) (int64, error) {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return 0, err
	}

	return tc.store.Decrement(tagKey, value)
}

// Forget forgets/evicts a given key-val pair from the store
func (tc *taggedCache) Forget(keys ...string) (bool, error) {
	var tagKeys = make([]string, len(keys))
	for i, key := range keys {
		tagKey, err := tc.taggedItemKey(key)
		if err != nil {
			return false, err
		}

		tagKeys[i] = tagKey
	}

	return tc.store.Forget(tagKeys...)
}

// Forever puts a val in the given store until it is forgotten/evicted
func (tc *taggedCache) Forever(key string, value interface{}) error {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return err
	}

	return tc.store.Forever(tagKey, value)
}

// Many gets many values from the store
func (tc *taggedCache) Many(keys ...string) (Items, error) {
	var (
		taggedKeys = make([]string, len(keys))
		tagKeyMap  = map[string]string{}
	)
	for i, key := range keys {
		tagKey, err := tc.taggedItemKey(key)
		if err != nil {
			return nil, err
		}

		taggedKeys[i] = tagKey
		tagKeyMap[tagKey] = key
	}

	results, err := tc.store.Many(taggedKeys...)
	if err != nil {
		return nil, err
	}

	items := Items{}
	for _, result := range results {
		key, valid := tagKeyMap[result.Key()]
		if !valid {
			return nil, errors.New("tag key not found")
		}

		result.tagKey = result.Key()
		result.key = key
		items[result.Key()] = result
	}

	return items, nil
}

// PutMany puts many values in the given store until they are forgotten/evicted
func (tc *taggedCache) PutMany(entries ...Entry) error {
	for i, entry := range entries {
		key, err := tc.taggedItemKey(entry.Key)
		if err != nil {
			return err
		}

		entry.Key = key
		entries[i] = entry
	}

	return tc.store.PutMany(entries...)
}

// Prefix gets the cache key val
func (tc *taggedCache) Prefix() string {
	return tc.store.Prefix()
}

// GetInt64 gets an int64 val from the store
func (tc *taggedCache) GetInt64(key string) (int64, error) {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return 0, err
	}

	return tc.store.GetInt64(tagKey)
}

// GetInt gets an int val from the store
func (tc *taggedCache) GetInt(key string) (int, error) {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return 0, err
	}

	return tc.store.GetInt(tagKey)
}

// GetUint64 gets an uint64 val from the store
func (tc *taggedCache) GetUint64(key string) (uint64, error) {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return 0, err
	}

	return tc.store.GetUint64(tagKey)
}

// GetBool gets a bool val from the store
func (tc *taggedCache) GetBool(key string) (bool, error) {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return false, err
	}

	return tc.store.GetBool(tagKey)
}

// GetFloat64 gets a float val from the store
func (tc *taggedCache) GetFloat64(key string) (float64, error) {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return 0, err
	}

	return tc.store.GetFloat64(tagKey)
}

// GetFloat32 gets an int val from the store
func (tc *taggedCache) GetFloat32(key string) (float32, error) {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return 0, err
	}

	return tc.store.GetFloat32(tagKey)
}

// Get gets the struct representation of a val from the store
func (tc *taggedCache) Get(key string, entity interface{}) error {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return err
	}

	return tc.store.Get(tagKey, entity)
}

func (tc *taggedCache) Close() error {
	return tc.store.Close()
}

func (tc *taggedCache) GetString(key string) (string, error) {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return "", err
	}

	return tc.store.GetString(tagKey)
}

// Flush flushes all the given tags' associated records
func (tc *taggedCache) Flush() (bool, error) {
	if err := tc.tags.reset(); err != nil {
		return false, err
	}

	return true, nil
}

// GetTags returns the taggedCache Tags
func (tc *taggedCache) GetTags() tagSet {
	return tc.tags
}

// Exists checks if an entry exists in the cache for the given key
func (tc *taggedCache) Exists(key string) (bool, error) {
	tagKey, err := tc.taggedItemKey(key)
	if err != nil {
		return false, err
	}

	return tc.store.Exists(tagKey)
}

func (tc *taggedCache) taggedItemKey(key string) (string, error) {
	namespace, err := tc.tags.getNamespace()
	if err != nil {
		return namespace, err
	}

	h := sha1.New()
	h.Write([]byte(namespace))

	return tc.Prefix() + hex.EncodeToString(h.Sum(nil)) + ":" + key, nil
}
