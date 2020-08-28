package gocache

import (
	"testing"
)

func TestMemcacheConnector(t *testing.T) {
	memcacheStore, err := new(memcacheConnector).connect(memcacheStore())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := memcacheStore.(store); !ok {
		t.Error("Expected StoreInterface got", memcacheStore)
	}
}

func TestRedisConnector(t *testing.T) {
	redisStore, err := new(redisConnector).connect(redisStore())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := redisStore.(store); !ok {
		t.Error("Expected StoreInterface got", redisStore)
	}
}

func TestArrayConnector(t *testing.T) {
	mapStore, err := new(mapConnector).connect(mapStore())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := mapStore.(store); !ok {
		t.Error("Expected StoreInterface got", mapStore)
	}
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

func mapStore() map[string]interface{} {
	params := make(map[string]interface{})

	params["prefix"] = "golavel:"

	return params
}
