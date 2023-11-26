package gocache

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLock(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache    = createStore(t, d, e)
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
}

func TestLock_Get(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache             = createStore(t, d, e)
					wg                sync.WaitGroup
					cnt               int
					acquiredResultSet = map[bool]bool{
						true:  true,
						false: false,
					}
					l = cache.Lock("test", "test", 200*time.Millisecond)
				)
				wg.Add(2)
				go func() {
					defer wg.Done()

					acquired, err := l.Get(func() error {
						cnt++
						time.Sleep(50 * time.Millisecond)

						return nil
					})
					require.NoError(t, err)
					delete(acquiredResultSet, acquired)
				}()
				go func() {
					defer wg.Done()

					acquired, err := l.Get(func() error {
						cnt++
						time.Sleep(50 * time.Millisecond)

						return nil
					})
					require.NoError(t, err)
					delete(acquiredResultSet, acquired)
				}()
				wg.Wait()

				require.Equal(t, 1, cnt)
				require.Len(t, acquiredResultSet, 0)
			})
		}
	}
}

func TestLock_Block(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache = createStore(t, d, e)
					wg    sync.WaitGroup

					cnt             int
					acquiredResults []bool
					errs            []error
					l               = cache.Lock("test", "test", 200*time.Millisecond)
				)
				wg.Add(3)
				go func() {
					defer wg.Done()

					acquired, err := l.Block(25*time.Millisecond, 80*time.Millisecond, func() error {
						cnt++
						time.Sleep(50 * time.Millisecond)

						return nil
					})
					if err != nil {
						errs = append(errs, err)
					}
					acquiredResults = append(acquiredResults, acquired)
				}()
				go func() {
					defer wg.Done()

					acquired, err := l.Block(25*time.Millisecond, 80*time.Millisecond, func() error {
						cnt++
						time.Sleep(50 * time.Millisecond)

						return nil
					})
					if err != nil {
						errs = append(errs, err)
					}
					acquiredResults = append(acquiredResults, acquired)
				}()
				go func() {
					defer wg.Done()

					acquired, err := l.Block(25*time.Millisecond, 80*time.Millisecond, func() error {
						cnt++
						time.Sleep(50 * time.Millisecond)

						return nil
					})
					if err != nil {
						errs = append(errs, err)
					}
					acquiredResults = append(acquiredResults, acquired)
				}()
				wg.Wait()

				require.Equal(t, 2, cnt)
				require.Len(t, acquiredResults, 3)
				require.Len(t, errs, 1)
				require.True(t, errors.Is(errs[0], ErrBlockWaitTimeout))

				var blocked int
				for _, res := range acquiredResults {
					if !res {
						blocked++
					}
				}
				require.Equal(t, 1, blocked)
			})
		}
	}
}

func TestLock_Expire(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t, localDriver, redisDriver) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					cache    = createStore(t, d, e)
					l        = cache.Lock("test", "test", 10*time.Second)
					got, err = l.Acquire()
				)
				require.NoError(t, err)
				require.True(t, got)

				got, err = l.Expire(5 * time.Second)
				require.NoError(t, err)
				require.True(t, got)

				got, err = l.Expire(0 * time.Second)
				require.NoError(t, err)
				require.True(t, got)

				got, err = l.Expire(10 * time.Second)
				require.NoError(t, err)
				require.False(t, got)
			})
		}
	}
}
