package cache

import (
	"testing"
)

func TestPutGetMemcacheWithTags(t *testing.T) {
	cache := getMemcacheCache()

	expected := "value"

	tags := tags()

	cache.Tags(tags).Put("key", "value", 10)

	got, err := cache.Tags(tags).Get("key")

	if got != expected || err != nil {
		t.Error("Expected value, got ", got)
	}

	cache.Tags(tags).Forget("key")
}

func TestPutGetIntMemcacheWithTags(t *testing.T) {
	cache := getMemcacheCache()

	tags := tags()

	cache.Tags(tags).Put("key", 100, 1)

	got, err := cache.Tags(tags).Get("key")

	if got != int64(100) || err != nil {
		t.Error("Expected 100, got ", got)
	}

	cache.Tags(tags).Forget("key")
}

func TestPutGetFloatMemcacheWithTags(t *testing.T) {
	cache := getMemcacheCache()

	var expected float64

	expected = 9.99

	tags := tags()

	cache.Tags(tags).Put("key", expected, 1)

	got, err := cache.Tags(tags).Get("key")

	if got != expected || err != nil {
		t.Error("Expected 9.99, got ", got)
	}

	cache.Tags(tags).Forget("key")
}

func TestIncrementMemcacheWithTags(t *testing.T) {
	cache := getMemcacheCache()

	tags := tags()

	cache.Tags(tags).Increment("increment_key", 1)
	cache.Tags(tags).Increment("increment_key", 1)
	got, err := cache.Tags(tags).Get("increment_key")

	var expected int64 = 2

	if got != expected || err != nil {
		t.Error("Expected 2, got ", got)
	}

	cache.Tags(tags).Forget("increment_key")
}

func TestDecrementMemcacheWithTags(t *testing.T) {
	cache := getMemcacheCache()

	tags := tags()

	cache.Tags(tags).Increment("decrement_key", 2)
	cache.Tags(tags).Decrement("decrement_key", 1)

	var expected int64 = 1

	got, err := cache.Tags(tags).Get("decrement_key")

	if got != expected || err != nil {
		t.Error("Expected "+string(expected)+", got ", got)
	}

	cache.Tags(tags).Forget("decrement_key")
}

func TestForeverMemcacheWithTags(t *testing.T) {
	cache := getMemcacheCache()

	expected := "value"

	tags := tags()

	cache.Tags(tags).Forever("key", expected)

	got, err := cache.Tags(tags).Get("key")

	if got != expected || err != nil {
		t.Error("Expected "+expected+", got ", got)
	}

	cache.Tags(tags).Forget("key")
}

func TestPutGetManyMemcacheWithTags(t *testing.T) {
	cache := getMemcacheCache()

	tags := tags()

	keys := make(map[string]interface{})

	keys["key_1"] = "value"
	keys["key_2"] = int64(100)
	keys["key_3"] = float64(9.99)

	cache.Tags(tags).PutMany(keys, 10)

	result_keys := make([]string, 3)

	result_keys[0] = "key_1"
	result_keys[1] = "key_2"
	result_keys[2] = "key_3"

	results, err := cache.Tags(tags).Many(result_keys)

	if err != nil {
		panic(err)
	}

	for i, _ := range results {
		if results[i] != keys[i] {
			t.Error(i, results[i])
		}
	}

	cache.Tags(tags).Flush()
}

func TestPutGetStructMemcacheWithTags(t *testing.T) {
	cache := getMemcacheCache()

	tags := tags()

	var example Example

	example.Name = "Alejandro"
	example.Description = "Whatever"

	cache.Tags(tags).Put("key", example, 10)

	var newExample Example

	cache.Tags(tags).GetStruct("key", &newExample)

	if newExample != example {
		t.Error("The structs are not the same", newExample)
	}

	cache.Forget("key")
}

func getMemcacheCache() StoreInterface {
	params := make(map[string]interface{})

	params["server 1"] = "127.0.0.1:11211"
	params["prefix"] = "golavel:"

	cache, err := Cache("memcache", params)

	if err != nil {
		panic(err)
	}

	return cache
}
