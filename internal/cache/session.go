package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// SessionCache 缓存用户的 login_session，避免每次请求都查询数据库
type SessionCache struct {
	rdb *redis.Client
	ttl time.Duration
}

// NewSessionCache 创建会话缓存实例
func NewSessionCache(rdb *redis.Client, ttl time.Duration) *SessionCache {
	return &SessionCache{rdb: rdb, ttl: ttl}
}

// key 生成缓存键
func (c *SessionCache) key(userID uint) string {
	return fmt.Sprintf("auth:session:%d", userID)
}

// Get 读取缓存中的 login_session
// 返回值：session、是否命中、错误。缓存未命中时 found=false 且 err=nil。
func (c *SessionCache) Get(userID uint) (string, bool, error) {
	val, err := c.rdb.Get(context.Background(), c.key(userID)).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

// Set 写入 login_session（带 TTL）
func (c *SessionCache) Set(userID uint, session string) error {
	return c.rdb.Set(context.Background(), c.key(userID), session, c.ttl).Err()
}

// Del 删除缓存，用于登出/改密等使 session 失效的场景
func (c *SessionCache) Del(userID uint) error {
	return c.rdb.Del(context.Background(), c.key(userID)).Err()
}
