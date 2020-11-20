package gocache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPutGetInt64WithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		tags := tag()
		require.NoError(t, cache.Tags(tags).Put("key", 100, 1))

		got, err := cache.Tags(tags).GetInt64("key")
		require.NoError(t, err)
		require.Equal(t, int64(100), got)

		_, err = cache.Tags(tags).Forget("key")
		require.NoError(t, err)
	}
}

func TestPutGetFloat64WithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		var expected = 9.99
		tags := tag()
		require.NoError(t, cache.Tags(tags).Put("key", expected, 1))

		got, err := cache.Tags(tags).GetFloat64("key")
		require.NoError(t, err)
		require.Equal(t, expected, got)

		_, err = cache.Tags(tags).Forget("key")
		require.NoError(t, err)
	}
}

func TestIncrementWithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		tags := tag()
		_, err := cache.Tags(tags).Increment("increment_key", 1)
		require.NoError(t, err)

		_, err = cache.Tags(tags).Increment("increment_key", 1)
		require.NoError(t, err)

		got, err := cache.Tags(tags).GetInt64("increment_key")
		require.NoError(t, err)

		var expected int64 = 2
		require.Equal(t, expected, got)

		_, err = cache.Tags(tags).Forget("increment_key")
		require.NoError(t, err)
	}
}

func TestDecrementWithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		tags := tag()
		_, err := cache.Tags(tags).Increment("decrement_key", 2)
		require.NoError(t, err)

		_, err = cache.Tags(tags).Decrement("decrement_key", 1)
		require.NoError(t, err)

		var expected int64 = 1
		got, err := cache.Tags(tags).GetInt64("decrement_key")
		require.NoError(t, err)
		require.Equal(t, expected, got)

		_, err = cache.Tags(tags).Forget("decrement_key")
		require.NoError(t, err)
	}
}

func TestForeverWithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		expected := "value"
		tags := tag()
		require.NoError(t, cache.Tags(tags).Forever("key", expected))

		got, err := cache.Tags(tags).GetString("key")
		require.NoError(t, err)
		require.Equal(t, expected, got)

		_, err = cache.Tags(tags).Forget("key")
		require.NoError(t, err)
	}
}

func TestPutGetManyWithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		keys := map[string]string{
			"key_1": "value",
			"key_2": "100",
			"key_3": "9.99",
		}
		tags := tag()
		require.NoError(t, cache.Tags(tags).PutMany(keys, 10))

		resultKeys := []string{
			"key_1",
			"key_2",
			"key_3",
		}
		results, err := cache.Tags(tags).Many(resultKeys)
		require.NoError(t, err)

		for i := range results {
			require.Equal(t, keys[i], results[i])
		}

		_, err = cache.Tags(tags).Flush()
		require.NoError(t, err)
	}
}

func TestPutGetWithTags(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		tags := make([]string, 3)
		tags[0] = "tag1"
		tags[1] = "tag2"
		tags[2] = "tag3"

		var firstExample example
		firstExample.Name = "Alejandro"
		firstExample.Description = "Whatever"
		require.NoError(t, cache.Tags(tags...).Put("key", firstExample, 10))

		var newExample example
		require.NoError(t, cache.Tags(tags...).Get("key", &newExample))
		require.Equal(t, firstExample, newExample)

		_, err := cache.Tags(tags...).Forget("key")
		require.NoError(t, err)
	}
}

func TestTagSet(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		tagSet := cache.Tags("Alejandro").GetTags()
		namespace, err := tagSet.getNamespace()
		require.NoError(t, err)
		require.Equal(t, 20, len([]rune(namespace)))
		require.Nil(t, tagSet.reset())
	}
}

func tag() string {
	return "tag"
}
