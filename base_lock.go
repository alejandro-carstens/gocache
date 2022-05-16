package gocache

import (
	"fmt"
	"time"
)

type baseLock struct {
	lock
}

// Get attempts to acquire a lock. If acquired, fn will be invoked and the lock will be safely release once
// the invocation either succeeds or errors
func (l *baseLock) Get(fn func() error) (acquired bool, err error) {
	acquired, err = l.Acquire()
	if err != nil {
		return false, err
	}
	if !acquired {
		return acquired, nil
	}
	defer func() {
		if _, releaseErr := l.Release(); releaseErr == nil {
			return
		} else if err != nil {
			err = fmt.Errorf("gocache: %v: %v", err.Error(), releaseErr)
		} else {
			err = releaseErr
		}
	}()

	err = fn()
	return acquired, err
}

// Block will attempt to acquire a lock for the specified "wait" time. If acquired, fn will be invoked and
// the lock will be safely release once the invocation either succeeds or errors. The interval variable
// will be used as the wait duration between attempts to acquire the lock
func (l *baseLock) Block(interval, wait time.Duration, fn func() error) (acquired bool, err error) {
	starting := time.Now()
	for {
		acquired, err = l.Acquire()
		if err != nil {
			return false, err
		}
		if acquired {
			break
		}

		time.Sleep(interval)
		if time.Now().Add(-wait).After(starting) {
			return false, ErrBlockWaitTimeout
		}
	}

	if !acquired {
		return acquired, nil
	}
	defer func() {
		if _, releaseErr := l.Release(); releaseErr == nil {
			return
		} else if err != nil {
			err = fmt.Errorf("gocache: %v: %v", err.Error(), releaseErr)
		} else {
			err = releaseErr
		}
	}()

	err = fn()
	return acquired, err
}
