package gocache

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPutGetInt64WithTags(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache = createStore(t, d)
				ts    = tag()
			)
			require.NoError(t, cache.Tags(ts).Put("key", 100, time.Second))

			got, err := cache.Tags(ts).GetInt64("key")
			require.NoError(t, err)
			require.Equal(t, int64(100), got)

			_, err = cache.Tags(ts).Forget("key")
			require.NoError(t, err)
		})
	}
}

func TestPutGetIntWithTags(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache = createStore(t, d)
				ts    = tag()
			)
			require.NoError(t, cache.Tags(ts).Put("key", 100, time.Second))

			got, err := cache.Tags(ts).GetInt("key")
			require.NoError(t, err)
			require.Equal(t, 100, got)

			_, err = cache.Tags(ts).Forget("key")
			require.NoError(t, err)
		})
	}
}

func TestPutGetFloat64WithTags(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache    = createStore(t, d)
				expected = 9.99
				ts       = tag()
			)
			require.NoError(t, cache.Tags(ts).Put("key", expected, time.Second))

			got, err := cache.Tags(ts).GetFloat64("key")
			require.NoError(t, err)
			require.Equal(t, expected, got)

			_, err = cache.Tags(ts).Forget("key")
			require.NoError(t, err)
		})
	}
}

func TestPutGetFloat32WithTags(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache    = createStore(t, d)
				expected = 9.99
				ts       = tag()
			)
			require.NoError(t, cache.Tags(ts).Put("key", expected, time.Second))

			got, err := cache.Tags(ts).GetFloat32("key")
			require.NoError(t, err)
			require.Equal(t, float32(expected), got)

			_, err = cache.Tags(ts).Forget("key")
			require.NoError(t, err)
		})
	}
}

func TestPutGetUint64WithTags(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache = createStore(t, d)
				ts    = tag()
			)
			require.NoError(t, cache.Tags(ts).Put("key", 100, time.Second))

			got, err := cache.Tags(ts).GetUint64("key")
			require.NoError(t, err)
			require.Equal(t, uint64(100), got)

			_, err = cache.Tags(ts).Forget("key")
			require.NoError(t, err)
		})
	}
}

func TestIncrementWithTags(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache  = createStore(t, d)
				ts     = tag()
				_, err = cache.Tags(ts).Increment("increment_key", 1)
			)
			require.NoError(t, err)

			_, err = cache.Tags(ts).Increment("increment_key", 1)
			require.NoError(t, err)

			got, err := cache.Tags(ts).GetInt64("increment_key")
			require.NoError(t, err)

			var expected int64 = 2
			require.Equal(t, expected, got)

			_, err = cache.Tags(ts).Forget("increment_key")
			require.NoError(t, err)
		})
	}
}

func TestDecrementWithTags(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache  = createStore(t, d)
				ts     = tag()
				_, err = cache.Tags(ts).Increment("decrement_key", 2)
			)
			require.NoError(t, err)

			_, err = cache.Tags(ts).Decrement("decrement_key", 1)
			require.NoError(t, err)

			got, err := cache.Tags(ts).GetInt64("decrement_key")
			require.NoError(t, err)
			require.Equal(t, int64(1), got)

			_, err = cache.Tags(ts).Forget("decrement_key")
			require.NoError(t, err)
		})
	}
}

func TestForeverWithTags(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache    = createStore(t, d)
				expected = "value"
				ts       = tag()
			)
			require.NoError(t, cache.Tags(ts).Forever("key", expected))

			got, err := cache.Tags(ts).GetString("key")
			require.NoError(t, err)
			require.Equal(t, expected, got)

			_, err = cache.Tags(ts).Forget("key")
			require.NoError(t, err)
		})
	}
}

func TestPutGetManyWithTags(t *testing.T) {
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
				}
				results, err = cache.Tags(ts).Many("string", "uint64", "int", "int64", "float64", "float32", "struct")
			)
			require.NoError(t, err)

			for _, result := range results {
				require.NotEmpty(t, result.TagKey())

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
				}
			}

			_, err = cache.Flush()
			require.NoError(t, err)
		})
	}
}

func TestPutGetWithTags(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache = createStore(t, d)
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

			_, err := cache.Tags(ts...).Forget("key")
			require.NoError(t, err)
		})
	}
}

func TestTagSet(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache          = createStore(t, d)
				ts             = cache.Tags("Alejandro").GetTags()
				namespace, err = ts.getNamespace()
			)
			require.NoError(t, err)
			require.Equal(t, 20, len([]rune(namespace)))
			require.Nil(t, ts.reset())
		})
	}
}

func TestFlushWithTags(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache = createStore(t, d)
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

			_, err = cache.Flush()
			require.NoError(t, err)
		})
	}
}

func tag() string {
	return "tag"
}
