package gocache

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPutGetInt64WithTags(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache = createStore(t, d, e)
					ts    = tag()
				)
				require.NoError(t, cache.Tags(ts).Put("key", 100, time.Second))

				got, err := cache.Tags(ts).GetInt64("key")
				require.NoError(t, err)
				require.EqualValues(t, 100, got)

				_, err = cache.Tags(ts).Forget("key")
				require.NoError(t, err)
				require.NoError(t, cache.Tags(ts).TagSet().Flush())
			})
		}
	}
}

func TestPutGetIntWithTags(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache = createStore(t, d, e)
					ts    = tag()
				)
				require.NoError(t, cache.Tags(ts).Put("key", 100, time.Second))

				got, err := cache.Tags(ts).GetInt("key")
				require.NoError(t, err)
				require.Equal(t, 100, got)

				_, err = cache.Tags(ts).Forget("key")
				require.NoError(t, err)
				require.NoError(t, cache.Tags(ts).TagSet().Flush())
			})
		}
	}
}

func TestPutGetFloat64WithTags(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache    = createStore(t, d, e)
					expected = 9.99
					ts       = tag()
				)
				require.NoError(t, cache.Tags(ts).Put("key", expected, time.Second))

				got, err := cache.Tags(ts).GetFloat64("key")
				require.NoError(t, err)
				require.Equal(t, expected, got)

				_, err = cache.Tags(ts).Forget("key")
				require.NoError(t, err)
				require.NoError(t, cache.Tags(ts).TagSet().Flush())
			})
		}
	}
}

func TestPutGetFloat32WithTags(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache    = createStore(t, d, e)
					expected = 9.99
					ts       = tag()
				)
				require.NoError(t, cache.Tags(ts).Put("key", expected, time.Second))

				got, err := cache.Tags(ts).GetFloat32("key")
				require.NoError(t, err)
				require.EqualValues(t, expected, got)

				_, err = cache.Tags(ts).Forget("key")
				require.NoError(t, err)
				require.NoError(t, cache.Tags(ts).TagSet().Flush())
			})
		}
	}

}

func TestPutGetUint64WithTags(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache = createStore(t, d, e)
					ts    = tag()
				)
				require.NoError(t, cache.Tags(ts).Put("key", 100, time.Second))

				got, err := cache.Tags(ts).GetUint64("key")
				require.NoError(t, err)
				require.EqualValues(t, 100, got)

				_, err = cache.Tags(ts).Forget("key")
				require.NoError(t, err)
				require.NoError(t, cache.Tags(ts).TagSet().Flush())
			})
		}
	}
}

func TestPutGetBoolWithTags(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache = createStore(t, d, e)
					ts    = tag()
				)
				require.NoError(t, cache.Tags(ts).Put("key", true, time.Second))

				got, err := cache.Tags(ts).GetBool("key")
				require.NoError(t, err)
				require.Equal(t, true, got)

				require.NoError(t, cache.Tags(ts).Put("key", "a", time.Second))

				got, err = cache.Tags(ts).GetBool("key")
				require.NoError(t, err)
				require.Equal(t, true, got)

				require.NoError(t, cache.Tags(ts).Put("key", 1, time.Second))

				got, err = cache.Tags(ts).GetBool("key")
				require.NoError(t, err)
				require.Equal(t, true, got)

				require.NoError(t, cache.Tags(ts).Put("key", false, time.Second))

				got, err = cache.Tags(ts).GetBool("key")
				require.NoError(t, err)
				require.Equal(t, false, got)

				require.NoError(t, cache.Tags(ts).Put("key", 0, time.Second))

				got, err = cache.Tags(ts).GetBool("key")
				require.NoError(t, err)
				require.Equal(t, false, got)

				require.NoError(t, cache.Tags(ts).Put("key", "", time.Second))

				got, err = cache.Tags(ts).GetBool("key")
				require.NoError(t, err)
				require.Equal(t, false, got)

				_, err = cache.Tags(ts).Forget("key")
				require.NoError(t, err)
				require.NoError(t, cache.Tags(ts).TagSet().Flush())
			})
		}
	}
}

func TestIncrementWithTags(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache  = createStore(t, d, e)
					ts     = tag()
					_, err = cache.Tags(ts).Increment("increment_key", 1)
				)
				require.NoError(t, err)

				_, err = cache.Tags(ts).Increment("increment_key", 1)
				require.NoError(t, err)

				got, err := cache.Tags(ts).GetInt64("increment_key")
				require.NoError(t, err)
				require.EqualValues(t, 2, got)

				_, err = cache.Tags(ts).Forget("increment_key")
				require.NoError(t, err)
				require.NoError(t, cache.Tags(ts).TagSet().Flush())
			})
		}
	}
}

