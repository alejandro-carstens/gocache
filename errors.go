package gocache

import (
	"errors"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis"
)

// ErrNotFound represents an agnostic cache entry not found error
var ErrNotFound = errors.New("gocache: not found")

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