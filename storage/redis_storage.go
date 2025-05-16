package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr string) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisStorage{client: client}, nil
}

func (r *RedisStorage) Increment(ctx context.Context, key string, expiration int64) (int64, error) {
	luaScript := `
	local current
	current = redis.call("incr", KEYS[1])
	if tonumber(current) == 1 then
		redis.call("expire", KEYS[1], ARGV[1])
	end
	return current
	`

	result, err := r.client.Eval(ctx, luaScript, []string{key}, expiration).Int64()
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (r *RedisStorage) Block(ctx context.Context, key string, expiration int64) error {
	return r.client.Set(ctx, "blocked:"+key, "1", time.Duration(expiration)*time.Second).Err()
}

func (r *RedisStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	res, err := r.client.Exists(ctx, "blocked:"+key).Result()
	return res > 0, err
}