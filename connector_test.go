package gocache

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemcacheConnector(t *testing.T) {
	s, err := New(&MemcacheConfig{
		Prefix:  "golavel:",
		Servers: []string{os.Getenv("MEMCACHE_SERVER")},
	})
	require.NoError(t, err)

	_, ok := s.(store)
	require.True(t, ok)
}

func TestRedisConnector(t *testing.T) {
	s, err := New(&RedisConfig{
		Prefix: "golavel:",
		Addr:   os.Getenv("REDIS_ADDR"),
	})
	require.NoError(t, err)

	_, ok := s.(store)
	require.True(t, ok)
}

func TestLocalConnector(t *testing.T) {
	s, err := New(&LocalConfig{
		Prefix: "golavel:",
	})
	require.NoError(t, err)

	_, ok := s.(store)
	require.True(t, ok)
}
