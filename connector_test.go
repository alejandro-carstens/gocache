package gocache

import (
	"testing"
)

func TestMemcacheConnector(t *testing.T) {
	memcacheStore, err := New(memcacheStore())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := memcacheStore.(store); !ok {
		t.Error("Expected StoreInterface got", memcacheStore)
	}
}

func TestRedisConnector(t *testing.T) {
	redisStore, err := New(redisStore())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := redisStore.(store); !ok {
		t.Error("Expected StoreInterface got", redisStore)
	}
}

func TestArrayConnector(t *testing.T) {
	mapStore, err := New(mapStore())
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := mapStore.(store); !ok {
		t.Error("Expected StoreInterface got", mapStore)
	}
}

func redisStore() *Config {
	return &Config{
		Redis: &RedisConfig{
			Prefix: "golavel:",
			Addr:   "localhost:6379",
		},
	}
}

func memcacheStore() *Config {
	return &Config{
		Memcache: &MemcacheConfig{
			Prefix:  "golavel:",
			Servers: []string{"127.0.0.1:11211"},
		},
	}
}

func mapStore() *Config {
	return &Config{
		Map: &MapConfig{
			Prefix: "golavel:",
		},
	}
}
