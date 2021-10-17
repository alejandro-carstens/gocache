package gocache

import (
	"errors"

	"github.com/go-redis/redis"
)

var _ cacheConnector = &redisConnector{}

// redisConnector is the representation of the redis store connector
type redisConnector struct{}

// connect is responsible for connecting with the caching store
func (rc *redisConnector) connect(config config) (Cache, error) {
	cnf, valid := config.(*RedisConfig)
	if !valid {
		return nil, errors.New("config is not of type *RedisConfig")
	}
	return &RedisStore{
		client: redis.NewClient(&redis.Options{
			Network:            cnf.Network,
			Addr:               cnf.Addr,
			Dialer:             cnf.Dialer,
			OnConnect:          cnf.OnConnect,
			Password:           cnf.Password,
			DB:                 cnf.DB,
			MaxRetries:         cnf.MaxRetries,
			MinRetryBackoff:    cnf.MinRetryBackoff,
			MaxRetryBackoff:    cnf.MaxRetryBackoff,
			DialTimeout:        cnf.DialTimeout,
			ReadTimeout:        cnf.ReadTimeout,
			WriteTimeout:       cnf.WriteTimeout,
			PoolSize:           cnf.PoolSize,
			MinIdleConns:       cnf.MinIdleConns,
			MaxConnAge:         cnf.MaxConnAge,
			PoolTimeout:        cnf.PoolTimeout,
			IdleTimeout:        cnf.IdleTimeout,
			IdleCheckFrequency: cnf.IdleCheckFrequency,
			TLSConfig:          cnf.TLSConfig,
		}),
		prefix: prefix{
			val: cnf.Prefix,
		},
	}, nil
}

func (rc *redisConnector) validate(config config) error {
	cnf, valid := config.(*RedisConfig)
	if !valid {
		return errors.New("config is not of type *RedisConfig")
	}
	if len(cnf.Addr) == 0 {
		return errors.New("a redis address needs to be specified")
	}

	return nil
}
