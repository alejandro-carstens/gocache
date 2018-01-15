package cache

import (
	"strings"
	"testing"
)

var drivers = []string{"array", "memcache", "redis"}

type Example struct {
	Name        string
	Description string
}

func TestPutGet(t *testing.T) {
	for _, driver := range drivers {
		cache := store(driver)

		cache.Put("key", "value", 1)

		got, err := cache.Get("key")

		if got != "value" || err != nil {
			t.Error("Expected value, got ", got)
		}

		cache.Put("key", 1, 1)

		got, err = cache.Get("key")

		if got != int64(1) || err != nil {
			t.Error("Expected 1, got ", got)
		}

		cache.Put("key", 2.99, 1)

		got, err = cache.Get("key")

		if got != float64(2.99) || err != nil {
			t.Error("Expected 2.99, got", got)
		}

		cache.Forget("key")
	}
}

func TestPutGetInt(t *testing.T) {
	for _, driver := range drivers {
		cache := store(driver)

		cache.Put("key", 100, 1)

		got, err := cache.GetInt("key")

		if got != int64(100) || err != nil {
			t.Error("Expected 100, got ", got)
		}

		cache.Forget("key")
	}
}

func TestPutGetFloat(t *testing.T) {
	for _, driver := range drivers {
		cache := store(driver)

		var expected float64

		expected = 9.99

		cache.Put("key", expected, 1)

		got, err := cache.GetFloat("key")

		if got != expected || err != nil {
			t.Error("Expected 9.99, got ", got)
		}

		cache.Forget("key")
	}
}

func TestForever(t *testing.T) {
	for _, driver := range drivers {
		cache := store(driver)

		expected := "value"

		cache.Forever("key", expected)

		got, err := cache.Get("key")

		if got != expected || err != nil {
			t.Error("Expected "+expected+", got ", got)
		}

		cache.Forget("key")
	}
}

func TestPutGetMany(t *testing.T) {
	for _, driver := range drivers {
		cache := store(driver)

		keys := make(map[string]interface{})

		keys["key_1"] = "value"
		keys["key_2"] = int64(100)
		keys["key_3"] = float64(9.99)

		cache.PutMany(keys, 10)

		result_keys := make([]string, 3)

		result_keys[0] = "key_1"
		result_keys[1] = "key_2"
		result_keys[2] = "key_3"

		results, err := cache.Many(result_keys)

		if err != nil {
			panic(err)
		}

		for i, result := range results {
			if result != keys[i] {
				t.Error("Expected got", results["key_1"])
			}
		}

		cache.Flush()
	}
}

func TestPutGetStruct(t *testing.T) {
	for _, driver := range drivers {
		cache := store(driver)

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
}

func TestIncrement(t *testing.T) {
	for _, driver := range drivers {
		cache := store(driver)

		cache.Increment("increment_key", 1)
		cache.Increment("increment_key", 1)
		got, err := cache.GetInt("increment_key")

		cache.Forget("increment_key")

		var expected int64 = 2

		if got != expected || err != nil {
			t.Error("Expected 2, got ", got)
		}
	}
}

func TestDecrement(t *testing.T) {
	for _, driver := range drivers {
		cache := store(driver)

		cache.Increment("decrement_key", 2)
		cache.Decrement("decrement_key", 1)

		var expected int64 = 1

		got, err := cache.GetInt("decrement_key")

		if got != expected || err != nil {
			t.Error("Expected "+string(expected)+", got ", got)
		}

		cache.Forget("decrement_key")
	}
}

func TesIncrementDecrement(t *testing.T) {
	for _, driver := range drivers {
		cache := store(driver)

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
}

func store(store string) StoreInterface {
	switch strings.ToLower(store) {
	case "redis":
		cache, err := New(store, redisStore())

		if err != nil {
			panic(err)
		}

		return cache
	case "memcache":
		cache, err := New(store, memcacheStore())

		if err != nil {
			panic(err)
		}

		return cache
	case "array":
		cache, err := New(store, arrayStore())

		if err != nil {
			panic(err)
		}

		return cache
	}

	panic("No valid driver provided.")
}

func redisStore() map[string]interface{} {
	params := make(map[string]interface{})

	params["address"] = "localhost:6379"
	params["password"] = ""
	params["database"] = 0
	params["prefix"] = "golavel:"

	return params
}

func memcacheStore() map[string]interface{} {
	params := make(map[string]interface{})

	params["server 1"] = "127.0.0.1:11211"
	params["prefix"] = "golavel:"

	return params
}

func arrayStore() map[string]interface{} {
	params := make(map[string]interface{})

	params["prefix"] = "golavel:"

	return params
}
