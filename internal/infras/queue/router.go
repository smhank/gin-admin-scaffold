package queue

import (
	"fmt"
	"sync"

	"gin-admin-base/internal/domain/model"

	"gorm.io/gorm"
)

// MessageRouter 消息路由器，管理多个主题的发布订阅
type MessageRouter struct {
	queue Queue
	mu    sync.RWMutex
	db    *gorm.DB
}

// NewMessageRouter 创建消息路由器
func NewMessageRouter(queue Queue, db *gorm.DB) *MessageRouter {
	return &MessageRouter{
		queue: queue,
		db:    db,
	}
}

// Publish 发布消息到指定主题
func (r *MessageRouter) Publish(topic string, payload map[string]interface{}) error {
	return r.queue.Publish(topic, payload)
}

// Subscribe 订阅指定主题
func (r *MessageRouter) Subscribe(topic string, handler Handler) error {
	return r.queue.Subscribe(topic, handler)
}

// Close 关闭消息队列
func (r *MessageRouter) Close() error {
	return r.queue.Close()
}

// 预定义主题常量
const (
	// TopicOperationLog 操作日志主题
	TopicOperationLog = "operation_log"
	// TopicNotification 通知主题
	TopicNotification = "notification"
	// TopicSystemEvent 系统事件主题
	TopicSystemEvent = "system_event"
)

// RegisterDefaultHandlers 注册默认的消息处理器
func (r *MessageRouter) RegisterDefaultHandlers() error {
	// 注册操作日志处理器
	if err := r.Subscribe(TopicOperationLog, r.handleOperationLog); err != nil {
		return fmt.Errorf("注册操作日志处理器失败: %w", err)
	}

	// 注册通知处理器
	if err := r.Subscribe(TopicNotification, handleNotification); err != nil {
		return fmt.Errorf("注册通知处理器失败: %w", err)
	}

	// 注册系统事件处理器
	if err := r.Subscribe(TopicSystemEvent, handleSystemEvent); err != nil {
		return fmt.Errorf("注册系统事件处理器失败: %w", err)
	}

	return nil
}

// handleOperationLog 处理操作日志消息，异步写入数据库
func (r *MessageRouter) handleOperationLog(msg Message) error {
	payload := msg.Payload

	log := model.OperationLog{
		Username:  toString(payload["username"]),
		Action:    toString(payload["action"]),
		Method:    toString(payload["method"]),
		Path:      toString(payload["path"]),
		Params:    toString(payload["params"]),
		Result:    toString(payload["result"]),
		Duration:  toInt64(payload["duration"]),
		IP:        toString(payload["ip"]),
		UserAgent: toString(payload["user_agent"]),
	}

	if r.db != nil {
		if err := r.db.Create(&log).Error; err != nil {
			return fmt.Errorf("保存操作日志失败: %w", err)
		}
	} else {
		fmt.Printf("[Queue] 操作日志(DB未初始化): %+v\n", payload)
	}

	return nil
}

// handleNotification 处理通知消息
func handleNotification(msg Message) error {
	fmt.Printf("[Queue] 通知: %+v\n", msg.Payload)
	// TODO: 实际业务处理，如发送邮件、推送通知等
	return nil
}

// handleSystemEvent 处理系统事件消息
func handleSystemEvent(msg Message) error {
	fmt.Printf("[Queue] 系统事件: %+v\n", msg.Payload)
	// TODO: 实际业务处理，如清理缓存、刷新配置等
	return nil
}

// toString 安全转换为字符串
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

// toInt64 安全转换为 int64
func toInt64(v interface{}) int64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case int64:
		return val
	case float64:
		return int64(val)
	case int:
		return int64(val)
	}
	return 0
}
