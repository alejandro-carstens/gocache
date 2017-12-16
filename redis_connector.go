package cache

import (
	"github.com/go-redis/redis"
)

type RedisConnector struct{}

func (this *RedisConnector) Connect(params map[string]interface{}) StoreInterface {
	params = this.validate(params)

	return &RedisStore{
		Client: this.client(params["address"].(string), params["password"].(string), params["default_db"].(int)),
		Prefix: params["prefix"].(string),
	}
}

func (this *RedisConnector) client(address string, password string, default_db int) redis.Client {
	return *redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       default_db,
	})
}

func (this *RedisConnector) validate(params map[string]interface{}) map[string]interface{} {
	if _, ok := params["address"]; !ok {
		panic("You need to specify an address for your redis server. Ex: localhost:6379")
	}

	if _, ok := params["database"]; !ok {
		panic("You need to specify a database for your redis server. From 1 to 16.")
	}

	if _, ok := params["password"]; !ok {
		panic("You need to specify a password for your redis server.")
	}

	return params
}
