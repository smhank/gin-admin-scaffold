package cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	tasks   []Task
	cancels []context.CancelFunc
	mu      sync.Mutex
	running bool
	logger  *zap.Logger
}

// NewScheduler 创建定时任务调度器
func NewScheduler(logger *zap.Logger) *Scheduler {
	return &Scheduler{
		tasks:   make([]Task, 0),
		cancels: make([]context.CancelFunc, 0),
		logger:  logger,
	}
}

// AddTask 添加定时任务
func (s *Scheduler) AddTask(task Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = append(s.tasks, task)
}

// Start 启动所有定时任务
func (s *Scheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("定时任务调度器已在运行")
	}

	s.running = true

	for _, task := range s.tasks {
		if err := s.startTask(task); err != nil {
			s.logger.Warn("启动定时任务失败",
				zap.String("task", task.Name()),
				zap.Error(err),
			)
			continue
		}
		s.logger.Info("定时任务已启动",
			zap.String("task", task.Name()),
			zap.String("spec", task.Spec()),
		)
	}

	return nil
}

// startTask 启动单个定时任务
func (s *Scheduler) startTask(task Task) error {
	interval, err := ParseSpec(task.Spec())
	if err != nil {
		return fmt.Errorf("解析定时表达式失败: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancels = append(s.cancels, cancel)

	go func() {
		// 首次执行等待一个间隔，让系统先完成初始化
		time.Sleep(interval)

		for {
			select {
			case <-ctx.Done():
				s.logger.Info("定时任务已停止",
					zap.String("task", task.Name()),
				)
				return
			default:
				s.executeTask(ctx, task)
				time.Sleep(interval)
			}
		}
	}()

	return nil
}

// executeTask 执行单个任务（带 panic 恢复）
func (s *Scheduler) executeTask(ctx context.Context, task Task) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("定时任务执行 panic",
				zap.String("task", task.Name()),
				zap.Any("panic", r),
			)
		}
	}()

	startTime := time.Now()
	s.logger.Debug("定时任务开始执行",
		zap.String("task", task.Name()),
	)

	if err := task.Run(ctx); err != nil {
		s.logger.Error("定时任务执行失败",
			zap.String("task", task.Name()),
			zap.Duration("elapsed", time.Since(startTime)),
			zap.Error(err),
		)
		return
	}

	s.logger.Debug("定时任务执行完成",
		zap.String("task", task.Name()),
		zap.Duration("elapsed", time.Since(startTime)),
	)
}

// Stop 停止所有定时任务
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	for _, cancel := range s.cancels {
		cancel()
	}

	s.cancels = make([]context.CancelFunc, 0)
	s.running = false

	s.logger.Info("所有定时任务已停止")
}

// IsRunning 检查调度器是否在运行
func (s *Scheduler) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}
