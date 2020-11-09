package gocache

import (
	"github.com/go-redis/redis"
	"time"
)

type redisLock struct {
	client  *redis.Client
	seconds int64
	name    string
	owner   string
}

func (rl *redisLock) Acquire() (bool, error) {
	return rl.client.SetXX(rl.name, rl.owner, time.Duration(rl.seconds)*time.Second).Result()
}

func (rl redisLock) Release() (bool, error) {
	res, err := rl.client.Eval(redisLuaReleaseLockScript, []string{"1", rl.name, rl.owner}).Int64()

	return res > 0, err
}

func (rl *redisLock) ForceRelease() error {
	_, err := rl.client.Del(rl.name).Result()

	return err
}

func (rl *redisLock) GetCurrentOwner() (string, error) {
	return rl.client.Get(rl.name).Result()
}
