package gocache

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

var _ Lock = &redisLock{}

type redisLock struct {
	client  *redis.Client
	seconds int64
	name    string
	owner   string
}

// Acquire implementation of the Lock interface
func (rl *redisLock) Acquire() (bool, error) {
	return rl.client.SetNX(rl.name, rl.owner, time.Duration(rl.seconds)*time.Second).Result()
}

// Release implementation of the Lock interface
func (rl *redisLock) Release() (bool, error) {
	res, err := rl.client.Eval(redisLuaReleaseLockScript, []string{rl.name}, rl.owner).Int64()

	return res > 0, err
}

// ForceRelease implementation of the Lock interface
func (rl *redisLock) ForceRelease() error {
	_, err := rl.client.Del(rl.name).Result()

	return err
}

// GetCurrentOwner implementation of the Lock interface
func (rl *redisLock) GetCurrentOwner() (string, error) {
	res, err := rl.client.Get(rl.name).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		return "", nil
	}

	return res, err
}
