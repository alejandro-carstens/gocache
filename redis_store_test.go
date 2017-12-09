package cache

import (
	"config"
	"testing"
)

type Example struct {
	Name        string
	Description string
}

func TestPutGet(t *testing.T) {
	cache := getCache()

	cache.Put("key", "value", 1)

	got := cache.Get("key")

	if got != "value" {
		t.Error("Expected value, got ", got)
	}

	cache.Put("key", 1, 1)

	got = cache.Get("key")

	if got != int64(1) {
		t.Error("Expected 1, got ", got)
	}

	cache.Put("key", 2.99, 1)

	got = cache.Get("key")

	if got != float64(2.99) {
		t.Error("Expected 2.99, got", got)
	}

	cache.Forget("key")
}

func TestPutGetInt(t *testing.T) {
	cache := getCache()

	cache.Put("key", 100, 1)

	got, err := cache.GetInt("key")

	if got != int64(100) || err != nil {
		t.Error("Expected 100, got ", got)
	}

	cache.Forget("key")
}

func TestPutGetFloat(t *testing.T) {
	cache := getCache()

	var expected float64

	expected = 9.99

	cache.Put("key", expected, 1)

	got, err := cache.GetFloat("key")

	if got != expected || err != nil {
		t.Error("Expected 9.99, got ", got)
	}

	cache.Forget("key")
}

func TestIncrement(t *testing.T) {
	cache := getCache()

	cache.Increment("increment_key", 1)
	cache.Increment("increment_key", 1)
	got, err := cache.GetInt("increment_key")

	cache.Forget("increment_key")

	var expected int64 = 2

	if got != expected || err != nil {
		t.Error("Expected 2, got ", got)
	}
}

func TestDecrement(t *testing.T) {
	cache := getCache()

	cache.Increment("decrement_key", 2)
	cache.Decrement("decrement_key", 1)

	var expected int64 = 1

	got, err := cache.GetInt("decrement_key")

	if got != expected || err != nil {
		t.Error("Expected "+string(expected)+", got ", got)
	}

	cache.Forget("decrement_key")
}

func TestForever(t *testing.T) {
	cache := getCache()

	expected := "value"

	cache.Forever("key", expected)

	got := cache.Get("key")

	if got != expected {
		t.Error("Expected "+expected+", got ", got)
	}

	cache.Forget("key")
}

func TestStoreStruct(t *testing.T) {
	cache := getCache()

	var example Example

	example.Name = "Alejandro"
	example.Description = "Whatever"

	cache.Put("key", example, 10)

	var newExample Example

	cache.GetStruct("key", &newExample)

	if newExample != example {
		t.Error("The structs are not the same", newExample)
	}
}

func getCache() RedisStore {
	redisClient := database.Redis{}

	return RedisStore{
		Client: redisClient.Client(),
		Prefix: "golavel:",
	}
}
