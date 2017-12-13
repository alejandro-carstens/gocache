package cache

import (
	"testing"
)

func TestPutGetWithTags(t *testing.T) {
	cache := getCache()

	expected := "value"

	tags := tags()

	cache.Tags(tags).Put("key", "value", 10)

	got := cache.Tags(tags).Get("key")

	if got != expected {
		t.Error("Expected value, got ", got)
	}

	cache.Tags(tags).Forget("key")
}

func TestPutGetIntWithTags(t *testing.T) {
	cache := getCache()

	tags := tags()

	cache.Tags(tags).Put("key", 100, 1)

	got := cache.Tags(tags).Get("key")

	if got != int64(100) {
		t.Error("Expected 100, got ", got)
	}

	cache.Tags(tags).Forget("key")
}

func TestPutGetFloatWithTags(t *testing.T) {
	cache := getCache()

	var expected float64

	expected = 9.99

	tags := tags()

	cache.Tags(tags).Put("key", expected, 1)

	got := cache.Tags(tags).Get("key")

	if got != expected {
		t.Error("Expected 9.99, got ", got)
	}

	cache.Tags(tags).Forget("key")
}

func TestIncrementWithTags(t *testing.T) {
	cache := getCache()

	tags := tags()

	cache.Tags(tags).Increment("increment_key", 1)
	cache.Tags(tags).Increment("increment_key", 1)
	got := cache.Tags(tags).Get("increment_key")

	var expected int64 = 2

	if got != expected {
		t.Error("Expected 2, got ", got)
	}

	cache.Tags(tags).Forget("increment_key")
}

func TestDecrementWithTags(t *testing.T) {
	cache := getCache()

	tags := tags()

	cache.Tags(tags).Increment("decrement_key", 2)
	cache.Tags(tags).Decrement("decrement_key", 1)

	var expected int64 = 1

	got := cache.Tags(tags).Get("decrement_key")

	if got != expected {
		t.Error("Expected "+string(expected)+", got ", got)
	}

	cache.Tags(tags).Forget("decrement_key")
}

func TestForeverWithTags(t *testing.T) {
	cache := getCache()

	expected := "value"

	tags := tags()

	cache.Tags(tags).Forever("key", expected)

	got := cache.Tags(tags).Get("key")

	if got != expected {
		t.Error("Expected "+expected+", got ", got)
	}

	cache.Tags(tags).Forget("key")
}

func TestPutGetManyWithTags(t *testing.T) {
	cache := getCache()

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

	results := cache.Tags(tags).Many(result_keys)

	for i, _ := range results {
		if results[i] != keys[i] {
			t.Error(i, results[i])
		}
	}

	cache.Tags(tags).Flush()
}

func tags() []string {
	tags := make([]string, 1)

	tags[0] = "tag"

	return tags
}
