package gocache

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var drivers = []driver{
	redisDriver,
	memcacheDriver,
	localDriver,
}

type example struct {
	Name        string
	Description string
}

func TestPutGetInt64(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			cache := createStore(t, d)
			require.NoError(t, cache.Put("key", 100, time.Second))

			got, err := cache.GetInt64("key")
			require.NoError(t, err)
			require.Equal(t, int64(100), got)

			_, err = cache.Forget("key")
			require.NoError(t, err)
		})
	}
}

func TestPutGetInt(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			cache := createStore(t, d)
			require.NoError(t, cache.Put("key", 100, time.Second))

			got, err := cache.GetInt("key")
			require.NoError(t, err)
			require.Equal(t, 100, got)

			_, err = cache.Forget("key")
			require.NoError(t, err)
		})
	}
}

func TestPutGetString(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			cache := createStore(t, d)
			require.NoError(t, cache.Put("key", "value", time.Second))

			got, err := cache.GetString("key")
			require.NoError(t, err)
			require.Equal(t, "value", got)

			_, err = cache.Forget("key")
			require.NoError(t, err)
		})
	}
}

func TestPutGetFloat64(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache    = createStore(t, d)
				expected = 9.99
			)
			require.NoError(t, cache.Put("key", expected, time.Second))

			got, err := cache.GetFloat64("key")
			require.NoError(t, err)
			require.Equal(t, expected, got)

			_, err = cache.Forget("key")
			require.NoError(t, err)
		})
	}
}

func TestPutGetFloat32(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache    = createStore(t, d)
				expected = 9.99
			)
			require.NoError(t, cache.Put("key", expected, time.Second))

			got, err := cache.GetFloat32("key")
			require.NoError(t, err)
			require.Equal(t, float32(expected), got)

			_, err = cache.Forget("key")
			require.NoError(t, err)
		})
	}
}

func TestPutGetUint64(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			cache := createStore(t, d)
			require.NoError(t, cache.Put("key", 100, time.Second))

			got, err := cache.GetUint64("key")
			require.NoError(t, err)
			require.Equal(t, uint64(100), got)

			_, err = cache.Forget("key")
			require.NoError(t, err)
		})
	}
}

func TestForever(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache    = createStore(t, d)
				expected = "value"
			)
			require.NoError(t, cache.Forever("key", expected))

			got, err := cache.GetString("key")
			require.NoError(t, err)
			require.Equal(t, expected, got)

			_, err = cache.Forget("key")
			require.NoError(t, err)
		})
	}
}

func TestPutGetMany(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache = createStore(t, d)
				keys  = map[string]string{
					"key_1": "value",
					"key_2": "100",
					"key_3": "9.99",
				}
			)
			require.NoError(t, cache.PutMany(keys, 10*time.Second))

			var (
				resultKeys = []string{
					"key_1",
					"key_2",
					"key_3",
				}
				results, err = cache.Many(resultKeys)
			)
			require.NoError(t, err)

			for i, result := range results {
				require.Equal(t, result, keys[i])
			}

			_, err = cache.Flush()
			require.NoError(t, err)
		})
	}
}

func TestPutGet(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache        = createStore(t, d)
				firstExample example
			)
			firstExample.Name = "Alejandro"
			firstExample.Description = "Whatever"
			require.NoError(t, cache.Put("key", firstExample, 10*time.Second))

			var newExample example
			require.NoError(t, cache.Get("key", &newExample))
			require.Equal(t, firstExample, newExample)

			_, err := cache.Forget("key")
			require.NoError(t, err)
		})
	}
}

func TestIncrement(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache  = createStore(t, d)
				_, err = cache.Increment("increment_key", 1)
			)
			require.NoError(t, err)

			_, err = cache.Increment("increment_key", 1)
			require.NoError(t, err)

			got, err := cache.GetInt64("increment_key")
			require.NoError(t, err)

			var expected int64 = 2
			require.Equal(t, expected, got)

			_, err = cache.Forget("increment_key")
			require.NoError(t, err)
		})
	}
}

func TestDecrement(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache  = createStore(t, d)
				_, err = cache.Increment("decrement_key", 2)
			)
			require.NoError(t, err)

			_, err = cache.Decrement("decrement_key", 1)
			require.NoError(t, err)

			got, err := cache.GetInt64("decrement_key")
			require.NoError(t, err)

			var expected int64 = 1
			require.Equal(t, expected, got)

			_, err = cache.Forget("decrement_key")
			require.NoError(t, err)
		})
	}
}

func TestIncrementDecrement(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache    = createStore(t, d)
				got, err = cache.Increment("key", 2)
			)
			require.NoError(t, err)
			require.Equal(t, int64(2), got)

			got, err = cache.Increment("key", 8)
			require.NoError(t, err)
			require.Equal(t, int64(10), got)

			got, err = cache.Decrement("key", 10)
			require.NoError(t, err)
			require.Equal(t, int64(0), got)

			got, err = cache.Decrement("key1", 0)
			require.NoError(t, err)
			require.Equal(t, int64(0), got)

			got, err = cache.Increment("key1", 10)
			require.NoError(t, err)
			require.Equal(t, int64(10), got)

			got, err = cache.Decrement("key1", 10)
			require.NoError(t, err)
			require.Equal(t, int64(0), got)

			_, err = cache.Flush()
			require.NoError(t, err)
		})
	}
}

func createStore(t *testing.T, d driver) Cache {
	var (
		cnf config
		err error
	)
	switch d {
	case redisDriver:
		cnf = &RedisConfig{
			Prefix: "golavel:",
			Addr:   os.Getenv("REDIS_ADDR"),
		}
	case memcacheDriver:
		cnf = &MemcacheConfig{
			Prefix:  "golavel:",
			Servers: []string{os.Getenv("MEMCACHE_SERVER")},
		}
	case localDriver:
		cnf = &LocalConfig{
			Prefix: "golavel:",
		}
	}

	cache, err := New(cnf)
	require.NoError(t, err)

	return cache
}
