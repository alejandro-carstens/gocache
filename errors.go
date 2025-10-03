package gocache

import (
	"errors"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/redis/go-redis/v9"
)

var (
	// ErrNotFound represents an agnostic cache entry not found error
	ErrNotFound = errors.New("gocache: not found")
	// ErrFailedToRetrieveEntry indicates that an entry was not able to be properly retrieved from the cache when
	// calling cache.Many
	ErrFailedToRetrieveEntry = errors.New("gocache: an error occurred while retrieving value, for more detail call Item.Error()")
	// ErrFailedToAddItemEntry is returned when we expected to add an entry to the cache but an entry was already
	// present for the given key
	ErrFailedToAddItemEntry = errors.New("gocache: failed to add entry to cache")
	// ErrBlockWaitTimeout is returned when the max wait for acquiring a lock during a Block call is exceeded
	ErrBlockWaitTimeout = errors.New("gocache: failed to acquire lock due to lock wait timeout")
	// ErrNotImplemented is returned for methods that have not been implemented for the Cache interface
	ErrNotImplemented = errors.New("gocache: method not implemented")
)

func checkErrNotFound(err error) error {
	if isErrNotFound(err) {
		return ErrNotFound
	}

	return err
}

func isErrNotFound(err error) bool {
	if errors.Is(err, ErrNotFound) {
		return true
	}
	if errors.Is(err, redis.Nil) {
		return true
	}
	if errors.Is(err, memcache.ErrCacheMiss) {
		return true
	}

	return false
}
