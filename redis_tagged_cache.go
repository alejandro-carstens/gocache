package gocache

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"reflect"
	"strings"
	"time"
)

const (
	referenceKeyForever  = ":forever"
	referenceKeyStandard = ":standard"
)

var _ TaggedCache = &redisTaggedCache{}

// redisTaggedCache is the representation of the redis tagged cache store
type redisTaggedCache struct {
	taggedCache
}

// Put implementation of the TaggedCache interface
func (tc *redisTaggedCache) Put(key string, value interface{}, duration time.Duration) error {
	if duration == 0 {
		return tc.Forever(key, value)
	}
	if err := tc.pushKeys(key, referenceKeyStandard); err != nil {
		return err
	}

	return tc.taggedCache.Put(key, value, duration)
}

// PutMany implementation of the TaggedCache interface
func (tc *redisTaggedCache) PutMany(entries ...Entry) error {
	for i, entry := range entries {
		reference := referenceKeyStandard
		if entry.Duration == 0 {
			reference = referenceKeyForever
		}
		if err := tc.pushKeys(entry.Key, reference); err != nil {
			return err
		}

		key, err := tc.taggedItemKey(entry.Key)
		if err != nil {
			return err
		}

		entry.Key = key
		entries[i] = entry
	}

	return tc.store.PutMany(entries...)
}

// Forever implementation of the TaggedCache interface
func (tc *redisTaggedCache) Forever(key string, value interface{}) error {
	if err := tc.pushKeys(key, referenceKeyForever); err != nil {
		return err
	}

	return tc.taggedCache.Forever(key, value)
}

// Increment implementation of the TaggedCache interface
func (tc *redisTaggedCache) Increment(key string, value int64) (int64, error) {
	if err := tc.pushKeys(key, referenceKeyForever); err != nil {
		return 0, err
	}

	return tc.taggedCache.Increment(key, value)
}

// Decrement implementation of the TaggedCache interface
func (tc *redisTaggedCache) Decrement(key string, value int64) (int64, error) {
	if err := tc.pushKeys(key, referenceKeyForever); err != nil {
		return 0, err
	}

	return tc.taggedCache.Decrement(key, value)
}

// Flush flushes all the given tags' associated records. Note that for Redis all forever keys associated with
// the tags will also be deleted. Standard or expiring keys will be left alone until they expire
func (tc *redisTaggedCache) Flush() (bool, error) {
	if err := tc.deleteForeverKeys(); err != nil {
		return false, err
	}
	if err := tc.deleteStandardKeys(); err != nil {
		return false, err
	}

	return tc.taggedCache.Flush()
}

func (tc *redisTaggedCache) pushKeys(key, reference string) error {
	namespace, err := tc.tags.getNamespace()
	if err != nil {
		return err
	}

	h := sha1.New()
	h.Write([]byte(namespace))

	fullKey := tc.Prefix() + hex.EncodeToString(h.Sum(nil)) + ":" + key
	for _, segment := range strings.Split(namespace, "|") {
		var (
			inputs = []reflect.Value{reflect.ValueOf(tc.referenceKey(segment, reference)), reflect.ValueOf(fullKey)}
			res    = reflect.ValueOf(tc.store).MethodByName("Lpush").Call(inputs)
		)
		for _, r := range res {
			if !r.IsNil() {
				return errors.New(r.String())
			}
		}
	}

	return nil
}

func (tc *redisTaggedCache) deleteStandardKeys() error {
	return tc.deleteKeysByReference(referenceKeyStandard)
}

func (tc *redisTaggedCache) deleteForeverKeys() error {
	return tc.deleteKeysByReference(referenceKeyForever)
}

func (tc *redisTaggedCache) deleteKeysByReference(reference string) error {
	namespace, err := tc.tags.getNamespace()
	if err != nil {
		return err
	}

	for _, segment := range strings.Split(namespace, "|") {
		key := tc.referenceKey(segment, reference)
		if err = tc.deleteValues(key); err != nil {
			return err
		}
		if _, err = tc.store.Forget(segment); err != nil {
			return err
		}
	}

	return nil
}

func (tc *redisTaggedCache) deleteValues(key string) error {
	var (
		inputs = []reflect.Value{reflect.ValueOf(key), reflect.ValueOf(int64(0)), reflect.ValueOf(int64(-1))}
		keys   = reflect.ValueOf(tc.store).MethodByName("Lrange").Call(inputs)
	)
	if len(keys) == 0 {
		return nil
	}

	for _, k := range keys {
		if k.Len() == 0 {
			continue
		}

		delKeys := make([]string, k.Len())
		for i := 0; i < k.Len(); i++ {
			delKeys[i] = k.Index(i).String()
		}

		if len(delKeys) == 0 {
			continue
		}
		if _, err := tc.store.Forget(delKeys...); err != nil {
			return err
		}
	}

	return nil
}

func (tc *redisTaggedCache) referenceKey(segment, suffix string) string {
	return tc.Prefix() + segment + suffix
}
