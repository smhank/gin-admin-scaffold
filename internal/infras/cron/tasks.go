package cron

import (
	"context"
	"fmt"
	"time"

	"gin-admin-base/internal/infras/global"
	"gin-admin-base/internal/infras/queue"

	"go.uber.org/zap"
)

// RegisterDefaultTasks 注册默认定时任务
func RegisterDefaultTasks(scheduler *Scheduler) {
	// 清理过期操作日志（每天凌晨 2 点执行）
	scheduler.AddTask(NewFuncTask(
		"clean_expired_logs",
		"@every 24h",
		CleanExpiredLogs,
	))

	// 发送系统心跳（每 5 分钟执行一次）
	scheduler.AddTask(NewFuncTask(
		"system_heartbeat",
		"@every 5m",
		SystemHeartbeat,
	))

	// 检查数据库连接（每 30 秒执行一次）
	scheduler.AddTask(NewFuncTask(
		"check_db_connection",
		"@every 30s",
		CheckDBConnection,
	))
}

// CleanExpiredLogs 清理过期操作日志（保留最近 30 天）
func CleanExpiredLogs(ctx context.Context) error {
	if global.DB == nil {
		return nil
	}

	cutoff := time.Now().AddDate(0, 0, -30)

	result := global.DB.Exec("DELETE FROM operation_logs WHERE created_at < ?", cutoff)
	if result.Error != nil {
		return fmt.Errorf("清理过期日志失败: %w", result.Error)
	}

	global.Logger.Info("清理过期操作日志",
		zap.Int64("deleted_count", result.RowsAffected),
		zap.Time("cutoff", cutoff),
	)

	return nil
}

// SystemHeartbeat 发送系统心跳
func SystemHeartbeat(ctx context.Context) error {
	if global.MsgRouter == nil {
		return nil
	}

	return global.MsgRouter.Publish(queue.TopicSystemEvent, map[string]interface{}{
		"event":     "heartbeat",
		"timestamp": time.Now().Unix(),
		"status":    "running",
	})
}

// CheckDBConnection 检查数据库连接
func CheckDBConnection(ctx context.Context) error {
	if global.DB == nil {
		return nil
	}

	sqlDB, err := global.DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		global.Logger.Error("数据库连接检查失败", zap.Error(err))
		return err
	}

	global.Logger.Debug("数据库连接正常")
	return nil
}
