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
		return connect(&RedisConnector{}, params)
	case MemcacheDriver:
		return connect(&MemcacheConnector{}, params)
	case MapDriver:
		return connect(&MapConnector{}, params)
	}

	return connect(&MapConnector{}, params)
}

func connect(connector CacheConnector, params map[string]interface{}) (Cache, error) {
	return connector.Connect(params)
}
