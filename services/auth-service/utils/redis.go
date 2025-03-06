package utils

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"time"

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

	RedisClient = redis.NewClient(&redis.Options{
		Addr:      redisAddr,
		Password:  redisPassword,
		DB:        0,
		TLSConfig: tlsConfig,
	})

	// ✅ Test Redis connection
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("❌ Redis connection error:", err)
		RedisClient = nil
		return
	}

	fmt.Println("✅ Connected to Redis at", redisAddr)
}

// ✅ Store Token in Redis (For Logout)
func SaveTokenToBlacklist(token string, expiration int) error {
	if RedisClient == nil {
		return fmt.Errorf("redis is not connected")
	}
	ctx := context.Background()
	err := RedisClient.Set(ctx, token, "blacklisted", time.Duration(expiration)*time.Second).Err()
	return err
}

// ✅ Check if Token is Blacklisted
func IsTokenBlacklisted(token string) bool {
	if RedisClient == nil {
		fmt.Println("❌ Redis is not connected")
		return false
	}

	ctx := context.Background()
	result, err := RedisClient.Get(ctx, token).Result()
	if err == redis.Nil {
		return false
	}
	return result == "blacklisted"
}

// ✅ Cache User Details in Redis
func CacheUser(email, password, role string) {
	if RedisClient == nil {
		fmt.Println("❌ Redis is not connected")
		return
	}

	userData := map[string]string{
		"password": password,
		"role":     role,
	}
	jsonData, _ := json.Marshal(userData)

	ctx := context.Background()
	RedisClient.Set(ctx, email, jsonData, 5*time.Minute) // Cache for 5 mins
}

// ✅ Retrieve Cached User from Redis
func GetCachedUser(email string) (*struct {
	Password string `json:"password"`
	Role     string `json:"role"`
}, error) {
	if RedisClient == nil {
		return nil, fmt.Errorf("❌ Redis is not connected")
	}

	ctx := context.Background()
	data, err := RedisClient.Get(ctx, email).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("❌ user not found in cache")
	} else if err != nil {
		return nil, err
	}

	var cachedUser struct {
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	json.Unmarshal([]byte(data), &cachedUser)
	return &cachedUser, nil
}
