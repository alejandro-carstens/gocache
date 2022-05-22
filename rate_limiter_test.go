package gocache

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRateLimiter_Hit(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers {
			t.Run(d.string(), func(t *testing.T) {
				var (
					c           = createStore(t, d, e)
					rateLimiter = NewRateLimiter(c)
				)
				defer func() {
					_, err := c.Flush()
					require.NoError(t, err)
				}()

				hit, err := rateLimiter.Hit("whatever", 1*time.Second)
				require.NoError(t, err)
				require.EqualValues(t, 1, hit)

				val, err := c.GetInt64("whatever:counter")
				require.NoError(t, err)
				require.EqualValues(t, 1, val)

				unix, err := c.GetInt64("whatever:timer")
				require.NoError(t, err)
				require.NotEmpty(t, unix)

				hit, err = rateLimiter.Hit("whatever", 1*time.Second)
				require.NoError(t, err)
				require.EqualValues(t, 2, hit)

				// Let's allow the count to expire
				time.Sleep(1050 * time.Millisecond)

				_, err = c.GetInt64("whatever:timer")
				require.ErrorIs(t, err, ErrNotFound)

				_, err = c.GetInt64("whatever:counter")
				require.ErrorIs(t, err, ErrNotFound)
			})
		}
	}
}

func TestRateLimiter_ConcurrentHits(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers {
			t.Run(d.string(), func(t *testing.T) {
				var (
					wg          sync.WaitGroup
					c           = createStore(t, d, e)
					rateLimiter = NewRateLimiter(c)
				)
				defer func() {
					_, err := c.Flush()
					require.NoError(t, err)
				}()

				wg.Add(100)
				for i := 0; i < 100; i++ {
					go func() {
						defer wg.Done()

						_, err := rateLimiter.Hit("concurrency", 2*time.Second)
						require.NoError(t, err)
					}()
				}
				wg.Wait()

				attemptsLeft, err := rateLimiter.AttemptsLeft("concurrency", 99)
				require.NoError(t, err)
				require.EqualValues(t, 0, attemptsLeft)
			})
		}
	}
}

func TestRateLimiter_Attempts(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers {
			t.Run(d.string(), func(t *testing.T) {
				var (
					c           = createStore(t, d, e)
					rateLimiter = NewRateLimiter(c)
				)
				defer func() {
					_, err := c.Flush()
					require.NoError(t, err)
				}()

				for i := 0; i < 100; i++ {
					hit, err := rateLimiter.Hit("whatever", 2*time.Second)
					require.NoError(t, err)
					require.EqualValues(t, i+1, hit)
				}

				attempts, err := rateLimiter.Attempts("whatever")
				require.NoError(t, err)
				require.EqualValues(t, 100, attempts)
			})
		}
	}
}

func TestRateLimiter_AttemptsLeft(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers {
			t.Run(d.string(), func(t *testing.T) {
				var (
					c           = createStore(t, d, e)
					rateLimiter = NewRateLimiter(c)
				)
				defer func() {
					_, err := c.Flush()
					require.NoError(t, err)
				}()

				for i := 0; i < 100; i++ {
					hit, err := rateLimiter.Hit("whatever", 2*time.Second)
					require.NoError(t, err)
					require.EqualValues(t, i+1, hit)
				}

				attemptsLeft, err := rateLimiter.AttemptsLeft("whatever", 1000)
				require.NoError(t, err)
				require.EqualValues(t, 900, attemptsLeft)

				attemptsLeft, err = rateLimiter.AttemptsLeft("whatever", 50)
				require.NoError(t, err)
				require.EqualValues(t, 0, attemptsLeft)
			})
		}
	}
}

func TestRateLimiter_AvailableIn(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers {
			t.Run(d.string(), func(t *testing.T) {
				var (
					c           = createStore(t, d, e)
					rateLimiter = NewRateLimiter(c)
				)
				defer func() {
					_, err := c.Flush()
					require.NoError(t, err)
				}()

				hit, err := rateLimiter.Hit("whatever", 3*time.Second)
				require.NoError(t, err)
				require.EqualValues(t, 1, hit)

				time.Sleep(1 * time.Second)

				duration, err := rateLimiter.AvailableIn("whatever")
				require.NoError(t, err)
				require.Equal(t, 2*time.Second, duration)
			})
		}
	}
}

func TestRateLimiter_Clear(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers {
			t.Run(d.string(), func(t *testing.T) {
				var (
					c           = createStore(t, d, e)
					rateLimiter = NewRateLimiter(c)
				)
				defer func() {
					_, err := c.Flush()
					require.NoError(t, err)
				}()

				hit, err := rateLimiter.Hit("whatever", 3*time.Second)
				require.NoError(t, err)
				require.EqualValues(t, 1, hit)

				require.NoError(t, rateLimiter.Clear("whatever"))

				_, err = c.GetInt64("whatever")
				require.Equal(t, ErrNotFound, err)

				_, err = c.GetInt64("whatever:timer")
				require.Equal(t, ErrNotFound, err)
			})
		}
	}
}

func TestRateLimiter_TooManyAttempts(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers {
			t.Run(d.string(), func(t *testing.T) {
				var (
					c           = createStore(t, d, e)
					rateLimiter = NewRateLimiter(c)
				)
				defer func() {
					_, err := c.Flush()
					require.NoError(t, err)
				}()

				for i := 0; i < 100; i++ {
					hit, err := rateLimiter.Hit("whatever", 2*time.Second)
					require.NoError(t, err)
					require.EqualValues(t, i+1, hit)
				}

				tooMany, err := rateLimiter.TooManyAttempts("whatever", 100)
				require.NoError(t, err)
				require.True(t, tooMany)

				tooMany, err = rateLimiter.TooManyAttempts("whatever", 1000)
				require.NoError(t, err)
				require.False(t, tooMany)

				tooMany, err = rateLimiter.TooManyAttempts("whatever", 10)
				require.NoError(t, err)
				require.True(t, tooMany)
			})
		}
	}
}
