package database

import (
	"context"
	"fmt"
	"time"

	"blog/pkg/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis(cfg *config.RedisConfig) error {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect redis: %w", err)
	}

	RedisClient = client
	return nil
}

// GetRedis 获取 Redis 客户端实例
func GetRedis() *redis.Client {
	return RedisClient
}
