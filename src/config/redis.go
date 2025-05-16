package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

var Ctx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", GetEnv("REDIS_HOST", "localhost"), GetEnv("REDIS_PORT", "6379")),
		Password: GetEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic("Redis connection error: " + err.Error())
	}

	fmt.Println("Redis connected")
}
