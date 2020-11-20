package gocache

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type mapLocker struct {
	mu      sync.Mutex
	waiters int32
}

func (l *mapLocker) inc() {
	atomic.AddInt32(&l.waiters, 1)
}

func (l *mapLocker) dec() {
	atomic.AddInt32(&l.waiters, -1)
}

func (l *mapLocker) count() int32 {
	return atomic.LoadInt32(&l.waiters)
}

func (l *mapLocker) lock() {
	l.mu.Lock()
}

func (l *mapLocker) unlock() {
	l.mu.Unlock()
}

type mapLock struct {
	mu      sync.Mutex
	locks   map[string]*mapLocker
	seconds int64
	name    string
	owner   string
}

// Acquire implementation of the Lock interface
func (ml *mapLock) Acquire() (bool, error) {
	ml.mu.Lock()
	if _, exists := ml.locks[ml.key()]; exists {
		ml.mu.Unlock()
		return false, nil
	}

	ml.locks[ml.key()] = new(mapLocker)
	ml.mu.Unlock()
	ml.locks[ml.key()].lock()

	return true, nil
}

// Release implementation of the Lock interface
func (ml *mapLock) Release() (bool, error) {
	ml.mu.Lock()
	defer ml.mu.Unlock()
	nameLock, exists := ml.locks[ml.key()]
	if !exists {
		return true, nil
	}

	delete(ml.locks, ml.key())
	nameLock.unlock()

	return true, nil
}

// ForceRelease implementation of the Lock interface
func (ml *mapLock) ForceRelease() error {
	_, err := ml.Release()

	return err
}

// GetCurrentOwner implementation of the Lock interface
func (ml *mapLock) GetCurrentOwner() (string, error) {
	ml.mu.Lock()
	defer ml.mu.Unlock()
	if _, exists := ml.locks[ml.key()]; exists {
		return ml.owner, nil
	}

	return "", nil
}

func (ml *mapLock) key() string {
	return fmt.Sprintf("%v:%v", ml.owner, ml.name)
}
