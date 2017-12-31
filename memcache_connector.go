package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheConnector struct{}

// To return StoreInterface
func (this *MemcacheConnector) Connect(params map[string]interface{}) *MemcacheStore {
	params = this.validate(params)

	prefix := params["prefix"].(string)

	delete(params, "prefix")

	return &MemcacheStore{
		Client: this.client(params),
		Prefix: prefix,
	}
}

func (this *MemcacheConnector) client(params map[string]interface{}) memcache.Client {
	servers := make([]string, len(params)-1)

	for _, param := range params {
		servers = append(servers, param.(string))
	}

	return *memcache.New(servers...)
}

func (this *MemcacheConnector) validate(params map[string]interface{}) map[string]interface{} {
	if _, ok := params["prefix"]; !ok {
		panic("You need to specify a caching prefix.")
	}

	for key, param := range params {
		if _, ok := param.(string); !ok {
			panic("The" + key + "parameter is not of type string.")
		}
	}

	return params
}
