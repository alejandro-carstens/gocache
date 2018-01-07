package cache

import (
	"testing"
)

func TestRedisConnection(t *testing.T) {
	params := make(map[string]interface{})

	params["address"] = "localhost:6379"
	params["password"] = ""
	params["database"] = 0
	params["prefix"] = "golavel:"

	redisConnector := new(RedisConnector)

	redisStore, err := redisConnector.Connect(params)

	if err != nil {
		panic(err)
	}

	_, ok := redisStore.(StoreInterface)

	if !ok {
		t.Error("Expected StoreInterface got", redisStore)
	}
}
