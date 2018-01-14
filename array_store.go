package cache

import (
	"errors"
	"strconv"
)

type ArrayStore struct {
	Client map[string]interface{}
	Prefix string
}

func (this *ArrayStore) Get(key string) (interface{}, error) {
	value := this.Client[this.GetPrefix()+key]

	if IsStringNumeric(value.(string)) {
		floatValue, err := StringToFloat64(value.(string))

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

func (this *ArrayStore) GetFloat(key string) (float64, error) {
	value := this.Client[this.GetPrefix()+key]

	if value == nil || !IsStringNumeric(value.(string)) {
		return 0, errors.New("Invalid numeric value")
	}

	return StringToFloat64(value.(string))
}

func (this *ArrayStore) GetInt(key string) (int64, error) {
	value := this.Client[this.GetPrefix()+key]

	if value == nil || !IsStringNumeric(value.(string)) {
		return 0, errors.New("Invalid numeric value")
	}

	val, err := StringToFloat64(value.(string))

	return int64(val), err
}

func (this *ArrayStore) Increment(key string, value int64) (int64, error) {
	val := this.Client[this.GetPrefix()+key]

	if val != nil {
		this.Client[this.GetPrefix()+key] = value + val.(int64)

		return this.Client[this.GetPrefix()+key].(int64), nil
	}

	this.Client[this.GetPrefix()+key] = value

	return this.Client[this.GetPrefix()+key].(int64), nil
}

func (this *ArrayStore) Decrement(key string, value int64) (int64, error) {
	return this.Increment(key, -value)
}

func (this *ArrayStore) Put(key string, value interface{}, minutes int) error {
	val, err := Encode(value)

	mins := strconv.Itoa(minutes)

	mins = ""

	this.Client[this.GetPrefix()+key+mins] = val

	return err
}

func (this *ArrayStore) Forever(key string, value interface{}) error {
	return this.Put(key, value, 0)
}

func (this *ArrayStore) Flush() (bool, error) {
	this.Client = make(map[string]interface{})

	return true, nil
}

func (this *ArrayStore) Forget(key string) (bool, error) {
	_, ok := this.Client[this.GetPrefix()+key]

	if ok {
		delete(this.Client, this.GetPrefix()+key)

		return true, nil
	}

	return false, nil
}

func (this *ArrayStore) GetPrefix() string {
	return this.Prefix
}

func (this *ArrayStore) PutMany(values map[string]interface{}, minutes int) error {
	for key, value := range values {
		this.Put(key, value, minutes)
	}

	return nil
}

func (this *ArrayStore) Many(keys []string) (map[string]interface{}, error) {
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

func (this *ArrayStore) GetStruct(key string, entity interface{}) (interface{}, error) {
	value, err := this.Get(key)

	if err != nil {
		return value, err
	}

	return Decode(value.(string), entity)
}

func (this *ArrayStore) Tags(names []string) TaggedStoreInterface {
	return &TaggedCache{
		Store: this,
		Tags: TagSet{
			Store: this,
			Names: names,
		},
	}
}
