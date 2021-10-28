package gocache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLock(t *testing.T) {
	for _, d := range drivers {
		t.Run(d.string(), func(t *testing.T) {
			var (
				cache    = createStore(t, d)
				got, err = cache.Lock("test", "test", 10*time.Second).Acquire()
			)
			require.NoError(t, err)
			require.True(t, got)

			got, err = cache.Lock("test", "test", 10*time.Second).Acquire()
			require.NoError(t, err)
			require.False(t, got)

			user, err := cache.Lock("test", "test", 10*time.Second).GetCurrentOwner()
			require.NoError(t, err)
			require.Equal(t, "test", user)

			got, err = cache.Lock("test", "test", 10*time.Second).Release()
			require.NoError(t, err)
			require.True(t, got)

			got, err = cache.Lock("test", "test", 10*time.Second).Acquire()
			require.NoError(t, err)
			require.True(t, got)
			require.NoError(t, cache.Lock("test", "test", 10*time.Second).ForceRelease())

			_, err = cache.Flush()
			require.NoError(t, err)
		})
	}
}
