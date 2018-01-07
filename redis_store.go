package cache

import (
	"errors"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

type RedisStore struct {
	Client redis.Client
	Prefix string
}

func (this *RedisStore) get(key string) *redis.StringCmd {
	return this.Client.Get(this.Prefix + key)
}

func (this *RedisStore) Get(key string) (interface{}, error) {
	intVal, err := this.get(key).Int64()

	if err != nil {
		floatVal, err := this.get(key).Float64()

		if err != nil {
			value, err := this.get(key).Result()

			if err != nil {
				if err.Error() == "redis: nil" {
					return "0", nil
				}

				return value, err
			}

			return SimpleDecode(value)
		}

		if &floatVal == nil {
			return floatVal, errors.New("Float value is nil.")
		}

		return floatVal, nil
	}

	if &intVal == nil {
		return intVal, errors.New("Int value is nil.")
	}

	return intVal, nil
}

func (this *RedisStore) GetFloat(key string) (float64, error) {
	return this.get(key).Float64()
}

func (this *RedisStore) GetInt(key string) (int64, error) {
	return this.get(key).Int64()
}

func (this *RedisStore) Increment(key string, value int64) (int64, error) {
	return this.Client.IncrBy(this.Prefix+key, value).Result()
}

func (this *RedisStore) Decrement(key string, value int64) (int64, error) {
	return this.Client.DecrBy(this.Prefix+key, value).Result()
}

func (this *RedisStore) Put(key string, value interface{}, minutes int) error {
	time, err := time.ParseDuration(strconv.Itoa(minutes) + "m")

	if err != nil {
		return err
	}

	if IsNumeric(value) {
		return this.Client.Set(this.Prefix+key, value, time).Err()
	}

	val, err := Encode(value)

	if err != nil {
		return err
	}

	return this.Client.Set(this.GetPrefix()+key, val, time).Err()
}

func (this *RedisStore) Forever(key string, value interface{}) error {
	if IsNumeric(value) {
		err := this.Client.Set(this.Prefix+key, value, 0).Err()

		if err != nil {
			return err
		}

		return this.Client.Persist(this.Prefix + key).Err()
	}

	val, err := Encode(value)

	if err != nil {
		return err
	}

	err = this.Client.Set(this.Prefix+key, val, 0).Err()

	if err != nil {
		return err
	}

	return this.Client.Persist(this.Prefix + key).Err()
}

func (this *RedisStore) Flush() (bool, error) {
	err := this.Client.FlushDB().Err()

	if err != nil {
		return false, err
	}

	return true, nil
}

func (this *RedisStore) Forget(key string) (bool, error) {
	err := this.Client.Del(this.Prefix + key).Err()

	if err != nil {
		return false, err
	}

	return true, nil
}

func (this *RedisStore) GetPrefix() string {
	return this.Prefix
}

func (this *RedisStore) PutMany(values map[string]interface{}, minutes int) {
	pipe := this.Client.TxPipeline()

	for key, value := range values {
		this.Put(key, value, minutes)
	}

	_, err := pipe.Exec()

	if err != nil {
		panic(err)
	}
}

func (this *RedisStore) Many(keys []string) (map[string]interface{}, error) {
	values := make(map[string]interface{})

	pipe := this.Client.TxPipeline()

	for _, key := range keys {
		val, err := this.Get(key)

		if err != nil {
			return values, err
		}

		values[key] = val
	}

	_, err := pipe.Exec()

	return values, err
}

func (this *RedisStore) Connection() interface{} {
	return this.Client
}

func (this *RedisStore) Lpush(segment string, key string) {
	this.Client.LPush(segment, key)
}

func (this *RedisStore) Lrange(key string, start int64, stop int64) []string {
	return this.Client.LRange(key, start, stop).Val()
}

func (this *RedisStore) Tags(names []string) TaggedStoreInterface {
	return &RedisTaggedCache{
		TaggedCache{
			Store: this,
			Tags: TagSet{
				Store: this,
				Names: names,
			},
		},
	}
}

func (this *RedisStore) GetStruct(key string, entity interface{}) (interface{}, error) {
	value, err := this.get(key).Result()

	if err != nil {
		panic(err)
	}

	return Decode(value, entity)
}
