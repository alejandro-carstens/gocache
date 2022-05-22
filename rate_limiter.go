package gocache

import (
	"errors"
	"fmt"
	"time"
)

// NewRateLimiter creates an instance of *RateLimiter
func NewRateLimiter(cache Cache) *RateLimiter {
	return &RateLimiter{
		cache: cache,
	}
}

type (
	// RateLimiter is a struct whose purpose is to limit the rate at which something is executed or accessed. The
	// underlying logic of this implementation allows for a max number of hits x for a given duration y. If x is
	// exceeded during the y timeframe the RateLimiter will limit further calls until duration y expires. Once
	// duration y is expired calls will reset back to 0
	RateLimiter struct {
		cache Cache
	}
	// ThrottleResponse contains relevant information with respect to calls that are meant to be rate limited
	ThrottleResponse struct {
		retryAfter        time.Duration
		remainingAttempts int64
		maxAttempts       int64
	}
)

// TooManyAttempts determines if the given key has been "accessed" too many times
func (l *RateLimiter) TooManyAttempts(key string, maxAttempts int64) (bool, error) {
	left, err := l.AttemptsLeft(key, maxAttempts)
	if err != nil {
		return false, err
	}

	return left == 0, nil
}

// Hit increments the counter for a given key for a given decay time
func (l *RateLimiter) Hit(key string, decay time.Duration) (hits int64, err error) {
	// The cache will manage the decay as the timer will be expired by the cache
	if _, err = l.cache.Add(l.formatKey(key, "timer"), l.availableAt(decay).Unix(), decay); err != nil {
		return 0, err
	}
	if _, err = l.cache.Add(l.formatKey(key, "counter"), int64(0), decay); err != nil {
		return 0, err
	}

	hits, err = l.cache.Increment(l.formatKey(key, "counter"), 1)
	if err != nil {
		return 0, err
	}

	return hits, nil
}

// Attempts gets the number of attempts for the given key
func (l *RateLimiter) Attempts(key string) (int64, error) {
	val, err := l.cache.GetInt64(l.formatKey(key, "counter"))
	if err != nil && !errors.Is(err, ErrNotFound) {
		return val, err
	}

	return val, nil
}

// AttemptsLeft gets the number of attempts left for a given key
func (l *RateLimiter) AttemptsLeft(key string, maxAttempts int64) (int64, error) {
	attempts, err := l.Attempts(key)
	if err != nil {
		return 0, err
	}

	left := maxAttempts - attempts
	if left > 0 {
		return left, nil
	}

	val, err := l.cache.GetInt64(l.formatKey(key, "timer"))
	if err != nil && !errors.Is(err, ErrNotFound) {
		return 0, err
	}
	if val > 0 {
		return 0, nil
	}
	// If the timer is already at 0 we can clear the counter and timer
	// and return max attempts
	if err = l.Clear(key); err != nil {
		return 0, err
	}

	return maxAttempts, nil
}

// Clear clears the hits and lockout timer for the given key
func (l *RateLimiter) Clear(key string) error {
	if _, err := l.cache.Forget(l.formatKey(key, "counter")); err != nil {
		return err
	}

	_, err := l.cache.Forget(l.formatKey(key, "timer"))

	return err
}

// AvailableIn gets the number of seconds until the "key" is accessible again
func (l *RateLimiter) AvailableIn(key string) (time.Duration, error) {
	unixTime, err := l.cache.GetInt64(l.formatKey(key, "timer"))
	if err != nil && !errors.Is(err, ErrNotFound) {
		return 0, err
	}

	in := unixTime - time.Now().Unix()
	if in < 0 {
		in = 0
	}

	return time.Duration(in) * time.Second, nil
}

// Throttle rate limits the calls that made to fn
func (l *RateLimiter) Throttle(key string, maxCalls int64, decay time.Duration, fn func() error) (*ThrottleResponse, error) {
	tooMany, err := l.TooManyAttempts(key, maxCalls)
	if err != nil {
		return nil, err
	}
	if tooMany {
		return l.buildResponse(key, maxCalls, tooMany)
	}

	attempts, err := l.Hit(key, decay)
	if err != nil {
		return nil, err
	}
	if tooMany = attempts > maxCalls; tooMany {
		return l.buildResponse(key, maxCalls, tooMany)
	}
	if err = fn(); err != nil {
		return nil, err
	}

	return l.buildResponse(key, maxCalls, tooMany)
}

func (l *RateLimiter) buildResponse(key string, maxAttempts int64, tooMany bool) (*ThrottleResponse, error) {
	retryAfter, err := l.AvailableIn(key)
	if err != nil {
		return nil, err
	}

	var remainingAttempts int64
	if !tooMany {
		attempts, err := l.AttemptsLeft(key, maxAttempts)
		if err != nil {
			return nil, err
		}

		remainingAttempts = attempts + 1
	}

	return &ThrottleResponse{
		retryAfter:        retryAfter,
		remainingAttempts: remainingAttempts,
		maxAttempts:       maxAttempts,
	}, nil
}

// availableAt returns the current time plus a decay duration
func (*RateLimiter) availableAt(decay time.Duration) time.Time {
	return time.Now().Add(decay)
}

func (*RateLimiter) formatKey(p1, p2 string) string {
	return fmt.Sprintf("%s:%s", p1, p2)
}

// RetryAfter returns the duration that one should wait before executing the next call
func (r *ThrottleResponse) RetryAfter() time.Duration {
	return r.retryAfter
}

// RemainingAttempts returns the remaining calls that can be made to a function before it is throttled
func (r *ThrottleResponse) RemainingAttempts() int64 {
	return r.remainingAttempts
}

// MaxAttempts returns the max attempts for a given call
func (r *ThrottleResponse) MaxAttempts() int64 {
	return r.maxAttempts
}

// IsThrottled returns if the number of calls have been throttled
func (r *ThrottleResponse) IsThrottled() bool {
	return r.remainingAttempts == 0
}
