package cache

import (
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheConnector struct{}

func (this *MemcacheConnector) Connect(params map[string]interface{}) (StoreInterface, error) {
	params, err := this.validate(params)

	if err != nil {
		return &MemcacheStore{}, err
	}

	prefix := params["prefix"].(string)

	delete(params, "prefix")

	return &MemcacheStore{
		Client: this.client(params),
		Prefix: prefix,
	}, nil
}

func (this *MemcacheConnector) client(params map[string]interface{}) memcache.Client {
	servers := make([]string, len(params)-1)

	for _, param := range params {
		servers = append(servers, param.(string))
	}

	return *memcache.New(servers...)
}

func (this *MemcacheConnector) validate(params map[string]interface{}) (map[string]interface{}, error) {
	if _, ok := params["prefix"]; !ok {
		return params, errors.New("You need to specify a caching prefix.")
	}

	for key, param := range params {
		if _, ok := param.(string); !ok {
			return params, errors.New("The" + key + "parameter is not of type string.")
		}
	}

	return params, nil
}
