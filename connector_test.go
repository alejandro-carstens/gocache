package gocache

import (
	"os"
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
			Addr:   os.Getenv("REDIS_ADDR"),
		},
	}
}

func memcacheStore() *Config {
	return &Config{
		Memcache: &MemcacheConfig{
			Prefix:  "golavel:",
			Servers: []string{os.Getenv("MEMCACHE_SERVER")},
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
