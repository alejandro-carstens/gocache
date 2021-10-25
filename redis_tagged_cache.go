package gocache

import (
	"crypto/sha1"
	"encoding/hex"
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

	tc.pushForever(namespace, key)

	h := sha1.New()
	h.Write([]byte(namespace))

	return tc.store.Forever(tc.Prefix()+hex.EncodeToString(h.Sum(nil))+":"+key, value)
}

// TagFlush flushes the tags of the TaggedCache
func (tc *redisTaggedCache) TagFlush() error {
	return tc.deleteForeverKeys()
}

func (tc *redisTaggedCache) pushForever(namespace string, key string) {
	h := sha1.New()
	h.Write([]byte(namespace))

	fullKey := tc.Prefix() + hex.EncodeToString(h.Sum(nil)) + ":" + key
	for _, segment := range strings.Split(namespace, "|") {
		inputs := []reflect.Value{
			reflect.ValueOf(tc.foreverKey(segment)),
			reflect.ValueOf(fullKey),
		}

		reflect.ValueOf(tc.store).MethodByName("Lpush").Call(inputs)
	}
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

		for i := 0; i < k.Len(); i++ {
			if _, err := tc.store.Forget(k.Index(i).String()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (tc *redisTaggedCache) foreverKey(segment string) string {
	return tc.Prefix() + segment + ":forever"
}
