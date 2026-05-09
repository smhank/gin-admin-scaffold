package persistence

import (
	"fmt"
	"time"

	"gin-admin-base/internal/infras/config"
	migrationlib "gin-admin-base/internal/infras/migration"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitDB() (*gorm.DB, error) {
	mysqlConfig := config.GetMySQLConfig()

	// 构建 GORM 配置
	gormConfig := &gorm.Config{
		// 跳过默认事务，非事务查询不包裹事务，提升约 30% 性能
		SkipDefaultTransaction: mysqlConfig.SkipDefaultTx,
		// 缓存预编译语句，减少 SQL 解析开销
		PrepareStmt: mysqlConfig.PrepareStmt,
		// 命名策略
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false, // 使用复数表名
		},
	}

	// 打开数据库连接
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       mysqlConfig.DSN(),
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式
		DontSupportRenameColumn:   true,  // 重命名列时采用删除并新建的方式
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), gormConfig)

	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层 sql.DB 实例，配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(mysqlConfig.MaxIdleConns)       // 最大空闲连接数
	sqlDB.SetMaxOpenConns(mysqlConfig.MaxOpenConns)       // 最大打开连接数
	sqlDB.SetConnMaxLifetime(mysqlConfig.MaxLifetime)     // 连接最大存活时间
	sqlDB.SetConnMaxIdleTime(mysqlConfig.ConnMaxIdleTime) // 空闲连接最大存活时间

	// 运行迁移（建表 + 数据初始化，类似 Yii2 的 migrate/up）
	manager := migrationlib.NewManager(db)
	registerMigrations(manager)
	if err := manager.Up(); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %w", err)
	}

	return db, nil
}

// registerMigrations 注册所有迁移
func registerMigrations(manager *migrationlib.Manager) {
	manager.RegisterAll([]migrationlib.Migration{
		migrationlib.NewM20250506000000InitTables(),          // 建表
		migrationlib.NewM20250506000001InitSeedData(),        // 初始化种子数据
		migrationlib.NewM20250506000002AddOperationLogMenu(), // 添加操作历史菜单
	})
}

// 用于避免 time 包导入未使用的警告
var _ = time.Now
