package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/eduardo-andrade/rate-limiter/config"
	"github.com/eduardo-andrade/rate-limiter/limiter"
)

func TestIPRateLimiter(l *limiter.Limiter, cfg *config.Config) string {
	ctx := context.Background()
	key := "ip:127.0.0.1"
	allowed, blocked := 0, 0

	for i := 0; i < 10; i++ {
		ok, err := l.AllowRequest(ctx, key, false)
		if err != nil {
			return fmt.Sprintf("Erro: %v", err)
		}
		if ok {
			allowed++
		} else {
			blocked++
		}
	}

	return fmt.Sprintf("Limite por IP - Permitidas: %d, Bloqueadas: %d", allowed, blocked)
}

func TestTokenRateLimiter(l *limiter.Limiter, cfg *config.Config) string {
	ctx := context.Background()
	key := "token:abc123"
	allowed, blocked := 0, 0

	for i := 0; i < 15; i++ {
		ok, err := l.AllowRequest(ctx, key, true)
		if err != nil {
			return fmt.Sprintf("Erro: %v", err)
		}
		if ok {
			allowed++
		} else {
			blocked++
		}
	}

	return fmt.Sprintf("Limite por Token - Permitidas: %d, Bloqueadas: %d", allowed, blocked)
}

func TestExpiration(l *limiter.Limiter, cfg *config.Config) string {
	ctx := context.Background()
	key := "ip:192.168.0.1"

	for i := 0; i < cfg.IPLimit+1; i++ {
		l.AllowRequest(ctx, key, false)
	}

	time.Sleep(cfg.IPExpiration + time.Second)

	ok, err := l.AllowRequest(ctx, key, false)
	if err != nil {
		return fmt.Sprintf("Erro: %v", err)
	}

	if ok {
		return "Após expiração: requisição foi PERMITIDA (correto)"
	} else {
		return "Após expiração: requisição foi BLOQUEADA (erro)"
	}
}
