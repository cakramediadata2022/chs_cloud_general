package config

import (
	"context"
	"log"
	"time"

	"github.com/cakramediadata2022/chs_cloud_general/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	redisHost := cfg.RedisAddr

	if redisHost == "" {
		redisHost = ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:         redisHost,
		MinIdleConns: cfg.MinIdleConns,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
		Password:     cfg.RedisPassword, // no password set
		DB:           cfg.DB,            // use default DB
	})

	// Test the Redis connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}