func TestDecrementWithTags(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache  = createStore(t, d, e)
					ts     = tag()
					_, err = cache.Tags(ts).Increment("decrement_key", 2)
				)
				require.NoError(t, err)

				_, err = cache.Tags(ts).Decrement("decrement_key", 1)
				require.NoError(t, err)

				got, err := cache.Tags(ts).GetInt64("decrement_key")
				require.NoError(t, err)
				require.EqualValues(t, 1, got)

				_, err = cache.Tags(ts).Forget("decrement_key")
				require.NoError(t, err)
				require.NoError(t, cache.Tags(ts).TagSet().Flush())
			})
		}
	}
}

func TestForeverWithTags(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache    = createStore(t, d, e)
					expected = "value"
					ts       = tag()
				)
				require.NoError(t, cache.Tags(ts).Forever("key", expected))

				got, err := cache.Tags(ts).GetString("key")
				require.NoError(t, err)
				require.Equal(t, expected, got)

				_, err = cache.Tags(ts).Forget("key")
				require.NoError(t, err)
				require.NoError(t, cache.Tags(ts).TagSet().Flush())
			})
		}
	}
}

func TestForgetWithTags(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache    = createStore(t, d, e)
					ts       = tag()
					res, err = cache.Tags(ts).Add("key", 2, time.Second)
				)
				require.NoError(t, err)
				require.True(t, res)

				res, err = cache.Tags(ts).Forget("key")
				require.NoError(t, err)
				require.True(t, res)

				_, err = cache.Tags(ts).GetInt("key")
				require.Equal(t, ErrNotFound, err)

				res, err = cache.Tags(ts).Forget("key")
				require.False(t, res)
				require.NoError(t, cache.Tags(ts).TagSet().Flush())
			})
		}
	}
}

func TestForgetManyWithTags(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache = createStore(t, d, e)
					ts    = tag()
					err   = cache.Tags(ts).PutMany(Entry{
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

				err = cache.Tags(ts).ForgetMany("key1", "key2")
				require.NoError(t, err)

				res, err := cache.Tags(ts).Many("key1", "key2", "key3")
				require.NoError(t, err)
				require.Equal(t, ErrNotFound, res["key1"].Error())
				require.Equal(t, ErrNotFound, res["key2"].Error())

				v, err := res["key3"].Int()
				require.NoError(t, err)
				require.Equal(t, 3, v)
				require.NoError(t, cache.ForgetMany("key3"))

				_, err = cache.Tags(ts).GetInt("key")
				require.Equal(t, ErrNotFound, err)
				require.NoError(t, cache.Tags(ts).TagSet().Flush())
			})
		}
	}
}

func TestPutGetManyWithTags(t *testing.T) {
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
					ts = tag()
				)
				require.NoError(t, cache.Tags(ts).PutMany(entries...))

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
					results, err = cache.Tags(ts).Many("string", "uint64", "int", "int64", "float64", "float32", "struct", "error", "bool")
				)
				require.NoError(t, err)

				for _, result := range results {
					require.NotEmpty(t, result.TagKey())

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
					case "error":
						require.Equal(t, expectedResults[result.Key()], result.Error())
						require.True(t, result.EntryNotFound())
					case "bool":
						res, err := result.Bool()
						require.NoError(t, err)
						require.Equal(t, expectedResults[result.Key()], res)
					}
					if result.Error() != nil {
						continue
					}

					_, err = cache.Tags(ts).Forget(result.Key())
					require.NoError(t, err)
				}

				_, err = cache.Tags(ts).Flush()
				require.NoError(t, err)
				require.NoError(t, cache.Tags(ts).TagSet().Flush())
			})
		}
	}
}

func TestPutGetWithTags(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache = createStore(t, d, e)
					ts    = []string{
						"tag1",
						"tag2",
						"tag3",
					}
					firstExample = example{
						Name:        "Alejandro",
						Description: "Whatever",
					}
				)
				require.NoError(t, cache.Tags(ts...).Put("key", firstExample, 10*time.Second))

				var newExample example
				require.NoError(t, cache.Tags(ts...).Get("key", &newExample))
				require.Equal(t, firstExample, newExample)

				type custom int
				require.NoError(t, cache.Tags(ts...).Put("key", custom(1), time.Second))

				var c custom
				require.NoError(t, cache.Tags(ts...).Get("key", &c))
				require.EqualValues(t, 1, c)

				_, err := cache.Tags(ts...).Forget("key")
				require.NoError(t, err)
				require.NoError(t, cache.Tags(ts...).TagSet().Flush())
			})
		}
	}
}

