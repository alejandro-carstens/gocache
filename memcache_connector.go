package gocache

import (
	"errors"

	"github.com/bradfitz/gomemcache/memcache"
)

var _ cacheConnector = &memcacheConnector{}

type memcacheConnector struct{}

func (mc *memcacheConnector) connect(config config) (Cache, error) {
	cnf, valid := config.(*MemcacheConfig)
	if !valid {
		return nil, errors.New("config is not of type *RedisConfig")
	}

	client := memcache.New(cnf.Servers...)
	if cnf.MaxIdleConns > 0 {
		client.MaxIdleConns = cnf.MaxIdleConns
	}

	client.Timeout = cnf.Timeout

	return &MemcacheStore{
		client: client,
		prefix: prefix{
			val: cnf.Prefix,
		},
	}, nil
}

func (mc *memcacheConnector) validate(config config) error {
	cnf, valid := config.(*MemcacheConfig)
	if !valid {
		return errors.New("config is not of type *RedisConfig")
	}
	if len(cnf.Servers) == 0 {
		return errors.New("memcache.servers cannot be empty")
	}

	return nil
}
