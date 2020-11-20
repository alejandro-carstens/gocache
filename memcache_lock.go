package gocache

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type memcacheLock struct {
	client  *memcache.Client
	name    string
	owner   string
	seconds int64
}

// Acquire implementation of the Lock interface
func (ml *memcacheLock) Acquire() (bool, error) {
	err := ml.client.Add(&memcache.Item{
		Key:        ml.name,
		Value:      []byte(ml.owner),
		Expiration: int32(ml.seconds),
	})
	if err == memcache.ErrNotStored {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
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
	if err == memcache.ErrCacheMiss {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return string(item.Value), nil
}
