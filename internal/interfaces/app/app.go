package app

import (
	"gin-admin-base/internal/application"
	"gin-admin-base/internal/infras/cache"
	"gin-admin-base/internal/infras/config"
	"gin-admin-base/internal/infras/cron"
	"gin-admin-base/internal/infras/global"
	"gin-admin-base/internal/infras/logger"
	"gin-admin-base/internal/infras/persistence"
	"gin-admin-base/internal/infras/queue"
	"gin-admin-base/internal/interfaces/handler"
	"gin-admin-base/internal/interfaces/middleware"
	"gin-admin-base/internal/interfaces/router"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// App 应用主结构体
type App struct {
	Router *gin.Engine
}

// New 创建并初始化应用
func New() *App {
	// 1. 初始化配置
	config.InitConfig()

	// 2. 设置运行模式
	global.AppMode = viper.GetString("server.mode")
	if global.AppMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 3. 初始化日志
	global.Logger = logger.InitLogger(config.GetLoggerConfig())

	// 4. 初始化基础设施
	db := initDatabase()
	initRedis()
	initMessageQueue(db)
	scheduler := initCron()
	if scheduler != nil {
		defer scheduler.Stop()
	}

	// 5. 初始化 Handler 和 AuthService
	handlers, authSvc := initHandlers(db)

	// 6. 创建 Gin 引擎并注册中间件
	r := gin.Default()
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.CORSMiddleware())

	// 7. 注册路由
	router.New(r, handlers, authSvc).Register()

	return &App{Router: r}
}

// Run 启动 HTTP 服务器
func (a *App) Run() error {
	port := viper.GetString("server.port")
	if port == "" {
		port = "9501"
	}
	addr := "0.0.0.0:" + port
	global.Logger.Info("Starting server", zap.String("addr", addr))

	if err := a.Router.Run(addr); err != nil {
		global.Logger.Fatal("Server failed to start", zap.Error(err))
		return err
	}
	return nil
}

// Shutdown 关闭应用，释放资源
func (a *App) Shutdown() {
	if global.MsgRouter != nil {
		global.MsgRouter.Close()
	}
	global.Logger.Sync()
}

// initDatabase 初始化数据库
func initDatabase() *gorm.DB {
	db, err := persistence.InitDB()
	if err != nil {
		global.Logger.Warn("Database init failed, using mock auth", zap.Error(err))
		return nil
	}
	global.DB = db
	return db
}

// initRedis 初始化 Redis
func initRedis() {
	cache.InitRedis()
}

// initMessageQueue 初始化消息队列
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

// initCron 初始化定时任务
func initCron() *cron.Scheduler {
	if global.DB == nil {
		return nil
	}

	scheduler := cron.NewScheduler(global.Logger)
	cron.RegisterDefaultTasks(scheduler)
	if err := scheduler.Start(); err != nil {
		global.Logger.Warn("Cron scheduler start failed", zap.Error(err))
		return nil
	}
	global.Logger.Info("Cron scheduler started successfully")
	return scheduler
}

// initHandlers 初始化所有 Handler 和 AuthService
func initHandlers(db *gorm.DB) (*handler.Handlers, *application.AuthService) {
	if db == nil {
		return &handler.Handlers{
			Auth:         handler.NewAuthHandler(nil),
			Permission:   handler.NewPermissionHandler(nil),
			Role:         handler.NewRoleHandler(nil),
			Menu:         handler.NewMenuHandler(nil),
			Path:         handler.NewPathHandler(nil),
			User:         handler.NewUserHandler(nil),
			OperationLog: handler.NewOperationLogHandler(nil),
			Migration:    handler.NewMigrationHandler(nil),
		}, nil
	}

	authRepo := persistence.NewAuthRepo(db)
	authSvc := application.NewAuthService(authRepo)

	return &handler.Handlers{
		Auth:         handler.NewAuthHandler(authSvc),
		Permission:   handler.NewPermissionHandler(db),
		Role:         handler.NewRoleHandler(db),
		Menu:         handler.NewMenuHandler(db),
		Path:         handler.NewPathHandler(db),
		User:         handler.NewUserHandler(db),
		OperationLog: handler.NewOperationLogHandler(db),
		Migration:    handler.NewMigrationHandler(db),
	}, authSvc
}
