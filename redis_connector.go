package gocache

import (
	"errors"

	"github.com/go-redis/redis"
)

// redisConnector is the representation of the redis store connector
type redisConnector struct{}

// connect is responsible for connecting with the caching store
func (rc *redisConnector) connect(params map[string]interface{}) (Cache, error) {
	params, err := rc.validate(params)
	if err != nil {
		return nil, err
	}

	return &RedisStore{
		client: rc.client(params["address"].(string), params["password"].(string), params["database"].(int)),
		prefix: params["prefix"].(string),
	}, nil
}

func (rc *redisConnector) client(address string, password string, database int) redis.Client {
	return *redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       database,
	})
}

func (rc *redisConnector) validate(params map[string]interface{}) (map[string]interface{}, error) {
	if _, ok := params["address"]; !ok {
		return params, errors.New("you need to specify an address for your redis server. Ex: localhost:6379")
	}
	if _, ok := params["database"]; !ok {
		return params, errors.New("you need to specify a database for your redis server. From 1 to 16 0-indexed")
	}
	if _, ok := params["password"]; !ok {
		return params, errors.New("you need to specify a password for your redis server")
	}
	if _, ok := params["prefix"]; !ok {
		return params, errors.New("you need to specify a caching prefix")
	}

	return params, nil
}
