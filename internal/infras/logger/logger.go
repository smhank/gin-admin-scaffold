package logger

import (
	"os"

	"gin-admin-base/internal/infras/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger 根据配置初始化日志
func InitLogger(cfg *config.LoggerConfig) *zap.Logger {
	// 设置日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// dev 模式使用可读性更好的编码
	var encoder zapcore.Encoder
	if cfg.Mode == "dev" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 构建日志输出目标
	var cores []zapcore.Core

	// 控制台输出
	consoleWriter := zapcore.AddSync(os.Stdout)
	cores = append(cores, zapcore.NewCore(encoder, consoleWriter, level))

	// 文件输出（如果配置了文件路径）
	if cfg.FilePath != "" {
		// 使用 lumberjack 实现日志轮转
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,    // MB
			MaxBackups: cfg.MaxBackups, // 保留旧文件数
			MaxAge:     cfg.MaxAge,     // 保留天数
			Compress:   cfg.Compress,   // 压缩归档
		})
		cores = append(cores, zapcore.NewCore(encoder, fileWriter, level))
	}

	// 合并所有输出目标
	core := zapcore.NewTee(cores...)

	// 构建 Logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return logger
}
