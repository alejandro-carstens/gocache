package gocache

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"reflect"
	"strings"
)

var _ TaggedCache = &redisTaggedCache{}

// redisTaggedCache is the representation of the redis tagged cache store
type redisTaggedCache struct {
	taggedCache
}

// Forever puts a value in the given store until it is forgotten/evicted
func (tc *redisTaggedCache) Forever(key string, value interface{}) error {
	namespace, err := tc.tags.getNamespace()
	if err != nil {
		return err
	}
	if err = tc.pushForever(namespace, key); err != nil {
		return err
	}

	h := sha1.New()
	h.Write([]byte(namespace))

	return tc.store.Forever(tc.Prefix()+hex.EncodeToString(h.Sum(nil))+":"+key, value)
}

// Flush flushes all the given tags' associated records. Note that for Redis all forever keys associated with
// the tags will also be deleted. Standard or expiring keys will be left alone until they expire
func (tc *redisTaggedCache) Flush() (bool, error) {
	if err := tc.deleteForeverKeys(); err != nil {
		return false, err
	}

	return tc.taggedCache.Flush()
}

func (tc *redisTaggedCache) pushForever(namespace, key string) error {
	h := sha1.New()
	h.Write([]byte(namespace))

	fullKey := tc.Prefix() + hex.EncodeToString(h.Sum(nil)) + ":" + key
	for _, segment := range strings.Split(namespace, "|") {
		inputs := []reflect.Value{reflect.ValueOf(tc.foreverKey(segment)), reflect.ValueOf(fullKey)}

		res := reflect.ValueOf(tc.store).MethodByName("Lpush").Call(inputs)
		for _, r := range res {
			if !r.IsNil() {
				return errors.New(r.String())
			}
		}
	}

	return nil
}

func (tc *redisTaggedCache) deleteForeverKeys() error {
	namespace, err := tc.tags.getNamespace()
	if err != nil {
		return err
	}

	for _, segment := range strings.Split(namespace, "|") {
		key := tc.foreverKey(segment)
		if err = tc.deleteForeverValues(key); err != nil {
			return err
		}
		if _, err = tc.store.Forget(segment); err != nil {
			return err
		}
	}

	return nil
}

func (tc *redisTaggedCache) deleteForeverValues(key string) error {
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

func (tc *redisTaggedCache) foreverKey(segment string) string {
	return tc.Prefix() + segment + ":forever"
}
