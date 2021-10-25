package gocache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLock(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache    = createStore(t, d)
				got, err = cache.Lock("test", "test", 10).Acquire()
			)
			require.NoError(t, err)
			require.True(t, got)

			got, err = cache.Lock("test", "test", 10).Acquire()
			require.NoError(t, err)
			require.False(t, got)

			user, err := cache.Lock("test", "test", 10).GetCurrentOwner()
			require.NoError(t, err)
			require.Equal(t, "test", user)

			got, err = cache.Lock("test", "test", 10).Release()
			require.NoError(t, err)
			require.True(t, got)

			got, err = cache.Lock("test", "test", 10).Acquire()
			require.NoError(t, err)
			require.True(t, got)
			require.NoError(t, cache.Lock("test", "test", 10).ForceRelease())

			_, err = cache.Flush()
			require.NoError(t, err)
		})
	}
}
