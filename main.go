package main

import (
	"log"
	"net/http"
	"strconv"
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

	mux.HandleFunc("/test/run", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		testType := r.URL.Query().Get("testType")
		token := r.URL.Query().Get("token")
		ip := r.URL.Query().Get("ip")
		requestsStr := r.URL.Query().Get("requests")
		intervalStr := r.URL.Query().Get("interval")
		maxAllowedStr := r.URL.Query().Get("maxAllowed")

		if testType == "" {
			http.Error(w, "Parâmetros do teste não informados", http.StatusBadRequest)
			return
		}
		requests, err := strconv.Atoi(requestsStr)
		if err != nil || requests < 1 {
			requests = 10
		}

		intervalMs, err := strconv.Atoi(intervalStr)
		if err != nil || intervalMs < 0 {
			intervalMs = 100
		}
		interval := time.Duration(intervalMs) * time.Millisecond

		maxAllowed, err := strconv.Atoi(maxAllowedStr)
		if err != nil || maxAllowed < 1 {
			maxAllowed = 5
		}

		ctx := r.Context()

		var key string
		if testType == "token" && token != "" {
			key = "token:" + token
		} else if testType == "ip" && ip != "" {
			key = "ip:" + ip
		} else if testType == "ip" {
			key = "ip:127.0.0.1"
		} else {
			http.Error(w, "Invalid test parameters", http.StatusBadRequest)
			return
		}

		allowedCount := 0
		blockedCount := 0
		var lastErr error

		for i := 0; i < requests; i++ {
			testLimiter := limiter.NewLimiter(
				redisStorage,
				func() int {
					if testType == "token" {
						return cfg.IPLimit // ou 0
					}
					return maxAllowed
				}(),
				func() int {
					if testType == "token" {
						return maxAllowed
					}
					return cfg.TokenLimit // ou 0
				}(),
				cfg.IPExpiration,
				cfg.TokenExpiration,
			)

			allowed, err := testLimiter.AllowRequest(ctx, key, testType == "token")

			if err != nil {
				lastErr = err
				break
			}
			if allowed {
				allowedCount++
			} else {
				blockedCount++
			}
			time.Sleep(interval)
		}

		if lastErr != nil {
			http.Error(w, "Error running test: "+lastErr.Error(), http.StatusInternalServerError)
			return
		}

		result := "Test Type: " + testType + "\n" +
			"Key: " + key + "\n" +
			"Requests: " + strconv.Itoa(requests) + "\n" +
			"Interval: " + interval.String() + "\n" +
			"Allowed: " + strconv.Itoa(allowedCount) + "\n" +
			"Blocked: " + strconv.Itoa(blockedCount) + "\n"

		if allowedCount > maxAllowed {
			result += "\n⚠️ Aviso: O número de requisições permitidas ultrapassou o limite máximo configurado (" + strconv.Itoa(maxAllowed) + ").\n"
		}

		w.Write([]byte(result))
	})

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/test.html")
	})

	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Você acessou a rota /api"))
	})

	handler := middleware.RateLimiterMiddleware(rateLimiter, cfg, []string{"/test", "/test/run"})(mux)
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
