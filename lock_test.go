package gocache

import "testing"

func TestLock(t *testing.T) {
	for _, driver := range drivers {
		cache := createStore(driver)

		got, err := cache.Lock("test", "test", 10).Acquire()
		if !got {
			t.Error("Expected to acquire lock", got)
		}

		got, err = cache.Lock("test", "test", 10).Acquire()
		if got {
			t.Error("Expected to not acquire lock", got)
		}
		if err != nil {
			t.Fatal(err)
		}

		user, err := cache.Lock("test", "test", 10).GetCurrentOwner()
		if user != "test" {
			t.Error("Expected to not acquire lock", user)
		}
		if err != nil {
			t.Fatal(err)
		}

		got, err = cache.Lock("test", "test", 10).Release()
		if got {
			t.Error("Expected to release lock", got)
		}

		got, err = cache.Lock("test", "test", 10).Acquire()
		if !got {
			t.Error("Expected to acquire lock", got)
		}
		if err != nil {
			t.Fatal(err)
		}
		if err = cache.Lock("test", "test", 10).ForceRelease(); err != nil {
			t.Fatal(err)
		}
		if _, err := cache.Flush(); err != nil {
			t.Fatal(err)
		}
	}
}
