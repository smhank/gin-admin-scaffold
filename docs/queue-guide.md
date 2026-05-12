# 消息队列使用指南

## 概述

本项目基于 Watermill 架构思想，内置了消息队列功能，支持 **Redis Streams** 和 **MySQL** 两种驱动。消息队列用于解耦业务逻辑，实现异步处理、事件驱动等场景。

## 架构设计

```
┌─────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  发布者     │────▶│  消息路由器      │────▶│  队列驱动       │
│ (Publisher) │     │ (MessageRouter)  │     │ (Redis/MySQL)   │
└─────────────┘     └──────────────────┘     └─────────────────┘
                            │
                            ▼
                    ┌──────────────────┐
                    │  消息处理器      │
                    │ (Handler)        │
                    └──────────────────┘
```

### 核心组件

| 组件 | 文件 | 说明 |
|------|------|------|
| `Queue` 接口 | `message.go` | 定义 Publish、Subscribe、Close 三个核心方法 |
| `Message` 结构体 | `message.go` | 消息体，包含 ID、主题、负载、创建时间 |
| `Handler` 类型 | `message.go` | 消息处理函数类型 `func(msg Message) error` |
| `RedisStreamQueue` | `redis_stream.go` | Redis Streams 驱动实现 |
| `MySQLQueue` | `mysql_queue.go` | MySQL 驱动实现 |
| `MessageRouter` | `router.go` | 消息路由器，管理多主题发布订阅 |
| `InitQueue` | `factory.go` | 队列工厂，根据配置自动选择驱动 |

## 配置

### 配置文件 (config.yaml)

在 `internal/infras/config/config.yaml` 中配置消息队列：

```yaml
queue:
  driver: redis                    # 驱动类型: redis (Redis Streams) / mysql (MySQL)
  consumer_group: Gonio-group      # 消费者组名称
  default_max_len: 20              # 默认 Stream 最大长度
  trim_interval: 3600              # Stream 修剪间隔（秒）
  topic_concurrency:               # 各主题并发数
    operation_log: 1
    notification: 1
    system_event: 1
  topic_max_len:                   # 各主题最大长度
    operation_log: 100
    notification: 50
    system_event: 20
```

### Redis Streams 驱动

当 `driver: redis` 时，使用 Redis Streams 作为消息队列后端。

**特点：**
- 高性能，基于内存
- 支持消费组，消息可持久化到 Redis
- 支持消息确认机制（ACK）
- 支持消息重新投递

**前置条件：**
- Redis 服务正常运行
- 已配置 `redis.addr`（默认 `127.0.0.1:6379`）

### MySQL 驱动

当 `driver: mysql` 时，使用 MySQL 数据库表作为消息队列后端。

**特点：**
- 无需额外中间件
- 消息持久化到数据库
- 支持消息状态管理（pending/consumed/failed）
- 轮询模式，每秒扫描待处理消息

**前置条件：**
- 数据库连接正常

**queue_messages 表结构：**

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint (PK) | 主键 |
| msg_id | varchar(36) | 消息唯一 ID |
| topic | varchar(255) | 消息主题 |
| payload | text | 消息内容 (JSON) |
| status | varchar(20) | 状态: pending/consumed/failed |
| created_at | datetime | 创建时间 |
| consumed_at | datetime | 消费时间 |

## 预定义主题

| 主题常量 | 值 | 说明 |
|----------|-----|------|
| `TopicOperationLog` | `operation_log` | 操作日志 |
| `TopicNotification` | `notification` | 通知消息 |
| `TopicSystemEvent` | `system_event` | 系统事件 |

## 快速开始

### 1. 初始化消息队列

消息队列在 `internal/interfaces/app/app.go` 中自动初始化，无需手动配置。初始化后的消息路由器存储在全局变量 `global.MsgRouter` 中，可在项目的任何地方直接使用：

```go
// app.go 中自动执行
func initMessageQueue(db *gorm.DB) {
    if db == nil {
        return
    }

    q, err := queue.InitQueue(cache.RedisClient, db)
    if err != nil {
        global.Logger.Warn("Message queue init failed", zap.Error(err))
        return
    }

    global.MsgRouter = queue.NewMessageRouter(q, db)
    if err := global.MsgRouter.RegisterDefaultHandlers(); err != nil {
        global.Logger.Warn("Register default handlers failed", zap.Error(err))
    } else {
        global.Logger.Info("Message queue initialized successfully")
    }
}
```

### 2. 发布消息

通过全局变量 `global.MsgRouter` 在任意位置发布消息：

```go
import (
    "gin-admin-base/internal/infras/global"
    "gin-admin-base/internal/infras/queue"
)

// 发布操作日志消息
global.MsgRouter.Publish(queue.TopicOperationLog, map[string]interface{}{
    "action":     "create_user",
    "user_id":    1,
    "username":   "admin",
    "ip":         "127.0.0.1",
    "timestamp":  time.Now().Unix(),
})

// 发布通知消息
global.MsgRouter.Publish(queue.TopicNotification, map[string]interface{}{
    "type":    "email",
    "to":      "user@example.com",
    "subject": "欢迎注册",
    "body":    "感谢您注册我们的系统",
})

// 发布系统事件
global.MsgRouter.Publish(queue.TopicSystemEvent, map[string]interface{}{
    "event": "config_changed",
    "key":   "site_name",
    "old":   "旧名称",
    "new":   "新名称",
})
```

### 3. 订阅消息

```go
import "gin-admin-base/internal/infras/global"

// 自定义订阅
global.MsgRouter.Subscribe("custom_topic", func(msg queue.Message) error {
    fmt.Printf("收到消息: ID=%s, Topic=%s, Payload=%+v\n",
        msg.ID, msg.Topic, msg.Payload)
    // 处理业务逻辑...
    return nil // 返回 nil 表示处理成功
})
```

