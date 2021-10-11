package gocache

import (
	"errors"

	"github.com/go-redis/redis"
)

var _ cacheConnector = &redisConnector{}

// redisConnector is the representation of the redis store connector
type redisConnector struct{}

// connect is responsible for connecting with the caching store
func (rc *redisConnector) connect(config *Config) (Cache, error) {
	return &RedisStore{
		client: redis.NewClient(&redis.Options{
			Network:            config.Redis.Network,
			Addr:               config.Redis.Addr,
			Dialer:             config.Redis.Dialer,
			OnConnect:          config.Redis.OnConnect,
			Password:           config.Redis.Password,
			DB:                 config.Redis.DB,
			MaxRetries:         config.Redis.MaxRetries,
			MinRetryBackoff:    config.Redis.MinRetryBackoff,
			MaxRetryBackoff:    config.Redis.MaxRetryBackoff,
			DialTimeout:        config.Redis.DialTimeout,
			ReadTimeout:        config.Redis.ReadTimeout,
			WriteTimeout:       config.Redis.WriteTimeout,
			PoolSize:           config.Redis.PoolSize,
			MinIdleConns:       config.Redis.MinIdleConns,
			MaxConnAge:         config.Redis.MaxConnAge,
			PoolTimeout:        config.Redis.PoolTimeout,
			IdleTimeout:        config.Redis.IdleTimeout,
			IdleCheckFrequency: config.Redis.IdleCheckFrequency,
			TLSConfig:          config.Redis.TLSConfig,
		}),
		prefix: prefix{val: config.Redis.Prefix},
	}, nil
}

func (rc *redisConnector) validate(config *Config) error {
	if config.Redis == nil {
		return errors.New("a redis config needs to be specified")
	}
	if len(config.Redis.Addr) == 0 {
		return errors.New("a redis address needs to be specified")
	}

	return nil
}
