package gocache

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

const (
	deleteLimit = 1000
	redisOk     = "OK"
)

var _ Cache = &RedisStore{}

// NewRedisStore validates the passed in config and creates a Cache implementation of type *RedisStore
func NewRedisStore(cnf *RedisConfig, encoder Encoder) (*RedisStore, error) {
	if err := cnf.validate(); err != nil {
		return nil, err
	}
	return &RedisStore{
		prefix: prefix{
			val: cnf.Prefix,
		},
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
		encoder: encoder,
	}, nil
}

// RedisStore is the representation of the redis caching store
type RedisStore struct {
	prefix
	client  *redis.Client
	encoder Encoder
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

// GetBool gets a bool value from the store
func (s *RedisStore) GetBool(key string) (bool, error) {
	value, err := s.get(key).Result()
	if err != nil {
		return false, checkErrNotFound(err)
	}
	if isStringNumeric(value) || isStringBool(value) {
		return stringToBool(value), nil
	}
	if err = s.encoder.Decode([]byte(value), &value); err != nil {
		return false, err
	}

	return stringToBool(value), nil
}

// GetString gets a string value from the store
func (s *RedisStore) GetString(key string) (string, error) {
	value, err := s.get(key).Result()
	if err != nil {
		return "", checkErrNotFound(err)
	}
	if err = s.encoder.Decode([]byte(value), &value); err != nil {
		return "", err
	}

	return value, nil
}

// Increment increments an integer counter by a given val
func (s *RedisStore) Increment(key string, value int64) (int64, error) {
	return s.client.IncrBy(s.k(key), value).Result()
}

// Decrement decrements an integer counter by a given val
func (s *RedisStore) Decrement(key string, value int64) (int64, error) {
	return s.client.DecrBy(s.k(key), value).Result()
}

// Put puts a value in the given store for a predetermined amount of time in seconds
func (s *RedisStore) Put(key string, value interface{}, duration time.Duration) error {
	if isNumeric(value) || isBool(value) {
		return s.client.Set(s.k(key), value, duration).Err()
	}

	val, err := s.encoder.Encode(value)
	if err != nil {
		return err
	}

	return s.client.Set(s.k(key), val, duration).Err()
}

// Add an item to the cache only if an item doesn't already exist for the given key, or if the existing item has
// expired. If the record was successfully added true will be returned else false will be returned
func (s *RedisStore) Add(key string, value interface{}, duration time.Duration) (bool, error) {
	if isNumeric(value) || isBool(value) {
		res, err := s.client.Eval(redisLuaAddScript, []string{s.k(key)}, value, duration.Seconds()).String()
		if err != nil && !errors.Is(err, redis.Nil) {
			return false, err
		}

		return res == redisOk, nil
	}

	val, err := s.encoder.Encode(value)
	if err != nil {
		return false, err
	}

	res, err := s.client.Eval(redisLuaAddScript, []string{s.k(key)}, val, duration.Seconds()).String()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, err
	}

	return res == redisOk, nil
}

// Forever puts a value in the given store until it is forgotten/evicted
func (s *RedisStore) Forever(key string, value interface{}) error {
	if isNumeric(value) || isBool(value) {
		if err := s.client.Set(s.k(key), value, 0).Err(); err != nil {
			return err
		}

		return s.client.Persist(s.k(key)).Err()
	}

	val, err := s.encoder.Encode(value)
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
func (s *RedisStore) Forget(keys ...string) (bool, error) {
	if len(keys) == 0 {
		return true, nil
	}

	var delKeys []string
	for _, key := range keys {
		delKeys = append(delKeys, s.k(key))
		if len(delKeys) < deleteLimit {
			continue
		}
		if err := s.client.Del(delKeys...).Err(); err != nil {
			return false, checkErrNotFound(err)
		}
	}

	if len(delKeys) == 0 {
		return true, nil
	}
	if err := s.client.Del(delKeys...).Err(); err != nil {
		return false, checkErrNotFound(err)
	}

	return true, nil
}

// PutMany puts many values in the given store until they are forgotten/evicted
func (s *RedisStore) PutMany(entries ...Entry) error {
	pipe := s.client.TxPipeline()
	for _, entry := range entries {
		if err := s.Put(entry.Key, entry.Value, entry.Duration); err != nil {
			return err
		}
	}

	_, err := pipe.Exec()

	return err
}

// Many gets many values from the store
func (s *RedisStore) Many(keys ...string) (Items, error) {
	var (
		items = Items{}
		pipe  = s.client.TxPipeline()
	)
	for _, key := range keys {
		val, err := s.get(key).Result()
		if err != nil {
			items[key] = Item{
				key:     key,
				err:     checkErrNotFound(err),
				encoder: s.encoder,
			}

			continue
		}
		if isStringNumeric(val) || isStringBool(val) {
			items[key] = Item{
				key:     key,
				value:   val,
				encoder: s.encoder,
			}

			continue
		}

		items[key] = Item{
			key:     key,
			value:   val,
			encoder: s.encoder,
		}
	}

	_, err := pipe.Exec()

	return items, err
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
	value, err := s.get(key).Bytes()
	if err != nil {
		return checkErrNotFound(err)
	}

	return s.encoder.Decode(value, entity)
}

// Lock returns a redis implementation of the Lock interface
func (s *RedisStore) Lock(name, owner string, duration time.Duration) Lock {
	return &redisLock{
		client:   s.client,
		name:     name,
		owner:    owner,
		duration: duration,
	}
}

// Exists checks if an entry exists in the cache for the given key
func (s *RedisStore) Exists(key string) (bool, error) {
	_, err := s.get(key).Result()
	if err == nil {
		return true, nil
	} else if err != nil && isErrNotFound(err) {
		return false, nil
	}

	return false, err
}

// Lpush runs the Redis lpush command (used via reflection, do not delete)
func (s *RedisStore) Lpush(segment, key string) error {
	return s.client.LPush(segment, key).Err()
}

// Lrange runs the Redis lrange command (used via reflection, do not delete)
func (s *RedisStore) Lrange(key string, start, stop int64) []string {
	return s.client.LRange(key, start, stop).Val()
}

func (s *RedisStore) get(key string) *redis.StringCmd {
	return s.client.Get(s.k(key))
}
