package cache

import (
	"testing"
)

func TestMemcachePutGet(t *testing.T) {
	cache := getMemcacheCache()

	cache.Put("foo", "bar", 1)

	got, err := cache.Get("foo")

	if got != "bar" || err != nil {
		t.Error("Expected bar got", got)
	}

	cache.Put("foo", 100, 1)

	gotInt, err := cache.Get("foo")

	if gotInt != int64(100) || err != nil {
		t.Error("Expected bar 100", gotInt)
	}

	cache.Put("foo", 10.1, 1)

	gotFloat, err := cache.Get("foo")

	if gotFloat != float64(10.1) || err != nil {
		t.Error("Expected bar 10.1", gotFloat)
	}

	cache.Forget("foo")
}

func TestMemcacheGetInt(t *testing.T) {
	cache := getMemcacheCache()

	cache.Put("foo", 10, 1)

	got, err := cache.GetInt("foo")

	if err != nil {
		panic(err)
	}

	if got != int64(10) {
		t.Error("Expected bar 10.0", got)
	}
}

func TestMemcacheGetFloat(t *testing.T) {
	cache := getMemcacheCache()

	cache.Put("foo", 10.0, 1)

	got, err := cache.GetFloat("foo")

	if err != nil {
		panic(err)
	}

	if got != float64(10.0) {
		t.Error("Expected bar 10.0", got)
	}

	cache.Forget("foo")
}

func TestMemcacheIncrementDecrement(t *testing.T) {
	cache := getMemcacheCache()

	got, err := cache.Increment("key", 2)

	if got != int64(2) {
		t.Error("Expected bar 2", got)
	}

	got, err = cache.Increment("key", 8)

	if got != int64(10) {
		t.Error("Expected bar 10", got)
	}

	got, err = cache.Decrement("key", 10)

	if got != int64(0) {
		t.Error("Expected bar 0", got)
	}

	got, err = cache.Decrement("key1", 0)

	if got != int64(0) {
		t.Error("Expected bar 0", got)
	}

	got, err = cache.Increment("key1", 10)

	if got != int64(10) {
		t.Error("Expected bar 10", got)
	}

	got, err = cache.Decrement("key1", 10)

	if got != int64(0) {
		t.Error("Expected bar 0", got)
	}

	if err != nil {
		panic(err)
	}

	cache.Flush()
}

func TestMemcachePutManyGetMany(t *testing.T) {
	cache := getMemcacheCache()

	keys := make(map[string]interface{})

	keys["foo_1"] = "value"
	keys["foo_2"] = int64(100)
	keys["foo_3"] = float64(9.99)

	cache.PutMany(keys, 10)

	result_keys := make([]string, 3)

	result_keys[0] = "foo_1"
	result_keys[1] = "foo_2"
	result_keys[2] = "foo_3"

	results, err := cache.Many(result_keys)

	if err != nil {
		panic(err)
	}

	for i, result := range results {
		if result != keys[i] {
			t.Error("Expected got", result)
		}
	}

	cache.Flush()
}

func TestMemcacheForever(t *testing.T) {
	cache := getMemcacheCache()

	expected := "value"

	cache.Forever("key", expected)

	got, err := cache.Get("key")

	if got != expected || err != nil {
		t.Error("Expected "+expected+", got ", got)
	}

	cache.Forget("key")
}

func TestMemcachePutGetStruct(t *testing.T) {
	cache := getMemcacheCache()

	var example Example

	example.Name = "Alejandro"
	example.Description = "Whatever"

	cache.Put("key", example, 10)

	var newExample Example

	cache.GetStruct("key", &newExample)

	if newExample != example {
		t.Error("The structs are not the same", newExample)
	}

	cache.Forget("key")
}

func getMemcacheCache() StoreInterface {
	params := make(map[string]interface{})

	params["server 1"] = "127.0.0.1:11211"
	params["prefix"] = "golavel:"

	return Cache("memcache", params)
}
