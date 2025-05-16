package tests

import (
	"context"
	"testing"
	"time"
	"net/http"
	"net/http/httptest"
	"github.com/eduardo-andrade/rate-limiter/limiter"
	"github.com/eduardo-andrade/rate-limiter/storage"
)

func setupTestLimiter(t *testing.T) *limiter.Limiter {
	redisAddr := "localhost:6379"
	store, err := storage.NewRedisStorage(redisAddr)
	if err != nil {
		t.Fatalf("failed to connect to Redis: %v", err)
	}

	return limiter.NewLimiter(store, 5, 10, 10*time.Second, 10*time.Second)
}

func TestRateLimitByIP(t *testing.T) {
	l := setupTestLimiter(t)
	ctx := context.Background()
	key := "ip:127.0.0.1"

	allowedCount := 0
	blockedCount := 0

	for i := 0; i < 10; i++ {
		allowed, err := l.AllowRequest(ctx, key, false)
		if err != nil {
			t.Fatalf("error in AllowRequest: %v", err)
		}
		if allowed {
			allowedCount++
		} else {
			blockedCount++
		}
	}

	if allowedCount > 5 {
		t.Errorf("expected max 5 allowed requests, got %d", allowedCount)
	}
	if blockedCount == 0 {
		t.Errorf("expected some blocked requests, got %d", blockedCount)
	}
}

func TestRateLimitByToken(t *testing.T) {
	l := setupTestLimiter(t)
	ctx := context.Background()
	key := "token:abc123"

	allowedCount := 0
	blockedCount := 0

	for i := 0; i < 15; i++ {
		allowed, err := l.AllowRequest(ctx, key, true)
		if err != nil {
			t.Fatalf("error in AllowRequest: %v", err)
		}
		if allowed {
			allowedCount++
		} else {
			blockedCount++
		}
	}

	if allowedCount > 10 {
		t.Errorf("expected max 10 allowed requests, got %d", allowedCount)
	}
	if blockedCount == 0 {
		t.Errorf("expected some blocked requests, got %d", blockedCount)
	}
}

func TestResetAfterExpiration(t *testing.T) {
	l := setupTestLimiter(t)
	ctx := context.Background()
	key := "ip:192.168.0.1"

	// Exceed limit
	for i := 0; i < 6; i++ {
		l.AllowRequest(ctx, key, false)
	}

	t.Log("Waiting for expiration...")
	time.Sleep(11 * time.Second) // wait for expiration

	allowed, err := l.AllowRequest(ctx, key, false)
	if err != nil {
		t.Fatalf("error in AllowRequest after expiration: %v", err)
	}
	if !allowed {
		t.Error("expected request to be allowed after expiration")
	}
}

func BenchmarkRateLimitByIP(b *testing.B) {
	l := setupTestLimiter(&testing.T{})
	ctx := context.Background()
	key := "ip:10.0.0.1"

	for i := 0; i < b.N; i++ {
		_, _ = l.AllowRequest(ctx, key, false)
	}
}

func BenchmarkRateLimitByToken(b *testing.B) {
	l := setupTestLimiter(&testing.T{})
	ctx := context.Background()
	key := "token:test123"

	for i := 0; i < b.N; i++ {
		_, _ = l.AllowRequest(ctx, key, true)
	}
}

func TestHTTPHandlerWithIP(t *testing.T) {
	l := setupTestLimiter(t)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	limiterFunc := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			allowed, _ := l.AllowRequest(ctx, "ip:1.2.3.4", false)
			if !allowed {
				http.Error(w, "rate limit", http.StatusTooManyRequests)
				return
			}
			h.ServeHTTP(w, r)
		})
	}

	ts := httptest.NewServer(limiterFunc(handler))
	defer ts.Close()

	for i := 0; i < 7; i++ {
		resp, err := http.Get(ts.URL)
		if err != nil {
			t.Fatalf("failed request: %v", err)
		}
		if i >= 5 && resp.StatusCode != http.StatusTooManyRequests {
			t.Errorf("expected 429, got %d", resp.StatusCode)
		}
	}
}