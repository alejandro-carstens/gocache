package gocache

import (
	"errors"

	"github.com/bradfitz/gomemcache/memcache"
)

var _ cacheConnector = &memcacheConnector{}

type memcacheConnector struct{}

func (mc *memcacheConnector) connect(config *Config) (Cache, error) {
	client := memcache.New(config.Memcache.Servers...)
	if config.Memcache.MaxIdleConns > 0 {
		client.MaxIdleConns = config.Memcache.MaxIdleConns
	}

	client.Timeout = config.Memcache.Timeout

	return &MemcacheStore{
		client: client,
		prefix: config.Memcache.Prefix,
	}, nil
}

func (mc *memcacheConnector) validate(config *Config) error {
	if config.Memcache == nil {
		return errors.New("memcache config not specified")
	}
	if len(config.Memcache.Servers) == 0 {
		return errors.New("memcache.servers cannot be empty")
	}

	return nil
}
