package cache

import (
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheStore struct {
	Client memcache.Client
	Prefix string
}

func (this *MemcacheStore) Put(key string, value interface{}, minutes int) error {
	return this.Client.Set(this.item(key, value, minutes))
}

func (this *MemcacheStore) Forever(key string, value interface{}) error {
	return this.Put(key, value, 0)
}

func (this *MemcacheStore) get(key string) (string, error) {
	item, err := this.Client.Get(this.GetPrefix() + key)

	if err != nil {
		if err.Error() == "memcache: cache miss" {
			return "", nil
		}

		return "", err
	}

	return this.getItemValue(item.Value), nil
}

func (this *MemcacheStore) Get(key string) (interface{}, error) {
	value, err := this.get(key)

	if err != nil {
		return value, err
	}

	return this.processValue(value), nil
}

func (this *MemcacheStore) GetFloat(key string) (float64, error) {
	value, err := this.get(key)

	if err != nil {
		return 0.0, err
	}

	if !IsStringNumeric(value) {
		return 0.0, errors.New("Invalid numeric value")
	}

	return StringToFloat64(value), nil
}

func (this *MemcacheStore) GetInt(key string) (int64, error) {
	value, err := this.get(key)

	if err != nil {
		return 0, err
	}

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

func (this *MemcacheStore) Increment(key string, value int64) (int64, error) {
	newValue, err := this.Client.Increment(this.GetPrefix()+key, uint64(value))

	if err != nil {
		if err.Error() != "memcache: cache miss" {
			return value, err
		}

		this.Put(key, value, 0)

		return value, nil
	}

	return int64(newValue), nil
}

func (this *MemcacheStore) Decrement(key string, value int64) (int64, error) {
	newValue, err := this.Client.Decrement(this.GetPrefix()+key, uint64(value))

	if err != nil {
		if err.Error() != "memcache: cache miss" {
			return value, err
		}

		this.Put(key, 0, 0)

		return int64(0), nil
	}

	return int64(newValue), nil
}

func (this *MemcacheStore) GetPrefix() string {
	return this.Prefix
}

func (this *MemcacheStore) PutMany(values map[string]interface{}, minutes int) {
	for key, value := range values {
		this.Put(key, value, minutes)
	}
}

func (this *MemcacheStore) Many(keys []string) (map[string]interface{}, error) {
	items := make(map[string]interface{})

	for _, key := range keys {
		val, err := this.Get(key)

		if err != nil {
			return items, err
		}

		items[key] = val
	}

	return items, nil
}

func (this *MemcacheStore) Forget(key string) (bool, error) {
	err := this.Client.Delete(this.GetPrefix() + key)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (this *MemcacheStore) Flush() (bool, error) {
	err := this.Client.DeleteAll()

	if err != nil {
		return false, err
	}

	return true, nil
}

func (this *MemcacheStore) GetStruct(key string, entity interface{}) (interface{}, error) {
	value, err := this.get(key)

	if err != nil {
		return value, err
	}

	return Decode(value, entity)
}

func (this *MemcacheStore) Tags(names []string) TaggedStoreInterface {
	return &TaggedCache{
		Store: this,
		Tags: TagSet{
			Store: this,
			Names: names,
		},
	}
}

func (this *MemcacheStore) getItemValue(itemValue []byte) string {
	value, err := SimpleDecode(string(itemValue))

	if err != nil {
		return string(itemValue)
	}

	return value
}

func (this *MemcacheStore) processValue(value string) interface{} {
	if IsStringNumeric(value) {
		floatValue := StringToFloat64(value)

		if IsFloat(floatValue) {
			return floatValue
		}

		return int64(floatValue)
	}

	return value
}
