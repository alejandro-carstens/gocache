package gocache

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type counter int64

func (c *counter) increment() int64 {
	return atomic.AddInt64((*int64)(c), 1)
}

func TestRateLimiter_Hit(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
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
		for _, d := range drivers(t) {
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
		for _, d := range drivers(t) {
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
		for _, d := range drivers(t) {
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
		for _, d := range drivers(t) {
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
		for _, d := range drivers(t) {
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
		for _, d := range drivers(t) {
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

func TestRateLimit_ThrottleSynchronous(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					c           = createStore(t, d, e)
					rateLimiter = NewRateLimiter(c)
				)
				defer func() {
					_, err := c.Flush()
					require.NoError(t, err)
				}()

				var cnt int
				for i := 0; i < 101; i++ {
					response, err := rateLimiter.Throttle("whatever", 100, 2*time.Second, func() error {
						cnt++

						return nil
					})
					require.NoError(t, err)
					require.EqualValues(t, 100, response.MaxAttempts())
					require.True(t, 1*time.Second <= response.RetryAfter())
					if i < 100 {
						require.Equal(t, i+1, cnt)
						require.EqualValues(t, 100-i, response.RemainingAttempts())
						require.False(t, response.IsThrottled())

						continue
					}

					require.True(t, response.IsThrottled())
					require.EqualValues(t, 0, response.RemainingAttempts())
					require.Equal(t, 100, cnt)
				}
			})
		}
	}
}

func TestRateLimit_ThrottleConcurrent(t *testing.T) {
	for _, e := range encoders {
		for _, d := range drivers(t) {
			t.Run(d.string(), func(t *testing.T) {
				var (
					c           = createStore(t, d, e)
					rateLimiter = NewRateLimiter(c)
				)
				defer func() {
					_, err := c.Flush()
					require.NoError(t, err)
				}()

				var (
					cnt counter
					wg  sync.WaitGroup
				)
				wg.Add(100)
				for i := 0; i < 100; i++ {
					go func() {
						defer wg.Done()
						_, err := rateLimiter.Throttle("whatever", 99, 2*time.Second, func() error {
							cnt.increment()

							return nil
						})
						require.NoError(t, err)
					}()
				}
				wg.Wait()

				response, err := rateLimiter.Throttle("whatever", 99, 2*time.Second, func() error {
					cnt++

					return nil
				})
				require.NoError(t, err)
				require.EqualValues(t, 99, response.MaxAttempts())
				require.True(t, 1*time.Second <= response.RetryAfter())
				require.True(t, response.IsThrottled())
				require.EqualValues(t, 0, response.RemainingAttempts())
				require.EqualValues(t, 99, cnt)
			})
		}
	}
}
