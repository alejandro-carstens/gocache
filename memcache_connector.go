package gocache

import (
	"errors"

	"github.com/bradfitz/gomemcache/memcache"
)

// memcacheConnector is the representation of the memcache store connector
type memcacheConnector struct{}

// connect is responsible for connecting with the caching store
func (mc *memcacheConnector) connect(params map[string]interface{}) (Cache, error) {
	params, err := mc.validate(params)
	if err != nil {
		return nil, err
	}

	prefix := params["prefix"].(string)

	delete(params, "prefix")

	return &MemcacheStore{
		client: mc.client(params),
		prefix: prefix,
	}, nil
}

func (mc *memcacheConnector) client(params map[string]interface{}) memcache.Client {
	servers := make([]string, len(params)-1)
	for _, param := range params {
		servers = append(servers, param.(string))
	}

	return *memcache.New(servers...)
}

func (mc *memcacheConnector) validate(params map[string]interface{}) (map[string]interface{}, error) {
	if _, ok := params["prefix"]; !ok {
		return params, errors.New("you need to specify a caching prefix")
	}

	for key, param := range params {
		if _, ok := param.(string); !ok {
			return params, errors.New("the" + key + "parameter is not of type string")
		}
	}

	return params, nil
}
