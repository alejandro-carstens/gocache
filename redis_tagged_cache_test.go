package cache

import (
	"testing"
)

func TestPutGetWithTags(t *testing.T) {
	cache := getCache()

	tags := make([]string, 1)

	tags[0] = "tag"

	expected := "value"

	cache.Tags(tags).Put("key", "value", 10)

	got := cache.Tags(tags).Get("key")

	if got != expected {
		t.Error("Expected value, got ", got)
	}

	cache.Tags(tags).Forget("key")
}