### 4. 自定义消息处理器

在 `router.go` 中注册新的处理器：

```go
// 在 RegisterDefaultHandlers 方法中添加
func (r *MessageRouter) RegisterDefaultHandlers() error {
    // ... 已有处理器 ...

    // 添加自定义处理器
    if err := r.Subscribe("my_topic", handleMyTopic); err != nil {
        return fmt.Errorf("注册 my_topic 处理器失败: %w", err)
    }

    return nil
}

// 实现处理器
func handleMyTopic(msg Message) error {
    // 从 Payload 中提取数据
    userID := msg.Payload["user_id"].(float64)
    action := msg.Payload["action"].(string)

    // 执行业务逻辑
    fmt.Printf("处理用户 %d 的 %s 操作\n", int(userID), action)
    return nil
}
```

## 默认消息处理器

### 1. 操作日志处理器 (`handleOperationLog`)

操作日志消息通过消息队列异步写入数据库，避免同步写入影响接口响应性能：

```go
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
```

### 2. 通知处理器 (`handleNotification`)

```go
func handleNotification(msg Message) error {
    fmt.Printf("[Queue] 通知: %+v\n", msg.Payload)
    // TODO: 实际业务处理，如发送邮件、推送通知等
    return nil
}
```

### 3. 系统事件处理器 (`handleSystemEvent`)

```go
func handleSystemEvent(msg Message) error {
    fmt.Printf("[Queue] 系统事件: %+v\n", msg.Payload)
    // TODO: 实际业务处理，如清理缓存、刷新配置等
    return nil
}
```

## 最佳实践

### 1. 消息幂等性

消息处理器应设计为幂等的，即同一条消息重复处理不会产生副作用：

```go
func handleOrderPaid(msg Message) error {
    orderID := msg.Payload["order_id"].(string)

    // 先检查是否已处理
    var count int64
    db.Model(&ProcessedOrder{}).Where("order_id = ?", orderID).Count(&count)
    if count > 0 {
        return nil // 已处理，跳过
    }

    // 处理订单
    // ...

    // 记录已处理
    db.Create(&ProcessedOrder{OrderID: orderID})
    return nil
}
```

### 2. 错误处理

- 处理器返回 `nil` 表示处理成功，消息将被确认（ACK）
- 处理器返回 `error` 表示处理失败，消息将不会被确认
  - Redis 驱动：消息会保留在 Stream 中，下次重新投递
  - MySQL 驱动：消息状态标记为 `failed`

### 3. 消息超时

对于耗时操作，建议在处理器内部使用超时控制：

```go
func handleLongRunningTask(msg Message) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    done := make(chan error, 1)
    go func() {
        // 执行耗时操作
        done <- processTask(msg.Payload)
    }()

    select {
    case err := <-done:
        return err
    case <-ctx.Done():
        return fmt.Errorf("任务处理超时")
    }
}
```

### 4. 消息大小限制

- Redis Streams：单条消息建议不超过 1MB
- MySQL：受 `text` 字段限制（最大 65KB），建议消息负载不超过 10KB

### 5. 主题命名规范

```
# 推荐格式
<业务模块>:<操作类型>

# 示例
user:created       # 用户创建
user:updated       # 用户更新
order:paid         # 订单支付
order:refunded     # 订单退款
notification:email # 邮件通知
system:config      # 系统配置变更
```

## 性能对比

| 特性 | Redis Streams | MySQL |
|------|---------------|-------|
| 吞吐量 | 高 (10万+/秒) | 中 (1万+/秒) |
| 延迟 | 毫秒级 | 秒级 (轮询) |
| 持久化 | 支持 (RDB/AOF) | 支持 |
| 消息确认 | 支持 (ACK) | 支持 (状态标记) |
| 消息回溯 | 支持 | 支持 |
| 额外依赖 | Redis | 无 |

## 故障排查

### 1. Redis Streams 连接失败

```
错误: Redis 客户端未初始化，无法创建 Redis Streams 队列
```

**解决：**
- 检查 Redis 服务是否正常运行

### 2. MySQL 连接失败

```
错误: 数据库未初始化，无法创建 MySQL 队列
```

**解决：**
- 检查数据库连接配置
- 确保 `queue_messages` 表已自动创建

### 3. 消息堆积

如果消息消费速度跟不上生产速度：

1. **Redis 驱动**：增加消费者数量（多个实例订阅同一消费组）
2. **MySQL 驱动**：缩短轮询间隔（修改 `time.Sleep` 时间）

### 4. 消息重复消费

- Redis 驱动：消费者宕机重启后，未 ACK 的消息会被重新投递
- MySQL 驱动：应用重启后，`pending` 状态的消息会被重新消费

**建议：** 处理器实现幂等性设计

## 扩展指南

### 添加新的队列驱动

1. 在 `message.go` 中定义 `Queue` 接口（已完成）
2. 创建新的驱动文件，实现 `Queue` 接口
3. 在 `factory.go` 的 `InitQueue` 函数中添加新的 case

```go
// factory.go
case QueueTypeKafka:
    return NewKafkaQueue(brokers), nil
```

### 添加新的预定义主题

在 `router.go` 中添加主题常量和默认处理器：

```go
const (
    TopicEmailNotification = "email:notification"
)

func (r *MessageRouter) RegisterDefaultHandlers() error {
    // ... 已有处理器 ...
    if err := r.Subscribe(TopicEmailNotification, handleEmailNotification); err != nil {
        return fmt.Errorf("注册邮件通知处理器失败: %w", err)
    }
    return nil
}
