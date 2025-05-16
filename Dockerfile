FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o rate-limiter .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/rate-limiter .
COPY .env .

EXPOSE 8080

CMD ["./rate-limiter"]