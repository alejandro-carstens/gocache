package cache

import (
	"crypto/sha1"
	"encoding/hex"
	"reflect"
	"strings"
)

type RedisTaggedCache struct {
	TaggedCache
}

func (this *RedisTaggedCache) Forever(key string, value interface{}) {
	namespace := this.Tags.GetNamespace()

	this.pushForever(namespace, key)

	h := sha1.New()

	h.Write(([]byte(namespace)))

	this.Store.Forever(this.GetPrefix()+hex.EncodeToString(h.Sum(nil))+":"+key, value)
}

func (this *RedisTaggedCache) pushForever(namespace string, key string) {
	h := sha1.New()

	h.Write(([]byte(namespace)))

	fullKey := this.GetPrefix() + hex.EncodeToString(h.Sum(nil)) + ":" + key

	segments := strings.Split(namespace, "|")

	for _, segment := range segments {

		inputs := []reflect.Value{
			reflect.ValueOf(this.foreverKey(segment)),
			reflect.ValueOf(fullKey),
		}

		reflect.ValueOf(this.Store).MethodByName("Lpush").Call(inputs)
	}
}

func (this *RedisTaggedCache) TagFlush() {
	this.deleteForeverKeys()
}

func (this *RedisTaggedCache) deleteForeverKeys() {
	segments := strings.Split(this.Tags.GetNamespace(), "|")

	for _, segment := range segments {
		key := this.foreverKey(segment)

		this.deleteForeverValues(key)

		this.Store.Forget(segment)
	}
}

func (this *RedisTaggedCache) deleteForeverValues(key string) {
	inputs := []reflect.Value{
		reflect.ValueOf(key),
		reflect.ValueOf(int64(0)),
		reflect.ValueOf(int64(-1)),
	}

	keys := reflect.ValueOf(this.Store).MethodByName("Lrange").Call(inputs)

	if len(keys) > 0 {
		for _, key := range keys {
			if key.Len() > 0 {
				for i := 0; i < key.Len(); i++ {
					this.Store.Forget(key.Index(i).String())
				}
			}
		}
	}
}

func (this *RedisTaggedCache) foreverKey(segment string) string {
	return this.GetPrefix() + segment + ":forever"
}
