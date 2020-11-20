package gocache

import (
	"testing"
)

func TestLock(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		lock := cache.Lock("test", "test", 10)
		got, err := lock.Acquire()
		if !got {
			t.Error("Expected to acquire lock", got)
		}

		if driver == mapDriver {
			got, err = lock.Acquire()
		} else {
			got, err = cache.Lock("test", "test", 10).Acquire()
		}
		if got {
			t.Error("Expected to not acquire lock", got)
		}
		if err != nil {
			t.Fatal(err)
		}

		var user string
		if driver == mapDriver {
			user, err = lock.GetCurrentOwner()
		} else {
			user, err = cache.Lock("test", "test", 10).GetCurrentOwner()
		}
		if user != "test" {
			t.Error("Expected to not acquire lock", user)
		}
		if err != nil {
			t.Fatal(err)
		}
		if driver == mapDriver {
			got, err = lock.Release()
		} else {
			got, err = cache.Lock("test", "test", 10).Release()
		}
		if !got {
			t.Error("Expected to release lock", got)
		}

		lock = cache.Lock("test", "test", 10)
		got, err = lock.Acquire()
		if !got {
			t.Error("Expected to acquire lock", got)
		}
		if err != nil {
			t.Fatal(err)
		}
		if driver == mapDriver {
			err = lock.ForceRelease()
		} else {
			err = cache.Lock("test", "test", 10).ForceRelease()
		}
		if err != nil {
			t.Fatal(err)
		}
		if _, err := cache.Flush(); err != nil {
			t.Fatal(err)
		}
	}
}
