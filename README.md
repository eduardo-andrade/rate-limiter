# Rate Limiter in Go

A configurable rate limiter middleware for Go web applications that can limit requests based on IP address or API tokens.

## Features

- IP-based rate limiting
- Token-based rate limiting (with higher priority)
- Redis-backed storage
- Configurable limits and expiration times
- Docker support

## Configuration

Configure the rate limiter using environment variables:

- `REDIS_ADDR`: Redis server address (default: `localhost:6379`)
- `ENABLE_IP_LIMITER`: Enable IP-based limiting (default: `true`)
- `IP_LIMIT`: Maximum requests per IP (default: `5`)
- `IP_EXPIRATION`: IP limit expiration in seconds (default: `300`)
- `ENABLE_TOKEN_LIMITER`: Enable token-based limiting (default: `true`)
- `TOKEN_LIMIT`: Maximum requests per token (default: `10`)
- `TOKEN_EXPIRATION`: Token limit expiration in seconds (default: `300`)

## Running with Docker

1. Copy `.env.example` to `.env` and adjust settings if needed
2. Run: `docker-compose up --build`

## Testing

- For IP limiting: Make multiple requests from the same IP
- For token limiting: Include `API_KEY: <TOKEN>` header in requests

When limits are exceeded, the server responds with HTTP 429 status code.
