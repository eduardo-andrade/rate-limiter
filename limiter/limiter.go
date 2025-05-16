package limiter

import (
	"context"
	"time"

	"github.com/eduardo-andrade/rate-limiter/storage"
)

type Limiter struct {
	storage    storage.Storage
	ipLimit    int
	ipExpiry   int64
	tokenLimit int
	tokenExpiry int64
}

func NewLimiter(storage storage.Storage, ipLimit, tokenLimit int, ipExpiry, tokenExpiry time.Duration) *Limiter {
	return &Limiter{
		storage:    storage,
		ipLimit:    ipLimit,
		ipExpiry:   int64(ipExpiry.Seconds()),
		tokenLimit: tokenLimit,
		tokenExpiry: int64(tokenExpiry.Seconds()),
	}
}

func (l *Limiter) AllowRequest(ctx context.Context, identifier string, isToken bool) (bool, error) {
	
	blocked, err := l.storage.IsBlocked(ctx, identifier)
	if err != nil {
		return false, err
	}
	if blocked {
		return false, nil
	}

	limit := l.ipLimit
	expiry := l.ipExpiry
	if isToken {
		limit = l.tokenLimit
		expiry = l.tokenExpiry
	}

	count, err := l.storage.Increment(ctx, identifier, expiry)
	if err != nil {
		return false, err
	}

	if count > int64(limit) {
		if err := l.storage.Block(ctx, identifier, expiry); err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}