package gocache

import (
	"time"

	"github.com/go-redis/redis"
)

const redisNilErrorResponse string = "redis: nil"

// RedisStore is the representation of the redis caching store
type RedisStore struct {
	client *redis.Client
	prefix string
}

// GetFloat64 gets a float value from the store
func (rs *RedisStore) GetFloat64(key string) (float64, error) {
	return rs.get(key).Float64()
}

// GetInt64 gets an int value from the store
func (rs *RedisStore) GetInt64(key string) (int64, error) {
	return rs.get(key).Int64()
}

// GetString gets a string value from the store
func (rs *RedisStore) GetString(key string) (string, error) {
	value, err := rs.get(key).Result()
	if err != nil {
		return "", err
	}

	return simpleDecode(value)
}

// GetRawString gets a raw string value from the store
func (rs *RedisStore) GetRawString(key string) (string, error) {
	value, err := rs.get(key).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

// Increment increments an integer counter by a given value
func (rs *RedisStore) Increment(key string, value int64) (int64, error) {
	return rs.client.IncrBy(rs.prefix+key, value).Result()
}

// Decrement decrements an integer counter by a given value
func (rs *RedisStore) Decrement(key string, value int64) (int64, error) {
	return rs.client.DecrBy(rs.prefix+key, value).Result()
}

// Put puts a value in the given store for a predetermined amount of time in seconds
func (rs *RedisStore) Put(key string, value interface{}, seconds int) error {
	duration := time.Duration(int64(seconds)) * time.Second
	if isNumeric(value) {
		return rs.client.Set(rs.GetPrefix()+key, value, duration).Err()
	}

	val, err := encode(value)
	if err != nil {
		return err
	}

	return rs.client.Set(rs.GetPrefix()+key, val, duration).Err()
}

// PutRawString puts a raw string value in the given store for a predetermined amount of time in seconds
func (rs *RedisStore) PutRawString(key, value string, seconds int) error {
	duration := time.Duration(int64(seconds)) * time.Second

	return rs.client.Set(rs.GetPrefix()+key, value, duration).Err()
}

// Forever puts a value in the given store until it is forgotten/evicted
func (rs *RedisStore) Forever(key string, value interface{}) error {
	if isNumeric(value) {
		if err := rs.client.Set(rs.GetPrefix()+key, value, 0).Err(); err != nil {
			return err
		}

		return rs.client.Persist(rs.GetPrefix() + key).Err()
	}

	val, err := encode(value)
	if err != nil {
		return err
	}
	if err = rs.client.Set(rs.GetPrefix()+key, val, 0).Err(); err != nil {
		return err
	}

	return rs.client.Persist(rs.GetPrefix() + key).Err()
}

// Flush flushes the store
func (rs *RedisStore) Flush() (bool, error) {
	if err := rs.client.FlushDB().Err(); err != nil {
		return false, err
	}

	return true, nil
}

// Forget forgets/evicts a given key-value pair from the store
func (rs *RedisStore) Forget(key string) (bool, error) {
	if err := rs.client.Del(rs.prefix + key).Err(); err != nil {
		return false, err
	}

	return true, nil
}

// GetPrefix gets the cache key prefix
func (rs *RedisStore) GetPrefix() string {
	return rs.prefix
}

// PutMany puts many values in the given store until they are forgotten/evicted
func (rs *RedisStore) PutMany(values map[string]string, seconds int) error {
	pipe := rs.client.TxPipeline()

	for key, value := range values {
		if err := rs.Put(key, value, seconds); err != nil {
			return err
		}
	}

	_, err := pipe.Exec()

	return err
}

// Many gets many values from the store
func (rs *RedisStore) Many(keys []string) (map[string]string, error) {
	pipe := rs.client.TxPipeline()

	values := map[string]string{}
	for _, key := range keys {
		val, err := rs.GetString(key)
		if err != nil {
			return values, err
		}

		values[key] = val
	}

	_, err := pipe.Exec()

	return values, err
}

// Connection returns the the store's client
func (rs *RedisStore) Connection() interface{} {
	return rs.client
}

// Tags returns the taggedCache for the given store
func (rs *RedisStore) Tags(names ...string) TaggedCache {
	return &redisTaggedCache{
		taggedCache{
			store: rs,
			tags: tagSet{
				store: rs,
				names: names,
			},
		},
	}
}

// Close closes the client releasing all open resources
func (rs *RedisStore) Close() error {
	return rs.client.Close()
}

// Get gets the struct representation of a value from the store
func (rs *RedisStore) Get(key string, entity interface{}) error {
	value, err := rs.get(key).Result()
	if err != nil {
		return err
	}
	_, err = decode(value, entity)

	return err
}

// Lock returns a redis implementation of the Lock interface
func (rs *RedisStore) Lock(name, owner string, seconds int64) Lock {
	return &redisLock{
		client:  rs.client,
		seconds: seconds,
		name:    name,
		owner:   owner,
	}
}

// Lpush runs the Redis lpush command (used via reflection, do not delete)
func (rs *RedisStore) Lpush(segment string, key string) {
	rs.client.LPush(segment, key)
}

// Lrange runs the Redis lrange command (used via reflection, do not delete)
func (rs *RedisStore) Lrange(key string, start int64, stop int64) []string {
	return rs.client.LRange(key, start, stop).Val()
}

func (rs *RedisStore) get(key string) *redis.StringCmd {
	return rs.client.Get(rs.GetPrefix() + key)
}
