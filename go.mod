module github.com/eduardo-andrade/rate-limiter

go 1.21

replace github.com/eduardo-andrade/rate-limiter/middleware => ./middleware

replace github.com/eduardo-andrade/rate-limiter/config => ./config

replace github.com/eduardo-andrade/rate-limiter/limiter => ./limiter

replace github.com/eduardo-andrade/rate-limiter/storage => ./storage

require (
	github.com/go-redis/redis/v8 v8.11.5
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)
