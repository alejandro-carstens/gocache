package cache

import (
	"strings"
)

// REDIS_DRIVER specifies the redis driver name
const REDIS_DRIVER = "redis"

// MEMCACHE_DRIVER specifies the memcache driver name
const MEMCACHE_DRIVER = "memcache"

// ARRAY_DRIVER specifies the array driver name
const ARRAY_DRIVER = "array"

// New new-ups an instance of StoreInterface
func New(driver string, params map[string]interface{}) (StoreInterface, error) {
	switch strings.ToLower(driver) {
	case REDIS_DRIVER:
		return connect(new(RedisConnector), params)
	case MEMCACHE_DRIVER:
		return connect(new(MemcacheConnector), params)
	case ARRAY_DRIVER:
		return connect(new(ArrayConnector), params)
	}

	return connect(new(ArrayConnector), params)
}

func connect(connector CacheConnectorInterface, params map[string]interface{}) (StoreInterface, error) {
	return connector.Connect(params)
}
