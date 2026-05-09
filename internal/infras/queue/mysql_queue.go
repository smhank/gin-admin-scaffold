package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MySQLQueue 基于 MySQL 的消息队列实现
type MySQLQueue struct {
	db        *gorm.DB
	handlers  map[string]Handler
	consumers map[string]context.CancelFunc
	mu        sync.RWMutex
}

// NewMySQLQueue 创建 MySQL 消息队列
func NewMySQLQueue(db *gorm.DB) *MySQLQueue {
	// 自动建表
	if err := db.AutoMigrate(&QueueMessage{}); err != nil {
		// 建表失败不影响队列创建，后续操作会报错
	}

	return &MySQLQueue{
		db:        db,
		handlers:  make(map[string]Handler),
		consumers: make(map[string]context.CancelFunc),
	}
}

// QueueMessage 数据库中的消息记录
type QueueMessage struct {
	ID         uint       `gorm:"primarykey"`
	MsgID      string     `gorm:"column:msg_id;type:varchar(36);index;not null"`
	Topic      string     `gorm:"column:topic;type:varchar(255);index;not null"`
	Payload    string     `gorm:"column:payload;type:text;not null"`
	Status     string     `gorm:"column:status;type:varchar(20);default:'pending';index"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
	ConsumedAt *time.Time `gorm:"column:consumed_at"`
}

// TableName 表名
func (QueueMessage) TableName() string {
	return "queue_messages"
}

// Publish 发布消息到 MySQL
func (q *MySQLQueue) Publish(topic string, payload map[string]interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %w", err)
	}

	msg := QueueMessage{
		MsgID:   uuid.New().String(),
		Topic:   topic,
		Payload: string(data),
		Status:  "pending",
	}

	if err := q.db.Create(&msg).Error; err != nil {
		return fmt.Errorf("保存消息到数据库失败: %w", err)
	}

	return nil
}

// Subscribe 订阅 MySQL 消息
func (q *MySQLQueue) Subscribe(topic string, handler Handler) error {
	q.mu.Lock()
	if _, exists := q.handlers[topic]; exists {
		q.mu.Unlock()
		return fmt.Errorf("主题 %s 已被订阅", topic)
	}
	q.handlers[topic] = handler
	q.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	q.mu.Lock()
	q.consumers[topic] = cancel
	q.mu.Unlock()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// 轮询待处理消息
				var messages []QueueMessage
				q.db.Where("topic = ? AND status = ?", topic, "pending").
					Order("created_at ASC").
					Limit(10).
					Find(&messages)

				for _, msg := range messages {
					var payload map[string]interface{}
					if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
						// 无效消息，标记为失败
						q.db.Model(&msg).Update("status", "failed")
						continue
					}

					message := Message{
						ID:        msg.MsgID,
						Topic:     msg.Topic,
						Payload:   payload,
						CreatedAt: msg.CreatedAt,
					}

					if err := handler(message); err != nil {
						// 处理失败，标记为失败
						q.db.Model(&msg).Update("status", "failed")
						continue
					}

					// 处理成功
					now := time.Now()
					q.db.Model(&msg).Updates(map[string]interface{}{
						"status":      "consumed",
						"consumed_at": &now,
					})
				}

				time.Sleep(1 * time.Second)
			}
		}
	}()

	return nil
}

// Close 关闭所有消费者
func (q *MySQLQueue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	for _, cancel := range q.consumers {
		cancel()
	}
	return nil
}
