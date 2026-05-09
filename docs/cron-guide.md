# 定时任务使用指南

## 概述

本项目内置了轻量级的定时任务调度器，支持基于时间间隔的任务调度。无需额外依赖，适用于周期性的后台任务，如数据清理、心跳检测、状态同步等。

## 架构设计

```
┌─────────────────┐
│  调度器         │
│  (Scheduler)    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐     ┌─────────────────┐
│  任务 1         │     │  任务 2         │
│  (Task)         │     │  (Task)         │
│  - clean_logs   │     │  - heartbeat    │
│  - @every 24h   │     │  - @every 5m    │
└─────────────────┘     └─────────────────┘
```

### 核心组件

| 组件 | 文件 | 说明 |
|------|------|------|
| `Task` 接口 | `task.go` | 定义 Name、Spec、Run 三个方法 |
| `FuncTask` 结构体 | `task.go` | 基于函数的任务实现 |
| `Scheduler` 结构体 | `manager.go` | 任务调度器，管理任务的启动/停止 |
| `RegisterDefaultTasks` | `tasks.go` | 注册默认定时任务 |

## 快速开始

### 1. 初始化调度器

定时任务在 `main.go` 中自动初始化，无需手动配置：

```go
// main.go 中自动执行
scheduler := cron.NewScheduler(global.Logger)
cron.RegisterDefaultTasks(scheduler)
scheduler.Start()
defer scheduler.Stop()
```

### 2. 创建自定义任务

```go
package mytask

import (
    "context"
    "fmt"
    "time"

    "gin-admin-base/internal/infras/cron"
)

// 方式一：使用 Task 接口
type MyTask struct{}

func (t *MyTask) Name() string { return "my_task" }
func (t *MyTask) Spec() string { return "@every 10s" }
func (t *MyTask) Run(ctx context.Context) error {
    fmt.Println("执行自定义任务")
    return nil
}

// 方式二：使用 FuncTask（推荐）
func RegisterMyTask(scheduler *cron.Scheduler) {
    scheduler.AddTask(cron.NewFuncTask(
        "my_task",           // 任务名称
        "@every 10s",        // 执行间隔
        func(ctx context.Context) error {
            fmt.Println("执行自定义任务")
            return nil
        },
    ))
}
```

### 3. 注册任务

在 `main.go` 中注册：

```go
import "your-project/internal/infras/cron"

scheduler := cron.NewScheduler(global.Logger)

// 注册默认任务
cron.RegisterDefaultTasks(scheduler)

// 注册自定义任务
RegisterMyTask(scheduler)

scheduler.Start()
```

## 定时表达式

### 支持的格式

| 格式 | 说明 | 示例 |
|------|------|------|
| `@every <duration>` | 每隔指定时间执行 | `@every 5s`, `@every 1h` |
| `<duration>` | 直接使用 time.Duration | `5s`, `1m`, `30s` |

### 常用间隔

| 表达式 | 说明 |
|--------|------|
| `@every 30s` | 每 30 秒 |
| `@every 1m` | 每 1 分钟 |
| `@every 5m` | 每 5 分钟 |
| `@every 1h` | 每 1 小时 |
| `@every 24h` | 每 24 小时（每天） |

## 默认定时任务

| 任务名称 | 间隔 | 说明 |
|----------|------|------|
| `clean_expired_logs` | 每 24 小时 | 清理 30 天前的操作日志 |
| `system_heartbeat` | 每 5 分钟 | 发送系统心跳事件到消息队列 |
| `check_db_connection` | 每 30 秒 | 检查数据库连接是否正常 |

### 1. 清理过期日志 (`clean_expired_logs`)

```go
func CleanExpiredLogs(ctx context.Context) error {
    // 删除 30 天前的操作日志
    cutoff := time.Now().AddDate(0, 0, -30)
    result := global.DB.Exec("DELETE FROM operation_logs WHERE created_at < ?", cutoff)
    global.Logger.Info("清理过期操作日志",
        zap.Int64("deleted_count", result.RowsAffected),
    )
    return nil
}
```

