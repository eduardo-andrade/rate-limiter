package main

import (
	"log"
	"net/http"
	"time"

	"github.com/eduardo-andrade/rate-limiter/config"
	"github.com/eduardo-andrade/rate-limiter/limiter"
	"github.com/eduardo-andrade/rate-limiter/middleware"
	"github.com/eduardo-andrade/rate-limiter/storage"
)

func main() {
	cfg := config.LoadConfig()

	redisStorage, err := storage.NewRedisStorage(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	rateLimiter := limiter.NewLimiter(
		redisStorage,
		cfg.IPLimit,
		cfg.TokenLimit,
		cfg.IPExpiration,
		cfg.TokenExpiration,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	handler := middleware.RateLimiterMiddleware(rateLimiter, cfg)(mux)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server starting on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}