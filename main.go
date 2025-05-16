package main

import (
	"log"
	"net/http"
	"time"

	"github.com/eduardo-andrade/rate-limiter/config"
	"github.com/eduardo-andrade/rate-limiter/limiter"
	"github.com/eduardo-andrade/rate-limiter/middleware"
	"github.com/eduardo-andrade/rate-limiter/storage"
	"github.com/eduardo-andrade/rate-limiter/tests"
	"github.com/eduardo-andrade/rate-limiter/web"
	
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

	mux.HandleFunc("/test/ip", func(w http.ResponseWriter, r *http.Request) {
	results := tests.TestIPRateLimiter(rateLimiter, cfg)
	w.Write([]byte(results))
})

mux.HandleFunc("/test/token", func(w http.ResponseWriter, r *http.Request) {
	results := tests.TestTokenRateLimiter(rateLimiter, cfg)
	w.Write([]byte(results))
})

mux.HandleFunc("/test/expiration", func(w http.ResponseWriter, r *http.Request) {
	results := tests.TestExpiration(rateLimiter, cfg)
	w.Write([]byte(results))
})

mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/test.html")
})

	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Você acessou a rota /api"))
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