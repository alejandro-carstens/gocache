package cache

import (
	"errors"
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

func (this *MemcacheStore) get(key string) string {
	item, err := this.Client.Get(this.GetPrefix() + key)

	if err != nil {
		panic(err)
	}

	value, err := SimpleDecode(string(item.Value))

	if err != nil {
		return string(item.Value)
	}

	return value
}

func (this *MemcacheStore) Get(key string) interface{} {
	value := this.get(key)

	if IsStringNumeric(value) {
		floatValue := StringToFloat64(value)

		if IsFloat(floatValue) {
			return floatValue
		}

		return int64(floatValue)
	}

	return value
}

func (this *MemcacheStore) GetFloat(key string) (float64, error) {
	value := this.get(key)

	if !IsStringNumeric(value) {
		return 0.0, errors.New("Invalid numeric value")
	}

	return StringToFloat64(value), nil
}

func (this *MemcacheStore) GetInt(key string) (int64, error) {
	value := this.get(key)

	if !IsStringNumeric(value) {
		return 0, errors.New("Invalid numeric value")
	}

	return int64(StringToFloat64(value)), nil
}

func (this *MemcacheStore) item(key string, value interface{}, minutes int) *memcache.Item {
	val, err := Encode(value)

	if err != nil {
		panic(err)
	}

	return &memcache.Item{
		Key:        this.GetPrefix() + key,
		Value:      []byte(val),
		Expiration: int32(minutes),
	}
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
