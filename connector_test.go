package gocache

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMemcacheConnector(t *testing.T) {
	memcacheStore, err := New(memcacheStore())
	require.NoError(t, err)

	_, ok := memcacheStore.(store)
	require.True(t, ok)
}

func TestRedisConnector(t *testing.T) {
	redisStore, err := New(redisStore())
	require.NoError(t, err)

	_, ok := redisStore.(store)
	require.True(t, ok)
}

func TestArrayConnector(t *testing.T) {
	mapStore, err := New(mapStore())
	require.NoError(t, err)

	_, ok := mapStore.(store)
	require.True(t, ok)
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
