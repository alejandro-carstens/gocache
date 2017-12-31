package cache

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheStore struct {
	Client memcache.Client
	Prefix string
}

func (this *MemcacheStore) Put(key string, value interface{}, minutes int) {
	val, err := Encode(value)

	if err != nil {
		panic(err)
	}

	err = this.Client.Add(&memcache.Item{Key: this.getKey(key), Value: []byte(val), Expiration: int32(minutes)})

	if err != nil {
		panic(err)
	}
}

func (this *MemcacheStore) get(key string) (*memcache.Item, error) {
	return this.Client.Get(this.GetPrefix() + key)
}

func (this *MemcacheStore) Get(key string) interface{} {
	item, err := this.get(this.getKey(key))

	if err != nil {
		panic(err)
	}

	return item.Value
}

func (this *MemcacheStore) getKey(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(this.GetPrefix() + key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (this *MemcacheStore) Increment(key string, value int64) int64 {
	newValue, err := this.Client.Increment(this.GetPrefix()+key, uint64(value))

	if err != nil {
		panic(err)
	}

	return int64(newValue)
}

func (this *MemcacheStore) Decrement(key string, value int64) int64 {
	newValue, err := this.Client.Decrement(this.GetPrefix()+key, uint64(value))

	if err != nil {
		panic(err)
	}

	return int64(newValue)
}

func (this *MemcacheStore) GetPrefix() string {
	return this.Prefix
}
