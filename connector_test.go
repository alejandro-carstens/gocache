package cache

import (
	"testing"
)

func TestMemcacheConnector(t *testing.T) {
	memcacheConnector := new(MemcacheConnector)

	memcacheStore, err := memcacheConnector.Connect(memcacheStore())

	if err != nil {
		panic(err)
	}

	_, ok := memcacheStore.(StoreInterface)

	if !ok {
		t.Error("Expected StoreInterface got", memcacheStore)
	}
}

func TestRedisConnector(t *testing.T) {
	redisConnector := new(RedisConnector)

	redisStore, err := redisConnector.Connect(redisStore())

	if err != nil {
		panic(err)
	}

	_, ok := redisStore.(StoreInterface)

	if !ok {
		t.Error("Expected StoreInterface got", redisStore)
	}
}
