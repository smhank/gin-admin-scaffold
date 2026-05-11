package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
}

// ClearAllPermissionCaches 清除所有用户权限缓存
// 当权限或角色被修改时调用，确保所有用户的权限缓存失效
func ClearAllPermissionCaches() {
	if RedisClient == nil {
		return
	}
	// 先 Ping 检查 Redis 是否可用
	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		// Redis 不可用时静默忽略
		return
	}
	// 使用 SCAN 命令遍历所有权限缓存键并删除
	iter := RedisClient.Scan(Ctx, 0, "user_permissions:*", 0).Iterator()
	for iter.Next(Ctx) {
		RedisClient.Del(Ctx, iter.Val())
	}
}
