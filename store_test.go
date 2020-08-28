package gocache

import (
	"strings"
	"testing"
)

var drivers = []string{
	"map",
	"memcache",
	"redis",
}

type example struct {
	Name        string
	Description string
}

func TestPutGetInt64(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		if err := cache.Put("key", 100, 1); err != nil {
			t.Fatal(err)
		}

		got, err := cache.GetInt64("key")
		if got != int64(100) || err != nil {
			t.Error("Expected 100, got ", got)
		}
		if _, err := cache.Forget("key"); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPutGetString(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		if err := cache.Put("key", "value", 1); err != nil {
			t.Fatal(err)
		}

		got, err := cache.GetString("key")
		if got != "value" || err != nil {
			t.Error("Expected value, got ", got)
		}
		if _, err := cache.Forget("key"); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPutGetFloat64(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		var expected float64

		expected = 9.99
		if err := cache.Put("key", expected, 1); err != nil {
			t.Fatal(err)
		}
		got, err := cache.GetFloat64("key")
		if got != expected || err != nil {
			t.Error("Expected 9.99, got ", got)
		}
		if _, err := cache.Forget("key"); err != nil {
			t.Fatal(err)
		}
	}
}

func TestForever(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		expected := "value"
		if err := cache.Forever("key", expected); err != nil {
			t.Fatal(err)
		}

		got, err := cache.GetString("key")
		if got != expected || err != nil {
			t.Error("Expected "+expected+", got ", got)
		}
		if _, err := cache.Forget("key"); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPutGetMany(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		keys := make(map[string]string)

		keys["key_1"] = "value"
		keys["key_2"] = "100"
		keys["key_3"] = "9.99"

		if err := cache.PutMany(keys, 10); err != nil {
			t.Fatal(err)
		}

		resultKeys := make([]string, 3)

		resultKeys[0] = "key_1"
		resultKeys[1] = "key_2"
		resultKeys[2] = "key_3"

		results, err := cache.Many(resultKeys)
		if err != nil {
			t.Fatal(err)
		}

		for i, result := range results {
			if result != keys[i] {
				t.Error("Expected got", results["key_1"])
			}
		}

		if _, err := cache.Flush(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestPutGet(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		var firstExample example

		firstExample.Name = "Alejandro"
		firstExample.Description = "Whatever"

		if err := cache.Put("key", firstExample, 10); err != nil {
			t.Fatal(err)
		}

		var newExample example

		if err := cache.Get("key", &newExample); err != nil {
			t.Fatal(err)
		}
		if newExample != firstExample {
			t.Error("The structs are not the same", newExample)
		}
		if _, err := cache.Forget("key"); err != nil {
			t.Fatal(err)
		}
	}
}

func TestIncrement(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		if _, err := cache.Increment("increment_key", 1); err != nil {
			t.Fatal(err)
		}
		if _, err := cache.Increment("increment_key", 1); err != nil {
			t.Fatal(err)
		}

		got, err := cache.GetInt64("increment_key")
		if _, err := cache.Forget("increment_key"); err != nil {
			t.Fatal(err)
		}

		var expected int64 = 2
		if got != expected || err != nil {
			t.Error("Expected 2, got ", got)
		}
	}
}

func TestDecrement(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		if _, err := cache.Increment("decrement_key", 2); err != nil {
			t.Fatal(err)
		}
		if _, err := cache.Decrement("decrement_key", 1); err != nil {
			t.Fatal(err)
		}

		var expected int64 = 1

		got, err := cache.GetInt64("decrement_key")
		if got != expected || err != nil {
			t.Error("Expected "+string(expected)+", got ", got)
		}
		if _, err := cache.Forget("decrement_key"); err != nil {
			t.Fatal(err)
		}
	}
}

func TestIncrementDecrement(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

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
			t.Fatal(err)
		}
		if _, err := cache.Flush(); err != nil {
			t.Fatal(err)
		}
	}
}

func createStore(store string) Cache {
	switch strings.ToLower(store) {
	case RedisDriver:
		cache, err := New(store, redisStore())
		if err != nil {
			panic(err)
		}

		return cache
	case MemcacheDriver:
		cache, err := New(store, memcacheStore())
		if err != nil {
			panic(err)
		}

		return cache
	case MapDriver:
		cache, err := New(store, mapStore())
		if err != nil {
			panic(err)
		}

		return cache
	}

	panic("No valid driver provided.")
}