func TestFlushWithTags(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache = createStore(t, d, e)
					ts1   = []string{"person", "dev"}
					ts2   = []string{"bot", "dev", "ai"}
					ts3   = []string{"person", "painter"}
					ts4   = []string{"person", "driver", "current"}
					ts5   = []string{"person", "driver", "legend"}
				)
				require.NoError(t, cache.Tags(ts1...).Put("joe", "doe", 0))
				require.NoError(t, cache.Tags(ts2...).Put("bot", "doe", time.Second))
				require.NoError(t, cache.Tags(ts3...).Forever("jane", "doe"))

				require.NoError(t, cache.Tags(ts4...).PutMany(Entry{
					Key:      "checo",
					Value:    "perez",
					Duration: 0,
				}, Entry{
					Key:      "lewis",
					Value:    "hamilton",
					Duration: time.Second,
				}))
				require.NoError(t, cache.Tags(ts5...).Put("ayrton", "senna", time.Second))

				val, err := cache.Tags(ts1...).GetString("joe")
				require.NoError(t, err)
				require.Equal(t, "doe", val)

				val, err = cache.Tags(ts2...).GetString("bot")
				require.NoError(t, err)
				require.Equal(t, "doe", val)
				// We flush dev, so we won't be able to access joe or bot anymore
				_, err = cache.Tags("dev").Flush()
				require.NoError(t, err)

				_, err = cache.Tags(ts1...).GetString("joe")
				require.True(t, errors.Is(err, ErrNotFound))

				_, err = cache.Tags(ts2...).GetString("bot")
				require.True(t, errors.Is(err, ErrNotFound))

				// We flush painter so jane should not be available
				val, err = cache.Tags(ts3...).GetString("jane")
				require.NoError(t, err)
				require.Equal(t, "doe", val)

				_, err = cache.Tags("painter").Flush()
				require.NoError(t, err)

				_, err = cache.Tags(ts3...).GetString("jane")
				require.True(t, errors.Is(err, ErrNotFound))

				// We flush all the current drivers so checo and lewis should not be available
				_, err = cache.Tags("current").Flush()
				require.NoError(t, err)

				_, err = cache.Tags(ts4...).GetString("checo")
				require.True(t, errors.Is(err, ErrNotFound))

				_, err = cache.Tags(ts4...).GetString("lewis")
				require.True(t, errors.Is(err, ErrNotFound))

				// We should still be able to access ayrton since he is a legend driver
				val, err = cache.Tags(ts5...).GetString("ayrton")
				require.NoError(t, err)
				require.Equal(t, "senna", val)
				require.NoError(t, cache.Tags(ts1...).TagSet().Flush())
				require.NoError(t, cache.Tags(ts2...).TagSet().Flush())
				require.NoError(t, cache.Tags(ts3...).TagSet().Flush())
				require.NoError(t, cache.Tags(ts4...).TagSet().Flush())
				require.NoError(t, cache.Tags(ts5...).TagSet().Flush())
			})
		}
	}
}

func TestTagExists(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache = createStore(t, d, e)
					ts    = tag()
					err   = cache.Tags(ts).Put("key", 2, time.Second)
				)
				require.NoError(t, err)

				exists, err := cache.Tags(ts).Exists("key")
				require.NoError(t, err)
				require.True(t, exists)

				_, err = cache.Tags(ts).Forget("key")
				require.NoError(t, err)

				exists, err = cache.Tags(ts).Exists("key")
				require.NoError(t, err)
				require.False(t, exists)

				_, err = cache.Tags(ts).Flush()
				require.NoError(t, err)
				require.NoError(t, cache.Tags(ts).TagSet().Flush())
			})
		}
	}
}

func TestTagAdd(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache    = createStore(t, d, e)
					ts       = tag()
					res, err = cache.Tags(ts).Add("key", 2, time.Second)
				)
				require.NoError(t, err)
				require.True(t, res)

				res, err = cache.Tags(ts).Add("key", 2, time.Second)
				require.NoError(t, err)
				require.False(t, res)

				i, err := cache.Tags(ts).GetInt("key")
				require.NoError(t, err)
				require.Equal(t, 2, i)

				res, err = cache.Tags(ts).Flush()
				require.NoError(t, err)
				require.True(t, res)

				res, err = cache.Tags(ts).Add("key", 2, time.Second)
				require.NoError(t, err)
				require.True(t, res)

				res, err = cache.Tags(ts).Add("other_key", "whatever", time.Second)
				require.NoError(t, err)
				require.True(t, res)

				res, err = cache.Tags(ts).Add("other_key", "whatever", time.Second)
				require.NoError(t, err)
				require.False(t, res)

				v, err := cache.Tags(ts).GetString("other_key")
				require.NoError(t, err)
				require.Equal(t, "whatever", v)

				_, err = cache.Tags(ts).Flush()
				require.NoError(t, err)
				require.NoError(t, cache.Tags(ts).TagSet().Flush())
			})
		}
	}
}

func tag() string {
	return "tag"
}
