package cache

import (
	"errors"
	"github.com/go-redis/redis"
)

type RedisConnector struct{}

func (this *RedisConnector) Connect(params map[string]interface{}) (StoreInterface, error) {
	params, err := this.validate(params)

	if err != nil {
		return &RedisStore{}, err
	}

	return &RedisStore{
		Client: this.client(params["address"].(string), params["password"].(string), params["database"].(int)),
		Prefix: params["prefix"].(string),
	}, nil
}

func (this *RedisConnector) client(address string, password string, database int) redis.Client {
	return *redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       database,
	})
}

func (this *RedisConnector) validate(params map[string]interface{}) (map[string]interface{}, error) {
	if _, ok := params["address"]; !ok {
		return params, errors.New("You need to specify an address for your redis server. Ex: localhost:6379")
	}

	if _, ok := params["database"]; !ok {
		return params, errors.New("You need to specify a database for your redis server. From 1 to 16 0-indexed")
	}

	if _, ok := params["password"]; !ok {
		return params, errors.New("You need to specify a password for your redis server.")
	}

	if _, ok := params["prefix"]; !ok {
		return params, errors.New("You need to specify a caching prefix.")
	}

	return params, nil
}
