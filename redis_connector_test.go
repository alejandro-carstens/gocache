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

	redis_connector := new(RedisConnector)

	redis_store := redis_connector.Connect(params)

	_, ok := redis_store.(StoreInterface)

	if !ok {
		t.Error("Expected StoreInterface got", redis_store)
	}
}
