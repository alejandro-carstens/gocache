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

func (ml *mapLock) Acquire() (bool, error) {
	nameLock, exists := ml.locks[ml.key()]
	if !exists {
		nameLock = new(mapLocker)
		ml.locks[ml.name] = nameLock
	}

	nameLock.inc()
	ml.mu.Unlock()

	nameLock.lock()
	nameLock.dec()

	return true, nil
}

func (ml *mapLock) Release() (bool, error) {
	ml.mu.Lock()
	nameLock, exists := ml.locks[ml.key()]
	if !exists {
		ml.mu.Unlock()
		return true, nil
	}

	if nameLock.count() == 0 {
		delete(ml.locks, ml.key())
	}
	nameLock.unlock()

	ml.mu.Unlock()
	return true, nil
}

func (ml *mapLock) ForceRelease() error {
	_, err := ml.Release()

	return err
}

func (ml *mapLock) GetCurrentOwner() (string, error) {
	ml.mu.Lock()
	if _, exists := ml.locks[ml.key()]; !exists {
		return ml.name, nil
	}
	ml.mu.Unlock()

	return "", nil
}

func (ml *mapLock) key() string {
	return fmt.Sprintf("%v:%v", ml.owner, ml.name)
}
