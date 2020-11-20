package gocache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLock(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		lock := cache.Lock("test", "test", 10)
		got, err := lock.Acquire()
		require.NoError(t, err)
		require.True(t, got)

		if driver == mapDriver {
			got, err = lock.Acquire()
		} else {
			got, err = cache.Lock("test", "test", 10).Acquire()
		}
		require.NoError(t, err)
		require.False(t, got)

		var user string
		if driver == mapDriver {
			user, err = lock.GetCurrentOwner()
		} else {
			user, err = cache.Lock("test", "test", 10).GetCurrentOwner()
		}
		require.NoError(t, err)
		require.Equal(t, "test", user)

		if driver == mapDriver {
			got, err = lock.Release()
		} else {
			got, err = cache.Lock("test", "test", 10).Release()
		}
		require.NoError(t, err)
		require.True(t, got)

		lock = cache.Lock("test", "test", 10)
		got, err = lock.Acquire()
		require.NoError(t, err)
		require.True(t, got)

		if driver == mapDriver {
			err = lock.ForceRelease()
		} else {
			err = cache.Lock("test", "test", 10).ForceRelease()
		}
		require.NoError(t, err)

		_, err = cache.Flush()
		require.NoError(t, err)
	}
}
