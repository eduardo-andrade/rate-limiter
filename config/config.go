package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr          string
	EnableIPLimiter    bool
	IPLimit            int
	IPExpiration       time.Duration
	EnableTokenLimiter bool
	TokenLimit         int
	TokenExpiration    time.Duration
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		RedisAddr:          getEnv("REDIS_ADDR", "localhost:6379"),
		EnableIPLimiter:    getEnvAsBool("ENABLE_IP_LIMITER", true),
		IPLimit:            getEnvAsInt("IP_LIMIT", 5),
		IPExpiration:       getEnvAsDuration("IP_EXPIRATION", 5*time.Minute),
		EnableTokenLimiter: getEnvAsBool("ENABLE_TOKEN_LIMITER", true),
		TokenLimit:         getEnvAsInt("TOKEN_LIMIT", 10),
		TokenExpiration:    getEnvAsDuration("TOKEN_EXPIRATION", 5*time.Minute),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
