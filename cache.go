package cache

import (
	"strings"
)

const REDIS_DRIVER = "redis"
const MEMCACHE_DRIVER = "memcache"
const ARRAY_DRIVER = "array"

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
