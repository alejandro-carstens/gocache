package cache

import (
	"testing"
)

func TestMemcachePutGet(t *testing.T) {
	cache := getMemcacheCache()

	cache.Put("foo", "bar", 1)

	got := cache.Get("foo")

	if got != "bar" {
		t.Error("Expected bar got", got)
	}

	cache.Put("foo", 100, 1)

	gotInt := cache.Get("foo")

	if gotInt != int64(100) {
		t.Error("Expected bar 100", gotInt)
	}

	cache.Put("foo", 10.1, 1)

	gotFloat := cache.Get("foo")

	if gotFloat != float64(10.1) {
		t.Error("Expected bar 10.1", gotFloat)
	}
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
}

func getMemcacheCache() *MemcacheStore {
	params := make(map[string]interface{})

	params["server 1"] = "127.0.0.1:11211"
	params["prefix"] = "golavel:"

	memcacheConnector := new(MemcacheConnector)

	return memcacheConnector.Connect(params)
}
