package cache

import (
	"strings"
)

const REDIS_DRIVER = "redis"
const MEMCACHE_DRIVER = "memcache"

func Cache(driver string, params map[string]interface{}) (StoreInterface, error) {
	switch strings.ToLower(driver) {
	case REDIS_DRIVER:
		return connect(new(RedisConnector), params)
	case MEMCACHE_DRIVER:
		return connect(new(MemcacheConnector), params)
	default:
		panic("The provided driver was not found.")
	}
}

func connect(connector CacheConnectorInterface, params map[string]interface{}) (StoreInterface, error) {
	return connector.Connect(params)
}
