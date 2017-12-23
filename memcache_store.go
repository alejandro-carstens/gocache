package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheStore struct {
	Client memcache.Client
	Prefix string
}

func (this *MemcacheStore) Put(key string, value interface{}, minutes int) {
	err := this.Client.Set(this.item(key, value, minutes))

	if err != nil {
		panic(err)
	}
}

func (this *MemcacheStore) item(key string, value interface{}, expiration int) *memcache.Item {
	val, err := Encode(value)

	if err != nil {
		panic(err)
	}

	return &memcache.Item{Key: this.Prefix + key, Value: []byte(val), Expiration: int32(expiration)}
}
