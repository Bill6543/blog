package cache

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// viewKey 浏览量计数器的 Redis Hash key
const viewKey = "blog:view"

// ViewCounter 文章浏览量计数器：先累加到 Redis，再定时批量落库
type ViewCounter struct {
	rdb *redis.Client
}

// NewViewCounter 创建浏览量计数器实例
func NewViewCounter(rdb *redis.Client) *ViewCounter {
	return &ViewCounter{rdb: rdb}
}

// Incr 增加一篇文章的浏览量
func (c *ViewCounter) Incr(articleID uint) error {
	field := strconv.FormatUint(uint64(articleID), 10)
	return c.rdb.HIncrBy(context.Background(), viewKey, field, 1).Err()
}

// Drain 读取并清空所有浏览量计数，返回 map[articleID]增量
// 使用 RENAME 原子地把计数搬到临时 key，避免与新的 Incr 竞争、也避免并发 Drain 重复读取
func (c *ViewCounter) Drain() (map[uint]int64, error) {
	ctx := context.Background()
	tmpKey := viewKey + ":drain"

	// 先检查 key 是否存在，避免在不存在的 key 上执行 RenameNX 导致错误
	exists, err := c.rdb.Exists(ctx, viewKey).Result()
	if err != nil {
		return nil, err
	}
	if exists == 0 {
		// 没有待落库的计数
		return map[uint]int64{}, nil
	}

	renamed, err := c.rdb.RenameNX(ctx, viewKey, tmpKey).Result()
	if err != nil {
		return nil, err
	}
	if !renamed {
		// 没有待落库的计数（理论上不会到这里，因为上面已检查）
		return map[uint]int64{}, nil
	}

	vals, err := c.rdb.HGetAll(ctx, tmpKey).Result()
	if err != nil {
		return nil, err
	}
	_ = c.rdb.Del(ctx, tmpKey).Err()

	counts := make(map[uint]int64, len(vals))
	for field, v := range vals {
		id, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			continue
		}
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n <= 0 {
			continue
		}
		counts[uint(id)] = n
	}
	return counts, nil
}
