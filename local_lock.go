package gocache

import (
	"errors"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

var _ Lock = &localLock{}

func newLocalLock(client *cache.Cache, name, owner string, duration time.Duration) *localLock {
	return (&localLock{
		c:        client,
		name:     name,
		owner:    owner,
		duration: duration,
	}).initBaseLock()
}

type localLock struct {
	baseLock
	c        *cache.Cache
	name     string
	owner    string
	duration time.Duration
}

// Acquire implementation of the Lock interface
func (l *localLock) Acquire() (bool, error) {
	err := l.c.Add(l.name, l.owner, l.duration)
	if err != nil && err.Error() == fmt.Sprintf("Item %s already exists", l.name) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// Release implementation of the Lock interface
func (l *localLock) Release() (bool, error) {
	currentOwner, err := l.GetCurrentOwner()
	if err != nil {
		return false, err
	}
	if currentOwner == l.owner {
		l.c.Delete(l.name)

		return true, nil
	}

	return false, nil
}

// ForceRelease implementation of the Lock interface
func (l *localLock) ForceRelease() error {
	l.c.Delete(l.name)

	return nil
}

// GetCurrentOwner implementation of the Lock interface
func (l *localLock) GetCurrentOwner() (string, error) {
	value, valid := l.c.Get(l.name)
	if !valid {
		return "", nil
	}

	owner, valid := value.(string)
	if !valid {
		return "", errors.New("owner is not of type string")
	}

	return owner, nil
}

// Expire implementation of the Lock interface
func (l *localLock) Expire(time.Duration) (bool, error) {
	return false, ErrNotImplemented
}

func (l *localLock) initBaseLock() *localLock {
	l.lock = l

	return l
}
