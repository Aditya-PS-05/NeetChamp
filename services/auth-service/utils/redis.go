package utils

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() {
	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	redisPassword := os.Getenv("REDIS_PASSWORD")
	useTLS := os.Getenv("REDIS_TLS") == "true"

	fmt.Println("Connecting to Redis at", redisAddr)

	var tlsConfig *tls.Config
	if useTLS {
		tlsConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		fmt.Println("Using TLS for Redis connection")
	} else {
		fmt.Println("⚠️ Not using TLS for Redis")
	}

	// Create Redis client
	RedisClient = redis.NewClient(&redis.Options{
		Addr:      redisAddr,
		Password:  redisPassword,
		DB:        0,
		TLSConfig: tlsConfig,
	})

	// Test Redis connection
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Redis connection error:", err)
		RedisClient = nil
		return
	}

	fmt.Println("Connected to Redis at", redisAddr)
}

// 🔹 Store Token in Redis (For Logout)
func SaveTokenToBlacklist(token string, expiration int) error {
	if RedisClient == nil {
		return fmt.Errorf("redis is not connected")
	}

	ctx := context.Background()
	err := RedisClient.Set(ctx, token, "blacklisted", 0).Err()
	if err != nil {
		return fmt.Errorf("failed to save token to Redis: %v", err)
	}
	return nil
}

// 🔹 Check if Token is Blacklisted
func IsTokenBlacklisted(token string) bool {
	if RedisClient == nil {
		fmt.Println("Redis is not connected")
		return false
	}

	ctx := context.Background()
	result, err := RedisClient.Get(ctx, token).Result()
	if err == redis.Nil {
		return false // Not blacklisted
	}
	return result == "blacklisted"
}
