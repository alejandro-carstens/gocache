package gocache

import (
	"errors"

	"github.com/patrickmn/go-cache"
)

var _ cacheConnector = &localConnector{}

type localConnector struct{}

func (c *localConnector) connect(config config) (Cache, error) {
	cnf, valid := config.(*LocalConfig)
	if !valid {
		return nil, errors.New("config is not of type *RedisConfig")
	}

	return &LocalStore{
		c:                 cache.New(cnf.DefaultExpiration, cnf.DefaultInterval),
		defaultExpiration: cnf.DefaultExpiration,
		defaultInterval:   cnf.DefaultInterval,
		prefix: prefix{
			val: cnf.Prefix,
		},
	}, nil
}

func (c *localConnector) validate(config config) error {
	if _, valid := config.(*LocalConfig); !valid {
		return errors.New("config is not of type *RedisConfig")
	}

	return nil
}
