package gocache

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type (
	driver  string
	example struct {
		Name        string
		Description string
	}
)

func (d driver) string() string {
	return string(d)
}

const (
	redisDriver    driver = "redis"
	memcacheDriver driver = "memcache"
	localDriver    driver = "local"
)

var drivers = []driver{
	redisDriver,
	memcacheDriver,
	localDriver,
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
				cache   = createStore(t, d)
				entries = []Entry{
					{
						Key:      "string",
						Value:    "string",
						Duration: 10 * time.Second,
					},
					{
						Key:      "uint64",
						Value:    uint64(100),
						Duration: 10 * time.Second,
					},
					{
						Key:      "int",
						Value:    100,
						Duration: 10 * time.Second,
					},
					{
						Key:      "int64",
						Value:    int64(100),
						Duration: 10 * time.Second,
					},
					{
						Key:      "float64",
						Value:    float64(100),
						Duration: 10 * time.Second,
					},
					{
						Key:      "float32",
						Value:    float32(100),
						Duration: 10 * time.Second,
					},
					{
						Key: "struct",
						Value: example{
							Name:        "hello",
							Description: "world",
						},
						Duration: 10 * time.Second,
					},
				}
			)
			require.NoError(t, cache.PutMany(entries...))

			var (
				expectedResults = map[string]interface{}{
					"string":  "string",
					"uint64":  uint64(100),
					"int":     100,
					"int64":   int64(100),
					"float64": float64(100),
					"float32": float32(100),
					"struct": example{
						Name:        "hello",
						Description: "world",
					},
					"error": ErrNotFound,
				}
				results, err = cache.Many("string", "uint64", "int", "int64", "float64", "float32", "struct", "error")
			)
			require.NoError(t, err)

			for _, result := range results {
				switch result.Key() {
				case "string":
					require.Equal(t, expectedResults[result.Key()], result.String())
				case "uint64":
					res, err := result.Uint64()
					require.NoError(t, err)
					require.Equal(t, expectedResults[result.Key()], res)
				case "int":
					res, err := result.Int()
					require.NoError(t, err)
					require.Equal(t, expectedResults[result.Key()], res)
				case "int64":
					res, err := result.Int64()
					require.NoError(t, err)
					require.Equal(t, expectedResults[result.Key()], res)
				case "float64":
					res, err := result.Float64()
					require.NoError(t, err)
					require.Equal(t, expectedResults[result.Key()], res)
				case "float32":
					res, err := result.Float32()
					require.NoError(t, err)
					require.Equal(t, expectedResults[result.Key()], res)
				case "struct":
					var res example
					require.NoError(t, result.Unmarshal(&res))
					require.Equal(t, expectedResults[result.Key()], res)
					require.False(t, result.EntryNotFound())
					require.NoError(t, result.Error())
				case "error":
					require.Equal(t, expectedResults[result.Key()], result.Error())
					require.True(t, result.EntryNotFound())
				}
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
	t.Helper()

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
