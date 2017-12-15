package cache

import (
	"github.com/go-redis/redis"
)

type RedisConnector struct {
	Store
}

const address = "localhost:6379"
const password = ""
const default_db = 0

func (this *RedisConnector) Connect(params map[string]interface{}) RedisStore {
	params, err := this.validate(params)

	if err != nil {
		panic("Invalid parameters for redis client.")
	}

	address := params["address"].(string)
	password := params["password"].(string)
	default_db := params["default_db"].(int)
	prefix := params["prefix"].(string)

	return RedisStore{
		Client: this.client(address, password, default_db),
		Prefix: prefix,
	}
}

func (this *RedisConnector) client(address string, password string, default_db int) redis.Client {
	return *redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       default_db,
	})
}

func (this *RedisConnector) validate(params map["string"]interface{}) map["string"]interface{} {
    // TODO implement this method
} 

