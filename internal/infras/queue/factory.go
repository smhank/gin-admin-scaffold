package queue

import (
	"fmt"

	"gin-admin-base/internal/infras/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// QueueType 队列类型
type QueueType string

const (
	// QueueTypeRedis Redis Streams 队列
	QueueTypeRedis QueueType = "redis"
	// QueueTypeMySQL MySQL 队列
	QueueTypeMySQL QueueType = "mysql"
)

// InitQueue 根据配置初始化消息队列
func InitQueue(redisClient *redis.Client, db *gorm.DB) (Queue, error) {
	queueCfg := config.GetQueueConfig()
	queueType := QueueType(queueCfg.Driver)
	if queueType == "" {
		queueType = QueueTypeRedis // 默认使用 Redis
	}

	switch queueType {
	case QueueTypeRedis:
		if redisClient == nil {
			return nil, fmt.Errorf("Redis 客户端未初始化，无法创建 Redis Streams 队列")
		}
		return NewRedisStreamQueue(redisClient, queueCfg), nil
	case QueueTypeMySQL:
		if db == nil {
			return nil, fmt.Errorf("数据库未初始化，无法创建 MySQL 队列")
		}
		return NewMySQLQueue(db), nil
	default:
		return nil, fmt.Errorf("不支持的消息队列类型: %s", queueType)
	}
}
