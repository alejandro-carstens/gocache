package gocache

import (
	"errors"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type localLock struct {
	c       *cache.Cache
	name    string
	owner   string
	seconds int64
}

// Acquire implementation of the Lock interface
func (l *localLock) Acquire() (bool, error) {
	err := l.c.Add(l.name, l.owner, time.Duration(l.seconds)*time.Second)
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
