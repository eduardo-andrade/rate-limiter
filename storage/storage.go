package storage

import "context"

type Storage interface {
	Increment(ctx context.Context, key string, expiration int64) (int64, error)
	Block(ctx context.Context, key string, expiration int64) error
	IsBlocked(ctx context.Context, key string) (bool, error)
}