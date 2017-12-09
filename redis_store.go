package cache

import (
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

func (this *RedisStore) Get(key string) interface{} {
	intVal, err := this.get(key).Int64()

	if err != nil {

		floatVal, err := this.get(key).Float64()

		if err != nil {
			value, err := this.get(key).Result()

			if err != nil {
				panic(err)
			}

			val, err := SimpleDecode(value)

			if err != nil {
				panic(err)
			}

			return val
		}

		if &floatVal == nil {
			panic("Float value is nil.")
		}

		return floatVal
	}

	if &intVal == nil {
		panic("Int value is nil.")
	}

	return intVal
}

func (this *RedisStore) GetFloat(key string) (float64, error) {
	return this.get(key).Float64()
}

func (this *RedisStore) GetInt(key string) (int64, error) {
	return this.get(key).Int64()
}

func (this *RedisStore) Increment(key string, value int64) int64 {
	value, err := this.Client.IncrBy(this.Prefix+key, value).Result()

	if err != nil {
		panic(err)
	}

	return value
}

func (this *RedisStore) Decrement(key string, value int64) int64 {
	value, err := this.Client.DecrBy(this.Prefix+key, value).Result()

	if err != nil {
		panic(err)
	}

	return value
}

func (this *RedisStore) Put(key string, value interface{}, minutes int) {
	time, err := time.ParseDuration(strconv.Itoa(minutes) + "m")

	if err != nil {
		panic(err)
	}

	if IsNumeric(value) {
		err = this.Client.Set(this.Prefix+key, value, time).Err()

		if err != nil {
			panic(err)
		}

		return
	}

	val, err := Encode(value)

	if err != nil {
		panic(err)
	}

	err = this.Client.Set(this.Prefix+key, val, time).Err()

	if err != nil {
		panic(err)
	}
}

func (this *RedisStore) Forever(key string, value interface{}) {
	if IsNumeric(value) {
		err := this.Client.Set(this.Prefix+key, value, 0).Err()

		if err != nil {
			panic(err)
		}

		err = this.Client.Persist(this.Prefix + key).Err()

		if err != nil {
			panic(err)
		}

		return
	}

	val, err := Encode(value)

	if err != nil {
		panic(err)
	}

	err = this.Client.Set(this.Prefix+key, val, 0).Err()

	if err != nil {
		panic(err)
	}

	err = this.Client.Persist(this.Prefix + key).Err()

	if err != nil {
		panic(err)
	}
}

func (this *RedisStore) Flush() bool {
	err := this.Client.FlushDB().Err()

	if err != nil {
		panic(err)
	}

	return true
}

func (this *RedisStore) Forget(key string) bool {
	err := this.Client.Del(this.Prefix + key).Err()

	if err != nil {
		panic(err)
	}

	return true
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

func (this *RedisStore) Many(keys []string) map[string]interface{} {
	values := make(map[string]interface{})

	pipe := this.Client.TxPipeline()

	for _, key := range keys {
		values[key] = this.Get(key)
	}

	_, err := pipe.Exec()

	if err != nil {
		panic(err)
	}

	return values
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

func (this *RedisStore) Tags(names []string) *RedisTaggedCache {
	taggedCache := &RedisTaggedCache{
		TaggedCache{
			Store: this,
			Tags: TagSet{
				Store: this,
				Names: names,
			},
		},
	}

	return taggedCache
}

func (this *RedisStore) GetStruct(key string, entity interface{}) (interface{}, error) {
	value, err := this.get(key).Result()

	if err != nil {
		panic(err)
	}

	return Decode(value, entity)
}
