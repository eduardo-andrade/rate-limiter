package middleware

import (
	"net/http"
	"strings"

	"github.com/eduardo-andrade/rate-limiter/config"
	"github.com/eduardo-andrade/rate-limiter/limiter"
)

func RateLimiterMiddleware(limiter *limiter.Limiter, cfg *config.Config, skipPaths []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Verifica se a rota atual está na lista de rotas a ignorar
			for _, p := range skipPaths {
				if strings.HasPrefix(r.URL.Path, p) {
					// pula o rate limiter para esta rota
					next.ServeHTTP(w, r)
					return
				}
			}

			ctx := r.Context()

			if cfg.EnableTokenLimiter {
				token := r.Header.Get("API_KEY")
				if token != "" {
					allowed, err := limiter.AllowRequest(ctx, "token:"+token, true)
					if err != nil {
						http.Error(w, "internal server error", http.StatusInternalServerError)
						return
					}
					if !allowed {
						http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
						return
					}
					next.ServeHTTP(w, r)
					return
				}
			}

			if cfg.EnableIPLimiter {
				ip := strings.Split(r.RemoteAddr, ":")[0]
				allowed, err := limiter.AllowRequest(ctx, "ip:"+ip, false)
				if err != nil {
					http.Error(w, "internal server error", http.StatusInternalServerError)
					return
				}
				if !allowed {
					http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
