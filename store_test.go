package gocache

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

const (
	redisDriver    string = "redis"
	memcacheDriver string = "memcache"
	mapDriver      string = "map"
)

var drivers = []string{
	"map",
	"memcache",
	"redis",
}

type example struct {
	Name        string
	Description string
}

func TestPutGetInt64(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)
		require.NoError(t, cache.Put("key", 100, 1))

		got, err := cache.GetInt64("key")
		require.NoError(t, err)
		require.Equal(t, int64(100), got)

		_, err = cache.Forget("key")
		require.NoError(t, err)
	}
}

func TestPutGetString(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)
		require.NoError(t, cache.Put("key", "value", 1))

		got, err := cache.GetString("key")
		require.NoError(t, err)
		require.Equal(t, "value", got)

		_, err = cache.Forget("key")
		require.NoError(t, err)
	}
}

func TestPutGetFloat64(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		expected := 9.99
		require.NoError(t, cache.Put("key", expected, 1))
		got, err := cache.GetFloat64("key")
		require.NoError(t, err)
		require.Equal(t, expected, got)

		_, err = cache.Forget("key")
		require.NoError(t, err)
	}
}

func TestForever(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		expected := "value"
		require.NoError(t, cache.Forever("key", expected))

		got, err := cache.GetString("key")
		require.NoError(t, err)
		require.Equal(t, expected, got)

		_, err = cache.Forget("key")
		require.NoError(t, err)
	}
}

func TestPutGetMany(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		keys := map[string]string{
			"key_1": "value",
			"key_2": "100",
			"key_3": "9.99",
		}
		require.NoError(t, cache.PutMany(keys, 10))

		resultKeys := []string{
			"key_1",
			"key_2",
			"key_3",
		}
		results, err := cache.Many(resultKeys)
		require.NoError(t, err)

		for i, result := range results {
			require.Equal(t, result, keys[i])
		}
		_, err = cache.Flush()
		require.NoError(t, err)
	}
}

func TestPutGet(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		var firstExample example
		firstExample.Name = "Alejandro"
		firstExample.Description = "Whatever"
		require.NoError(t, cache.Put("key", firstExample, 10))

		var newExample example
		require.NoError(t, cache.Get("key", &newExample))
		require.Equal(t, firstExample, newExample)

		_, err := cache.Forget("key")
		require.NoError(t, err)
	}
}

func TestIncrement(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		_, err := cache.Increment("increment_key", 1)
		require.NoError(t, err)

		_, err = cache.Increment("increment_key", 1)
		require.NoError(t, err)

		got, err := cache.GetInt64("increment_key")
		require.NoError(t, err)

		var expected int64 = 2
		require.Equal(t, expected, got)

		_, err = cache.Forget("increment_key")
		require.NoError(t, err)
	}
}

func TestDecrement(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		_, err := cache.Increment("decrement_key", 2)
		require.NoError(t, err)
		_, err = cache.Decrement("decrement_key", 1)
		require.NoError(t, err)

		got, err := cache.GetInt64("decrement_key")
		require.NoError(t, err)

		var expected int64 = 1
		require.Equal(t, expected, got)

		_, err = cache.Forget("decrement_key")
		require.NoError(t, err)
	}
}

func TestIncrementDecrement(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		got, err := cache.Increment("key", 2)
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
	}
}

func createStore(store string) Cache {
	var (
		cache Cache
		err   error
	)
	switch strings.ToLower(store) {
	case redisDriver:
		cache, err = New(redisStore())
	case memcacheDriver:
		cache, err = New(memcacheStore())
	case mapDriver:
		cache, err = New(mapStore())
	}
	if err != nil {
		panic(err)
	}
	if cache == nil {
		panic("No valid driver provided.")
	}

	return cache
}
