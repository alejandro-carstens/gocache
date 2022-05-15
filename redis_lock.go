package gocache

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

var _ Lock = &redisLock{}

func newRedisLock(client *redis.Client, name, owner string, duration time.Duration) *redisLock {
	return (&redisLock{
		client:   client,
		name:     name,
		owner:    owner,
		duration: duration,
	}).initBaseLock()
}

type redisLock struct {
	baseLock
	client   *redis.Client
	name     string
	owner    string
	duration time.Duration
}

// Acquire implementation of the Lock interface
func (rl *redisLock) Acquire() (bool, error) {
	return rl.client.SetNX(rl.name, rl.owner, rl.duration).Result()
}

// Release implementation of the Lock interface
func (rl *redisLock) Release() (bool, error) {
	res, err := rl.client.Eval(redisLuaReleaseLockScript, []string{rl.name}, rl.owner).Int64()

	return res > 0, err
}

// ForceRelease implementation of the Lock interface
func (rl *redisLock) ForceRelease() error {
	if _, err := rl.client.Del(rl.name).Result(); err != nil {
		return checkErrNotFound(err)
	}

	return nil
}

// GetCurrentOwner implementation of the Lock interface
func (rl *redisLock) GetCurrentOwner() (string, error) {
	res, err := rl.client.Get(rl.name).Result()
	if err != nil && errors.Is(err, redis.Nil) {
		return "", nil
	}

	return res, err
}

func (rl *redisLock) initBaseLock() *redisLock {
	rl.lock = rl

	return rl
}
