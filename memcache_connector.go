package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheConnector struct{}

// To return StoreInterface
func (this *MemcacheConnector) Connect(params map[string]interface{}) memcache.Client {
	this.validate(params)

	return this.client(params)
}

func (this *MemcacheConnector) client(params map[string]interface{}) memcache.Client {
	servers := make([]string, len(params))

	for _, param := range params {
		servers = append(servers, param.(string))
	}

	return *memcache.New(servers...)
}

func (this *MemcacheConnector) validate(params map[string]interface{}) map[string]interface{} {
	for key, param := range params {
		if _, ok := param.(string); !ok {
			panic("The" + key + "parameter is not of type string.")
		}
	}

	return params
}