### 2. 系统心跳 (`system_heartbeat`)

```go
func SystemHeartbeat(ctx context.Context) error {
    // 通过消息队列发送心跳事件
    return global.MsgRouter.Publish(queue.TopicSystemEvent, map[string]interface{}{
        "event":     "heartbeat",
        "timestamp": time.Now().Unix(),
        "status":    "running",
    })
}
```

### 3. 数据库连接检查 (`check_db_connection`)

```go
func CheckDBConnection(ctx context.Context) error {
    sqlDB, err := global.DB.DB()
    if err != nil {
        return fmt.Errorf("获取数据库实例失败: %w", err)
    }
    return sqlDB.PingContext(ctx)
}
```

## 最佳实践

### 1. 错误处理

任务返回 `error` 会被调度器记录日志，但不会影响后续执行：

```go
func MyTask(ctx context.Context) error {
    if err := doSomething(); err != nil {
        // 记录错误，任务会继续在下一个周期执行
        return fmt.Errorf("执行失败: %w", err)
    }
    return nil
}
```

### 2. Panic 恢复

调度器会自动捕获任务中的 panic，防止整个程序崩溃：

```go
func MyTask(ctx context.Context) error {
    // 即使这里 panic，调度器也会恢复并记录日志
    panic("意外错误")
}
```

### 3. 任务幂等性

定时任务应设计为幂等的，多次执行不会产生副作用：

```go
func CleanExpiredLogs(ctx context.Context) error {
    // 使用 WHERE 条件确保只删除过期数据
    // 多次执行是安全的
    global.DB.Where("created_at < ?", cutoff).Delete(&OperationLog{})
    return nil
}
```

### 4. 资源清理

在 `main.go` 中使用 `defer` 确保程序退出时停止所有任务：

```go
scheduler := cron.NewScheduler(global.Logger)
scheduler.Start()
defer scheduler.Stop() // 程序退出时自动停止
```

### 5. 避免任务重叠

如果任务执行时间可能超过间隔时间，建议在任务内部加锁：

```go
var mu sync.Mutex

func MyTask(ctx context.Context) error {
    if !mu.TryLock() {
        global.Logger.Warn("任务正在执行中，跳过本次")
        return nil
    }
    defer mu.Unlock()

    // 执行耗时操作
    time.Sleep(10 * time.Second)
    return nil
}
```

## 调度器 API

### 创建调度器

```go
scheduler := cron.NewScheduler(logger)
```

### 添加任务

```go
scheduler.AddTask(task)
```

### 启动所有任务

```go
err := scheduler.Start()
```

### 停止所有任务

```go
scheduler.Stop()
```

### 检查运行状态

```go
if scheduler.IsRunning() {
    fmt.Println("调度器正在运行")
}
```

## 扩展指南

### 添加新的默认定时任务

在 `tasks.go` 的 `RegisterDefaultTasks` 函数中添加：

```go
func RegisterDefaultTasks(scheduler *Scheduler) {
    // ... 已有任务 ...

    // 添加新任务
    scheduler.AddTask(NewFuncTask(
        "sync_user_status",    // 任务名称
        "@every 1h",           // 每小时执行一次
        SyncUserStatus,        // 任务函数
    ))
}

// 实现任务函数
func SyncUserStatus(ctx context.Context) error {
    // 同步用户状态逻辑
    return nil
}
```

### 创建带参数的任务

```go
func NewCleanLogsTask(retentionDays int) *FuncTask {
    return NewFuncTask(
        "clean_logs",
        "@every 24h",
        func(ctx context.Context) error {
            cutoff := time.Now().AddDate(0, 0, -retentionDays)
            // 清理逻辑
            return nil
        },
    )
}

// 使用
scheduler.AddTask(NewCleanLogsTask(30)) // 保留 30 天
scheduler.AddTask(NewCleanLogsTask(7))  // 保留 7 天
