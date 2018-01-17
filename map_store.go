package cache

import (
	"errors"
	"strconv"
)

// MapStore is the representation of an array caching store
type MapStore struct {
	Client map[string]interface{}
	Prefix string
}

// Get gets a value from the store
func (as *MapStore) Get(key string) (interface{}, error) {
	value := as.Client[as.GetPrefix()+key]

	if value == nil {
		return "", nil
	}

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

	return SimpleDecode(value.(string))
}

// GetFloat gets a float value from the store
func (as *MapStore) GetFloat(key string) (float64, error) {
	value := as.Client[as.GetPrefix()+key]

	if value == nil || !IsStringNumeric(value.(string)) {
		return 0, errors.New("Invalid numeric value")
	}

	return StringToFloat64(value.(string))
}

// GetInt gets an int value from the store
func (as *MapStore) GetInt(key string) (int64, error) {
	value := as.Client[as.GetPrefix()+key]

	if value == nil || !IsStringNumeric(value.(string)) {
		return 0, errors.New("Invalid numeric value")
	}

	val, err := StringToFloat64(value.(string))

	return int64(val), err
}

// Increment increments an integer counter by a given value
func (as *MapStore) Increment(key string, value int64) (int64, error) {
	val := as.Client[as.GetPrefix()+key]

	if val != nil {
		if IsStringNumeric(val.(string)) {
			floatValue, err := StringToFloat64(val.(string))

			if err != nil {
				return 0, err
			}

			result := value + int64(floatValue)

			err = as.Put(key, result, 0)

			return result, err
		}

	}

	err := as.Put(key, value, 0)

	return value, err
}

// Decrement decrements an integer counter by a given value
func (as *MapStore) Decrement(key string, value int64) (int64, error) {
	return as.Increment(key, -value)
}

// Put puts a value in the given store for a predetermined amount of time in mins.
func (as *MapStore) Put(key string, value interface{}, minutes int) error {
	val, err := Encode(value)

	mins := strconv.Itoa(minutes)

	mins = ""

	as.Client[as.GetPrefix()+key+mins] = val

	return err
}

// Forever puts a value in the given store until it is forgotten/evicted
func (as *MapStore) Forever(key string, value interface{}) error {
	return as.Put(key, value, 0)
}

// Flush flushes the store
func (as *MapStore) Flush() (bool, error) {
	as.Client = make(map[string]interface{})

	return true, nil
}

// Forget forgets/evicts a given key-value pair from the store
func (as *MapStore) Forget(key string) (bool, error) {
	_, ok := as.Client[as.GetPrefix()+key]

	if ok {
		delete(as.Client, as.GetPrefix()+key)

		return true, nil
	}

	return false, nil
}

// GetPrefix gets the cache key prefix
func (as *MapStore) GetPrefix() string {
	return as.Prefix
}

// PutMany puts many values in the given store until they are forgotten/evicted
func (as *MapStore) PutMany(values map[string]interface{}, minutes int) error {
	for key, value := range values {
		as.Put(key, value, minutes)
	}

	return nil
}

// Many gets many values from the store
func (as *MapStore) Many(keys []string) (map[string]interface{}, error) {
	items := make(map[string]interface{})

	for _, key := range keys {
		val, err := as.Get(key)

		if err != nil {
			return items, err
		}

		items[key] = val
	}

	return items, nil
}

// GetStruct gets the struct representation of a value from the store
func (as *MapStore) GetStruct(key string, entity interface{}) (interface{}, error) {
	value := as.Client[as.GetPrefix()+key]

	return Decode(value.(string), entity)
}

// Tags returns the TaggedCache for the given store
func (as *MapStore) Tags(names []string) TaggedStoreInterface {
	return &TaggedCache{
		Store: as,
		Tags: TagSet{
			Store: as,
			Names: names,
		},
	}
}
