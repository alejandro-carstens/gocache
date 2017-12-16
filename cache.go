package cache

import (
	"strings"
)

const REDIS_DRIVER = "redis"

func Cache(driver string, params map[string]interface{}) StoreInterface {
	switch strings.ToLower(driver) {
	case REDIS_DRIVER:
		return new(RedisConnector).Connect(params)
	default:
		panic("The provided driver was not found")
	}
}
