package gocache

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/alejandro-carstens/gocache/encoder"
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

var (
	driverList = []driver{
		redisDriver,
		memcacheDriver,
		localDriver,
	}
	encoders = []encoder.Encoder{
		encoder.JSON{},
		encoder.Msgpack{},
	}
)

func drivers(t *testing.T, excludeList ...driver) []driver {
	t.Helper()

	var list []driver
	for _, d := range driverList {
		var exclude bool
		for _, e := range excludeList {
			if d == e {
				exclude = true

				break
			}
		}

		if !exclude {
			list = append(list, d)
		}
	}

	return list
}

func TestPutGetInt64(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				cache := createStore(t, d, e)
				require.NoError(t, cache.Put("key", 100, time.Second))

				got, err := cache.GetInt64("key")
				require.NoError(t, err)
				require.EqualValues(t, 100, got)

				_, err = cache.Forget("key")
				require.NoError(t, err)
			})
		}
	}
}

func TestPutGetInt(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				cache := createStore(t, d, e)
				require.NoError(t, cache.Put("key", 100, time.Second))

				got, err := cache.GetInt("key")
				require.NoError(t, err)
				require.Equal(t, 100, got)

				_, err = cache.Forget("key")
				require.NoError(t, err)
			})
		}
	}
}

func TestPutGetString(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				cache := createStore(t, d, e)
				require.NoError(t, cache.Put("key", "value", time.Second))

				got, err := cache.GetString("key")
				require.NoError(t, err)
				require.Equal(t, "value", got)

				_, err = cache.Forget("key")
				require.NoError(t, err)
			})
		}
	}
}

func TestPutGetFloat64(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache    = createStore(t, d, e)
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
}

func TestPutGetFloat32(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache    = createStore(t, d, e)
					expected = 9.99
				)
				require.NoError(t, cache.Put("key", expected, time.Second))

				got, err := cache.GetFloat32("key")
				require.NoError(t, err)
				require.EqualValues(t, expected, got)

				_, err = cache.Forget("key")
				require.NoError(t, err)
			})
		}
	}
}

func TestPutGetUint64(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				cache := createStore(t, d, e)
				require.NoError(t, cache.Put("key", 100, time.Second))

				got, err := cache.GetUint64("key")
				require.NoError(t, err)
				require.EqualValues(t, 100, got)

				_, err = cache.Forget("key")
				require.NoError(t, err)
			})
		}
	}
}

func TestPutGetBool(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				cache := createStore(t, d, e)
				require.NoError(t, cache.Put("key", true, time.Second))

				got, err := cache.GetBool("key")
				require.NoError(t, err)
				require.Equal(t, true, got)

				require.NoError(t, cache.Put("key", "a", time.Second))

				got, err = cache.GetBool("key")
				require.NoError(t, err)
				require.Equal(t, true, got)

				require.NoError(t, cache.Put("key", 1, time.Second))

				got, err = cache.GetBool("key")
				require.NoError(t, err)
				require.Equal(t, true, got)

				require.NoError(t, cache.Put("key", false, time.Second))

				got, err = cache.GetBool("key")
				require.NoError(t, err)
				require.Equal(t, false, got)

				require.NoError(t, cache.Put("key", 0, time.Second))

				got, err = cache.GetBool("key")
				require.NoError(t, err)
				require.Equal(t, false, got)

				require.NoError(t, cache.Put("key", "", time.Second))

				got, err = cache.GetBool("key")
				require.NoError(t, err)
				require.Equal(t, false, got)

				_, err = cache.Forget("key")
				require.NoError(t, err)
			})
		}
	}
}

func TestForever(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache    = createStore(t, d, e)
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
}

func TestPutGetMany(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache   = createStore(t, d, e)
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
						{
							Key:      "bool",
							Value:    false,
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
						"bool":  false,
					}
					results, err = cache.Many("string", "uint64", "int", "int64", "float64", "float32", "struct", "error", "bool")
				)
				require.NoError(t, err)

				for _, result := range results {
					switch result.Key() {
					case "string":
						res, err := result.String()
						require.NoError(t, err)
						require.Equal(t, expectedResults[result.Key()], res)
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
					case "error":
						require.Equal(t, expectedResults[result.Key()], result.Error())
						require.True(t, result.EntryNotFound())
					case "bool":
						res, err := result.Bool()
						require.NoError(t, err)
						require.Equal(t, expectedResults[result.Key()], res)
					}
				}

				_, err = cache.Flush()
				require.NoError(t, err)
			})
		}
	}
}

func TestPutGet(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache        = createStore(t, d, e)
					firstExample example
				)
				firstExample.Name = "Alejandro"
				firstExample.Description = "Whatever"
				require.NoError(t, cache.Put("key", firstExample, 10*time.Second))

				var newExample example
				require.NoError(t, cache.Get("key", &newExample))
				require.Equal(t, firstExample, newExample)

				type custom int
				require.NoError(t, cache.Put("key", custom(1), time.Second))

				var c custom
				require.NoError(t, cache.Get("key", &c))
				require.EqualValues(t, 1, c)

				_, err := cache.Forget("key")
				require.NoError(t, err)
			})
		}
	}
}

func TestIncrement(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache  = createStore(t, d, e)
					_, err = cache.Increment("increment_key", 1)
				)
				require.NoError(t, err)

				_, err = cache.Increment("increment_key", 1)
				require.NoError(t, err)

				got, err := cache.GetInt64("increment_key")
				require.NoError(t, err)
				require.EqualValues(t, 2, got)

				_, err = cache.Forget("increment_key")
				require.NoError(t, err)
			})
		}
	}
}

