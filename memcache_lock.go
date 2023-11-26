package gocache

import (
	"errors"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

var _ Lock = &memcacheLock{}

func newMemcacheLock(client *memcache.Client, name, owner string, duration time.Duration) *memcacheLock {
	return (&memcacheLock{
		client:   client,
		name:     name,
		owner:    owner,
		duration: duration,
	}).initBaseLock()
}

type memcacheLock struct {
	baseLock
	client   *memcache.Client
	name     string
	owner    string
	duration time.Duration
}

// Acquire implementation of the Lock interface
func (ml *memcacheLock) Acquire() (bool, error) {
	return ml.acquire(ml.duration)
}

// Release implementation of the Lock interface
func (ml *memcacheLock) Release() (bool, error) {
	currentOwner, err := ml.GetCurrentOwner()
	if err != nil {
		return false, err
	}
	if currentOwner == ml.owner {
		return true, ml.client.Delete(ml.name)
	}

	return false, nil
}

// ForceRelease implementation of the Lock interface
func (ml *memcacheLock) ForceRelease() error {
	return ml.client.Delete(ml.name)
}

// GetCurrentOwner implementation of the Lock interface
func (ml *memcacheLock) GetCurrentOwner() (string, error) {
	item, err := ml.client.Get(ml.name)
	if errors.Is(err, memcache.ErrCacheMiss) {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return string(item.Value), nil
}

// Expire implementation of the Lock interface
func (ml *memcacheLock) Expire(duration time.Duration) (bool, error) {
	if duration.Seconds() == 0 {
		if err := ml.ForceRelease(); err != nil {
			return false, err
		}

		return true, nil
	}
	if err := ml.client.Touch(ml.name, int32(duration.Seconds())); err != nil && errors.Is(err, memcache.ErrCacheMiss) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (ml *memcacheLock) acquire(duration time.Duration) (bool, error) {
	if err := ml.client.Add(&memcache.Item{
		Key:        ml.name,
		Value:      []byte(ml.owner),
		Expiration: int32(duration.Seconds()),
	}); errors.Is(err, memcache.ErrNotStored) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (ml *memcacheLock) initBaseLock() *memcacheLock {
	ml.lock = ml

	return ml
}
