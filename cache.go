package cache

import (
	"strings"
)

const REDIS_DRIVER = "redis"

func Cache(driver string, params map[string]interface{}) StoreInterface {
	switch strings.ToLower(driver) {
	case REDIS_DRIVER:
		return connect(new(RedisConnector), params)
	default:
		panic("The provided driver was not found")
	}
}

func connect(connector CacheConnectorInterface, params map[string]interface{}) StoreInterface {
	return connector.Connect(params)
}
