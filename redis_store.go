package gocache

import (
	"time"

	"github.com/go-redis/redis"
)

var _ Cache = &RedisStore{}

// RedisStore is the representation of the redis caching store
type RedisStore struct {
	prefix
	client *redis.Client
}

// GetFloat64 gets a float64 value from the store
func (s *RedisStore) GetFloat64(key string) (float64, error) {
	res, err := s.get(key).Float64()

	return res, checkErrNotFound(err)
}

// GetFloat32 gets a float32 value from the store
func (s *RedisStore) GetFloat32(key string) (float32, error) {
	res, err := s.get(key).Float32()

	return res, checkErrNotFound(err)
}

// GetInt64 gets an int64 value from the store
func (s *RedisStore) GetInt64(key string) (int64, error) {
	res, err := s.get(key).Int64()

	return res, checkErrNotFound(err)
}

// GetInt gets an int value from the store
func (s *RedisStore) GetInt(key string) (int, error) {
	res, err := s.get(key).Int()

	return res, checkErrNotFound(err)
}

// GetUint64 gets an uint64 value from the store
func (s *RedisStore) GetUint64(key string) (uint64, error) {
	res, err := s.get(key).Uint64()

	return res, checkErrNotFound(err)
}

// GetString gets a string value from the store
func (s *RedisStore) GetString(key string) (string, error) {
	value, err := s.get(key).Result()
	if err != nil {
		return "", checkErrNotFound(err)
	}

	return simpleDecode(value)
}

// Increment increments an integer counter by a given value
func (s *RedisStore) Increment(key string, value int64) (int64, error) {
	return s.client.IncrBy(s.k(key), value).Result()
}

// Decrement decrements an integer counter by a given value
func (s *RedisStore) Decrement(key string, value int64) (int64, error) {
	return s.client.DecrBy(s.k(key), value).Result()
}

// Put puts a value in the given store for a predetermined amount of time in seconds
func (s *RedisStore) Put(key string, value interface{}, seconds int) error {
	duration := time.Duration(int64(seconds)) * time.Second
	if isNumeric(value) {
		return s.client.Set(s.k(key), value, duration).Err()
	}

	val, err := encode(value)
	if err != nil {
		return err
	}

	return s.client.Set(s.k(key), val, duration).Err()
}

// Forever puts a value in the given store until it is forgotten/evicted
func (s *RedisStore) Forever(key string, value interface{}) error {
	if isNumeric(value) {
		if err := s.client.Set(s.k(key), value, 0).Err(); err != nil {
			return err
		}

		return s.client.Persist(s.k(key)).Err()
	}

	val, err := encode(value)
	if err != nil {
		return err
	}
	if err = s.client.Set(s.k(key), val, 0).Err(); err != nil {
		return err
	}

	return s.client.Persist(s.k(key)).Err()
}

// Flush flushes the store
func (s *RedisStore) Flush() (bool, error) {
	if err := s.client.FlushDB().Err(); err != nil {
		return false, err
	}

	return true, nil
}

// Forget forgets/evicts a given key-value pair from the store
func (s *RedisStore) Forget(key string) (bool, error) {
	if err := s.client.Del(s.k(key)).Err(); err != nil {
		return false, checkErrNotFound(err)
	}

	return true, nil
}

// PutMany puts many values in the given store until they are forgotten/evicted
func (s *RedisStore) PutMany(values map[string]string, seconds int) error {
	pipe := s.client.TxPipeline()

	for key, value := range values {
		if err := s.Put(key, value, seconds); err != nil {
			return err
		}
	}

	_, err := pipe.Exec()

	return err
}

// Many gets many values from the store
func (s *RedisStore) Many(keys []string) (map[string]string, error) {
	pipe := s.client.TxPipeline()

	values := map[string]string{}
	for _, key := range keys {
		val, err := s.GetString(key)
		if err != nil {
			return values, err
		}

		values[key] = val
	}

	_, err := pipe.Exec()

	return values, err
}

// Connection returns the the store's c
func (s *RedisStore) Connection() interface{} {
	return s.client
}

// Tags returns the taggedCache for the given store
func (s *RedisStore) Tags(names ...string) TaggedCache {
	return &redisTaggedCache{
		taggedCache{
			store: s,
			tags: tagSet{
				store: s,
				names: names,
			},
		},
	}
}

// Close closes the c releasing all open resources
func (s *RedisStore) Close() error {
	return s.client.Close()
}

// Get gets the struct representation of a value from the store
func (s *RedisStore) Get(key string, entity interface{}) error {
	value, err := s.get(key).Result()
	if err != nil {
		return checkErrNotFound(err)
	}

	_, err = decode(value, entity)

	return err
}

// Lock returns a redis implementation of the Lock interface
func (s *RedisStore) Lock(name, owner string, seconds int64) Lock {
	return &redisLock{
		client:  s.client,
		seconds: seconds,
		name:    name,
		owner:   owner,
	}
}

// Lpush runs the Redis lpush command (used via reflection, do not delete)
func (s *RedisStore) Lpush(segment string, key string) {
	s.client.LPush(segment, s.k(key))
}

// Lrange runs the Redis lrange command (used via reflection, do not delete)
func (s *RedisStore) Lrange(key string, start int64, stop int64) []string {
	return s.client.LRange(s.k(key), start, stop).Val()
}

func (s *RedisStore) get(key string) *redis.StringCmd {
	return s.client.Get(s.k(key))
}
