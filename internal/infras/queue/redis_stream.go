package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gin-admin-base/internal/infras/config"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RedisStreamQueue 基于 Redis Streams 的消息队列实现
type RedisStreamQueue struct {
	client    *redis.Client
	ctx       context.Context
	handlers  map[string]Handler
	consumers map[string]context.CancelFunc
	cfg       *config.QueueConfig
}

// NewRedisStreamQueue 创建 Redis Streams 消息队列
func NewRedisStreamQueue(client *redis.Client, cfg *config.QueueConfig) *RedisStreamQueue {
	q := &RedisStreamQueue{
		client:    client,
		ctx:       context.Background(),
		handlers:  make(map[string]Handler),
		consumers: make(map[string]context.CancelFunc),
		cfg:       cfg,
	}

	// 如果配置了 trim_interval，启动定期清理协程
	if cfg.TrimInterval > 0 {
		go q.trimLoop()
	}

	return q
}

// trimLoop 定期清理过长的 Stream
func (q *RedisStreamQueue) trimLoop() {
	ticker := time.NewTicker(time.Duration(q.cfg.TrimInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-q.ctx.Done():
			return
		case <-ticker.C:
			q.trimAllStreams()
		}
	}
}

// trimAllStreams 清理所有已订阅的 Stream
func (q *RedisStreamQueue) trimAllStreams() {
	for topic := range q.handlers {
		maxLen := q.cfg.GetTopicMaxLen(topic)
		if maxLen > 0 {
			// 使用 XTRIM 保留最新的 N 条消息
			q.client.XTrimMaxLen(q.ctx, topic, maxLen)
		}
	}
}

// Publish 发布消息到 Redis Stream
func (q *RedisStreamQueue) Publish(topic string, payload map[string]interface{}) error {
	msg := Message{
		ID:        uuid.New().String(),
		Topic:     topic,
		Payload:   payload,
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %w", err)
	}

	// 使用 XADD 添加到 Stream，并限制最大长度
	args := &redis.XAddArgs{
		Stream: topic,
		Values: map[string]interface{}{
			"data": string(data),
		},
	}

	// 设置 stream 最大长度
	if maxLen := q.cfg.GetTopicMaxLen(topic); maxLen > 0 {
		args.MaxLen = maxLen
		args.Approx = true // 使用近似裁剪，性能更好
	}

	_, err = q.client.XAdd(q.ctx, args).Result()
	if err != nil {
		return fmt.Errorf("发布消息到 Redis Stream 失败: %w", err)
	}

	return nil
}

// Subscribe 订阅 Redis Stream
func (q *RedisStreamQueue) Subscribe(topic string, handler Handler) error {
	if _, exists := q.handlers[topic]; exists {
		return fmt.Errorf("主题 %s 已被订阅", topic)
	}

	q.handlers[topic] = handler

	ctx, cancel := context.WithCancel(q.ctx)
	q.consumers[topic] = cancel

	consumerGroup := q.cfg.ConsumerGroup
	if consumerGroup == "" {
		consumerGroup = "Gonio-group"
	}

	// 先创建消费组（幂等操作，只执行一次）
	if err := q.client.XGroupCreateMkStream(ctx, topic, consumerGroup, "0").Err(); err != nil {
		// 如果消费组已存在，忽略错误
		if err.Error() != "BUSYGROUP Consumer Group name already exists" {
			return fmt.Errorf("创建消费组失败: %w", err)
		}
	}

	// 获取该 topic 的并发数
	concurrency := q.cfg.GetTopicConcurrency(topic)

	// 启动多个消费者协程
	for i := 0; i < concurrency; i++ {
		consumerID := fmt.Sprintf("consumer-%s-%d", uuid.New().String()[:8], i)
		go q.consumeMessages(ctx, topic, consumerGroup, consumerID, handler)
	}

	return nil
}

// consumeMessages 消费者协程，持续读取并处理消息
func (q *RedisStreamQueue) consumeMessages(ctx context.Context, topic, consumerGroup, consumerID string, handler Handler) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// 读取消息
			streams, err := q.client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    consumerGroup,
				Consumer: consumerID,
				Streams:  []string{topic, ">"},
				Count:    10,
				Block:    2 * time.Second,
			}).Result()

			if err != nil {
				// redis.Nil 表示超时无消息，属于正常情况
				if err == redis.Nil {
					continue
				}
				// 其他错误记录日志（通过 fmt 简单记录，避免循环依赖）
				continue
			}

			for _, stream := range streams {
				for _, message := range stream.Messages {
					dataStr, ok := message.Values["data"].(string)
					if !ok {
						continue
					}

					var msg Message
					if err := json.Unmarshal([]byte(dataStr), &msg); err != nil {
						continue
					}

					// 处理消息
					if err := handler(msg); err != nil {
						// 处理失败，不确认消息（下次会重新投递）
						continue
					}

					// 处理成功，确认消息
					q.client.XAck(ctx, topic, consumerGroup, message.ID)
				}
			}
		}
	}
}

// Close 关闭所有消费者
func (q *RedisStreamQueue) Close() error {
	for _, cancel := range q.consumers {
		cancel()
	}
	return nil
}
