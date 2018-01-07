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
	item, err := this.item(key, value, minutes)

	if err != nil {
		return err
	}

	return this.Client.Set(item)
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

	return this.processValue(value)
}

func (this *MemcacheStore) GetFloat(key string) (float64, error) {
	value, err := this.get(key)

	if err != nil {
		return 0.0, err
	}

	if !IsStringNumeric(value) {
		return 0.0, errors.New("Invalid numeric value")
	}

	return StringToFloat64(value)
}

func (this *MemcacheStore) GetInt(key string) (int64, error) {
	value, err := this.get(key)

	if err != nil {
		return 0, err
	}

	if !IsStringNumeric(value) {
		return 0, errors.New("Invalid numeric value")
	}

	val, err := StringToFloat64(value)

	return int64(val), err
}

func (this *MemcacheStore) item(key string, value interface{}, minutes int) (*memcache.Item, error) {
	val, err := Encode(value)

	return &memcache.Item{
		Key:        this.GetPrefix() + key,
		Value:      []byte(val),
		Expiration: int32(minutes),
	}, err
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

func (this *MemcacheStore) PutMany(values map[string]interface{}, minutes int) error {
	for key, value := range values {
		err := this.Put(key, value, minutes)

		if err != nil {
			return err
		}
	}

	return nil
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

func (this *MemcacheStore) processValue(value string) (interface{}, error) {
	if IsStringNumeric(value) {
		floatValue, err := StringToFloat64(value)

		if err != nil {
			return floatValue, err
		}

		if IsFloat(floatValue) {
			return floatValue, err
		}

		return int64(floatValue), err
	}

	return value, nil
}