func TestDecrement(t *testing.T) {
	var expectedLoweBoundSet = map[string]int64{
		memcacheDriver.string(): 0,
		redisDriver.string():    -2,
		localDriver.string():    -2,
	}
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache  = createStore(t, d, e)
					_, err = cache.Increment("decrement_key", 2)
				)
				require.NoError(t, err)

				_, err = cache.Decrement("decrement_key", 1)
				require.NoError(t, err)

				got, err := cache.GetInt64("decrement_key")
				require.NoError(t, err)
				require.EqualValues(t, 1, got)

				_, err = cache.Forget("decrement_key")
				require.NoError(t, err)

				got, err = cache.Decrement("decrement_key", 2)
				require.NoError(t, err)
				require.Equal(t, expectedLoweBoundSet[d.string()], got)

				_, err = cache.Forget("decrement_key")
				require.NoError(t, err)
			})
		}
	}
}

func TestIncrementDecrement(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache    = createStore(t, d, e)
					got, err = cache.Increment("key", 2)
				)
				require.NoError(t, err)
				require.EqualValues(t, 2, got)

				got, err = cache.Increment("key", 8)
				require.NoError(t, err)
				require.EqualValues(t, 10, got)

				got, err = cache.Decrement("key", 10)
				require.NoError(t, err)
				require.EqualValues(t, 0, got)

				got, err = cache.Decrement("key1", 0)
				require.NoError(t, err)
				require.EqualValues(t, 0, got)

				got, err = cache.Increment("key1", 10)
				require.NoError(t, err)
				require.EqualValues(t, 10, got)

				got, err = cache.Decrement("key1", 10)
				require.NoError(t, err)
				require.EqualValues(t, 0, got)

				_, err = cache.Flush()
				require.NoError(t, err)
			})
		}
	}
}

func TestExists(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache = createStore(t, d, e)
					err   = cache.Put("key", 2, time.Second)
				)
				require.NoError(t, err)

				exists, err := cache.Exists("key")
				require.NoError(t, err)
				require.True(t, exists)

				_, err = cache.Forget("key")
				require.NoError(t, err)

				exists, err = cache.Exists("key")
				require.NoError(t, err)
				require.False(t, exists)

				_, err = cache.Flush()
				require.NoError(t, err)
			})
		}
	}
}

func TestAdd(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache    = createStore(t, d, e)
					res, err = cache.Add("key", 2, time.Second)
				)
				require.NoError(t, err)
				require.True(t, res)

				res, err = cache.Add("key", 2, time.Second)
				require.NoError(t, err)
				require.False(t, res)

				i, err := cache.GetInt("key")
				require.NoError(t, err)
				require.Equal(t, 2, i)

				res, err = cache.Forget("key")
				require.NoError(t, err)
				require.True(t, res)

				res, err = cache.Add("key", 2, time.Second)
				require.NoError(t, err)
				require.True(t, res)

				res, err = cache.Add("other_key", "whatever", time.Second)
				require.NoError(t, err)
				require.True(t, res)

				res, err = cache.Add("other_key", "whatever", time.Second)
				require.NoError(t, err)
				require.False(t, res)

				v, err := cache.GetString("other_key")
				require.NoError(t, err)
				require.Equal(t, "whatever", v)

				_, err = cache.Flush()
				require.NoError(t, err)
			})
		}
	}
}

func TestForget(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache    = createStore(t, d, e)
					res, err = cache.Add("key", 2, time.Second)
				)
				require.NoError(t, err)
				require.True(t, res)

				res, err = cache.Forget("key")
				require.NoError(t, err)
				require.True(t, res)

				_, err = cache.GetInt("key")
				require.Equal(t, ErrNotFound, err)

				res, err = cache.Forget("key")
				require.False(t, res)
			})
		}
	}
}

func TestForgetMany(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache = createStore(t, d, e)
					err   = cache.PutMany(Entry{
						Key:      "key1",
						Value:    1,
						Duration: time.Second,
					}, Entry{
						Key:      "key2",
						Value:    2,
						Duration: time.Second,
					}, Entry{
						Key:      "key3",
						Value:    3,
						Duration: time.Second,
					})
				)
				require.NoError(t, err)

				err = cache.ForgetMany("key1", "key2")
				require.NoError(t, err)

				res, err := cache.Many("key1", "key2", "key3")
				require.NoError(t, err)
				require.Equal(t, ErrNotFound, res["key1"].Error())
				require.Equal(t, ErrNotFound, res["key2"].Error())

				v, err := res["key3"].Int()
				require.NoError(t, err)
				require.Equal(t, 3, v)
				require.NoError(t, cache.ForgetMany("key3"))

				_, err = cache.GetInt("key")
				require.Equal(t, ErrNotFound, err)
			})
		}
	}
}

func TestExpire(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t, localDriver) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache = createStore(t, d, e)
					err   = cache.Put("key1", 1, time.Second)
				)
				require.NoError(t, err)
				require.NoError(t, cache.Expire("key1", 0*time.Second))
				require.Error(t, ErrNotFound, cache.Expire("key1", time.Second))
			})
		}
	}
}

func createStore(t *testing.T, d driver, encoder encoder.Encoder) Cache {
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

	cache, err := New(cnf, encoder)
	require.NoError(t, err)

	return cache
}
