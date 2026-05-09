package queue

import (
	"time"
)

// Message 消息体
type Message struct {
	ID        string                 `json:"id"`
	Topic     string                 `json:"topic"`
	Payload   map[string]interface{} `json:"payload"`
	CreatedAt time.Time              `json:"created_at"`
}

// Handler 消息处理函数
type Handler func(msg Message) error

// Queue 消息队列接口
type Queue interface {
	// Publish 发布消息
	Publish(topic string, payload map[string]interface{}) error
	// Subscribe 订阅消息
	Subscribe(topic string, handler Handler) error
	// Close 关闭队列
	Close() error
}
