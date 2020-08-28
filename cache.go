package gocache

import (
	"strings"
)

// RedisDriver specifies the redis driver name
const RedisDriver = "redis"

// MemcacheDriver specifies the memcache driver name
const MemcacheDriver = "memcache"

// MAP_DRIVER specifies the map driver name
const MapDriver = "map"

// New new-ups an instance of Store
func New(driver string, params map[string]interface{}) (Cache, error) {
	switch strings.ToLower(driver) {
	case RedisDriver:
		return connect(&redisConnector{}, params)
	case MemcacheDriver:
		return connect(&memcacheConnector{}, params)
	case MapDriver:
		break
	}

	return connect(&mapConnector{}, params)
}

func connect(connector cacheConnector, params map[string]interface{}) (Cache, error) {
	return connector.connect(params)
}
