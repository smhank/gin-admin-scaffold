package cron

import (
	"context"
	"fmt"
	"time"
)

// Task 定时任务接口
type Task interface {
	// Name 任务名称
	Name() string
	// Spec 定时表达式 (cron 格式: 秒 分 时 日 月 周)
	Spec() string
	// Run 执行任务
	Run(ctx context.Context) error
}

// TaskFunc 任务函数类型
type TaskFunc func(ctx context.Context) error

// FuncTask 基于函数的定时任务
type FuncTask struct {
	name string
	spec string
	fn   TaskFunc
}

// NewFuncTask 创建函数任务
func NewFuncTask(name, spec string, fn TaskFunc) *FuncTask {
	return &FuncTask{
		name: name,
		spec: spec,
		fn:   fn,
	}
}

func (t *FuncTask) Name() string                  { return t.name }
func (t *FuncTask) Spec() string                  { return t.spec }
func (t *FuncTask) Run(ctx context.Context) error { return t.fn(ctx) }

// ParseSpec 解析 cron 表达式为间隔
// 支持格式: "@every 5s", "@every 1m", "@every 1h"
// 也支持简单格式: "5s", "1m", "1h"
func ParseSpec(spec string) (time.Duration, error) {
	// 如果已经是 @every 格式
	if len(spec) > 7 && spec[:7] == "@every " {
		return time.ParseDuration(spec[7:])
	}

	// 尝试直接解析为 duration
	if d, err := time.ParseDuration(spec); err == nil {
		return d, nil
	}

	return 0, fmt.Errorf("无法解析定时表达式: %s", spec)
}
