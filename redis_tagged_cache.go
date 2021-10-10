package gocache

import (
	"crypto/sha1"
	"encoding/hex"
	"reflect"
	"strings"
)

// redisTaggedCache is the representation of the redis tagged cache store
type redisTaggedCache struct {
	taggedCache
}

// Forever puts a value in the given store until it is forgotten/evicted
func (rtc *redisTaggedCache) Forever(key string, value interface{}) error {
	namespace, err := rtc.tags.getNamespace()
	if err != nil {
		return err
	}

	rtc.pushForever(namespace, key)

	h := sha1.New()
	h.Write([]byte(namespace))

	return rtc.store.Forever(rtc.GetPrefix()+hex.EncodeToString(h.Sum(nil))+":"+key, value)
}

// TagFlush flushes the tags of the TaggedCache
func (rtc *redisTaggedCache) TagFlush() error {
	return rtc.deleteForeverKeys()
}

func (rtc *redisTaggedCache) pushForever(namespace string, key string) {
	h := sha1.New()
	h.Write([]byte(namespace))

	fullKey := rtc.GetPrefix() + hex.EncodeToString(h.Sum(nil)) + ":" + key
	for _, segment := range strings.Split(namespace, "|") {
		inputs := []reflect.Value{
			reflect.ValueOf(rtc.foreverKey(segment)),
			reflect.ValueOf(fullKey),
		}

		reflect.ValueOf(rtc.store).MethodByName("Lpush").Call(inputs)
	}
}

func (rtc *redisTaggedCache) deleteForeverKeys() error {
	namespace, err := rtc.tags.getNamespace()
	if err != nil {
		return err
	}

	for _, segment := range strings.Split(namespace, "|") {
		key := rtc.foreverKey(segment)
		if err = rtc.deleteForeverValues(key); err != nil {
			return err
		}
		if _, err = rtc.store.Forget(segment); err != nil {
			return err
		}
	}

	return nil
}

func (rtc *redisTaggedCache) deleteForeverValues(key string) error {
	var (
		inputs = []reflect.Value{reflect.ValueOf(key), reflect.ValueOf(int64(0)), reflect.ValueOf(int64(-1))}
		keys   = reflect.ValueOf(rtc.store).MethodByName("Lrange").Call(inputs)
	)
	if len(keys) == 0 {
		return nil
	}

	for _, k := range keys {
		if k.Len() == 0 {
			continue
		}

		for i := 0; i < k.Len(); i++ {
			if _, err := rtc.store.Forget(k.Index(i).String()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (rtc *redisTaggedCache) foreverKey(segment string) string {
	return rtc.GetPrefix() + segment + ":forever"
}
