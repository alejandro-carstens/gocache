package gocache

import "github.com/patrickmn/go-cache"

type localConnector struct{}

func (c *localConnector) connect(config *Config) (Cache, error) {
	return &LocalStore{
		c:                 cache.New(config.Local.DefaultExpiration, config.Local.DefaultInterval),
		defaultExpiration: config.Local.DefaultExpiration,
		defaultInterval:   config.Local.DefaultInterval,
		prefix:            config.Local.Prefix,
	}, nil
}

func (c *localConnector) validate(_ *Config) error {
	return nil
}
