package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
)

// Lua script to be executed on the redis side
const (
	lockScript = `
		return redis.call('SET', KEYS[1], ARGV[1], 'NX', 'PX', ARGV[2])
	`
	unlockScript = `
		if redis.call("get",KEYS[1]) == ARGV[1] then
		    return redis.call("del",KEYS[1])
		else
		    return 0
		end
	`
)

type Lock struct {
	RedisClient *redis.Client
}

func (rl *Lock) Lock(key, value string, timeoutMs int) (bool, error) {
	ctx := context.Background()

	resp := rl.RedisClient.Eval(ctx, lockScript, []string{key}, []string{value, strconv.Itoa(timeoutMs)})

	if res, err := resp.Result(); err != nil {
		return false, err
	} else {
		return res == "OK", nil
	}
}

func (rl *Lock) Unlock(key, value string) {
	ctx := context.Background()

	_ = rl.RedisClient.Eval(ctx, unlockScript, []string{key}, []string{value})
}
