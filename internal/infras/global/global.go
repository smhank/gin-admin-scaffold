package global

import (
	"gin-admin-base/internal/infras/queue"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB        *gorm.DB
	Redis     *redis.Client
	Logger    *zap.Logger
	MsgRouter *queue.MessageRouter
	AppMode   string // debug / release
)
